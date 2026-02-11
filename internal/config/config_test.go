package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.True(t, cfg.Enabled)
	assert.Equal(t, "manual-commit", cfg.Strategy)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.False(t, cfg.Telemetry)
}

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load("")
	require.NoError(t, err)
	assert.Equal(t, "manual-commit", cfg.Strategy)
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	entireDir := filepath.Join(dir, ".open-entire")
	require.NoError(t, os.MkdirAll(entireDir, 0o755))

	configJSON := `{"strategy": "auto-commit", "log_level": "debug"}`
	require.NoError(t, os.WriteFile(filepath.Join(entireDir, "settings.json"), []byte(configJSON), 0o644))

	cfg, err := Load(dir)
	require.NoError(t, err)
	assert.Equal(t, "auto-commit", cfg.Strategy)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.True(t, cfg.Enabled) // default
}

func TestLoadLocalOverridesProject(t *testing.T) {
	dir := t.TempDir()
	entireDir := filepath.Join(dir, ".open-entire")
	require.NoError(t, os.MkdirAll(entireDir, 0o755))

	projectJSON := `{"strategy": "manual-commit"}`
	localJSON := `{"strategy": "auto-commit"}`
	require.NoError(t, os.WriteFile(filepath.Join(entireDir, "settings.json"), []byte(projectJSON), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(entireDir, "settings.local.json"), []byte(localJSON), 0o644))

	cfg, err := Load(dir)
	require.NoError(t, err)
	assert.Equal(t, "auto-commit", cfg.Strategy) // local wins
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	cfg := &Config{
		Enabled:  true,
		Strategy: "auto-commit",
		LogLevel: "debug",
	}

	require.NoError(t, Save(dir, cfg))

	loaded, err := Load(dir)
	require.NoError(t, err)
	assert.Equal(t, "auto-commit", loaded.Strategy)
	assert.Equal(t, "debug", loaded.LogLevel)
}
