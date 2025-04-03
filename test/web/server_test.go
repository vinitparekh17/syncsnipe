package web_test

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
		testFunc func(t *testing.T, frontendDir string)
	}{
		{
			name: "TestWebCmdStructure",
			testFunc: func(t *testing.T, frontendDir string) {
				webCmd := test.GetWebCmd(t, frontendDir)
				assert.NotNil(t, webCmd)
				assert.Equal(t, "web", webCmd.Use)
			},
		},
		{
			name: "TestDefaultPort",
			port: webPort,
			testFunc: func(t *testing.T, frontendDir string) {
				cleanup := test.StartWebServer(t, frontendDir)
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
			testFunc: func(t *testing.T, frontendDir string) {
				cleanup := test.StartWebServer(t, frontendDir, "-p", "8081")
				defer cleanup()
				res, err := http.Get(fmt.Sprintf("http://%s:%s", webHost, "8081"))
				require.NoError(t, err)
				defer res.Body.Close()
				assert.Equal(t, http.StatusOK, res.StatusCode)
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			frontendDir := test.FrontendDir
			_, err := os.Stat(frontendDir)
			if os.IsNotExist(err) {
				t.Logf("frontend directory '%s' not found, creating temporary directory", frontendDir)
				tempFrontendDir := t.TempDir()
				dummyIndexPath := filepath.Join(tempFrontendDir, "index.html")
				dummyIndexContent := []byte("<html><body>Hello, World!</body></html>")
				err = os.WriteFile(dummyIndexPath, dummyIndexContent, 0644)
				require.NoError(t, err)
				frontendDir = tempFrontendDir
				t.Logf("created dummy index.html in %s", tempFrontendDir)
			} else if err != nil {
				t.Fatalf("failed to check frontend directory: %v", err)
			}
			tc.testFunc(t, frontendDir)
		})
	}
	t.Cleanup(func() {
		test.CleanupTest(t, test.MockDB)
	})
}
