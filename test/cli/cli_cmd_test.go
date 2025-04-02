package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinitparekh17/syncsnipe/test"
)

func TestSetupTest(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() func()
		cleanup    func()
		expectErr  bool
		errMessage string
	}{
		{
			name: "MissingSchema_Error",
			setup: func() func() {
				// Store original and set new value properly
				originalSchema := test.SchemaFile
				test.SchemaFile = "non-existent.sql"
				// Return a function that restores the original
				return func() {
					test.SchemaFile = originalSchema
				}
			},
			expectErr:  true,
			errMessage: "failed to load schema: error initializing local file system: stat non-existent.sql: no such file or directory",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var cleanup func()
			if tc.setup != nil {
				cleanup = tc.setup()
			}

			// Ensure cleanup runs after test completes
			if cleanup != nil {
				defer cleanup()
			}

			q, err := test.SetupTest(t)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMessage)
				assert.Nil(t, q)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, q)
			}
		})
	}
}
