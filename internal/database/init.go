package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/stuffbin"
)

type DB struct {
	*sql.DB
}

func GetDatabase(dbFile string) *DB {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("failed to open db connection: %v", err)
	}

	return &DB{db}
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
