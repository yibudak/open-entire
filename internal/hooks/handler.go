package hooks

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/yibudak/open-entire/internal/config"
	"github.com/yibudak/open-entire/internal/strategy"
)

// Handler dispatches hook events to the configured strategy.
type Handler struct {
	repoDir  string
	cfg      *config.Config
	strategy strategy.Strategy
}

// NewHandler creates a new hook event handler.
func NewHandler(repoDir string, cfg *config.Config, strat strategy.Strategy) *Handler {
	return &Handler{
		repoDir:  repoDir,
		cfg:      cfg,
		strategy: strat,
	}
}

// HandlePostCommit handles the post-commit hook event.
func (h *Handler) HandlePostCommit(ctx context.Context) error {
	if !h.cfg.Enabled {
		slog.Debug("entire is disabled, skipping post-commit")
		return nil
	}

	event := &strategy.CommitEvent{
		RepoDir: h.repoDir,
	}

	if err := h.strategy.OnCommit(ctx, event); err != nil {
		return fmt.Errorf("strategy post-commit failed: %w", err)
	}

	return nil
}

// HandlePrePush handles the pre-push hook event.
func (h *Handler) HandlePrePush(ctx context.Context) error {
	if !h.cfg.Enabled {
		return nil
	}

	event := &strategy.PushEvent{
		RepoDir: h.repoDir,
	}

	return h.strategy.OnPush(ctx, event)
}
