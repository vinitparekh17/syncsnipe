package handler

import (
	"database/sql"
	"net/http"

	"github.com/knadh/stuffbin"
	"github.com/vinitparekh17/syncsnipe/internal/colorlog"
	"github.com/vinitparekh17/syncsnipe/internal/sync"
)

type Handler struct {
	db          *sql.DB
	syncWatcher *sync.SyncWatcher
	syncWorker  *sync.SyncWorker
}

func HandleStaticFiles(staticFile string, fs stuffbin.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		if _, err := w.Write(file.ReadBytes()); err != nil {
			colorlog.Error("failed to write file bytes in response")
		}
	}
}
