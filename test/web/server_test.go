package web_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/vinitparekh17/syncsnipe/cmd/web"
	"github.com/vinitparekh17/syncsnipe/test"
)

func getWebCmd(t *testing.T) *cobra.Command {
	q, err := test.SetupTest(t)
	assert.NoError(t, err)
	webCmd, err := web.NewWebCmd(q)
	assert.NoError(t, err)
	return webCmd
}

func TestWebCmd(t *testing.T) {

	t.Run("TestWebCmd", func(t *testing.T) {
		webCmd := getWebCmd(t)
		assert.NotNil(t, webCmd)
		assert.Equal(t, "web", webCmd.Use)
		defer test.CleanupTest(t, test.MockDB)
	})

	// t.Run("TestWebCmdRun", func(t *testing.T) {
	// 	webCmd := getWebCmd(t)
	// 	err := webCmd.RunE(webCmd, []string{})
	// 	assert.NoError(t, err)
	// 	defer test.CleanupTest(t, test.MockDB)
	// })
}
