package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/vinitparekh17/syncsnipe/internal/database"
)

type SyncService interface {
	AddSyncRule(ctx context.Context, profileName, sourceDir, targetDir string) error
	ListSyncRules(ctx context.Context) ([]database.ListSyncRulesGroupByProfileRow, error)
	RemoveSyncRuleByProfile(ctx context.Context, profileName, sourceDir string) error
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

	profile, err := s.DB.GetProfileByName(ctx, profileName)
	if err != nil {
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
		return fmt.Errorf("no sync rule found on '%s' profile", profileName)
	}

	return nil
}
