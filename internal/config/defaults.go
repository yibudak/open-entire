package config

// DefaultConfig returns the base configuration defaults.
func DefaultConfig() Config {
	return Config{
		Enabled:   true,
		Strategy:  "manual-commit",
		LogLevel:  "info",
		Telemetry: false,
		StrategyOptions: StrategyOptions{
			Summarize: SummarizeOptions{
				Enabled: false,
			},
		},
	}
}
