package strategy

import (
	"context"
	"fmt"
)

// Strategy defines how and when checkpoints are created.
type Strategy interface {
	Name() string
	OnAgentResponse(ctx context.Context, event *AgentResponseEvent) error
	OnCommit(ctx context.Context, event *CommitEvent) error
	OnPush(ctx context.Context, event *PushEvent) error
}

// AgentResponseEvent is fired when an AI agent produces a response.
type AgentResponseEvent struct {
	RepoDir   string
	SessionID string
	AgentName string
}

// CommitEvent is fired on git commit.
type CommitEvent struct {
	RepoDir    string
	CommitHash string
	Message    string
	Branch     string
}

// PushEvent is fired on git push.
type PushEvent struct {
	RepoDir string
	Remote  string
	Branch  string
}

// New creates a strategy by name.
func New(name string, repoDir string) (Strategy, error) {
	switch name {
	case "manual-commit":
		return NewManualCommit(repoDir), nil
	case "auto-commit":
		return NewAutoCommit(repoDir), nil
	default:
		return nil, fmt.Errorf("unknown strategy: %s", name)
	}
}
