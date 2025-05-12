package sync

import (
	"context"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

type SyncOperation struct {
	profileID int64
	path      string
	operation string
	hash      string
	timeStamp time.Time
}

type SyncWorker struct {
	db          *database.Queries
	syncQueue   chan *SyncOperation
	workerCount int
	stopCh      chan struct{}
	rules       map[string]string
	wg          sync.WaitGroup
}

func NewSyncWorker(db *database.Queries, syncQueue chan *SyncOperation, workerCount int) (*SyncWorker, error) {
	if workerCount <= 0 {
		workerCount = runtime.NumCPU() // Default to CPU count
	}

	return &SyncWorker{
		db:          db,
		syncQueue:   syncQueue,
		rules:       make(map[string]string),
		workerCount: workerCount,
		stopCh:      make(chan struct{}),
	}, nil
}

func (s *SyncWorker) Start(ctx context.Context) {
	// Launch the worker pool
	for i := range s.workerCount {
		s.wg.Add(1)
		go func(workerID int) {
			defer s.wg.Done()
			s.workerLoop(ctx, workerID)
		}(i)
	}

	colorlog.Success("Started sync worker pool with %d workers", s.workerCount)
}

func (s *SyncWorker) workerLoop(ctx context.Context, workerID int) {
	for {
		select {
		case op, ok := <-s.syncQueue:
			if !ok {
				colorlog.Info("Worker %d stopping: queue closed", workerID)
				return
			}
			s.processOperation(op)
		case <-ctx.Done():
			colorlog.Info("Worker %d stopping: context canceled", workerID)
			return
		case <-s.stopCh:
			colorlog.Info("Worker %d stopping: stop signal received", workerID)
			return
		}
	}
}

// TODO: Complete the logic, must support all file events
func (s *SyncWorker) processOperation(op *SyncOperation) {
	colorlog.Info("Processing operation %s on %s", op.operation, op.path)

	if err := s.loadRules(op.profileID); err != nil {
		colorlog.Error("error loading rules for profile %d: %v", op.profileID, err)
		return
	}

	targetPath, err := s.getTargetPath(op.path)
	if err != nil {
		colorlog.Error("error getting target path for %s: %v", op.path, err)
		return
	}

	switch op.operation {
	case CreateOrModifyEvent:
		s.handleCreateOrModify(op, targetPath)
	default:
		colorlog.Warn("Unsupported operation %s", op.operation)
	}
}

func (s *SyncWorker) handleCreateOrModify(op *SyncOperation, targetPath string) {
	existing, err := s.db.GetFile(context.Background(), database.GetFileParams{
		SourcePath: op.path,
		TargetPath: targetPath,
	})

	if err != nil {
		colorlog.Error("error getting file from db: %v", err)
		return
	}

	if existing.Hash != op.hash && existing.LastSynced > op.timeStamp.Unix() {
		_, err := s.db.AddConflict(context.Background(), database.AddConflictParams{
			SourcePath: op.path,
			TargetPath: targetPath,
			SourceHash: op.hash,
			TargetHash: existing.Hash,
			SourceTime: op.timeStamp.Unix(),
			TargetTime: existing.ModTime,
			DetectedAt: time.Now().Unix(),
		})

		if err != nil {
			colorlog.Error("error adding conflict to db: %v", err)
		}
		colorlog.Warn("Conflict detected for %s", op.path)
		return
	}

	if err := copyFile(op.path, targetPath); err != nil {
		colorlog.Error("error copying file %s to %s: %v", op.path, targetPath, err)
		return
	}

	err = s.db.UpsertFile(context.Background(), database.UpsertFileParams{
		SourcePath: op.path,
		TargetPath: targetPath,
		Hash:       op.hash,
		Size:       getFileSize(op.path),
		ModTime:    op.timeStamp.Unix(),
		LastSynced: time.Now().Unix(),
	})

	if err != nil {
		colorlog.Error("error upserting file to db: %v", err)
	}

}

func (s *SyncWorker) getTargetPath(sourcePath string) (string, error) {
	dir := filepath.Dir(sourcePath)
	targetDir, ok := s.rules[dir]
	if !ok {
		return "", nil
	}
	return filepath.Join(targetDir, filepath.Base(sourcePath)), nil
}

func (s *SyncWorker) loadRules(profileID int64) error {
	rules, err := s.db.ListSyncRules(context.Background(), profileID)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		s.rules[rule.SourceDir] = rule.TargetDir
	}
	return nil
}

func (s *SyncWorker) close() {
	close(s.stopCh)
	s.wg.Wait()
	close(s.syncQueue)

	colorlog.Success("SyncWorker pool shutdown complete")
}
