package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/gobwas/glob"
	"github.com/vinitparekh17/syncsnipe/internal/database"
	"github.com/vinitparekh17/syncsnipe/internal/types"
)

const (
	maxPatternLength     = 260 // Windows MAX_PATH limit
	maxComponentLength   = 255 // Common filesystem limit
	maxNestedDirectories = 32  // Practical nesting limit
)

var (
	windowsSpecialFiles = map[string]struct{}{
		"CON": {}, "PRN": {}, "AUX": {}, "NUL": {},
		"COM1": {}, "COM2": {}, "COM3": {}, "COM4": {}, "COM5": {},
		"COM6": {}, "COM7": {}, "COM8": {}, "COM9": {},
		"LPT1": {}, "LPT2": {}, "LPT3": {}, "LPT4": {}, "LPT5": {},
		"LPT6": {}, "LPT7": {}, "LPT8": {}, "LPT9": {},
	}

	unixSpecialFiles = map[string]struct{}{
		"dev/null": {}, "dev/zero": {}, "dev/random": {}, "dev/urandom": {},
		"etc/passwd": {}, "etc/shadow": {}, "etc/hosts": {},
		".bashrc": {}, ".bash_profile": {}, ".profile": {}, ".zshrc": {},
		".ssh/id_rsa": {}, ".ssh/id_dsa": {}, ".ssh/authorized_keys": {},
		".env": {}, ".git": {}, ".svn": {}, ".hg": {},
	}

	macSpecialFiles = map[string]struct{}{
		".DS_Store": {}, ".Spotlight-V100": {}, ".Trashes": {}, ".fseventsd": {},
		"._": {},
	}
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

// validateExactPattern checks if a pattern is valid and secure for use in ignore rules.
// It performs multiple security and format validations to prevent potential issues.
func validateExactPattern(pattern string) error {
	// Normalized path for cross-platform compatibility
	pattern = filepath.ToSlash(filepath.Clean(pattern))

	if pattern == "" {
		return errors.New("invalid exact pattern: empty pattern not allowed")
	}

	if !utf8.ValidString(pattern) {
		return errors.New("invalid exact pattern: contains invalid UTF-8 sequences")
	}

	if strings.ContainsAny(pattern, "\x00\n\r\t") {
		return errors.New("invalid exact pattern: contains control characters")
	}

	if strings.HasPrefix(pattern, "/") {
		return errors.New("invalid exact pattern: absolute paths not allowed")
	}

	if strings.Contains(pattern, "../") {
		return errors.New("invalid exact pattern: directory traversal not allowed")
	}

	if len(pattern) > maxPatternLength {
		return fmt.Errorf("invalid exact pattern: exceeds maximum length of %d characters", maxPatternLength)
	}

	parts := strings.Split(pattern, "/")
	if len(parts) > maxNestedDirectories {
		return fmt.Errorf("invalid exact pattern: exceeds %d nested directories", maxNestedDirectories)
	}

	for _, part := range parts {
		if len(part) > maxComponentLength {
			return fmt.Errorf("invalid exact pattern: component exceeds %d characters", maxComponentLength)
		}
	}

	upperBase := strings.ToUpper(filepath.Base(pattern))
	if _, isSpecial := windowsSpecialFiles[upperBase]; isSpecial {
		return fmt.Errorf("invalid exact pattern: Windows reserved name '%s' not allowed", pattern)
	}

	for reserved := range windowsSpecialFiles {
		if strings.HasPrefix(upperBase, reserved+".") {
			return fmt.Errorf("invalid exact pattern: Windows reserved name '%s' not allowed", reserved)
		}
	}

	basePattern := filepath.Base(pattern)
	if _, isUnixSpecial := unixSpecialFiles[basePattern]; isUnixSpecial {
		return fmt.Errorf("invalid exact pattern: Unix/Linux special file '%s' not allowed", basePattern)
	}

	if _, isMacSpecial := macSpecialFiles[basePattern]; isMacSpecial {
		return fmt.Errorf("invalid exact pattern: macOS special file '%s' not allowed", basePattern)
	}

	return nil
}
