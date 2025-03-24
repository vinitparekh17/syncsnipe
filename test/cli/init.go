package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

// Test variables shared across tests
var (
	testProfileName = "test-profile" // test profile name to be used in add, rename, and delete tests

	mockDB     *database.DB
	schemaFile = filepath.Join("..", "..", "sql", "schema.sql")
)

const mockDBFile = "mock.db"

// Setup function to prepare test environment
func setupTest(t *testing.T) (*database.Queries, error) {
	mockDB = database.GetDatabase(mockDBFile)

	if err := mockDB.LoadSchema(schemaFile); err != nil {
		return nil, fmt.Errorf("unable to load schema: %v", err)
	}

	colorlog.Success("successfully Connected to sqlite")

	return database.New(mockDB), nil
}

// Helper function to execute command and capture output
func executeCommand(cmd *cobra.Command, args ...string) error {
	cmd.SetArgs(args)
	return cmd.Execute()
}

func cleanupTest(t *testing.T, mockDB *database.DB) {
	t.Helper()

	pattern := "*.db*"
	files, err := filepath.Glob(pattern)
	if err != nil && err == filepath.ErrBadPattern {
		return
	}

	// Ensure files exist to trigger file deletion loop
	if len(files) == 0 {
		t.Log("No matching .db files found.")
		return
	}

	// Attempt to remove each file
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			t.Errorf("Failed to delete %s: %v", file, err) // Continue loop
		}
	}
}
