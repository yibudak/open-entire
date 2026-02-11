package web

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yibudak/open-entire/internal/checkpoint"
)

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	store := checkpoint.NewStore(s.repo)
	checkpoints, err := store.List()
	if err != nil {
		slog.Debug("failed to list checkpoints", "error", err)
		checkpoints = nil
	}

	data := map[string]interface{}{
		"Title":       "Entire — Dashboard",
		"RepoDir":     s.repoDir,
		"Checkpoints": checkpoints,
	}

	s.renderTemplate(w, "dashboard.html", data)
}

func (s *Server) handleCheckpointsList(w http.ResponseWriter, r *http.Request) {
	store := checkpoint.NewStore(s.repo)
	checkpoints, err := store.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":       "Entire — Checkpoints",
		"Checkpoints": checkpoints,
	}

	s.renderTemplate(w, "checkpoints.html", data)
}

func (s *Server) handleCheckpointDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	store := checkpoint.NewStore(s.repo)

	cp, err := store.Get(id)
	if err != nil {
		http.Error(w, "Checkpoint not found", http.StatusNotFound)
		return
	}

	// Get diff if commit hash exists
	var diff string
	if cp.CommitHash != "" {
		diff, _ = s.repo.DiffContent(cp.CommitHash)
	}

	data := map[string]interface{}{
		"Title":      "Entire — Checkpoint " + id[:8],
		"Checkpoint": cp,
		"Diff":       diff,
	}

	s.renderTemplate(w, "checkpoint_detail.html", data)
}

func (s *Server) handleSessionDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idx := chi.URLParam(r, "idx")
	store := checkpoint.NewStore(s.repo)

	cp, err := store.Get(id)
	if err != nil {
		http.Error(w, "Checkpoint not found", http.StatusNotFound)
		return
	}

	var sessionIdx int
	if _, err := parseIdx(idx); err == nil {
		sessionIdx = parseIdxInt(idx)
	}

	transcript, _ := store.FormattedTranscript(id, sessionIdx)

	data := map[string]interface{}{
		"Title":        "Entire — Session",
		"Checkpoint":   cp,
		"SessionIndex": sessionIdx,
		"Transcript":   transcript,
	}

	s.renderTemplate(w, "session_detail.html", data)
}

func (s *Server) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseFS(templatesFS, "templates/base.html", "templates/"+name)
	if err != nil {
		slog.Error("template parse error", "template", name, "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		slog.Error("template execute error", "template", name, "error", err)
	}
}

func parseIdx(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

func parseIdxInt(s string) int {
	i, _ := parseIdx(s)
	return i
}
