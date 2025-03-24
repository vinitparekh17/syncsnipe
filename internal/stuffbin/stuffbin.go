package stuffbin

import (
	"fmt"
	"os"

	"github.com/knadh/stuffbin"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
)

// TODO: reconsider behaviour to tackle errors rather than log.Fetal

// LoadFile loads the file from linked binary / from local file system at given filePath parameter
func LoadFile(filePath string) (stuffbin.FileSystem, error) {
	path, err := os.Executable() // get self executable path
	if err != nil {
		return nil, fmt.Errorf("error while getting self executable path: %v", err)
	}

	fs, err := stuffbin.UnStuff(path)
	if err != nil {
		if err == stuffbin.ErrNoID {
			colorlog.Warn("unstuff failed in binary, attempting to use local file system for %s path", filePath)

			fs, err = stuffbin.NewLocalFS("/", filePath)
			if err != nil {
				return nil, fmt.Errorf("error initializing local file system: %v", err)
			}
		} else {
			return nil, fmt.Errorf("unable to unstuff %s path err: %v", filePath, err)
		}
	}
	return fs, nil
}
