package sync

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeHash(t *testing.T) {

	tmpFile, _ := os.CreateTemp("", "test")

	// Normal case
	t.Run("ComputeHash", func(t *testing.T) {
		os.WriteFile(tmpFile.Name(), []byte("hello"), 0644)

		hash, err := ComputeHash(tmpFile.Name())
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.Equal(t, 64, len(hash)) // SHA-256 hex length
	})

	// Permission error case
	t.Run("PermissionError", func(t *testing.T) {
		os.Chmod(tmpFile.Name(), 0000) // No read perms
		hash, err := ComputeHash(tmpFile.Name())
		assert.Error(t, err)
		assert.Empty(t, hash)
		assert.Contains(t, err.Error(), "permission denied")
	})

	// Large file case
	t.Run("LargeFile", func(t *testing.T) {
		largeFile, _ := os.CreateTemp("", "large")
		defer os.Remove(largeFile.Name())
		data := make([]byte, 1*1024*1024*1024) // 1GB data
		largeFile.Write(data)
		hash, err := ComputeHash(largeFile.Name())
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.Equal(t, 64, len(hash)) // SHA-256 hex length
	})

	// Non-existent file case
	t.Run("NonExistentFile", func(t *testing.T) {
		hash, err := ComputeHash("non-existent")
		assert.Error(t, err)
		assert.Empty(t, hash)
		assert.Contains(t, err.Error(), "no such file or directory")
	})

	// Directory case
	t.Run("Directory", func(t *testing.T) {
		hash, err := ComputeHash(".")
		assert.Error(t, err)
		assert.Empty(t, hash)
		assert.Contains(t, err.Error(), "is a directory")
	})

	defer os.Remove(tmpFile.Name())
}
