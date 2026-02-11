package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yibudak/open-entire/internal/checkpoint"
)

func (s *Server) apiListCheckpoints(w http.ResponseWriter, r *http.Request) {
	store := checkpoint.NewStore(s.repo)
	checkpoints, err := store.List()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, checkpoints)
}

func (s *Server) apiGetCheckpoint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	store := checkpoint.NewStore(s.repo)

	cp, err := store.Get(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	result := map[string]interface{}{
		"checkpoint": cp,
	}

	if cp.CommitHash != "" {
		diff, _ := s.repo.DiffContent(cp.CommitHash)
		result["diff"] = diff

		files, _ := s.repo.DiffFiles(cp.CommitHash)
		result["files"] = files
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) apiGetSession(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idx := parseIdxInt(chi.URLParam(r, "idx"))

	store := checkpoint.NewStore(s.repo)

	cp, err := store.Get(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	if idx >= len(cp.Sessions) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "session index out of range"})
		return
	}

	transcript, _ := store.FormattedTranscript(id, idx)
	rawTranscript, _ := store.RawTranscript(id, idx)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"session":    cp.Sessions[idx],
		"transcript": transcript,
		"raw":        rawTranscript,
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
