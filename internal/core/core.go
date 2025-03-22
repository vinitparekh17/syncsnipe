package core

import (
	"context"
	"errors"

	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/sync"
)

type SyncEngine struct {
	DB           *database.Queries
	Watcher      *sync.SyncWatcher
	Worker       *sync.SyncWorker
	ShutdownChan chan struct{}
}

// TODO: additional validation panding, dir specific

// AddSyncRule inserts sync rule records after input sanitization
func AddSyncRule(app *SyncEngine, profileName, sourceDir, targetDir string) error {
	if sourceDir == targetDir {
		return errors.New("source directory and target directory must be different")
	}

	profile, err := app.DB.GetProfileByName(context.Background(), profileName)
	if err != nil {
		return err
	}

	if _, err = app.DB.AddSyncRule(context.Background(), database.AddSyncRuleParams{
		ProfileID: profile.ID,
		SourceDir: sourceDir,
		TargetDir: targetDir,
	}); err != nil {
		return err
	}
	return nil
}
