package cli

import (
	"fmt"
	"log"
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
func setupTest() *database.Queries {
	mockDB = database.GetDatabase(mockDBFile)

	if err := mockDB.Ping(); err != nil {
		log.Fatal("error pinging db: %w", err)
	}

	if err := mockDB.LoadSchema(schemaFile); err != nil {
		log.Fatal("unable to load schema: %w", err)
	}

	colorlog.Success("successfully Connected to sqlite")

	return database.New(mockDB)
}

// Helper function to execute command and capture output
func executeCommand(cmd *cobra.Command, args ...string) error {
	cmd.SetArgs(args)
	return cmd.Execute()
}

func cleanupTest(t *testing.T, mockDB *database.DB) {
	t.Helper() // Marks this as a helper function
	defer mockDB.Close()
	// Delete files created during test with this *.db* file pattern
	pattern := "*.db*"
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("Error while matching files: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("No matching .db files found.")
		return
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			log.Printf("Failed to delete %s: %v", file, err)
		}
	}
}
