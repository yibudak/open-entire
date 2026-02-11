package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/yibudak/open-entire/internal/config"
	"github.com/yibudak/open-entire/internal/git"
	"github.com/yibudak/open-entire/internal/hooks"
)

func newEnableCmd() *cobra.Command {
	var (
		strategy string
		agent    string
		local    bool
		force    bool
	)

	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Initialize Open-Entire in a Git repository",
		Long:  "Install Git hooks and configure Open-Entire to capture AI coding sessions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoDir, err := findRepoRoot()
			if err != nil {
				return fmt.Errorf("not a git repository (or any parent): %w", err)
			}

			repo, err := git.Open(repoDir)
			if err != nil {
				return fmt.Errorf("failed to open repository: %w", err)
			}

			// Check if already enabled
			if repo.IsEntireEnabled() && !force {
				fmt.Println("Open-Entire is already enabled in this repository.")
				fmt.Println("Use --force to re-initialize.")
				return nil
			}

			// Create .open-entire directory
			entireDir := repoDir + "/.open-entire"
			if err := os.MkdirAll(entireDir, 0o755); err != nil {
				return fmt.Errorf("failed to create .open-entire directory: %w", err)
			}

			// Save config
			cfg := config.DefaultConfig()
			if strategy != "" {
				cfg.Strategy = strategy
			}
			if err := config.Save(repoDir, &cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			// Install hooks
			if err := hooks.Install(repoDir, force); err != nil {
				return fmt.Errorf("failed to install hooks: %w", err)
			}

			// Initialize checkpoints branch
			if err := repo.EnsureCheckpointsBranch(); err != nil {
				slog.Warn("could not create checkpoints branch", "error", err)
			}

			_ = local
			_ = agent

			fmt.Println("Open-Entire enabled successfully!")
			fmt.Printf("  Strategy: %s\n", cfg.Strategy)
			fmt.Printf("  Config:   %s/.open-entire/settings.json\n", repoDir)
			fmt.Println("\nYour AI coding sessions will now be captured as checkpoints.")
			return nil
		},
	}

	cmd.Flags().StringVar(&strategy, "strategy", "", "capture strategy (manual-commit or auto-commit)")
	cmd.Flags().StringVar(&agent, "agent", "", "AI agent to detect (default: auto)")
	cmd.Flags().BoolVar(&local, "local", false, "store data locally only")
	cmd.Flags().BoolVar(&force, "force", false, "force re-initialization")

	return cmd
}
