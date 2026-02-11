package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yibudak/open-entire/internal/git"
)

func newResetCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Delete shadow branch and session state",
		Long:  "Remove Entire's shadow branch and local session state for a clean start.",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoDir, err := findRepoRoot()
			if err != nil {
				return fmt.Errorf("not a git repository: %w", err)
			}

			if !force {
				fmt.Println("This will delete all local Entire state including active sessions.")
				fmt.Println("Use --force to confirm.")
				return nil
			}

			repo, err := git.Open(repoDir)
			if err != nil {
				return fmt.Errorf("failed to open repository: %w", err)
			}

			// Delete shadow branches
			branches, _ := repo.OrphanedShadowBranches()
			for _, b := range branches {
				_ = repo.DeleteBranch(b)
			}

			// Remove local state
			statePath := filepath.Join(repoDir, ".entire", "state.json")
			if err := os.Remove(statePath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove state: %w", err)
			}

			fmt.Println("Entire state reset successfully.")
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "confirm reset")

	return cmd
}
