package config

import (
	"os"
	"strings"
)

func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("ENTIRE_ENABLED"); v != "" {
		cfg.Enabled = strings.EqualFold(v, "true") || v == "1"
	}
	if v := os.Getenv("ENTIRE_STRATEGY"); v != "" {
		cfg.Strategy = v
	}
	if v := os.Getenv("ENTIRE_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("ENTIRE_TELEMETRY"); v != "" {
		cfg.Telemetry = strings.EqualFold(v, "true") || v == "1"
	}
}
