package core

import (
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/sync"
)

type App struct {
	DbQuery *database.Queries
	Watcher *sync.SyncWatcher
	Worker  *sync.SyncWorker
}
