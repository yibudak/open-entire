package session

// AgentDetector detects if an AI agent is currently active.
type AgentDetector interface {
	// Name returns the agent name.
	Name() string
	// Detect checks if the agent is active and returns the session ID.
	Detect(repoDir string) (sessionID string, err error)
}
