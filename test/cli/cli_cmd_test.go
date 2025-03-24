package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupTest_MissingSchema_Error(t *testing.T) {
	oldSchema := schemaFile
	schemaFile = "non-existent.sql"
	defer func() {
		schemaFile = oldSchema
	}()

	q, err := setupTest(t)
	assert.ErrorContains(t, err, "unable to load schema")
	assert.Nil(t, q)
}

func TestNewCliCmd(t *testing.T) {
	t.Run("NewCliCmd", func(t *testing.T) {
		cliCmd := getCliCmd(t)
		assert.NotNil(t, cliCmd)
	})

	defer cleanupTest(t, mockDB)
}
