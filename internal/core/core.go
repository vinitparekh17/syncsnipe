package core

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/sync"
)

type App struct {
	DB           *database.Queries
	Watcher      *sync.SyncWatcher
	Worker       *sync.SyncWorker
	ShutdownChan chan struct{}
}

func AddProfile(app *App, profileName string) error {
	profileName = strings.TrimSpace(profileName)
	profileName = regexp.MustCompile(`\s+`).ReplaceAllString(profileName, " ")

	if profileName == "" {
		return errors.New("empty profile's name is not allowed")
	}

	if len(profileName) < 2 || len(profileName) > 50 {
		return errors.New("profile's name must be between 2 and 50 characters")
	}

	var validProfileName = regexp.MustCompile(`^[a-zA-Z0-9 _-]+$`)
	if !validProfileName.MatchString(profileName) {
		return errors.New("profile's name can only contain letters, numbers, spaces, dashes (-), and underscores (_)")
	}

	count, err := app.DB.IsProfileExists(context.Background(), profileName)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("profile with %s name already exists", profileName)
	}

	_, err = app.DB.CreateProfile(context.Background(), profileName)
	return err
}

// TODO: additional validation panding, dir specific

// AddSyncRule inserts sync rule records after input sanitization
func AddSyncRule(app *App, profileName, sourceDir, targetDir string) error {
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
