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

func ShouldIgnore(path string, ignoreList []string) bool {
    for _, pattern := range ignoreList {
        if matched, err := filepath.Match(pattern, filepath.Base(path)); matched && err == nil {
            return true
        }
    }
    return false
}
