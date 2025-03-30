package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/vinitparekh17/syncsnipe/cmd/cli"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

const (
	mockDBFile = "mock.db"
)

var (
	TestProfileName = "test-profile" // Test profile name used across tests
	MockDB          *database.DB
	SchemaFile      = filepath.Join("..", "..", "sql", "schema.sql")
)

// SetupTest initializes test database and loads schema
func SetupTest(t *testing.T) (*database.Queries, error) {
	MockDB, err := database.GetDatabase(mockDBFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %v", err)
	}

	if err := MockDB.LoadSchema(SchemaFile); err != nil {
		return nil, fmt.Errorf("failed to load schema: %v", err)
	}

	colorlog.Success("successfully connected to SQLite database")
	return database.New(MockDB), nil
}

// GetCliCmd returns a configured CLI command for testing
func GetCliCmd(t *testing.T) *cobra.Command {
	q, err := SetupTest(t)
	assert.NoError(t, err)
	return cli.NewCliCmd(q)
}

// ExecuteCommand runs a CLI command with given arguments
func ExecuteCommand(cmd *cobra.Command, args ...string) error {
	cmd.SetArgs(args)
	return cmd.Execute()
}

// CleanupTest removes all test database files
func CleanupTest(t *testing.T, mockDB *database.DB) {
	t.Helper()

	files, err := filepath.Glob("*.db*")
	if err != nil {
		if err == filepath.ErrBadPattern {
			return
		}
		t.Errorf("Failed to glob db files: %v", err)
		return
	}

	if len(files) == 0 {
		t.Log("No database files found to cleanup")
		return
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			t.Errorf("Failed to delete %s: %v", file, err)
		}
	}
}

// CreateTempDir creates and returns path to a temporary directory
func CreateTempDir(t *testing.T) string {
	return t.TempDir()
}

// CreateTempFile creates and returns path to a temporary file
func CreateTempFile(t *testing.T) string {
	tempFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return tempFile.Name()
}

// CreateTestProfile creates a test profile using the CLI
func CreateTestProfile(t *testing.T) {
	cliCmd := GetCliCmd(t)
	err := ExecuteCommand(cliCmd, "profile", "add", TestProfileName)
	assert.NoError(t, err)
}
