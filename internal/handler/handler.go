package handler

import (
	"database/sql"
	"net/http"
	"path"
	"path/filepath"

	"github.com/knadh/stuffbin"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/sync"
)

type Handler struct {
	db          *sql.DB
	syncWatcher *sync.SyncWatcher
	syncWorker  *sync.SyncWorker
}

var FrontendDir = filepath.Join("frontend", "build")

func ServeIndexPage(fs stuffbin.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

    // prevent page caching in order to get latest content
    r.Header.Add("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
    r.Header.Add("Pragma", "no-cache")
    r.Header.Add("Expires", "-1")

		file, err := fs.Get(path.Join(FrontendDir, "index.html"))
		if err != nil {
			colorlog.Error("error at fs.Get: %v", err)
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}

		r.Header.Set("Content-Type", "text/html")
		if _, err := w.Write(file.ReadBytes()); err != nil {
			colorlog.Error("failed to write file bytes in response")
		}
	}
}
