package claude

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yibudak/open-entire/internal/agent"
	"github.com/yibudak/open-entire/pkg/types"
)

// ClaudeAgent implements the Agent interface for Claude Code.
type ClaudeAgent struct{}

func init() {
	agent.Register(&ClaudeAgent{})
}

func (a *ClaudeAgent) Name() string {
	return agentName
}

func (a *ClaudeAgent) Detect(repoDir string) (string, error) {
	return Detect(repoDir)
}

func (a *ClaudeAgent) ParseSession(sessionID string, repoDir string) (*types.SessionData, error) {
	projDir := ProjectDir(repoDir)
	path := filepath.Join(projDir, sessionID+".jsonl")

	session, err := ParseJSONL(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Claude session %s: %w", sessionID, err)
	}

	session.ID = sessionID

	// Parse subagent sessions
	subFiles, err := SubagentFiles(repoDir, sessionID)
	if err == nil {
		for _, sf := range subFiles {
			subSession, err := ParseJSONL(sf)
			if err != nil {
				continue
			}
			subSession.ID = strings.TrimSuffix(filepath.Base(sf), ".jsonl")
			session.NestedSessions = append(session.NestedSessions, *subSession)
		}
	}

	return session, nil
}

func (a *ClaudeAgent) SessionPaths(repoDir string) types.AgentPaths {
	return types.AgentPaths{
		SessionDir: ProjectDir(repoDir),
		Pattern:    "*.jsonl",
	}
}
