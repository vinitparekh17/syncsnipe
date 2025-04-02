package cli

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinitparekh17/syncsnipe/test"
)

func TestSyncCommands(t *testing.T) {
	cliCmd := test.GetCliCmd(t)
	test.CreateTestProfile(t)
	tempSrcDir := test.CreateTempDir(t)
	tempTargetDir := test.CreateTempDir(t)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "SyncCmd",
			args:    []string{"sync"},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "AddSyncWithWrongProfileName",
			args:    []string{"sync", "add", "invalid-profile-name", tempSrcDir, tempTargetDir},
			wantErr: true,
			errMsg:  "profile with name 'invalid-profile-name' does not exist",
		},
		{
			name:    "AddSyncWithWrongTargetDir",
			args:    []string{"sync", "add", test.TestProfileName, tempSrcDir, "invalid-target-dir"},
			wantErr: true,
			errMsg:  "target directory 'invalid-target-dir' does not exist",
		},
		{
			name:    "AddSyncWithSameSourceAndTargetDir",
			args:    []string{"sync", "add", test.TestProfileName, tempSrcDir, tempSrcDir},
			wantErr: true,
			errMsg:  "source directory and target directory must be different",
		},
		{
			name:    "AddSyncWithValidParams",
			args:    []string{"sync", "add", test.TestProfileName, tempSrcDir, tempTargetDir},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "StatusSyncRuleWithWrongProfileName",
			args:    []string{"sync", "status", "invalid-profile-name"},
			wantErr: true,
			errMsg:  "no sync rule found on 'invalid-profile-name' profile",
		},
		{
			name:    "StatusSyncRule",
			args:    []string{"sync", "status", test.TestProfileName},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "RemoveSyncRuleWithWrongProfileName",
			args:    []string{"sync", "remove", "invalid-profile-name", tempSrcDir},
			wantErr: true,
			errMsg:  fmt.Sprintf("no sync rule found on '%s' profile for source directory '%s'", "invalid-profile-name", tempSrcDir),
		},
		{
			name:    "RemoveSyncRuleWithWrongSourceDir",
			args:    []string{"sync", "remove", test.TestProfileName, "invalid-source-dir"},
			wantErr: true,
			errMsg:  fmt.Sprintf("no sync rule found on '%s' profile for source directory '%s'", test.TestProfileName, "invalid-source-dir"),
		},
		{
			name:    "RemoveSyncRule",
			args:    []string{"sync", "remove", test.TestProfileName, tempSrcDir},
			wantErr: false,
			errMsg:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := test.ExecuteCommand(cliCmd, tc.args...)

			if tc.wantErr {
				assert.Error(t, err)
				if tc.errMsg != "" {
					assert.ErrorContains(t, err, tc.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}

	// Cleanup after all tests
	defer test.CleanupTest(t, test.MockDB)
}
