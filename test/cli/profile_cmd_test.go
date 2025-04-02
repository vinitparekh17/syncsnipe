package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinitparekh17/syncsnipe/test"
)

func TestProfileCommands(t *testing.T) {
	cliCmd := test.GetCliCmd(t) // Setup CLI command
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "ProfileCmd",
			args:    []string{"profile"},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "AddProfileWithoutName",
			args:    []string{"profile", "add"},
			wantErr: true,
			errMsg:  "accepts 1 arg(s), received 0",
		},
		{
			name:    "AddProfile",
			args:    []string{"profile", "add", test.TestProfileName},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "RenameProfile",
			args:    []string{"profile", "rename", test.TestProfileName, "new-profile"},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "DeleteProfile",
			args:    []string{"profile", "delete", "new-profile"},
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "DeleteNonExistentProfile",
			args:    []string{"profile", "delete", "FakeProfile"},
			wantErr: true,
			errMsg:  "profile with FakeProfile name does not exist",
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

	t.Cleanup(func() {
		test.CleanupTest(t, test.MockDB)
	})
}
