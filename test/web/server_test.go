package web_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vinitparekh17/syncsnipe/test"
)

const (
	webHost = "localhost"
	webPort = "8080"
)

func TestWebCmd(t *testing.T) {
	tests := []struct {
		name     string
		port     string
		testFunc func(t *testing.T)
	}{
		{
			name: "TestWebCmdStructure",
			testFunc: func(t *testing.T) {
				webCmd := test.GetWebCmd(t)
				assert.NotNil(t, webCmd)
				assert.Equal(t, "web", webCmd.Use)
			},
		},
		{
			name: "TestDefaultPort",
			port: webPort,
			testFunc: func(t *testing.T) {
				cleanup := test.StartWebServer(t)
				defer cleanup()

				res, err := http.Get(fmt.Sprintf("http://%s:%s", webHost, webPort))
				require.NoError(t, err)
				defer res.Body.Close()
				assert.Equal(t, http.StatusOK, res.StatusCode)
			},
		},
		{
			name: "TestCustomPort",
			port: "8081",
			testFunc: func(t *testing.T) {
				cleanup := test.StartWebServer(t, "-p", "8081")
				defer cleanup()

				res, err := http.Get(fmt.Sprintf("http://%s:%s", webHost, "8081"))
				require.NoError(t, err)
				defer res.Body.Close()
				assert.Equal(t, http.StatusOK, res.StatusCode)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.testFunc)
	}

	t.Cleanup(func() {
		test.CleanupTest(t, test.MockDB)
	})
}
