package git

import (
	"fmt"
	"strings"
)

// DiffFiles returns the list of files changed in a commit.
func (r *Repository) DiffFiles(commitHash string) ([]string, error) {
	out, err := r.run("git", "diff", "--name-only", commitHash+"^", commitHash)
	if err != nil {
		// Initial commit
		out, err = r.run("git", "diff", "--name-only", "--root", commitHash)
		if err != nil {
			return nil, err
		}
	}

	var files []string
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

// DiffContent returns the full diff content for a commit.
func (r *Repository) DiffContent(commitHash string) (string, error) {
	out, err := r.run("git", "diff", commitHash+"^", commitHash)
	if err != nil {
		out, err = r.run("git", "diff", "--root", commitHash)
		if err != nil {
			return "", err
		}
	}
	return out, nil
}

// DiffLinesChanged returns the number of added and removed lines in a commit.
func (r *Repository) DiffLinesChanged(commitHash string) (added, removed int, err error) {
	out, e := r.run("git", "diff", "--numstat", commitHash+"^", commitHash)
	if e != nil {
		out, e = r.run("git", "diff", "--numstat", "--root", commitHash)
		if e != nil {
			return 0, 0, e
		}
	}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line == "" {
			continue
		}
		var a, r int
		if _, err := fmt.Sscanf(line, "%d\t%d", &a, &r); err == nil {
			added += a
			removed += r
		}
	}
	return added, removed, nil
}
