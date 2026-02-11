package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the merged configuration.
type Config struct {
	Enabled         bool            `json:"enabled"`
	Strategy        string          `json:"strategy"`
	LogLevel        string          `json:"log_level"`
	Telemetry       bool            `json:"telemetry"`
	StrategyOptions StrategyOptions `json:"strategy_options"`
}

// StrategyOptions holds strategy-specific configuration.
type StrategyOptions struct {
	Summarize SummarizeOptions `json:"summarize"`
}

// SummarizeOptions controls automatic summarization.
type SummarizeOptions struct {
	Enabled bool `json:"enabled"`
}

// Load reads and merges configuration from all layers.
// Priority: local > project > global > defaults
func Load(repoDir string) (*Config, error) {
	cfg := DefaultConfig()

	// Layer 1: Global config (~/.config/open-entire/settings.json)
	globalPath := globalConfigPath()
	if err := mergeFromFile(&cfg, globalPath); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if repoDir != "" {
		// Layer 2: Project config (.open-entire/settings.json)
		projectPath := filepath.Join(repoDir, ".open-entire", "settings.json")
		if err := mergeFromFile(&cfg, projectPath); err != nil && !os.IsNotExist(err) {
			return nil, err
		}

		// Layer 3: Local config (.open-entire/settings.local.json)
		localPath := filepath.Join(repoDir, ".open-entire", "settings.local.json")
		if err := mergeFromFile(&cfg, localPath); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

	// Layer 4: Env var overrides (highest precedence)
	applyEnvOverrides(&cfg)

	return &cfg, nil
}

// Save writes config to the project settings file.
func Save(repoDir string, cfg *Config) error {
	dir := filepath.Join(repoDir, ".open-entire")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(dir, "settings.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func globalConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "open-entire", "settings.json")
}

func mergeFromFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, cfg)
}
