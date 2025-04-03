package core

import (
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/sync"
)

const (
	profileNotFoundErr = "profile with name '%s' does not exist"
)

type SyncEngine struct {
	DB           *database.Queries
	Watcher      *sync.SyncWatcher
	Worker       *sync.SyncWorker
	ShutdownChan chan struct{}
}
