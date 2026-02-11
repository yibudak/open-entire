package checkpoint

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/yibudak/open-entire/pkg/types"
)

// GenerateID creates a 12-character hex checkpoint ID.
func GenerateID() (string, error) {
	b := make([]byte, 6) // 6 bytes = 12 hex chars
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate checkpoint ID: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// ShardPath returns the sharded storage path for a checkpoint ID.
// Format: <first-2-chars>/<remaining-10-chars>/
func ShardPath(id string) string {
	if len(id) < 12 {
		return id + "/"
	}
	return id[:2] + "/" + id[2:] + "/"
}

// MetadataPath returns the full path to metadata.json for a checkpoint.
func MetadataPath(id string) string {
	return ShardPath(id) + "metadata.json"
}

// SessionPath returns the path to a session folder within a checkpoint.
func SessionPath(id string, index int) string {
	return fmt.Sprintf("%s%d/", ShardPath(id), index)
}

// SessionFiles returns all standard file paths for a session within a checkpoint.
func SessionFiles(id string, index int) map[string]string {
	base := SessionPath(id, index)
	return map[string]string{
		"content_hash": base + "content_hash.txt",
		"context":      base + "context.md",
		"full":         base + "full.jsonl",
		"metadata":     base + "metadata.json",
		"prompt":       base + "prompt.txt",
	}
}

// NewMetadata creates checkpoint metadata from the given parameters.
func NewMetadata(id, commitHash, branch, author, message, strategyName string) *types.CheckpointMetadata {
	return &types.CheckpointMetadata{
		ID:         id,
		CommitHash: commitHash,
		Branch:     branch,
		Author:     author,
		Message:    message,
		Strategy:   strategyName,
	}
}
