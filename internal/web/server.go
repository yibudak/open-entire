package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yibudak/open-entire/internal/git"
)

// Server is the local web viewer.
type Server struct {
	repo    *git.Repository
	repoDir string
	router  chi.Router
}

// NewServer creates a new web server.
func NewServer(repo *git.Repository, repoDir string) *Server {
	s := &Server{
		repo:    repo,
		repoDir: repoDir,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServerFS(staticFS)))

	// HTML pages
	r.Get("/", s.handleDashboard)
	r.Get("/checkpoints", s.handleCheckpointsList)
	r.Get("/checkpoints/{id}", s.handleCheckpointDetail)
	r.Get("/checkpoints/{id}/sessions/{idx}", s.handleSessionDetail)

	// JSON API
	r.Route("/api", func(r chi.Router) {
		r.Get("/checkpoints", s.apiListCheckpoints)
		r.Get("/checkpoints/{id}", s.apiGetCheckpoint)
		r.Get("/checkpoints/{id}/sessions/{idx}", s.apiGetSession)
	})

	s.router = r
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
