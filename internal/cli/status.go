package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yibudak/open-entire/internal/config"
	"github.com/yibudak/open-entire/internal/git"
	"github.com/yibudak/open-entire/internal/session"
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show Entire status for current repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoDir, err := findRepoRoot()
			if err != nil {
				return fmt.Errorf("not a git repository: %w", err)
			}

			cfg, err := config.Load(repoDir)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			repo, err := git.Open(repoDir)
			if err != nil {
				return fmt.Errorf("failed to open repository: %w", err)
			}

			enabled := repo.IsEntireEnabled()
			branch, _ := repo.CurrentBranch()

			fmt.Printf("Repository: %s\n", repoDir)
			fmt.Printf("Enabled:    %t\n", enabled)
			fmt.Printf("Strategy:   %s\n", cfg.Strategy)
			fmt.Printf("Branch:     %s\n", branch)

			if !enabled {
				fmt.Println("\nRun 'open-entire enable' to start capturing sessions.")
				return nil
			}

			// Show active sessions
			store, err := session.NewStore(repoDir)
			if err == nil {
				sessions := store.ActiveSessions()
				if len(sessions) > 0 {
					fmt.Printf("\nActive Sessions: %d\n", len(sessions))
					for _, s := range sessions {
						fmt.Printf("  - %s (%s)\n", s.ID, s.AgentName)
					}
				}
			}

			// Show checkpoint count
			count, err := repo.CheckpointCount()
			if err == nil {
				fmt.Printf("Checkpoints: %d\n", count)
			}

			if cfgDetailed {
				data, _ := json.MarshalIndent(cfg, "", "  ")
				fmt.Printf("\nConfig:\n%s\n", data)
			}

			return nil
		},
	}
	return cmd
}
