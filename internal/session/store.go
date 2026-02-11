package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/yibudak/open-entire/pkg/types"
)

// State represents the persisted session state.
type State struct {
	Sessions []Session `json:"sessions"`
}

// Store manages session state persistence.
type Store struct {
	repoDir string
	path    string
	state   State
}

// NewStore creates a store backed by .open-entire/state.json.
func NewStore(repoDir string) (*Store, error) {
	s := &Store{
		repoDir: repoDir,
		path:    filepath.Join(repoDir, ".open-entire", "state.json"),
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// ActiveSessions returns sessions in the ACTIVE phase.
func (s *Store) ActiveSessions() []Session {
	var active []Session
	for _, sess := range s.state.Sessions {
		if sess.Phase == types.SessionActive {
			active = append(active, sess)
		}
	}
	return active
}

// StuckSessions returns sessions that appear stuck (ACTIVE for > 1 hour with no agent process).
func (s *Store) StuckSessions() []Session {
	var stuck []Session
	threshold := time.Now().Add(-1 * time.Hour)
	for _, sess := range s.state.Sessions {
		if sess.Phase == types.SessionActive && sess.StartedAt.Before(threshold) {
			stuck = append(stuck, sess)
		}
	}
	return stuck
}

// StartSession records a new active session.
func (s *Store) StartSession(id, agentName string) error {
	sess := Session{
		ID:        id,
		AgentName: agentName,
		RepoDir:   s.repoDir,
		Phase:     types.SessionActive,
		StartedAt: time.Now(),
	}
	s.state.Sessions = append(s.state.Sessions, sess)
	return s.save()
}

// EndSession marks a session as ended.
func (s *Store) EndSession(id string) error {
	now := time.Now()
	for i, sess := range s.state.Sessions {
		if sess.ID == id {
			s.state.Sessions[i].Phase = types.SessionEnded
			s.state.Sessions[i].EndedAt = &now
			return s.save()
		}
	}
	return nil
}

// GetSession returns a session by ID.
func (s *Store) GetSession(id string) (*Session, bool) {
	for _, sess := range s.state.Sessions {
		if sess.ID == id {
			return &sess, true
		}
	}
	return nil, false
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.state)
}

func (s *Store) save() error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
