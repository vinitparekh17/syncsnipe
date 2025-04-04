package cli

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinitparekh17/syncsnipe/test"
)

func TestIgnoreCommands(t *testing.T) {
	cliCmd := test.GetCliCmd(t) // Setup CLI command
	test.CreateTestProfile(t)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		// Base command test
		{
			name:    "Base ignore command shows help",
			args:    []string{"ignore"},
			wantErr: false,
		},

		// Add command tests
		{
			name:    "Add ignore pattern succeeds",
			args:    []string{"ignore", "add", test.TestProfileName, ".md"},
			wantErr: false,
		},
		{
			name:    "Add without arguments fails",
			args:    []string{"ignore", "add"},
			wantErr: true,
			errMsg:  "accepts 2 arg(s), received 0",
		},
		{
			name:    "Add with missing pattern fails",
			args:    []string{"ignore", "add", test.TestProfileName},
			wantErr: true,
			errMsg:  "accepts 2 arg(s), received 1",
		},
		{
			name:    "Add with invalid profile fails",
			args:    []string{"ignore", "add", "invalid", ".md"},
			wantErr: true,
			errMsg:  "profile with name 'invalid' does not exist",
		},
		{
			name:    "Add with invalid pattern (control character) fails",
			args:    []string{"ignore", "add", test.TestProfileName, "\x00"},
			wantErr: true,
			errMsg:  "invalid exact pattern: contains control characters",
		},
		{
			name:    "Add with invalid pattern (invalid UTF-8) fails",
			args:    []string{"ignore", "add", test.TestProfileName, string([]byte{0xFF, 0xFE, 0xFD})},
			wantErr: true,
			errMsg:  "invalid exact pattern: contains invalid UTF-8 sequences",
		},

		// List command tests
		{
			name:    "List ignore patterns succeeds",
			args:    []string{"ignore", "list", test.TestProfileName},
			wantErr: false,
		},
		{
			name:    "List without profile fails",
			args:    []string{"ignore", "list"},
			wantErr: true,
			errMsg:  "accepts 1 arg(s), received 0",
		},
		{
			name:    "List with invalid profile fails",
			args:    []string{"ignore", "list", "invalid"},
			wantErr: true,
			errMsg:  "profile with name 'invalid' does not exist",
		},

		// Remove command tests
		{
			name:    "Remove ignore pattern succeeds",
			args:    []string{"ignore", "remove", test.TestProfileName, ".md"},
			wantErr: false,
		},
		{
			name:    "Remove without arguments fails",
			args:    []string{"ignore", "remove"},
			wantErr: true,
			errMsg:  "accepts 2 arg(s), received 0",
		},
		{
			name:    "Remove with missing pattern fails",
			args:    []string{"ignore", "remove", test.TestProfileName},
			wantErr: true,
			errMsg:  "accepts 2 arg(s), received 1",
		},
		{
			name:    "Remove with invalid profile fails",
			args:    []string{"ignore", "remove", "invalid", ".md"},
			wantErr: true,
			errMsg:  "profile with name 'invalid' does not exist",
		},
		{
			name:    "Remove with nonexistent pattern fails",
			args:    []string{"ignore", "remove", test.TestProfileName, "invalid"},
			wantErr: true,
			errMsg:  fmt.Sprintf("no ignore pattern 'invalid' found for profile '%s' to remove", test.TestProfileName),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := test.ExecuteCommand(cliCmd, tc.args...)

			if tc.wantErr {
				assert.Error(t, err, "Expected error for test case: %s", tc.name)
				if tc.errMsg != "" {
					assert.ErrorContains(t, err, tc.errMsg, "Error message mismatch for test case: %s", tc.name)
				}
			} else {
				assert.NoError(t, err, "Unexpected error for test case: %s", tc.name)
			}
		})
	}

	t.Cleanup(func() {
		test.CleanupTest(t, test.MockDB)
	})
}
