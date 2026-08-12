package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vinitparekh17/syncsnipe/internal/core"
	"github.com/vinitparekh17/syncsnipe/internal/database"
)

type APIHandler struct {
	db             *database.Queries
	profileService core.ProfileService
	syncService    core.SyncService
}

func NewAPIHandler(db *database.Queries) *APIHandler {
	return &APIHandler{
		db:             db,
		profileService: core.NewProfile(db),
		syncService:    core.NewSync(db),
	}
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func jsonError(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}

func (h *APIHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	stats, err := h.db.GetDashboardStats(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, fmt.Sprintf("failed to fetch stats: %v", err))
		return
	}

	jsonResponse(w, http.StatusOK, stats)
}

func (h *APIHandler) HandleProfiles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		profiles, err := h.db.ListProfiles(r.Context())
		if err != nil {
			jsonError(w, http.StatusInternalServerError, fmt.Sprintf("failed to list profiles: %v", err))
			return
		}

		type ProfileWithRules struct {
			database.Profile
			Rules []database.SyncRule `json:"syncRules"`
		}

		result := make([]ProfileWithRules, 0, len(profiles))
		for _, p := range profiles {
			rules, err := h.db.ListSyncRules(r.Context(), p.ID)
			if err != nil {
				rules = []database.SyncRule{}
			}
			result = append(result, ProfileWithRules{
				Profile: p,
				Rules:   rules,
			})
		}
		jsonResponse(w, http.StatusOK, result)

	case http.MethodPost:
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := h.profileService.AddProfile(r.Context(), req.Name); err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}
		jsonResponse(w, http.StatusCreated, map[string]string{"status": "success", "message": "profile created successfully"})

	case http.MethodDelete:
		name := r.URL.Query().Get("name")
		if name == "" {
			var req struct {
				Name string `json:"name"`
			}
			_ = json.NewDecoder(r.Body).Decode(&req)
			name = req.Name
		}

		if name == "" {
			jsonError(w, http.StatusBadRequest, "profile name is required")
			return
		}

		if err := h.profileService.DeleteProfile(r.Context(), name); err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"status": "success", "message": "profile deleted successfully"})

	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *APIHandler) HandleSyncRules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rules, err := h.db.ListSyncRulesGroupByProfile(r.Context())
		if err != nil {
			jsonError(w, http.StatusInternalServerError, fmt.Sprintf("failed to list sync rules: %v", err))
			return
		}
		jsonResponse(w, http.StatusOK, rules)

	case http.MethodPost:
		var req struct {
			Profile         string `json:"profile"`
			SourceDir       string `json:"sourceDir"`
			TargetDir       string `json:"targetDir"`
			InitialSyncMode string `json:"initialSyncMode"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.InitialSyncMode == "" {
			req.InitialSyncMode = "mirror"
		}

		if err := h.syncService.AddSyncRule(r.Context(), req.Profile, req.SourceDir, req.TargetDir, req.InitialSyncMode); err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}
		jsonResponse(w, http.StatusCreated, map[string]string{"status": "success", "message": "sync rule added successfully"})

	case http.MethodDelete:
		var req struct {
			Profile   string `json:"profile"`
			SourceDir string `json:"sourceDir"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			req.Profile = r.URL.Query().Get("profile")
			req.SourceDir = r.URL.Query().Get("sourceDir")
		}

		if req.Profile == "" || req.SourceDir == "" {
			jsonError(w, http.StatusBadRequest, "profile and sourceDir are required")
			return
		}

		if err := h.syncService.RemoveSyncRuleByProfile(r.Context(), req.Profile, req.SourceDir); err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}
		jsonResponse(w, http.StatusOK, map[string]string{"status": "success", "message": "sync rule removed successfully"})

	default:
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *APIHandler) GetRecentActivities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	recentFiles, err := h.db.ListRecentFiles(r.Context())
	if err != nil {
		recentFiles = []database.ListRecentFilesRow{}
	}

	type ActivityItem struct {
		Label   string `json:"label"`
		Details string `json:"details"`
		Color   string `json:"color"`
	}

	activities := make([]ActivityItem, 0, len(recentFiles))
	for _, f := range recentFiles {
		activities = append(activities, ActivityItem{
			Label:   fmt.Sprintf("Synced %s", f.SourcePath),
			Details: fmt.Sprintf("%d bytes transferred", f.Size),
			Color:   "green",
		})
	}

	jsonResponse(w, http.StatusOK, activities)
}
