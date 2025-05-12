package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/vinitparekh17/syncsnipe/internal/database"
)

type SyncService interface {
	AddSyncRule(ctx context.Context, profileName, sourceDir, targetDir string) error
	ListSyncRules(ctx context.Context) ([]database.ListSyncRulesGroupByProfileRow, error)
	RemoveSyncRuleByProfile(ctx context.Context, profileName, sourceDir string) error
	GetSyncStatusByProfileName(ctx context.Context, profileName string) (database.GetSyncStatusByProfileNameRow, error)
	// RunSyncRule(ctx context.Context, profileName string) error
}

type Sync struct {
	DB *database.Queries
}

func NewSync(q *database.Queries) SyncService {
	return &Sync{DB: q}
}

// AddSyncRule inserts sync rule records after input sanitization
func (s *Sync) AddSyncRule(ctx context.Context, profileName, sourceDir, targetDir string) error {
	if sourceDir == targetDir {
		return errors.New("source directory and target directory must be different")
	}

	if !checkIsDirExists(sourceDir) {
		return fmt.Errorf("source directory '%s' does not exist", sourceDir)
	}

	if !checkIsDirExists(targetDir) {
		return fmt.Errorf("target directory '%s' does not exist", targetDir)
	}

	profileID, err := s.DB.GetProfileIDByName(ctx, profileName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf(profileNotFoundErr, profileName)
		}
		return err
	}

	if _, err = s.DB.AddSyncRule(context.Background(), database.AddSyncRuleParams{
		ProfileID: profileID,
		SourceDir: sourceDir,
		TargetDir: targetDir,
	}); err != nil {
		return err
	}
	return nil
}

// ListSyncRules returns all sync rules grouped by profile
func (s *Sync) ListSyncRules(ctx context.Context) ([]database.ListSyncRulesGroupByProfileRow, error) {
	syncRules, err := s.DB.ListSyncRulesGroupByProfile(ctx)
	if err != nil {
		return nil, err
	}
	return syncRules, nil
}

func (s *Sync) RemoveSyncRuleByProfile(ctx context.Context, profileName, sourceDir string) error {
	rows, err := s.DB.DeleteSyncRuleByProfileName(context.Background(), database.DeleteSyncRuleByProfileNameParams{
		SourceDir: sourceDir,
		Name:      profileName,
	})

	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no sync rule found on '%s' profile for source directory '%s'", profileName, sourceDir)
	}

	return nil
}

func (s *Sync) GetSyncStatusByProfileName(ctx context.Context, profileName string) (database.GetSyncStatusByProfileNameRow, error) {
	syncRule, err := s.DB.GetSyncStatusByProfileName(ctx, profileName)
	if err != nil {
		if err == sql.ErrNoRows {
			return database.GetSyncStatusByProfileNameRow{}, fmt.Errorf("no sync rule found on '%s' profile", profileName)
		}
		return database.GetSyncStatusByProfileNameRow{}, err
	}
	return syncRule, nil
}

func checkIsDirExists(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}

// func (s *Sync) RunSyncRule(ctx context.Context, profileName string) error {
// 	profileID, err := s.DB.GetProfileIDByName(ctx, profileName)
// 	if err != nil {
// 		return fmt.Errorf("failed to get profile ID: %w", err)
// 	}

// 	syncRules, err := s.DB.ListSyncRulesByProfileID(ctx, profileID)
// 	if err != nil {
// 		return fmt.Errorf("failed to list sync rules: %w", err)
// 	}
// 	if len(syncRules) == 0 {
// 		return fmt.Errorf("no sync rules found for profile '%s'", profileName)
// 	}

// 	watcher, err := sync.NewSyncWatcher(s.DB)
// 	if err != nil {
// 		return fmt.Errorf("failed to create sync watcher: %w", err)
// 	}
// 	defer func() {
// 		colorlog.Info("cleaning up watcher resources")
// 		watcher.Close()
// 	}()

// 	var watchErrors []string
// 	for _, rule := range syncRules {
// 		if err := watcher.AddDirectory(rule.SourceDir); err != nil {
// 			errMsg := fmt.Sprintf("failed to watch '%s': %v", rule.SourceDir, err)
// 			colorlog.Error(errMsg)
// 			watchErrors = append(watchErrors, errMsg)
// 		} else {
// 			colorlog.Info("watching directory: %s", rule.SourceDir)
// 		}
// 	}
// 	if len(watchErrors) == len(syncRules) {
// 		return fmt.Errorf("no directories could be watched: %s", strings.Join(watchErrors, "; "))
// 	}

// 	colorlog.Info("starting sync watcher for profile '%s'", profileName)
// 	watcher.Start(ctx)
// 	done := make(chan struct{})

// 	// Start watcher with proper shutdown coordination
// 	go func() {
// 		watcher.Run(ctx) // Blocking version within goroutine
// 		close(done)      // Signal completion
// 	}()

// 	// Wait for either context cancellation or watcher completion
// 	select {
// 	case <-ctx.Done():
// 		colorlog.Info("received shutdown signal, stopping watcher...")
// 		watcher.Stop() // Graceful shutdown
// 		// Wait for watcher to finish cleanup with timeout
// 		select {
// 		case <-done:
// 			colorlog.Info("watcher stopped gracefully")
// 		case <-time.After(5 * time.Second):
// 			colorlog.Warn("watcher shutdown timed out")
// 		}
// 		return ctx.Err()
// 	case <-done:
// 		colorlog.Info("watcher completed naturally")
// 		return nil
// 	}
// }
