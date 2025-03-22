package sync

import (
	"context"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

type SyncWatcher struct {
	watcher     *fsnotify.Watcher
	paths       map[string]bool
	syncQueue   chan *SyncOperation
	ignoreList  map[int64][]string // list per profileId
	db          *database.Queries
	mu          sync.Mutex
	fileHashMap map[string]string // for tracking file hashes
	Worker      *SyncWorker
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

	if err := sw.loadIgnoreList(); err != nil {
		return nil, err
	}

	worker, err := NewSyncWorker(db, sw.syncQueue)
	if err != nil {
		return nil, err
	}

	sw.Worker = worker
	return sw, nil
}

func (sw *SyncWatcher) Start() {
	sw.wg.Add(1)
	go func() {
		defer sw.wg.Done()
		debounce := make(map[string]*time.Timer)
		const debounceTime = 250 * time.Millisecond

		for {
			select {
			case event, ok := <-sw.watcher.Events:
				if !ok {
					return
				}

				if timer, exists := debounce[event.Name]; exists {
					timer.Reset(debounceTime)
				} else {
					debounce[event.Name] = time.AfterFunc(debounceTime, func() {
						sw.handleEvent(event)
						delete(debounce, event.Name)
					})
				}

			case err, ok := <-sw.watcher.Errors:
				if !ok {
					return
				}
				colorlog.Fatal("watcher err in event loop: %v", err)
			}
		}
	}()
}

func (sw *SyncWatcher) handleEvent(event fsnotify.Event) {
	profileID, err := sw.getProfileIDForPath(event.Name)
	if err != nil {
		colorlog.Error("no profile found for this file: %s", event.Name)
		return
	}

	if sw.shouldIgnore(event.Name, profileID) {
		return
	}

	op := SyncOperation{path: event.Name, timeStamp: time.Now()}
	switch {
	case event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write:
		op.operation = "create_or_modify"
		// TODO: create hash
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		op.operation = "remove"
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		op.operation = "rename"
		// op.OldPath = event.Name
	default:
		return
	}
}

func (sw *SyncWatcher) AddDirectory(path string) error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	if sw.paths[path] {
		return nil
	}

	err := WatchRecursive(sw.watcher, path)
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

func (sw *SyncWatcher) Close() {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	err := sw.watcher.Close()
	if err != nil {
		colorlog.Error("failed to close fsnotify watcher: %v", err)
	}
	close(sw.syncQueue)
	sw.wg.Wait()
	colorlog.Success("SyncWatcher closed successfully.")
}

func (sw *SyncWatcher) loadIgnoreList() error {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	profiles, err := sw.db.ListProfiles(context.Background())
	if err != nil {
		return err
	}

	for _, profile := range profiles {
		petterns, err := sw.db.ListIgnorePattern(context.Background(), profile.ID)
		if err != nil {
			return err
		}

		if len(petterns) == 0 {
			continue
		}

		list := make([]string, len(petterns))
		for i, p := range petterns {
			list[i] = p.Pattern
		}

		sw.ignoreList[profile.ID] = list
	}

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
