package sync

import (
	"context"
	"path/filepath"
	"sync"
	"time"

	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

type SyncOperation struct {
	path      string
	operation string
	hash      string
	timeStamp time.Time
}

type SyncWorker struct {
	db        *database.Queries
	syncQueue chan *SyncOperation
	rules     map[string]string
	wg        sync.WaitGroup
}

func NewSyncWorker(db *database.Queries, syncQueue chan *SyncOperation) (*SyncWorker, error) {
	return &SyncWorker{
		db:        db,
		syncQueue: syncQueue,
		rules:     make(map[string]string),
	}, nil
}

func (s *SyncWorker) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for op := range s.syncQueue {
			s.processOperation(op)
		}
		colorlog.Info("SyncWorker stopped")
	}()
}

// TODO: Complete the logic, must support all file events
func (s *SyncWorker) processOperation(op *SyncOperation) {
	colorlog.Info("Processing operation %s on %s", op.operation, op.path)
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
