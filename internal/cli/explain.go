package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yibudak/open-entire/internal/checkpoint"
	"github.com/yibudak/open-entire/internal/git"
)

func newExplainCmd() *cobra.Command {
	var (
		cpID          string
		commitHash    string
		generate      bool
		full          bool
		short         bool
		rawTranscript bool
	)

	cmd := &cobra.Command{
		Use:   "explain",
		Short: "Explain a session or checkpoint",
		Long:  "Read session data and display a formatted transcript or AI-generated summary.",
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

			// Resolve checkpoint ID
			id := cpID
			if commitHash != "" {
				resolved, err := repo.CheckpointFromCommit(commitHash)
				if err != nil {
					return fmt.Errorf("no checkpoint found for commit %s: %w", commitHash, err)
				}
				id = resolved
			}

			if id == "" {
				return fmt.Errorf("specify --checkpoint or --commit")
			}

			cp, err := store.Get(id)
			if err != nil {
				return fmt.Errorf("checkpoint %s not found: %w", id, err)
			}

			if rawTranscript {
				transcript, err := store.RawTranscript(id, 0)
				if err != nil {
					return err
				}
				fmt.Print(transcript)
				return nil
			}

			// Display checkpoint info
			fmt.Printf("Checkpoint: %s\n", cp.ID)
			fmt.Printf("Commit:     %s\n", cp.CommitHash)
			fmt.Printf("Branch:     %s\n", cp.Branch)
			fmt.Printf("Created:    %s\n", cp.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("Message:    %s\n", cp.Message)

			if cp.Attribution != nil {
				fmt.Printf("Attribution: %.0f%% agent (%d/%d lines)\n",
					cp.Attribution.AgentPercent, cp.Attribution.AgentLines, cp.Attribution.TotalLines)
			}

			if full {
				for i, s := range cp.Sessions {
					fmt.Printf("\n--- Session %d (%s) ---\n", i, s.AgentName)
					transcript, err := store.FormattedTranscript(id, i)
					if err != nil {
						fmt.Printf("  (error reading transcript: %v)\n", err)
						continue
					}
					fmt.Print(transcript)
				}
			}

			if short {
				fmt.Printf("\nSessions: %d\n", len(cp.Sessions))
				for _, s := range cp.Sessions {
					fmt.Printf("  - %s: %d input, %d output tokens\n",
						s.AgentName, s.TokenUsage.InputTokens, s.TokenUsage.OutputTokens)
				}
			}

			_ = generate
			return nil
		},
	}

	cmd.Flags().StringVar(&cpID, "checkpoint", "", "checkpoint ID")
	cmd.Flags().StringVar(&commitHash, "commit", "", "commit hash")
	cmd.Flags().BoolVar(&generate, "generate", false, "generate AI summary")
	cmd.Flags().BoolVar(&full, "full", false, "show full transcript")
	cmd.Flags().BoolVarP(&short, "short", "s", false, "show summary only")
	cmd.Flags().BoolVar(&rawTranscript, "raw-transcript", false, "show raw JSONL transcript")

	return cmd
}
