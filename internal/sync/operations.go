package sync

import (
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
	wg        sync.WaitGroup
}

type SyncWorker struct {
	db        *database.Queries
	syncQueue chan *SyncOperation
	rules     map[string]string
	mu        sync.RWMutex
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
}
