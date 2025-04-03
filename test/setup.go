package test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vinitparekh17/syncsnipe/cmd/cli"
	"github.com/vinitparekh17/syncsnipe/cmd/web"
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
	FrontendDir     = filepath.Join("..", "..", "frontend", "build")
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
	require.NoError(t, err)
	return cli.NewCliCmd(q)
}

func StartWebServer(t *testing.T, frontendDir string, args ...string) func() {
	t.Helper()
	webCmd := GetWebCmd(t, frontendDir)
	const webHost = "localhost"
	webPort := "8080" // Default port

	// Create a cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	// Start the web server in a goroutine
	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		if err := executeCommandWithContext(ctx, webCmd, args...); err != nil && !errors.Is(err, context.Canceled) {
			t.Errorf("Server execution failed: %v", err)
		}
	}()

	// Extract custom port from arguments
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-p" || args[i] == "--port" {
			webPort = args[i+1]
			break
		}
	}

	// Wait for server to be ready using exponential backoff
	backoff := 50 * time.Millisecond
	maxWait := 5 * time.Second
	deadline := time.Now().Add(maxWait)
	serverURL := fmt.Sprintf("http://%s:%s/health", webHost, webPort)

	for time.Until(deadline) > 0 {
		resp, err := http.Get(serverURL)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				break
			}
		}
		time.Sleep(backoff)
		backoff *= 2 // Exponential backoff
	}

	// Return a cleanup function
	return func() {
		cancel() // Cancel the context to stop the server

		// Wait for the server to exit with timeout
		select {
		case <-serverDone:
			t.Log("Server stopped successfully")
		case <-time.After(5 * time.Second):
			t.Errorf("Server failed to stop within timeout")
		}
	}
}

func executeCommandWithContext(ctx context.Context, cmd *cobra.Command, args ...string) error {
	cmd.SetArgs(args)

	// Create a channel to receive command execution result
	done := make(chan error, 1)
	go func() {
		done <- cmd.Execute()
	}()

	// Wait for either context cancellation or command completion
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// GetWebCmd returns a configured web command for testing
func GetWebCmd(t *testing.T, frontendDir string) *cobra.Command {
	q, err := SetupTest(t)
	require.NoError(t, err)
	webCmd, err := web.NewWebCmd(q, frontendDir)
	require.NoError(t, err)
	return webCmd
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

// FrontendFilesAvailable checks if the frontend build files exist
func FrontendFilesAvailable() bool {
	_, err := os.Stat(FrontendDir)
	return err == nil
}
