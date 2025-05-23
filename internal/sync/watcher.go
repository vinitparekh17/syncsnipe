package sync

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

type SyncWatcher struct {
	// Core components
	watcher   *fsnotify.Watcher
	worker    *SyncWorker
	db        *database.Queries
	syncQueue chan *SyncOperation

	// State management
	paths       map[string]bool
	fileHashMap map[string]string  // for tracking file hashes
	ignoreList  map[int64][]string // list per profileId
	mu          sync.Mutex
	wg          sync.WaitGroup
}

func NewSyncWatcher(db *database.Queries) (*SyncWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	sw := &SyncWatcher{
		watcher:     watcher,
		paths:       make(map[string]bool),
		syncQueue:   make(chan *SyncOperation, 100),
		db:          db,
		fileHashMap: make(map[string]string),
	}

	worker, err := NewSyncWorker(db, sw.syncQueue, runtime.NumCPU())
	if err != nil {
		return nil, err
	}

	sw.worker = worker
	return sw, nil
}

func (sw *SyncWatcher) Start(ctx context.Context) {
	sw.wg.Add(1)
	go func() {
		defer sw.wg.Done()
		sw.eventLoop(ctx)
	}()
	colorlog.Success("SyncWatcher started successfully.")
}

func (sw *SyncWatcher) Close() {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if err := sw.watcher.Close(); err != nil {
		colorlog.Error("failed to close fsnotify watcher: %v", err)
	}

	sw.worker.close()

	done := make(chan struct{})
	go func() {
		sw.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		colorlog.Success("All workers exited successfully")
	case <-time.After(5 * time.Second):
		colorlog.Warn("Timeout waiting for workers to exit - forcing shutdown")
	}

	colorlog.Success("SyncWatcher closed successfully.")
}

func (sw *SyncWatcher) eventLoop(ctx context.Context) {
	debounceHandler := sw.debounce(DebounceTime, sw.handleEvent)
	for {
		select {
		case event, ok := <-sw.watcher.Events:
			if !ok {
				return
			}
			debounceHandler(event)
		case err, ok := <-sw.watcher.Errors:
			if !ok {
				return
			}
			colorlog.Error("watcher err in event loop: %v", err)
			return
		case <-ctx.Done():
			colorlog.Info("SyncWatcher event loop terminated.")
			return
		}
	}
}

// debounce creates a debounced version of the given handler function.
// It ensures that the handler is called at most once within the specified duration.
func (sw *SyncWatcher) debounce(debounceTime time.Duration, handle func(event fsnotify.Event)) func(event fsnotify.Event) {
	type debounceEntry struct {
		timer  *time.Timer
		active bool
	}

	debounceMap := make(map[string]*debounceEntry)
	return func(event fsnotify.Event) {
		sw.mu.Lock()
		defer sw.mu.Unlock()

		entry, exists := debounceMap[event.Name]
		if !exists {
			entry = &debounceEntry{
				timer:  time.NewTimer(debounceTime),
				active: true,
			}
			debounceMap[event.Name] = entry
			go func() {
				<-entry.timer.C
				sw.mu.Lock()
				if entry.active {
					handle(event)
				}
				delete(debounceMap, event.Name)
				sw.mu.Unlock()
			}()
		} else {
			entry.timer.Reset(debounceTime)
		}
	}
}

func (sw *SyncWatcher) handleEvent(event fsnotify.Event) {
	profileID, err := sw.getProfileIDForPath(event.Name)
	if err != nil {
		colorlog.Error("no profile found for this file: %s", event.Name)
		return
	}

	if err := sw.loadIgnoreList(profileID); err != nil {
		colorlog.Error("failed to load ignore list for profile %d: %v", profileID, err)
		return
	}

	if sw.shouldIgnore(event.Name, profileID) {
		return
	}

	op := sw.createOperation(event)
	if op != nil {
		sw.syncQueue <- op
	}
}

func (sw *SyncWatcher) createOperation(event fsnotify.Event) *SyncOperation {
	op := &SyncOperation{
		path:      event.Name,
		timeStamp: time.Now(),
	}

	switch {
	case event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write:
		op.operation = CreateOrModifyEvent
		hash, err := ComputeHash(event.Name)
		if err != nil {
			colorlog.Error("skipping event for file %s: %v", filepath.Base(event.Name), err)
			return nil
		}
		op.hash = hash
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		op.operation = DeleteEvent
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		op.operation = RenameEvent
		// op.OldPath = event.Name
	default:
		return nil
	}

	return op
}

func (sw *SyncWatcher) AddDirectory(path string) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if sw.paths[path] {
		return nil
	}

	err := watchRecursive(sw.watcher, path)
	if err == nil {
		sw.paths[path] = true
	}

	return err
}

func (sw *SyncWatcher) RemoveDirectory(path string) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	err := sw.watcher.Remove(path)
	if err == nil {
		delete(sw.paths, path)
	}
	return err
}

func (sw *SyncWatcher) loadIgnoreList(profileID int64) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	sw.ignoreList = make(map[int64][]string)
	if err := sw.loadIgnorePatternsForProfile(profileID); err != nil {
		return fmt.Errorf("failed to load ignore patterns for profile %d: %w", profileID, err)
	}

	return nil
}

func (sw *SyncWatcher) loadIgnorePatternsForProfile(profileID int64) error {
	patterns, err := sw.db.ListIgnorePattern(context.Background(), profileID)
	if err != nil {
		return fmt.Errorf("failed to list ignore patterns: %w", err)
	}

	if len(patterns) == 0 {
		return nil
	}

	list := make([]string, len(patterns))
	for i, p := range patterns {
		list[i] = p.Pattern
	}

	sw.ignoreList[profileID] = list
	return nil
}

func (sw *SyncWatcher) getProfileIDForPath(path string) (int64, error) {
	dir := filepath.Dir(path)
	profileID, err := sw.db.GetProfileIDBySourceDir(context.Background(), dir)
	if err != nil {
		return 0, err
	}
	return profileID, nil
}

func (sw *SyncWatcher) shouldIgnore(path string, profileID int64) bool {
	petterns, exists := sw.ignoreList[profileID]
	if !exists {
		return false
	}

	for _, pattern := range petterns {
		if matched, err := filepath.Match(pattern, filepath.Base(path)); matched && err == nil {
			return true
		}
	}
	return false
}
