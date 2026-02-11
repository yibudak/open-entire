package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/yibudak/open-entire/internal/git"
)

func newCleanCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Remove orphaned Entire data",
		Long:  "Clean up orphaned shadow branches and incomplete checkpoint data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoDir, err := findRepoRoot()
			if err != nil {
				return fmt.Errorf("not a git repository: %w", err)
			}

			repo, err := git.Open(repoDir)
			if err != nil {
				return fmt.Errorf("failed to open repository: %w", err)
			}

			// Find orphaned shadow branches
			branches, err := repo.OrphanedShadowBranches()
			if err != nil {
				return fmt.Errorf("failed to list shadow branches: %w", err)
			}

			if len(branches) == 0 {
				fmt.Println("No orphaned data found.")
				return nil
			}

			fmt.Printf("Found %d orphaned shadow branch(es):\n", len(branches))
			for _, b := range branches {
				fmt.Printf("  - %s\n", b)
			}

			if dryRun {
				fmt.Println("\nDry run — no changes made.")
				return nil
			}

			for _, b := range branches {
				if err := repo.DeleteBranch(b); err != nil {
					slog.Warn("failed to delete branch", "branch", b, "error", err)
				} else {
					slog.Debug("deleted orphaned branch", "branch", b)
				}
			}

			fmt.Printf("Cleaned %d orphaned branch(es).\n", len(branches))
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be cleaned without making changes")

	return cmd
}
