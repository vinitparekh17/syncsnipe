package handler

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"strings"
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

func ServeIndexPage(fs stuffbin.FileSystem, frontendDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// prevent page caching in order to get latest content
		SetNoCache(w)
		file, err := fs.Get(filepath.Join(frontendDir, "index.html"))

		if err != nil {
			colorlog.Error("error at fs.Get: %v", err)
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}

		http.ServeContent(w, r, "index.html", time.Now(), file)
	}
}

// ServeApp serves the _app directory (sveltekit build output)
func ServeApp(fs stuffbin.FileSystem, frontendDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// safety check to ensure the request is for a file within the _app directory
		if !strings.HasPrefix(r.URL.Path, "/_app/") {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		// prevent page caching in order to get latest content
		SetNoCache(w)

		file, err := fs.Get(filepath.Join(frontendDir, r.URL.Path))
		if err != nil {
			colorlog.Error("error at fs.Get: %v", err)
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}
		http.ServeContent(w, r, r.URL.Path, time.Now(), file)
	}
}

func SetNoCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "-1")
	w.Header().Set("Content-Type", "text/html")
}
