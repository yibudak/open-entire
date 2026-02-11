package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvOverrides(t *testing.T) {
	cfg := DefaultConfig()

	t.Setenv("ENTIRE_ENABLED", "false")
	t.Setenv("ENTIRE_STRATEGY", "auto-commit")
	t.Setenv("ENTIRE_LOG_LEVEL", "debug")
	t.Setenv("ENTIRE_TELEMETRY", "true")

	applyEnvOverrides(&cfg)

	assert.False(t, cfg.Enabled)
	assert.Equal(t, "auto-commit", cfg.Strategy)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.True(t, cfg.Telemetry)
}

func TestEnvOverridesBoolVariants(t *testing.T) {
	cfg := DefaultConfig()
	t.Setenv("ENTIRE_ENABLED", "0")
	applyEnvOverrides(&cfg)
	assert.False(t, cfg.Enabled)

	cfg2 := DefaultConfig()
	t.Setenv("ENTIRE_ENABLED", "1")
	applyEnvOverrides(&cfg2)
	assert.True(t, cfg2.Enabled)
}
