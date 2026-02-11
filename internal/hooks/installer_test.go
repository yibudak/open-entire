package hooks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T) string {
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, ".git", "hooks")
	require.NoError(t, os.MkdirAll(hooksDir, 0o755))
	return dir
}

func TestInstallAndRemove(t *testing.T) {
	dir := setupTestRepo(t)

	// Install
	err := Install(dir, false)
	require.NoError(t, err)
	assert.True(t, IsInstalled(dir))

	// Verify files exist
	postCommit := filepath.Join(dir, ".git", "hooks", "post-commit")
	data, err := os.ReadFile(postCommit)
	require.NoError(t, err)
	assert.Contains(t, string(data), "managed by open-entire")

	// Remove
	err = Remove(dir)
	require.NoError(t, err)
	assert.False(t, IsInstalled(dir))
}

func TestInstallDoesNotOverwriteExistingHook(t *testing.T) {
	dir := setupTestRepo(t)

	// Write an existing hook
	hookPath := filepath.Join(dir, ".git", "hooks", "post-commit")
	require.NoError(t, os.WriteFile(hookPath, []byte("#!/bin/sh\necho existing"), 0o755))

	// Should fail without force
	err := Install(dir, false)
	assert.Error(t, err)

	// Should succeed with force
	err = Install(dir, true)
	assert.NoError(t, err)
}

func TestInstallForce(t *testing.T) {
	dir := setupTestRepo(t)

	// Install
	require.NoError(t, Install(dir, false))

	// Re-install with force
	err := Install(dir, true)
	assert.NoError(t, err)
	assert.True(t, IsInstalled(dir))
}
