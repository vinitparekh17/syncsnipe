package sync

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gobwas/glob"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/types"
)

const configReloadInterval = 5 * time.Second

type ruleRuntime struct {
	id, profileID  int64
	source, target string
	updatedAt      int64
	initialMode    string
	ignore         []func(string) bool
	dirs           map[string]struct{}
}

type debounceEntry struct {
	timer *time.Timer
	event fsnotify.Event
}

// SyncWatcher owns fsnotify intake, rule configuration, debounce timers, and workers.
type SyncWatcher struct {
	watcher              *fsnotify.Watcher
	worker               *SyncWorker
	db                   *database.Queries
	syncQueue            chan *SyncOperation
	mu                   sync.Mutex
	rules                map[int64]*ruleRuntime
	debounce             map[string]*debounceEntry
	ctx                  context.Context
	cancel               context.CancelFunc
	startOnce, closeOnce sync.Once
	wg                   sync.WaitGroup
}

func NewSyncWatcher(db *database.Queries) (*SyncWatcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	queue := make(chan *SyncOperation, 4096)
	return &SyncWatcher{watcher: fw, db: db, syncQueue: queue, rules: make(map[int64]*ruleRuntime), debounce: make(map[string]*debounceEntry), worker: mustWorker(db, queue)}, nil
}

func mustWorker(db *database.Queries, q chan *SyncOperation) *SyncWorker {
	w, _ := NewSyncWorker(db, q, runtime.NumCPU())
	return w
}

func (sw *SyncWatcher) Start(parent context.Context) {
	sw.startOnce.Do(func() {
		sw.ctx, sw.cancel = context.WithCancel(parent)
		sw.worker.Start(sw.ctx)
		if err := sw.reloadRules(true); err != nil {
			colorlog.Error("initial rule load failed: %v", err)
		}
		sw.wg.Add(2)
		go func() { defer sw.wg.Done(); sw.eventLoop() }()
		go func() { defer sw.wg.Done(); sw.configLoop() }()
		colorlog.Success("SyncWatcher started successfully")
	})
}

func (sw *SyncWatcher) Close() {
	sw.closeOnce.Do(func() {
		if sw.cancel != nil {
			sw.cancel()
		}
		sw.mu.Lock()
		for _, entry := range sw.debounce {
			entry.timer.Stop()
		}
		sw.debounce = map[string]*debounceEntry{}
		sw.mu.Unlock()
		if err := sw.watcher.Close(); err != nil {
			colorlog.Error("closing filesystem watcher: %v", err)
		}
		sw.wg.Wait()
		close(sw.syncQueue)
		sw.worker.Wait()
		colorlog.Success("SyncWatcher closed successfully")
	})
}

func (sw *SyncWatcher) eventLoop() {
	for {
		select {
		case <-sw.ctx.Done():
			return
		case event, ok := <-sw.watcher.Events:
			if !ok {
				return
			}
			sw.debounceEvent(event)
		case err, ok := <-sw.watcher.Errors:
			if !ok {
				return
			}
			colorlog.Error("filesystem watcher error: %v; scheduling reconciliation", err)
			go sw.reconcileAll()
		}
	}
}

func (sw *SyncWatcher) configLoop() {
	ticker := time.NewTicker(configReloadInterval)
	defer ticker.Stop()
	for {
		select {
		case <-sw.ctx.Done():
			return
		case <-ticker.C:
			if err := sw.reloadRules(false); err != nil {
				colorlog.Error("reloading sync rules: %v", err)
			}
		}
	}
}

func (sw *SyncWatcher) debounceEvent(event fsnotify.Event) {
	rule, ok := sw.ruleForPath(event.Name)
	if !ok || sw.isIgnored(rule, event.Name) {
		return
	}

	// Immediately watch newly created directories so files created inside them are not missed
	if event.Op&(fsnotify.Create|fsnotify.Write) != 0 {
		if info, err := os.Lstat(event.Name); err == nil && info.IsDir() {
			if err := sw.addTree(rule, event.Name); err != nil {
				colorlog.Error("watching new directory %s: %v", event.Name, err)
			}
		}
	}

	key := fmt.Sprintf("%d:%s", rule.id, event.Name)
	sw.mu.Lock()
	defer sw.mu.Unlock()
	if old := sw.debounce[key]; old != nil {
		old.event = event
		old.timer.Reset(DebounceTime)
		return
	}
	entry := &debounceEntry{event: event}
	entry.timer = time.AfterFunc(DebounceTime, func() { sw.fireDebounced(key, rule, entry) })
	sw.debounce[key] = entry
}

func (sw *SyncWatcher) fireDebounced(key string, rule *ruleRuntime, entry *debounceEntry) {
	sw.mu.Lock()
	if sw.debounce[key] != entry {
		sw.mu.Unlock()
		return
	}
	event := entry.event
	delete(sw.debounce, key)
	sw.mu.Unlock()

	if event.Op&fsnotify.Create != 0 {
		if info, err := os.Lstat(event.Name); err == nil && info.IsDir() {
			if err := sw.addTree(rule, event.Name); err != nil {
				colorlog.Error("watching new directory %s: %v", event.Name, err)
			}
			go sw.reconcileRule(rule)
			return
		}
	}
	sw.enqueue(rule, event.Name)
}

func (sw *SyncWatcher) enqueue(rule *ruleRuntime, source string) {
	rel, err := filepath.Rel(rule.source, source)
	if err != nil || rel == "." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return
	}
	op := &SyncOperation{profileID: rule.profileID, sourcePath: source, targetPath: filepath.Join(rule.target, rel), operation: CreateOrModifyEvent, timeStamp: time.Now()}
	select {
	case <-sw.ctx.Done():
		return
	case sw.syncQueue <- op:
	case <-time.After(100 * time.Millisecond):
		colorlog.Warn("sync queue full for rule %d; scheduling reconciliation", rule.id)
		go sw.reconcileRule(rule)
	}
}

func (sw *SyncWatcher) enqueueBlocking(rule *ruleRuntime, source string) {
	rel, err := filepath.Rel(rule.source, source)
	if err != nil || rel == "." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return
	}
	op := &SyncOperation{profileID: rule.profileID, sourcePath: source, targetPath: filepath.Join(rule.target, rel), operation: CreateOrModifyEvent, timeStamp: time.Now()}
	select {
	case <-sw.ctx.Done():
		return
	case sw.syncQueue <- op:
	}
}

func (sw *SyncWatcher) reloadRules(initial bool) error {
	dbRules, err := sw.db.ListAllActiveRules(sw.ctx)
	if err != nil {
		return err
	}
	desired := make(map[int64]database.SyncRule)
	for _, rule := range dbRules {
		desired[rule.ID] = rule
	}
	sw.mu.Lock()
	current := make(map[int64]*ruleRuntime, len(sw.rules))
	for id, r := range sw.rules {
		current[id] = r
	}
	sw.mu.Unlock()
	for id, active := range current {
		wanted, ok := desired[id]
		if !ok || wanted.SourceDir != active.source || wanted.TargetDir != active.target || wanted.UpdatedAt != active.updatedAt {
			sw.removeRule(active)
		}
	}
	for id, wanted := range desired {
		sw.mu.Lock()
		_, active := sw.rules[id]
		sw.mu.Unlock()
		if !active {
			if err := sw.addRule(wanted, initial); err != nil {
				colorlog.Error("activating rule %d: %v", id, err)
			}
		} else {
			sw.refreshIgnore(id, wanted.ProfileID)
		}
	}
	return nil
}

func (sw *SyncWatcher) addRule(rule database.SyncRule, initial bool) error {
	source, err := filepath.Abs(filepath.Clean(rule.SourceDir))
	if err != nil {
		return err
	}
	target, err := filepath.Abs(filepath.Clean(rule.TargetDir))
	if err != nil {
		return err
	}
	if source == target || isWithin(source, target) || isWithin(target, source) {
		return fmt.Errorf("source and target directories must not overlap")
	}
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("source is not a directory: %s", source)
	}
	rt := &ruleRuntime{id: rule.ID, profileID: rule.ProfileID, source: source, target: target, updatedAt: rule.UpdatedAt, initialMode: rule.InitialSyncMode, dirs: make(map[string]struct{})}
	if err := sw.populateIgnore(rt); err != nil {
		return err
	}
	if err := sw.addTree(rt, source); err != nil {
		return err
	}
	sw.mu.Lock()
	sw.rules[rt.id] = rt
	sw.mu.Unlock()
	if initial && rt.initialMode != "events_only" {
		go sw.reconcileRule(rt)
	}
	return nil
}

func (sw *SyncWatcher) removeRule(rule *ruleRuntime) {
	sw.mu.Lock()
	delete(sw.rules, rule.id)
	for d := range rule.dirs {
		_ = sw.watcher.Remove(d)
	}
	sw.mu.Unlock()
}

func (sw *SyncWatcher) refreshIgnore(id, profileID int64) {
	sw.mu.Lock()
	r := sw.rules[id]
	sw.mu.Unlock()
	if r != nil {
		if err := sw.populateIgnore(r); err != nil {
			colorlog.Error("loading ignores for rule %d: %v", id, err)
		}
	}
}

func (sw *SyncWatcher) populateIgnore(rule *ruleRuntime) error {
	patterns, err := sw.db.ListIgnorePattern(sw.ctx, rule.profileID)
	if err != nil {
		return err
	}
	matchers := make([]func(string) bool, 0, len(patterns))
	for _, p := range patterns {
		m, err := compileIgnore(p)
		if err != nil {
			colorlog.Warn("invalid ignore pattern %q: %v", p.Pattern, err)
			continue
		}
		matchers = append(matchers, m)
	}
	sw.mu.Lock()
	rule.ignore = matchers
	sw.mu.Unlock()
	return nil
}
func compileIgnore(p database.IgnorePattern) (func(string) bool, error) {
	switch p.Type {
	case types.Exact:
		return func(path string) bool {
			return filepath.ToSlash(path) == filepath.ToSlash(p.Pattern) || filepath.Base(path) == p.Pattern
		}, nil
	case types.Regex:
		r, err := regexp.Compile(p.Pattern)
		if err != nil {
			return nil, err
		}
		return r.MatchString, nil
	default:
		g, err := glob.Compile(filepath.ToSlash(p.Pattern))
		if err != nil {
			return nil, err
		}
		return func(path string) bool { return g.Match(filepath.ToSlash(path)) || g.Match(filepath.Base(path)) }, nil
	}
}
func (sw *SyncWatcher) isIgnored(r *ruleRuntime, path string) bool {
	rel, err := filepath.Rel(r.source, path)
	if err != nil {
		return true
	}
	sw.mu.Lock()
	matchers := append([]func(string) bool(nil), r.ignore...)
	sw.mu.Unlock()
	for _, match := range matchers {
		if match(rel) {
			return true
		}
	}
	return false
}
func (sw *SyncWatcher) ruleForPath(path string) (*ruleRuntime, bool) {
	abs, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return nil, false
	}
	sw.mu.Lock()
	defer sw.mu.Unlock()
	var best *ruleRuntime
	for _, r := range sw.rules {
		if (abs == r.source || isWithin(r.source, abs)) && (best == nil || len(r.source) > len(best.source)) {
			best = r
		}
	}
	return best, best != nil
}

func isWithin(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != "." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) && rel != ".."
}

func (sw *SyncWatcher) addTree(rule *ruleRuntime, root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type()&os.ModeSymlink != 0 {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() || sw.isIgnored(rule, path) {
			return nil
		}
		sw.mu.Lock()
		_, exists := rule.dirs[path]
		sw.mu.Unlock()
		if exists {
			return nil
		}
		if err := sw.watcher.Add(path); err != nil {
			return err
		}
		sw.mu.Lock()
		rule.dirs[path] = struct{}{}
		sw.mu.Unlock()
		return nil
	})
}

func (sw *SyncWatcher) reconcileAll() {
	sw.mu.Lock()
	rules := make([]*ruleRuntime, 0, len(sw.rules))
	for _, r := range sw.rules {
		rules = append(rules, r)
	}
	sw.mu.Unlock()
	for _, r := range rules {
		sw.reconcileRule(r)
	}
}

func (sw *SyncWatcher) reconcileRule(rule *ruleRuntime) {
	if err := os.MkdirAll(rule.target, 0755); err != nil {
		colorlog.Error("creating target root %s: %v", rule.target, err)
		return
	}
	seen := map[string]struct{}{}
	err := filepath.WalkDir(rule.source, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == rule.source {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		rel, _ := filepath.Rel(rule.source, path)
		if sw.isIgnored(rule, path) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		seen[rel] = struct{}{}
		if d.IsDir() {
			_ = sw.addTree(rule, path)
		}
		sw.enqueueBlocking(rule, path)
		return nil
	})
	if err != nil {
		colorlog.Error("reconciling source %s: %v", rule.source, err)
		return
	}
	var stale []string
	_ = filepath.WalkDir(rule.target, func(path string, d os.DirEntry, err error) error {
		if err != nil || path == rule.target {
			return nil
		}
		rel, _ := filepath.Rel(rule.target, path)
		if _, ok := seen[rel]; !ok {
			stale = append(stale, path)
		}
		return nil
	})
	sort.Slice(stale, func(i, j int) bool { return len(stale[i]) > len(stale[j]) })
	for _, path := range stale {
		if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
			colorlog.Error("removing stale target %s: %v", path, err)
		}
	}
}
