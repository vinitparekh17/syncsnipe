package sync

import (
	"io"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
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

func copyFile(source, target string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}

	defer func() {
		if err := src.Close(); err != nil {
			colorlog.Warn("failed to close source file %s: %v", source, err)
		}
	}()

	dst, err := os.Create(target)
	if err != nil {
		return err
	}

	defer func() {
		if err := dst.Close(); err != nil {
			colorlog.Warn("failed to close target file %s: %v", target, err)
		}
	}()

	_, err = io.Copy(dst, src)
	return err
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}
