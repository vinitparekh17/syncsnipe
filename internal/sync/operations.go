package sync

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

// SyncOperation is a fully resolved, source-authoritative filesystem action.
// The watcher resolves paths before queueing so workers never need to query rules.
type SyncOperation struct {
	profileID  int64
	sourcePath string
	targetPath string
	operation  string
	timeStamp  time.Time
}

type SyncWorker struct {
	db        *database.Queries
	syncQueue <-chan *SyncOperation
	workers   int
	wg        sync.WaitGroup
}

func NewSyncWorker(db *database.Queries, queue <-chan *SyncOperation, workers int) (*SyncWorker, error) {
	if workers < 1 {
		workers = runtime.NumCPU()
	}
	return &SyncWorker{db: db, syncQueue: queue, workers: workers}, nil
}

func (w *SyncWorker) Start(ctx context.Context) {
	for i := 0; i < w.workers; i++ {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case op, ok := <-w.syncQueue:
					if !ok {
						return
					}
					w.processOperation(ctx, op)
				}
			}
		}()
	}
}

func (w *SyncWorker) Wait() { w.wg.Wait() }

func (w *SyncWorker) processOperation(ctx context.Context, op *SyncOperation) {
	info, err := os.Lstat(op.sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.RemoveAll(op.targetPath); err != nil && !os.IsNotExist(err) {
				colorlog.Error("removing target %s: %v", op.targetPath, err)
			} else {
				colorlog.Info("deleted target %s", op.targetPath)
			}
			_ = w.db.DeleteFile(ctx, database.DeleteFileParams{SourcePath: op.sourcePath, TargetPath: op.targetPath})
			return
		}
		colorlog.Error("stat source %s: %v", op.sourcePath, err)
		return
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() && !info.IsDir() {
		return
	}
	if info.IsDir() {
		if err := os.MkdirAll(op.targetPath, info.Mode().Perm()); err != nil {
			colorlog.Error("creating target directory %s: %v", op.targetPath, err)
		}
		return
	}

	hash, err := ComputeHash(op.sourcePath)
	if err != nil {
		return
	}
	if err := copyFileAtomic(op.sourcePath, op.targetPath, info.Mode().Perm()); err != nil {
		colorlog.Error("copying %s to %s: %v", op.sourcePath, op.targetPath, err)
		return
	}
	colorlog.Success("synced %s -> %s (%d bytes)", op.sourcePath, op.targetPath, info.Size())
	if err := w.db.UpsertFile(ctx, database.UpsertFileParams{
		SourcePath: op.sourcePath, TargetPath: op.targetPath, Hash: hash,
		Size: info.Size(), ModTime: info.ModTime().Unix(), LastSynced: time.Now().Unix(),
	}); err != nil {
		colorlog.Error("recording sync for %s: %v", op.sourcePath, err)
	}
}

func copyFileAtomic(source, target string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()
	tmp, err := os.CreateTemp(filepath.Dir(target), ".syncsnipe-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err = io.Copy(tmp, src); err == nil {
		err = tmp.Chmod(mode)
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(tmpName, target)
}
