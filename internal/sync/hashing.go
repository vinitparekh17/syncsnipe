package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
)

func ComputeHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		colorlog.Error("error opening file %s while computing hash: %v", filepath.Base(path), err)
		return "", err
	}
	defer func() {
		if err := f.Close(); err != nil {
			colorlog.Warn("failed to close file %s at computeHash: %v", filepath.Base(path), err)
		}
	}()
	h := sha256.New()
	buf := make([]byte, 64*1024)
	if _, err := io.CopyBuffer(h, f, buf); err != nil {
		colorlog.Error("error computing hash for file %s: %v", filepath.Base(path), err)
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
