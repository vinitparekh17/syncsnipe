package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/stuffbin"
)

type DB struct {
	*sql.DB
}

func GetDatabase(dbFile string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbFile+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) LoadSchema(filePath string) error {
	fs, err := stuffbin.LoadFile(filePath)
	if err != nil {
		colorlog.Error("error loading schema file: %v", err)
		return err
	}

	if !tableExists(db, "files") {
		file, err := fs.Get(filePath)
		if err != nil {
			colorlog.Error("error getting schema.sql: %v", err)
			return err
		}

		_, err = db.Exec(string(file.ReadBytes()))
		return err
	}
	// Lightweight forward migration for databases created before initial sync
	// modes were introduced. SQLite does not support ADD COLUMN IF NOT EXISTS.
	if !columnExists(db, "sync_rules", "initial_sync_mode") {
		if _, err := db.Exec("ALTER TABLE sync_rules ADD COLUMN initial_sync_mode TEXT NOT NULL DEFAULT 'mirror' CHECK(initial_sync_mode IN ('mirror', 'events_only'))"); err != nil {
			return err
		}
	}
	colorlog.Info("schema already loaded, skipping")
	return nil
}

func tableExists(db *DB, tableName string) bool {
	query := "SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name=?);"
	var exists int
	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		colorlog.Error("%v", err)
	}
	return exists == 1
}

func columnExists(db *DB, tableName, columnName string) bool {
	rows, err := db.Query("PRAGMA table_info(" + tableName + ")")
	if err != nil {
		colorlog.Error("checking columns for %s: %v", tableName, err)
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, columnType string
		var notNull int
		var defaultValue interface{}
		var primaryKey int
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err == nil && name == columnName {
			return true
		}
	}
	return false
}
