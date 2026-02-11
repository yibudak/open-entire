package session

import (
	"time"

	"github.com/yibudak/open-entire/pkg/types"
)

// Session represents a tracked AI agent session.
type Session struct {
	ID        string             `json:"id"`
	AgentName string             `json:"agent_name"`
	RepoDir   string             `json:"repo_dir"`
	Phase     types.SessionPhase `json:"phase"`
	StartedAt time.Time          `json:"started_at"`
	EndedAt   *time.Time         `json:"ended_at,omitempty"`
}
