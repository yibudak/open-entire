package types

import "time"

// SessionPhase represents the lifecycle phase of a session.
type SessionPhase string

const (
	SessionActive    SessionPhase = "ACTIVE"
	SessionEnded     SessionPhase = "ENDED"
	SessionCondensed SessionPhase = "CONDENSED"
)

// SessionData holds parsed data from an AI agent session.
type SessionData struct {
	ID            string          `json:"id"`
	AgentName     string          `json:"agent_name"`
	StartedAt     time.Time       `json:"started_at"`
	EndedAt       *time.Time      `json:"ended_at,omitempty"`
	Phase         SessionPhase    `json:"phase"`
	Prompts       []Prompt        `json:"prompts"`
	Responses     []Response      `json:"responses"`
	ToolCalls     []ToolCall      `json:"tool_calls"`
	TokenUsage    TokenUsage      `json:"token_usage"`
	FilesChanged  []string        `json:"files_changed"`
	NestedSessions []SessionData  `json:"nested_sessions,omitempty"`
}

// Prompt represents a user prompt in a session.
type Prompt struct {
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id"`
}

// Response represents an AI response in a session.
type Response struct {
	Content    string     `json:"content"`
	Timestamp  time.Time  `json:"timestamp"`
	RequestID  string     `json:"request_id"`
	TokenUsage TokenUsage `json:"token_usage"`
}

// ToolCall represents a tool invocation during a session.
type ToolCall struct {
	Name      string    `json:"name"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id"`
}

// TokenUsage tracks token consumption.
type TokenUsage struct {
	InputTokens       int `json:"input_tokens"`
	OutputTokens      int `json:"output_tokens"`
	CacheCreation     int `json:"cache_creation"`
	CacheReads        int `json:"cache_reads"`
	APICalls          int `json:"api_calls"`
}

// CheckpointMetadata is stored on the entire/checkpoints/v1 branch.
type CheckpointMetadata struct {
	ID           string           `json:"id"`
	CommitHash   string           `json:"commit_hash"`
	Branch       string           `json:"branch"`
	Author       string           `json:"author"`
	Message      string           `json:"message"`
	CreatedAt    time.Time        `json:"created_at"`
	Strategy     string           `json:"strategy"`
	Sessions     []SessionSummary `json:"sessions"`
	Attribution  *Attribution     `json:"attribution,omitempty"`
}

// SessionSummary is a lightweight view of a session within a checkpoint.
type SessionSummary struct {
	Index      int        `json:"index"`
	AgentName  string     `json:"agent_name"`
	SessionID  string     `json:"session_id"`
	TokenUsage TokenUsage `json:"token_usage"`
}

// SessionMetadata stored per-session within a checkpoint.
type SessionMetadata struct {
	AgentName    string      `json:"agent_name"`
	SessionID    string      `json:"session_id"`
	TokenUsage   TokenUsage  `json:"token_usage"`
	Attribution  Attribution `json:"attribution"`
	StartedAt    time.Time   `json:"started_at"`
	EndedAt      *time.Time  `json:"ended_at,omitempty"`
}

// Attribution tracks AI vs human line contribution.
type Attribution struct {
	AgentPercent  float64 `json:"agent_percent"`
	AgentLines    int     `json:"agent_lines"`
	TotalLines    int     `json:"total_lines"`
}

// AgentPaths holds filesystem paths for an agent's data.
type AgentPaths struct {
	SessionDir string `json:"session_dir"`
	Pattern    string `json:"pattern"`
}
