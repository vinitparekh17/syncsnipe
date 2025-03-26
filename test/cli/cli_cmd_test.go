package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinitparekh17/syncsnipe/test"
)

func TestSetupTest_MissingSchema_Error(t *testing.T) {
	oldSchema := test.SchemaFile
	test.SchemaFile = "non-existent.sql"
	defer func() {
		test.SchemaFile = oldSchema
	}()

	q, err := test.SetupTest(t)
	assert.ErrorContains(t, err, "unable to load schema")
	assert.Nil(t, q)
}

func TestNewCliCmd(t *testing.T) {
	t.Run("NewCliCmd", func(t *testing.T) {
		cliCmd := getCliCmd(t)
		assert.NotNil(t, cliCmd)
	})

	defer test.CleanupTest(t, test.MockDB)
}
