package claude

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/yibudak/open-entire/pkg/types"
)

// JSONLEvent represents a single event in a Claude Code JSONL transcript.
type JSONLEvent struct {
	Type       string          `json:"type"`
	Timestamp  string          `json:"timestamp,omitempty"`
	RequestID  string          `json:"requestId,omitempty"`
	Message    json.RawMessage `json:"message,omitempty"`
	Content    json.RawMessage `json:"content,omitempty"`
	Usage      *UsageData      `json:"usage,omitempty"`
	Role       string          `json:"role,omitempty"`
	CostUSD    float64         `json:"costUSD,omitempty"`
}

// UsageData represents token usage from Claude's API.
type UsageData struct {
	InputTokens        int `json:"input_tokens"`
	OutputTokens       int `json:"output_tokens"`
	CacheCreationInput int `json:"cache_creation_input_tokens"`
	CacheReadInput     int `json:"cache_read_input_tokens"`
}

// ParseJSONL parses a Claude Code JSONL transcript file.
func ParseJSONL(path string) (*types.SessionData, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	session := &types.SessionData{
		AgentName: "claude-code",
	}

	// Track request IDs for deduplication
	requestLastEvent := make(map[string]*JSONLEvent)
	var allEvents []JSONLEvent

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 10*1024*1024) // 10MB line buffer

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var event JSONLEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue // Skip malformed lines
		}

		allEvents = append(allEvents, event)

		// Deduplicate by requestId — keep last occurrence for token usage
		if event.RequestID != "" {
			requestLastEvent[event.RequestID] = &event
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading JSONL: %w", err)
	}

	// Process events
	for _, event := range allEvents {
		ts := parseTimestamp(event.Timestamp)

		switch event.Type {
		case "user":
			content := extractContent(event.Message)
			if content != "" {
				session.Prompts = append(session.Prompts, types.Prompt{
					Content:   content,
					Timestamp: ts,
					RequestID: event.RequestID,
				})
			}

		case "assistant":
			// Only process the last event per requestId for responses
			if event.RequestID != "" {
				last := requestLastEvent[event.RequestID]
				if last != nil && last.Timestamp == event.Timestamp {
					content := extractContent(event.Message)
					resp := types.Response{
						Content:   content,
						Timestamp: ts,
						RequestID: event.RequestID,
					}
					if event.Usage != nil {
						resp.TokenUsage = types.TokenUsage{
							InputTokens:   event.Usage.InputTokens,
							OutputTokens:  event.Usage.OutputTokens,
							CacheCreation: event.Usage.CacheCreationInput,
							CacheReads:    event.Usage.CacheReadInput,
						}
					}
					session.Responses = append(session.Responses, resp)
				}
			}
		}

		// Extract tool calls
		if event.Type == "assistant" {
			toolCalls := extractToolCalls(event.Message)
			for _, tc := range toolCalls {
				tc.Timestamp = ts
				tc.RequestID = event.RequestID
				session.ToolCalls = append(session.ToolCalls, tc)
			}
		}
	}

	// Aggregate token usage from deduplicated responses
	for _, last := range requestLastEvent {
		if last.Usage != nil {
			session.TokenUsage.InputTokens += last.Usage.InputTokens
			session.TokenUsage.OutputTokens += last.Usage.OutputTokens
			session.TokenUsage.CacheCreation += last.Usage.CacheCreationInput
			session.TokenUsage.CacheReads += last.Usage.CacheReadInput
			session.TokenUsage.APICalls++
		}
	}

	// Set timestamps
	if len(allEvents) > 0 {
		session.StartedAt = parseTimestamp(allEvents[0].Timestamp)
		lastTS := parseTimestamp(allEvents[len(allEvents)-1].Timestamp)
		session.EndedAt = &lastTS
	}

	return session, nil
}

func parseTimestamp(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	// Try ISO 8601 formats
	for _, layout := range []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func extractContent(raw json.RawMessage) string {
	if raw == nil {
		return ""
	}

	// Try as string
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}

	// Try as object with content field
	var obj struct {
		Content interface{} `json:"content"`
	}
	if json.Unmarshal(raw, &obj) == nil && obj.Content != nil {
		switch v := obj.Content.(type) {
		case string:
			return v
		case []interface{}:
			var parts []string
			for _, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					if text, ok := m["text"].(string); ok {
						parts = append(parts, text)
					}
				}
			}
			return strings.Join(parts, "\n")
		}
	}

	return string(raw)
}

func extractToolCalls(raw json.RawMessage) []types.ToolCall {
	if raw == nil {
		return nil
	}

	var msg struct {
		Content []struct {
			Type  string          `json:"type"`
			Name  string          `json:"name"`
			Input json.RawMessage `json:"input"`
		} `json:"content"`
	}

	if json.Unmarshal(raw, &msg) != nil {
		return nil
	}

	var calls []types.ToolCall
	for _, c := range msg.Content {
		if c.Type == "tool_use" {
			calls = append(calls, types.ToolCall{
				Name:  c.Name,
				Input: string(c.Input),
			})
		}
	}
	return calls
}
