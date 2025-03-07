package sync

import (
	"database/sql"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
)

type SyncWatcher struct {
	watcher     *fsnotify.Watcher
	paths       map[string]bool
	syncQueue   chan SyncOperation
	ignoreList  []string
	db          *sql.DB
	mu          sync.Mutex
	fileHashMap map[string]string // for tracking file hashes
}

func NewSyncWatcher() (*SyncWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &SyncWatcher{
		watcher:     watcher,
		paths:       make(map[string]bool),
		syncQueue:   make(chan SyncOperation, 100),
		fileHashMap: make(map[string]string),
	}, err
}

func (sw *SyncWatcher) Start() {
	go func() {
		debounce := make(map[string]*time.Timer)
		const debounceTime = 250 * time.Millisecond

		for {
			select {
			case event, ok := <-sw.watcher.Events:
				if !ok {
					return
				}

				if ShouldIgnore(event.Name, sw.ignoreList) {
					continue
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
				colorlog.Error("watcher err in event loop: %v", err)
				os.Exit(1)
			}
		}
	}()
}

func (sw *SyncWatcher) handleEvent(event fsnotify.Event) {
	var op string

	switch {
	case event.Op&fsnotify.Create == fsnotify.Create:
		op = "create"
		if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
			if err := sw.AddDirectory(event.Name); err != nil {
				colorlog.Error("failed to watch new dir %s: %v", event.Name, err)
			}
		}
	case event.Op&fsnotify.Write == fsnotify.Write:
		op = "modify"
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		op = "delete"
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		op = "rename"
	}
	if op == "" {
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

func (sw *SyncWatcher) Close() error {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	err := sw.watcher.Close()
	if err == nil {
		close(sw.syncQueue)
	}
	return err
}
