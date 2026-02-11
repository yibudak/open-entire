package claude

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseJSONL(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test-session.jsonl")

	jsonl := `{"type":"user","timestamp":"2025-01-15T10:00:00Z","requestId":"req-1","message":"Write a hello world function"}
{"type":"assistant","timestamp":"2025-01-15T10:00:05Z","requestId":"req-1","message":{"content":[{"type":"text","text":"Here is the function:"}]},"usage":{"input_tokens":100,"output_tokens":50,"cache_creation_input_tokens":10,"cache_read_input_tokens":5}}
{"type":"user","timestamp":"2025-01-15T10:01:00Z","requestId":"req-2","message":"Add tests"}
{"type":"assistant","timestamp":"2025-01-15T10:01:10Z","requestId":"req-2","message":{"content":[{"type":"text","text":"Here are the tests:"},{"type":"tool_use","name":"Write","input":{"path":"test.go"}}]},"usage":{"input_tokens":200,"output_tokens":100,"cache_creation_input_tokens":0,"cache_read_input_tokens":20}}
`
	require.NoError(t, os.WriteFile(path, []byte(jsonl), 0o644))

	session, err := ParseJSONL(path)
	require.NoError(t, err)

	assert.Equal(t, "claude-code", session.AgentName)
	assert.Len(t, session.Prompts, 2)
	assert.Len(t, session.Responses, 2)

	// Token aggregation
	assert.Equal(t, 300, session.TokenUsage.InputTokens)
	assert.Equal(t, 150, session.TokenUsage.OutputTokens)
	assert.Equal(t, 10, session.TokenUsage.CacheCreation)
	assert.Equal(t, 25, session.TokenUsage.CacheReads)
	assert.Equal(t, 2, session.TokenUsage.APICalls)

	// Tool calls
	assert.Len(t, session.ToolCalls, 1)
	assert.Equal(t, "Write", session.ToolCalls[0].Name)
}

func TestParseJSONLEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.jsonl")
	require.NoError(t, os.WriteFile(path, []byte(""), 0o644))

	session, err := ParseJSONL(path)
	require.NoError(t, err)
	assert.Len(t, session.Prompts, 0)
	assert.Len(t, session.Responses, 0)
}

func TestParseJSONLSkipsMalformed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.jsonl")
	content := `not json
{"type":"user","timestamp":"2025-01-15T10:00:00Z","requestId":"req-1","message":"Hello"}
also not json
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	session, err := ParseJSONL(path)
	require.NoError(t, err)
	assert.Len(t, session.Prompts, 1)
}

func TestParseJSONLDeduplicatesByRequestID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stream.jsonl")

	// Multiple assistant events with same requestId (streaming)
	content := `{"type":"assistant","timestamp":"2025-01-15T10:00:01Z","requestId":"req-1","message":"partial"}
{"type":"assistant","timestamp":"2025-01-15T10:00:02Z","requestId":"req-1","message":"partial more"}
{"type":"assistant","timestamp":"2025-01-15T10:00:03Z","requestId":"req-1","message":"complete","usage":{"input_tokens":100,"output_tokens":50,"cache_creation_input_tokens":0,"cache_read_input_tokens":0}}
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	session, err := ParseJSONL(path)
	require.NoError(t, err)

	// Should only count token usage once (from the last event with this requestId)
	assert.Equal(t, 100, session.TokenUsage.InputTokens)
	assert.Equal(t, 1, session.TokenUsage.APICalls)
}
