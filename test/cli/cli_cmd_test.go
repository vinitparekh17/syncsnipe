package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCliCmd(t *testing.T) {
	t.Run("NewCliCmd", func(t *testing.T) {
		cliCmd := getCliCmd()
		assert.NotNil(t, cliCmd)
	})

	defer cleanupTest(t, mockDB)
}
