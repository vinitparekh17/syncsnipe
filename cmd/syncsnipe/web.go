package syncsnipe

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/knadh/stuffbin"
	"github.com/spf13/cobra"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/server"
)

var (
  staticFile = filepath.Join("frontend", "build", "index.html")
)
func initFS() stuffbin.FileSystem {
  path, err := os.Executable()
  if err != nil {
    colorlog.Error("%v", err)
    os.Exit(1)
  }

  fs, err := stuffbin.UnStuff(path)
  if err != nil {
    if err == stuffbin.ErrNoID {
      colorlog.Error("unstuff failed in binary, using local file system for static files")

      fs, err = stuffbin.NewLocalFS("/", staticFile)
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

var webCmd = &cobra.Command{
  Use: "web",
  Short: "run web interface",
  Run: func(cmd *cobra.Command, args []string) {
    mux := server.NewMuxRouter()
    fs := initFS()
    mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
      }
      file, err := fs.Get(staticFile)
      if err != nil {
        colorlog.Info("%s", staticFile)
        colorlog.Error("error at fs.Get: %v", err)
        http.Error(w, "unable to find and serve static files", http.StatusInternalServerError)
        return
      }
      r.Header.Set("Content-Type", "text/html")
      w.Write(file.ReadBytes())
     
    }))

    server := server.NewServer(mux)
    if err := server.Run(); err != nil {
      colorlog.Error("%v",err)
      os.Exit(1)
    }
  },
}

