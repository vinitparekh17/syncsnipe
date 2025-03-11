package handler

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"time"

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
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "-1")
    w.Header().Set("Content-Type", "text/html")

		file, err := fs.Get(filepath.Join(FrontendDir, "index.html"))

		if err != nil {
			colorlog.Error("error at fs.Get: %v", err)
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}

    http.ServeContent(w, r, "index.html", time.Now(), file)
	}
}
