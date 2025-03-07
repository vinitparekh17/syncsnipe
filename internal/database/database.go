package database

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/knadh/stuffbin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
)

const dbFile = "syncsnipe.db"
var schemaFile = filepath.Join("internal", "database", "schema.sql")


type Db struct {
  *sql.DB
}

func GetDatabase() *Db {
  db, err := sql.Open("sqlite3", dbFile)
  if err != nil {
    colorlog.Fetal("failed to open db connection: %v", err)
  }

  return &Db{db}
}

func (db *Db) LoadSchema() error {
  fs, err := stuffbin.UnStuff(schemaFile)
   if err != nil {
    if err == stuffbin.ErrNoID {
      colorlog.Error("unstuff failed in binary, using local file system for static files")

      fs, err = stuffbin.NewLocalFS("/", schemaFile)
      if err != nil {
        colorlog.Error("error loading schema.sql in local fs: %v", err)
        return err
      }
    } else {
      colorlog.Error("error initializing FS: %v", err)
      return err
    }
  }

  if !tableExists(db, "file") {
  file, err := fs.Get(schemaFile)
  if err != nil {
    colorlog.Error("error getting schema.sql: %v", err)
    return err
  }

  _, err = db.Exec(string(file.ReadBytes()))
  return err 
  } else {
    return nil
  }
}


func tableExists(db *Db, tableName string) bool {
	query := "SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name=?);"
	var exists int
	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		colorlog.Error("%v",err)
	}
	return exists == 1
}
