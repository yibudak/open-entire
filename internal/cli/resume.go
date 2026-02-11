package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yibudak/open-entire/internal/git"
)

func newResumeCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "resume [branch]",
		Short: "Resume a session from a branch",
		Long:  "Checkout a branch and find the associated session from commit trailers.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			branch := args[0]

			repoDir, err := findRepoRoot()
			if err != nil {
				return fmt.Errorf("not a git repository: %w", err)
			}

			repo, err := git.Open(repoDir)
			if err != nil {
				return fmt.Errorf("failed to open repository: %w", err)
			}

			// Checkout the branch
			if err := repo.Checkout(branch, force); err != nil {
				return fmt.Errorf("failed to checkout branch %q: %w", branch, err)
			}

			// Find checkpoint trailer in recent commits
			cpID, err := repo.FindCheckpointTrailer(branch)
			if err != nil {
				fmt.Printf("Checked out %s (no Entire checkpoint found on this branch).\n", branch)
				return nil
			}

			fmt.Printf("Resumed on branch %s\n", branch)
			fmt.Printf("  Last checkpoint: %s\n", cpID)
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "force checkout (discard changes)")

	return cmd
}
