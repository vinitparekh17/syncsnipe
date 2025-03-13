package sync

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func WatchRecursive(watcher *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
}

func (sw *SyncWatcher) shouldIgnore(path string, profileID int64) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

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
