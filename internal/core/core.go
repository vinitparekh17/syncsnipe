package core

import (
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/sync"
)

type App struct {
	DBQuery *database.Queries
	Watcher *sync.SyncWatcher
	Worker  *sync.SyncWorker
}
