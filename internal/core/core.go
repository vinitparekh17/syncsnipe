package core

import (
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/sync"
)

type SyncEngine struct {
	DB           *database.Queries
	Watcher      *sync.SyncWatcher
	Worker       *sync.SyncWorker
	ShutdownChan chan struct{}
}
