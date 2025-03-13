package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/stuffbin"
)

const dbFile = "syncsnipe.db"

type DB struct {
	*sql.DB
}

func GetDatabase() *DB {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		colorlog.Fetal("failed to open db connection: %v", err)
	}

	return &DB{db}
}

func (db *DB) LoadSchema(filePath string) error {
	fs := stuffbin.LoadFile(filePath)

	if !tableExists(db, "files") {
		file, err := fs.Get(filePath)
		if err != nil {
			colorlog.Error("error getting schema.sql: %v", err)
			return err
		}

		_, err = db.Exec(string(file.ReadBytes()))
		return err
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
