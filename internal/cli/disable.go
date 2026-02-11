package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yibudak/open-entire/internal/hooks"
)

func newDisableCmd() *cobra.Command {
	var (
		uninstall bool
		force     bool
		project   bool
	)

	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Remove Entire hooks from repository",
		Long:  "Remove Git hooks installed by Entire. Data is preserved unless --uninstall is used.",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoDir, err := findRepoRoot()
			if err != nil {
				return fmt.Errorf("not a git repository: %w", err)
			}

			if err := hooks.Remove(repoDir); err != nil {
				return fmt.Errorf("failed to remove hooks: %w", err)
			}

			fmt.Println("Entire hooks removed.")

			if uninstall {
				fmt.Println("Session data preserved in .entire/ directory.")
				fmt.Println("To fully remove, delete the .entire/ directory manually.")
			}

			_ = force
			_ = project
			return nil
		},
	}

	cmd.Flags().BoolVar(&uninstall, "uninstall", false, "fully uninstall (preserves data)")
	cmd.Flags().BoolVar(&force, "force", false, "skip confirmation")
	cmd.Flags().BoolVar(&project, "project", false, "disable for project only")

	return cmd
}
