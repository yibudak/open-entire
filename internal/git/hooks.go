package git

// Hook-related constants and utilities are in internal/hooks package.
// This file provides git-level hook path resolution.

import (
	"os"
	"path/filepath"
)

// HooksDir returns the Git hooks directory for the repository.
func (r *Repository) HooksDir() string {
	return filepath.Join(r.Dir, ".git", "hooks")
}

// HookExists checks if a specific hook file exists.
func (r *Repository) HookExists(name string) bool {
	_, err := os.Stat(filepath.Join(r.HooksDir(), name))
	return err == nil
}
