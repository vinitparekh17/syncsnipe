package cli

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/vinitparekh17/syncsnipe/cmd/cli"
)

func getCliCmd() *cobra.Command {
	q := setupTest()
	return cli.NewCliCmd(q)
}

func TestProfileCommands(t *testing.T) {
	cliCmd := getCliCmd()

	t.Run("ProfileCmd", func(t *testing.T) {
		err := executeCommand(cliCmd, "profile")
		assert.NoError(t, err)
	})

	t.Run("AddProfileWithoutName", func(t *testing.T) {
		err := executeCommand(cliCmd, "profile", "add")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "accepts 1 arg(s), received 0")
	})

	t.Run("AddProfile", func(t *testing.T) {
		err := executeCommand(cliCmd, "profile", "add", testProfileName)
		assert.NoError(t, err)
	})

	t.Run("ListProfiles", func(t *testing.T) {
		err := executeCommand(cliCmd, "profile", "list")
		assert.NoError(t, err)
	})

	t.Run("RenameProfile", func(t *testing.T) {
		err := executeCommand(cliCmd, "profile", "rename", testProfileName, "new-profile")
		testProfileName = "new-profile" // reassigning profile name in order to delete it in the next test case
		assert.NoError(t, err)
	})

	t.Run("DeleteProfile", func(t *testing.T) {
		err := executeCommand(cliCmd, "profile", "delete", testProfileName)
		assert.NoError(t, err)
	})

	t.Run("EditNonExistentProfile", func(t *testing.T) {
		err := executeCommand(cliCmd, "profile", "rename", "FakeProfile", "new-profile")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "profile with FakeProfile name does not exist")
	})

	t.Run("DeleteNonExistentProfile", func(t *testing.T) {
		err := executeCommand(cliCmd, "profile", "delete", "FakeProfile")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "profile with FakeProfile name does not exist")
	})

	defer cleanupTest(t, mockDB)
}
