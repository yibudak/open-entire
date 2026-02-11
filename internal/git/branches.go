package git

import (
	"fmt"
	"strings"
)

// CreateShadowBranch creates a temporary shadow branch for a session.
func (r *Repository) CreateShadowBranch(sessionID, worktreeID string) (string, error) {
	name := fmt.Sprintf("entire/%s-%s", sessionID, worktreeID)
	_, err := r.run("git", "branch", name)
	if err != nil {
		return "", fmt.Errorf("failed to create shadow branch %s: %w", name, err)
	}
	return name, nil
}

// ShadowBranches returns all Entire shadow branches.
func (r *Repository) ShadowBranches() ([]string, error) {
	out, err := r.run("git", "branch", "--list", "entire/*")
	if err != nil {
		return nil, err
	}

	var branches []string
	for _, line := range strings.Split(out, "\n") {
		b := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "*"))
		b = strings.TrimSpace(b)
		if b != "" && b != CheckpointsBranch {
			branches = append(branches, b)
		}
	}
	return branches, nil
}

// HasCheckpointsBranch checks if the checkpoints branch exists.
func (r *Repository) HasCheckpointsBranch() bool {
	_, err := r.run("git", "rev-parse", "--verify", CheckpointsBranch)
	return err == nil
}
