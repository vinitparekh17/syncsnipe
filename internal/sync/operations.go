package sync

import (
	"database/sql"
	"time"
)

type SyncOperation struct {
	path      string
	operation string
	hash      string
	timeStamp time.Time
}

type SyncWorker struct {
  db *sql.DB
  syncQueue chan SyncOperation
  rules map[string]string
}

func NewSyncWorker(db *sql.DB, syncQueue chan SyncOperation) *SyncWorker {
  return &SyncWorker{
    db: db,
    syncQueue: syncQueue,
    rules: make(map[string]string),
  }
}
