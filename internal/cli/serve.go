package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/yibudak/open-entire/internal/git"
	"github.com/yibudak/open-entire/internal/web"
)

func newServeCmd() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start local web viewer",
		Long:  "Launch a local web server to browse checkpoints and sessions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoDir, err := findRepoRoot()
			if err != nil {
				return fmt.Errorf("not a git repository: %w", err)
			}

			repo, err := git.Open(repoDir)
			if err != nil {
				return fmt.Errorf("failed to open repository: %w", err)
			}

			addr := fmt.Sprintf(":%d", port)
			slog.Info("starting web viewer", "addr", addr, "repo", repoDir)
			fmt.Printf("Entire web viewer running at http://localhost:%d\n", port)

			srv := web.NewServer(repo, repoDir)
			return srv.ListenAndServe(addr)
		},
	}

	cmd.Flags().IntVar(&port, "port", 8080, "port to listen on")

	return cmd
}
