package attribution

import (
	"log/slog"
	"strings"

	"github.com/yibudak/open-entire/internal/git"
	"github.com/yibudak/open-entire/pkg/types"
)

// Tracker tracks line attribution for a repository.
type Tracker struct {
	repo *git.Repository
}

// NewTracker creates an attribution tracker.
func NewTracker(repo *git.Repository) *Tracker {
	return &Tracker{repo: repo}
}

// ForCommit calculates attribution for a specific commit.
func (t *Tracker) ForCommit(commitHash string, agentFiles []string) types.Attribution {
	added, _, err := t.repo.DiffLinesChanged(commitHash)
	if err != nil {
		slog.Debug("could not compute diff stats", "commit", commitHash, "error", err)
		return types.Attribution{}
	}

	// Heuristic: files touched by the agent are agent-authored lines
	agentFileSet := make(map[string]bool)
	for _, f := range agentFiles {
		agentFileSet[f] = true
	}

	changedFiles, err := t.repo.DiffFiles(commitHash)
	if err != nil {
		return Calculate(added, 0)
	}

	agentLines := 0
	humanLines := 0
	for _, f := range changedFiles {
		if agentFileSet[f] {
			// Count lines from this file's diff as agent-authored
			agentLines += countFileLines(t.repo, commitHash, f)
		} else {
			humanLines += countFileLines(t.repo, commitHash, f)
		}
	}

	if agentLines+humanLines == 0 {
		return Calculate(added, 0)
	}

	return Calculate(agentLines, humanLines)
}

func countFileLines(repo *git.Repository, commitHash, file string) int {
	diff, err := repo.DiffContent(commitHash)
	if err != nil {
		return 0
	}

	count := 0
	inFile := false
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "diff --git") {
			inFile = strings.Contains(line, file)
		}
		if inFile && strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			count++
		}
	}
	return count
}
