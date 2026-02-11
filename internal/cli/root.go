package cli

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/yibudak/open-entire/internal/config"
	"github.com/yibudak/open-entire/internal/logging"
)

var (
	cfgQuiet    bool
	cfgDetailed bool
)

// NewRootCmd creates the root cobra command.
func NewRootCmd(version, commit, date string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "open-entire",
		Short: "Open-source AI session capture CLI",
		Long: `Open-Entire captures AI coding agent sessions and links them to Git workflows.
It creates versioned, searchable checkpoints of your AI interactions.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level := "info"
			if cfgDetailed {
				level = "debug"
			}

			// Try to load config for log level
			repoDir, _ := findRepoRoot()
			if cfg, err := config.Load(repoDir); err == nil {
				if !cfgDetailed {
					level = cfg.LogLevel
				}
			}

			logging.Setup(level, cfgQuiet)
			slog.Debug("open-entire starting", "version", version)
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&cfgQuiet, "quiet", "q", false, "suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&cfgDetailed, "detailed", "v", false, "verbose output")

	rootCmd.AddCommand(
		newVersionCmd(version, commit, date),
		newEnableCmd(),
		newDisableCmd(),
		newStatusCmd(),
		newRewindCmd(),
		newResumeCmd(),
		newExplainCmd(),
		newCleanCmd(),
		newDoctorCmd(),
		newResetCmd(),
		newServeCmd(),
	)

	return rootCmd
}

// findRepoRoot walks up from cwd to find a .git directory.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(dir + "/.git"); err == nil {
			return dir, nil
		}
		parent := dir[:max(0, len(dir)-len("/"+dir[lastSlash(dir)+1:]))]
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func lastSlash(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			return i
		}
	}
	return -1
}
