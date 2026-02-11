package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yibudak/open-entire/internal/checkpoint"
	"github.com/yibudak/open-entire/internal/git"
)

func newRewindCmd() *cobra.Command {
	var (
		to       string
		list     bool
		reset    bool
		logsOnly bool
	)

	cmd := &cobra.Command{
		Use:   "rewind",
		Short: "Rewind to a previous checkpoint",
		Long:  "List and restore to a previous checkpoint state.",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoDir, err := findRepoRoot()
			if err != nil {
				return fmt.Errorf("not a git repository: %w", err)
			}

			repo, err := git.Open(repoDir)
			if err != nil {
				return fmt.Errorf("failed to open repository: %w", err)
			}

			store := checkpoint.NewStore(repo)

			if list {
				checkpoints, err := store.List()
				if err != nil {
					return fmt.Errorf("failed to list checkpoints: %w", err)
				}
				if len(checkpoints) == 0 {
					fmt.Println("No checkpoints found.")
					return nil
				}
				for _, cp := range checkpoints {
					fmt.Printf("  %s  %s  %s\n", cp.ID, cp.CreatedAt.Format("2006-01-02 15:04"), cp.Message)
				}
				return nil
			}

			if to == "" {
				return fmt.Errorf("specify --to <checkpoint-id> or --list to see available checkpoints")
			}

			if logsOnly {
				fmt.Printf("Restoring logs from checkpoint %s...\n", to)
				return store.RestoreLogs(repoDir, to)
			}

			if reset {
				fmt.Printf("Hard resetting to checkpoint %s...\n", to)
			} else {
				fmt.Printf("Rewinding to checkpoint %s...\n", to)
			}

			return store.Rewind(repoDir, to, reset)
		},
	}

	cmd.Flags().StringVar(&to, "to", "", "checkpoint ID to rewind to")
	cmd.Flags().BoolVar(&list, "list", false, "list available checkpoints")
	cmd.Flags().BoolVar(&reset, "reset", false, "hard reset to checkpoint")
	cmd.Flags().BoolVar(&logsOnly, "logs-only", false, "restore logs only")

	return cmd
}
