package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinitparekh17/syncsnipe/test"
)

func TestProfileCommands(t *testing.T) {
	cliCmd := test.GetCliCmd(t)

	t.Run("ProfileCmd", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "profile")
		assert.NoError(t, err)
	})

	t.Run("AddProfileWithoutName", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "profile", "add")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "accepts 1 arg(s), received 0")
	})

	t.Run("AddProfile", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "profile", "add", test.TestProfileName)
		assert.NoError(t, err)
	})

	t.Run("ListProfiles", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "profile", "list")
		assert.NoError(t, err)
	})

	t.Run("RenameProfile", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "profile", "rename", test.TestProfileName, "new-profile")
		test.TestProfileName = "new-profile" // reassigning profile name in order to delete it in the next test case
		assert.NoError(t, err)
	})

	t.Run("DeleteProfile", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "profile", "delete", test.TestProfileName)
		assert.NoError(t, err)
	})

	t.Run("EditNonExistentProfile", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "profile", "rename", "FakeProfile", "new-profile")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "profile with FakeProfile name does not exist")
	})

	t.Run("DeleteNonExistentProfile", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "profile", "delete", "FakeProfile")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "profile with FakeProfile name does not exist")
	})

	defer test.CleanupTest(t, test.MockDB)
}
