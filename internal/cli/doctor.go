package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/yibudak/open-entire/internal/session"
)

func newDoctorCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Scan and fix stuck sessions",
		Long:  "Find sessions stuck in ACTIVE state and offer to fix them.",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoDir, err := findRepoRoot()
			if err != nil {
				return fmt.Errorf("not a git repository: %w", err)
			}

			store, err := session.NewStore(repoDir)
			if err != nil {
				return fmt.Errorf("failed to open session store: %w", err)
			}

			stuck := store.StuckSessions()
			if len(stuck) == 0 {
				fmt.Println("No stuck sessions found. Everything looks good!")
				return nil
			}

			fmt.Printf("Found %d stuck session(s):\n", len(stuck))
			for _, s := range stuck {
				fmt.Printf("  - %s (%s, started %s)\n", s.ID, s.AgentName, s.StartedAt.Format("2006-01-02 15:04"))
			}

			if !force {
				fmt.Println("\nUse --force to automatically end stuck sessions.")
				return nil
			}

			for _, s := range stuck {
				if err := store.EndSession(s.ID); err != nil {
					slog.Warn("failed to end session", "id", s.ID, "error", err)
				} else {
					fmt.Printf("  Fixed: %s\n", s.ID)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "automatically fix stuck sessions")

	return cmd
}
