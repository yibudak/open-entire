package strategy

import (
	"context"
	"log/slog"

	"github.com/yibudak/open-entire/internal/checkpoint"
	"github.com/yibudak/open-entire/internal/git"
	"github.com/yibudak/open-entire/pkg/types"
)

// ManualCommit creates checkpoints only when the user makes a git commit.
type ManualCommit struct {
	repoDir string
}

// NewManualCommit creates a manual-commit strategy.
func NewManualCommit(repoDir string) *ManualCommit {
	return &ManualCommit{repoDir: repoDir}
}

func (s *ManualCommit) Name() string { return "manual-commit" }

func (s *ManualCommit) OnAgentResponse(ctx context.Context, event *AgentResponseEvent) error {
	// Manual strategy does not checkpoint on agent responses
	slog.Debug("manual-commit: ignoring agent response", "session", event.SessionID)
	return nil
}

func (s *ManualCommit) OnCommit(ctx context.Context, event *CommitEvent) error {
	slog.Info("manual-commit: creating checkpoint on commit", "commit", event.CommitHash)

	repo, err := git.Open(s.repoDir)
	if err != nil {
		return err
	}

	if !repo.HasCheckpointsBranch() {
		slog.Debug("checkpoints branch does not exist, skipping")
		return nil
	}

	id, err := checkpoint.GenerateID()
	if err != nil {
		return err
	}

	// Resolve commit info
	commitHash, _ := repo.HeadCommitHash()
	branch, _ := repo.CurrentBranch()
	message, _ := repo.LastCommitMessage()
	author := repo.Author()

	meta := checkpoint.NewMetadata(id, commitHash, branch, author, message, s.Name())

	store := checkpoint.NewStore(repo)

	// Create checkpoint with empty session bundle (agent parser will enrich)
	bundle := checkpoint.SessionBundle{
		Metadata: &types.SessionMetadata{
			AgentName: "unknown",
		},
	}

	if err := store.Create(meta, []checkpoint.SessionBundle{bundle}); err != nil {
		return err
	}

	// Add trailer to commit
	_ = repo.AddTrailer(git.TrailerCheckpoint, id)

	return nil
}

func (s *ManualCommit) OnPush(ctx context.Context, event *PushEvent) error {
	slog.Debug("manual-commit: push event (sync checkpoints)")
	return nil
}
