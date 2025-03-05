package stuffbin

import (
	"os"

	"github.com/knadh/stuffbin"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
)

func LoadFile(filePath string) stuffbin.FileSystem {
	path, err := os.Executable()
	if err != nil {
		colorlog.Error("%v", err)
		os.Exit(1)
	}

	fs, err := stuffbin.UnStuff(path)
	if err != nil {
		if err == stuffbin.ErrNoID {
			colorlog.Warn("unstuff failed in binary, using local file system for %s path", filePath)

			fs, err = stuffbin.NewLocalFS("/", filePath)
			if err != nil {
				colorlog.Error("error initializing local file system: %v", err)
				os.Exit(1)
			}
		} else {
			colorlog.Error("error initializing FS: %v", err)
		}
	}
	return fs
}
