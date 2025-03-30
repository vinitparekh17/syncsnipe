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

	profile, err := s.DB.GetProfileByName(ctx, profileName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("profile with name '%s' does not exist", profileName)
		}
		return err
	}

	if _, err = s.DB.AddSyncRule(context.Background(), database.AddSyncRuleParams{
		ProfileID: profile.ID,
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
