package agent

import (
	"fmt"

	"github.com/yibudak/open-entire/pkg/types"
)

// Agent is the interface for AI agent integrations.
type Agent interface {
	Name() string
	Detect(repoDir string) (sessionID string, err error)
	ParseSession(sessionID string, repoDir string) (*types.SessionData, error)
	SessionPaths(repoDir string) types.AgentPaths
}

var registry = map[string]Agent{}

// Register adds an agent to the registry.
func Register(a Agent) {
	registry[a.Name()] = a
}

// Get returns a registered agent by name.
func Get(name string) (Agent, error) {
	a, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown agent: %s", name)
	}
	return a, nil
}

// DetectAny tries all registered agents and returns the first match.
func DetectAny(repoDir string) (Agent, string, error) {
	for _, a := range registry {
		sessionID, err := a.Detect(repoDir)
		if err == nil && sessionID != "" {
			return a, sessionID, nil
		}
	}
	return nil, "", fmt.Errorf("no active agent detected")
}

// All returns all registered agents.
func All() map[string]Agent {
	return registry
}
