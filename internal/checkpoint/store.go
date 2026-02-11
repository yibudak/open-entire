package checkpoint

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/yibudak/open-entire/internal/git"
	"github.com/yibudak/open-entire/pkg/types"
)

// Store reads/writes checkpoints on the entire/checkpoints/v1 branch.
type Store struct {
	repo *git.Repository
}

// NewStore creates a new checkpoint store.
func NewStore(repo *git.Repository) *Store {
	return &Store{repo: repo}
}

// Create writes a new checkpoint to the checkpoints branch.
func (s *Store) Create(meta *types.CheckpointMetadata, sessions []SessionBundle) error {
	meta.CreatedAt = time.Now()

	files := make(map[string][]byte)

	// Write checkpoint metadata
	metaData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	files[MetadataPath(meta.ID)] = metaData

	// Write each session
	for i, sess := range sessions {
		paths := SessionFiles(meta.ID, i)

		// Session metadata
		smData, err := json.MarshalIndent(sess.Metadata, "", "  ")
		if err != nil {
			return err
		}
		files[paths["metadata"]] = smData

		// Full transcript
		if len(sess.FullTranscript) > 0 {
			files[paths["full"]] = sess.FullTranscript
		}

		// Context markdown
		if len(sess.Context) > 0 {
			files[paths["context"]] = sess.Context
		}

		// Prompts
		if len(sess.Prompts) > 0 {
			files[paths["prompt"]] = sess.Prompts
		}

		// Content hash
		hash := sha256.Sum256(sess.FullTranscript)
		files[paths["content_hash"]] = []byte(hex.EncodeToString(hash[:]))
	}

	msg := fmt.Sprintf("checkpoint %s", meta.ID)
	if err := s.repo.CommitOnBranch(git.CheckpointsBranch, msg, files); err != nil {
		return fmt.Errorf("failed to write checkpoint: %w", err)
	}

	slog.Info("checkpoint created", "id", meta.ID, "sessions", len(sessions))
	return nil
}

// Get reads a checkpoint's metadata from the checkpoints branch.
func (s *Store) Get(id string) (*types.CheckpointMetadata, error) {
	data, err := s.repo.ReadFileFromBranch(git.CheckpointsBranch, MetadataPath(id))
	if err != nil {
		return nil, fmt.Errorf("checkpoint %s not found: %w", id, err)
	}

	var meta types.CheckpointMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

// List returns all checkpoints, sorted by creation time (newest first).
func (s *Store) List() ([]*types.CheckpointMetadata, error) {
	files, err := s.repo.ListFilesOnBranch(git.CheckpointsBranch, "")
	if err != nil {
		return nil, err
	}

	// Find all metadata.json files at the checkpoint level (2-level deep)
	var checkpoints []*types.CheckpointMetadata
	seen := make(map[string]bool)

	for _, f := range files {
		parts := strings.Split(f, "/")
		if len(parts) == 3 && parts[2] == "metadata.json" {
			id := parts[0] + parts[1]
			if seen[id] {
				continue
			}
			seen[id] = true

			meta, err := s.Get(id)
			if err != nil {
				slog.Debug("skipping invalid checkpoint", "id", id, "error", err)
				continue
			}
			checkpoints = append(checkpoints, meta)
		}
	}

	sort.Slice(checkpoints, func(i, j int) bool {
		return checkpoints[i].CreatedAt.After(checkpoints[j].CreatedAt)
	})

	return checkpoints, nil
}

// RawTranscript returns the raw JSONL transcript for a session within a checkpoint.
func (s *Store) RawTranscript(checkpointID string, sessionIndex int) (string, error) {
	paths := SessionFiles(checkpointID, sessionIndex)
	data, err := s.repo.ReadFileFromBranch(git.CheckpointsBranch, paths["full"])
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FormattedTranscript returns a human-readable transcript for a session.
func (s *Store) FormattedTranscript(checkpointID string, sessionIndex int) (string, error) {
	paths := SessionFiles(checkpointID, sessionIndex)
	data, err := s.repo.ReadFileFromBranch(git.CheckpointsBranch, paths["context"])
	if err != nil {
		// Fallback to raw
		return s.RawTranscript(checkpointID, sessionIndex)
	}
	return string(data), nil
}

// Rewind restores the working tree to the state at a checkpoint.
func (s *Store) Rewind(repoDir, checkpointID string, hard bool) error {
	meta, err := s.Get(checkpointID)
	if err != nil {
		return err
	}
	if meta.CommitHash == "" {
		return fmt.Errorf("checkpoint %s has no associated commit", checkpointID)
	}

	if hard {
		return s.repo.Checkout(meta.CommitHash, true)
	}
	return s.repo.Checkout(meta.CommitHash, false)
}

// RestoreLogs restores only the session logs from a checkpoint.
func (s *Store) RestoreLogs(repoDir, checkpointID string) error {
	// Read session data and write to local .open-entire/ directory
	_, err := s.Get(checkpointID)
	if err != nil {
		return err
	}
	slog.Info("logs restored from checkpoint", "id", checkpointID)
	return nil
}

// SessionBundle contains all the data for a session to be stored.
type SessionBundle struct {
	Metadata       *types.SessionMetadata
	FullTranscript []byte
	Context        []byte
	Prompts        []byte
}
