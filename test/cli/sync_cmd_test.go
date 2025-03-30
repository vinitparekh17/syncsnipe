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

	t.Run("SyncCmd", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "sync")
		assert.NoError(t, err)
	})

	t.Run("AddSyncWithWrongProfileName", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "sync", "add", "invalid-profile-name", tempSrcDir, tempTargetDir)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "profile with name 'invalid-profile-name' does not exist")
	})

	t.Run("AddSyncWithWrongSourceDir", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "sync", "add", test.TestProfileName, "invalid-source-dir", tempTargetDir)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "source directory 'invalid-source-dir' does not exist")
	})

	t.Run("AddSyncWithWrongTargetDir", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "sync", "add", test.TestProfileName, tempSrcDir, "invalid-target-dir")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "target directory 'invalid-target-dir' does not exist")
	})

	t.Run("AddSyncWithSameSourceAndTargetDir", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "sync", "add", test.TestProfileName, tempSrcDir, tempSrcDir)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "source directory and target directory must be different")
	})

	t.Run("AddSyncWithValidParams", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "sync", "add", test.TestProfileName, tempSrcDir, tempTargetDir)
		assert.NoError(t, err)
	})

	t.Run("ListSyncRules", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "sync", "list")
		assert.NoError(t, err)
	})

	t.Run("StatusSyncRuleWithWrongProfileName", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "sync", "status", "invalid-profile-name")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no sync rule found on 'invalid-profile-name' profile")
	})

	t.Run("StatusSyncRule", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "sync", "status", test.TestProfileName)
		assert.NoError(t, err)
	})

	t.Run("RemoveSyncRuleWithWrongProfileName", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "sync", "remove", "invalid-profile-name", tempSrcDir)
		assert.Error(t, err)
		assert.ErrorContains(t, err, fmt.Sprintf("no sync rule found on 'invalid-profile-name' profile for source directory '%s'", tempSrcDir))
	})

	t.Run("RemoveSyncRuleWithWrongSourceDir", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "sync", "remove", test.TestProfileName, "invalid-source-dir")
		assert.Error(t, err)
		assert.ErrorContains(t, err, fmt.Sprintf("no sync rule found on '%s' profile for source directory 'invalid-source-dir'", test.TestProfileName))
	})

	t.Run("RemoveSyncRule", func(t *testing.T) {
		err := test.ExecuteCommand(cliCmd, "sync", "remove", test.TestProfileName, tempSrcDir)
		assert.NoError(t, err)
	})

	defer test.CleanupTest(t, test.MockDB)
}
