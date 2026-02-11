package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yibudak/open-entire/pkg/types"
)

func TestSessionStore(t *testing.T) {
	dir := t.TempDir()

	store, err := NewStore(dir)
	require.NoError(t, err)

	// Start session
	require.NoError(t, store.StartSession("sess-123", "claude-code"))

	// Should be active
	active := store.ActiveSessions()
	assert.Len(t, active, 1)
	assert.Equal(t, "sess-123", active[0].ID)
	assert.Equal(t, types.SessionActive, active[0].Phase)

	// Get session
	s, found := store.GetSession("sess-123")
	assert.True(t, found)
	assert.Equal(t, "claude-code", s.AgentName)

	// End session
	require.NoError(t, store.EndSession("sess-123"))

	active = store.ActiveSessions()
	assert.Len(t, active, 0)

	s, found = store.GetSession("sess-123")
	assert.True(t, found)
	assert.Equal(t, types.SessionEnded, s.Phase)
	assert.NotNil(t, s.EndedAt)
}

func TestSessionStorePersistence(t *testing.T) {
	dir := t.TempDir()

	store1, err := NewStore(dir)
	require.NoError(t, err)
	require.NoError(t, store1.StartSession("sess-abc", "claude-code"))

	// Load in a new store instance
	store2, err := NewStore(dir)
	require.NoError(t, err)

	s, found := store2.GetSession("sess-abc")
	assert.True(t, found)
	assert.Equal(t, "sess-abc", s.ID)
}
