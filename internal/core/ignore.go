package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gobwas/glob"
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/types"
)

type IgnoreService interface {
	AddIgnore(ctx context.Context, profileName, pattern string) error
	ListIgnorePatterns(ctx context.Context, profileName string) ([]database.IgnorePattern, error)
	DeleteIgnorePattern(ctx context.Context, profileName, pattern string) error
}

type Ignore struct {
	DB *database.Queries
}

func NewIgnore(q *database.Queries) IgnoreService {
	return &Ignore{DB: q}
}

func (i *Ignore) AddIgnore(ctx context.Context, profileName, pattern string) error {
	ignoreParams, err := validateIgnorePattern(pattern)
	if err != nil {
		return err
	}

	profileID, err := i.DB.GetProfileIDByName(ctx, profileName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf(profileNotFoundErr, profileName)
		}
		return err
	}
	ignoreParams.ProfileID = profileID

	_, err = i.DB.AddIgnorePattern(ctx, ignoreParams)
	if err != nil {
		return err
	}

	return nil
}

func (i *Ignore) ListIgnorePatterns(ctx context.Context, profileName string) ([]database.IgnorePattern, error) {
	profileID, err := i.DB.GetProfileIDByName(ctx, profileName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf(profileNotFoundErr, profileName)
		}
		return nil, err
	}

	ignorePatterns, err := i.DB.ListIgnorePattern(ctx, profileID)
	if err != nil {
		return nil, err
	}
	return ignorePatterns, nil
}

func (i *Ignore) DeleteIgnorePattern(ctx context.Context, profileName, pattern string) error {
	profileID, err := i.DB.GetProfileIDByName(ctx, profileName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf(profileNotFoundErr, profileName)
		}
		return err
	}

	rows, err := i.DB.RemoveIgnorePatternByProfileName(ctx, database.RemoveIgnorePatternByProfileNameParams{
		ProfileID: profileID,
		Pattern:   pattern,
	})

	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("no ignore pattern '%s' found for profile '%s' to remove", pattern, profileName)
	}

	return nil
}

// validateIgnorePattern validates the ignore pattern and returns the ignore pattern params
func validateIgnorePattern(pattern string) (database.AddIgnorePatternParams, error) {
	if pattern == "" {
		return database.AddIgnorePatternParams{}, errors.New("pattern is required")
	}

	// Detect pattern type and validate it
	patternType, err := detectPatternType(pattern)
	if err != nil {
		return database.AddIgnorePatternParams{}, err
	}

	// Convert string to enum type
	pt, err := types.StringToIgnoreType(patternType)
	if err != nil {
		return database.AddIgnorePatternParams{}, err
	}

	return database.AddIgnorePatternParams{
		Pattern: pattern,
		Type:    pt,
	}, nil
}

// detectPatternType determines whether the given pattern is regex, glob, or exact.
func detectPatternType(pattern string) (string, error) {
	switch {
	case isRegex(pattern):
		if _, err := regexp.Compile(pattern); err != nil {
			return "", fmt.Errorf("invalid regex pattern: %w", err)
		}
		return "regex", nil

	case isGlob(pattern):
		if _, err := glob.Compile(pattern); err != nil {
			return "", fmt.Errorf("invalid glob pattern: %w", err)
		}
		return "glob", nil

	default:
		if err := validateExactPattern(pattern); err != nil {
			return "", err
		}
		return "exact", nil
	}
}

// isRegex checks if a pattern is likely a regex
func isRegex(pattern string) bool {
	return strings.HasPrefix(pattern, "^") || strings.HasSuffix(pattern, "$") ||
		strings.ContainsAny(pattern, "()[]|+?{}")
}

// isGlob checks if a pattern is likely a glob
func isGlob(pattern string) bool {
	return strings.ContainsAny(pattern, "*?[{}")
}

// validateExactPattern ensures the pattern is a valid exact match
func validateExactPattern(pattern string) error {
	if !utf8.ValidString(pattern) {
		return errors.New("invalid exact pattern: contains invalid UTF-8 sequences")
	}
	if strings.ContainsAny(pattern, "\x00\n\r") {
		return errors.New("invalid exact pattern: contains control characters")
	}
	return nil
}
