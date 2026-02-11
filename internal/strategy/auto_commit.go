package strategy

import (
	"context"
	"log/slog"

	"github.com/yibudak/open-entire/internal/checkpoint"
	"github.com/yibudak/open-entire/internal/git"
	"github.com/yibudak/open-entire/pkg/types"
)

// AutoCommit creates checkpoints after each AI agent response.
type AutoCommit struct {
	repoDir string
}

// NewAutoCommit creates an auto-commit strategy.
func NewAutoCommit(repoDir string) *AutoCommit {
	return &AutoCommit{repoDir: repoDir}
}

func (s *AutoCommit) Name() string { return "auto-commit" }

func (s *AutoCommit) OnAgentResponse(ctx context.Context, event *AgentResponseEvent) error {
	slog.Info("auto-commit: creating checkpoint on agent response", "session", event.SessionID)

	repo, err := git.Open(s.repoDir)
	if err != nil {
		return err
	}

	if !repo.HasCheckpointsBranch() {
		return nil
	}

	id, err := checkpoint.GenerateID()
	if err != nil {
		return err
	}

	commitHash, _ := repo.HeadCommitHash()
	branch, _ := repo.CurrentBranch()
	author := repo.Author()

	meta := checkpoint.NewMetadata(id, commitHash, branch, author, "[entire] auto checkpoint", s.Name())

	store := checkpoint.NewStore(repo)
	bundle := checkpoint.SessionBundle{
		Metadata: &types.SessionMetadata{
			AgentName: event.AgentName,
			SessionID: event.SessionID,
		},
	}

	return store.Create(meta, []checkpoint.SessionBundle{bundle})
}

func (s *AutoCommit) OnCommit(ctx context.Context, event *CommitEvent) error {
	// Auto-commit also checkpoints on explicit commits
	return (&ManualCommit{repoDir: s.repoDir}).OnCommit(ctx, event)
}

func (s *AutoCommit) OnPush(ctx context.Context, event *PushEvent) error {
	slog.Debug("auto-commit: push event")
	return nil
}
