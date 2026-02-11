package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	CheckpointsBranch = "entire/checkpoints/v1"
	ShadowBranchPrefix = "entire/"
)

// Repository wraps Git operations for a repository.
type Repository struct {
	Dir string
}

// Open opens a git repository at the given path.
func Open(dir string) (*Repository, error) {
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		return nil, fmt.Errorf("not a git repository: %s", dir)
	}
	return &Repository{Dir: dir}, nil
}

// IsEntireEnabled checks if Open-Entire hooks are installed.
func (r *Repository) IsEntireEnabled() bool {
	hookPath := filepath.Join(r.Dir, ".git", "hooks", "post-commit")
	data, err := os.ReadFile(hookPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "open-entire")
}

// CurrentBranch returns the current branch name.
func (r *Repository) CurrentBranch() (string, error) {
	out, err := r.run("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// Checkout checks out a branch.
func (r *Repository) Checkout(branch string, force bool) error {
	args := []string{"checkout", branch}
	if force {
		args = append(args, "--force")
	}
	_, err := r.run("git", args...)
	return err
}

// EnsureCheckpointsBranch creates the orphan checkpoints branch if it doesn't exist.
func (r *Repository) EnsureCheckpointsBranch() error {
	// Check if branch already exists
	_, err := r.run("git", "rev-parse", "--verify", CheckpointsBranch)
	if err == nil {
		return nil // Already exists
	}

	// Create orphan branch with an empty commit
	currentBranch, _ := r.CurrentBranch()

	_, err = r.run("git", "checkout", "--orphan", CheckpointsBranch)
	if err != nil {
		return fmt.Errorf("failed to create orphan branch: %w", err)
	}

	// Remove all files from index
	_ = r.runSilent("git", "rm", "-rf", "--cached", ".")
	_ = r.runSilent("git", "clean", "-fd")

	// Create initial commit
	_, err = r.run("git", "commit", "--allow-empty", "--no-verify", "-m", "Initialize entire checkpoints")
	if err != nil {
		// Try to go back
		_, _ = r.run("git", "checkout", currentBranch)
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	// Switch back to original branch
	if currentBranch != "" && currentBranch != "HEAD" {
		_, _ = r.run("git", "checkout", currentBranch)
	}

	return nil
}

// CheckpointCount returns the number of checkpoints on the checkpoints branch.
func (r *Repository) CheckpointCount() (int, error) {
	out, err := r.run("git", "log", "--oneline", CheckpointsBranch)
	if err != nil {
		return 0, err
	}
	if out == "" {
		return 0, nil
	}
	return len(strings.Split(strings.TrimSpace(out), "\n")), nil
}

// FindCheckpointTrailer searches recent commits on a branch for an Entire-Checkpoint trailer.
func (r *Repository) FindCheckpointTrailer(branch string) (string, error) {
	out, err := r.run("git", "log", "--format=%B", "-10", branch)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Entire-Checkpoint:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Entire-Checkpoint:")), nil
		}
	}
	return "", fmt.Errorf("no checkpoint trailer found")
}

// CheckpointFromCommit finds the checkpoint ID from a specific commit's trailers.
func (r *Repository) CheckpointFromCommit(hash string) (string, error) {
	out, err := r.run("git", "log", "--format=%B", "-1", hash)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Entire-Checkpoint:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Entire-Checkpoint:")), nil
		}
	}
	return "", fmt.Errorf("no checkpoint trailer on commit %s", hash)
}

// OrphanedShadowBranches returns Entire shadow branches that no longer have active sessions.
func (r *Repository) OrphanedShadowBranches() ([]string, error) {
	out, err := r.run("git", "branch", "--list", ShadowBranchPrefix+"*")
	if err != nil {
		return nil, err
	}

	var branches []string
	for _, line := range strings.Split(out, "\n") {
		b := strings.TrimSpace(strings.TrimPrefix(line, "*"))
		b = strings.TrimSpace(b)
		if b != "" && b != CheckpointsBranch && strings.HasPrefix(b, ShadowBranchPrefix) {
			branches = append(branches, b)
		}
	}
	return branches, nil
}

// DeleteBranch deletes a local branch.
func (r *Repository) DeleteBranch(name string) error {
	_, err := r.run("git", "branch", "-D", name)
	return err
}

// CommitOnBranch creates a commit on the specified branch without changing the working tree.
func (r *Repository) CommitOnBranch(branch, message string, files map[string][]byte) error {
	// Save current branch
	currentBranch, err := r.CurrentBranch()
	if err != nil {
		return err
	}

	// Checkout target branch
	if _, err := r.run("git", "checkout", branch); err != nil {
		return fmt.Errorf("failed to checkout %s: %w", branch, err)
	}

	// Write files
	for path, data := range files {
		fullPath := filepath.Join(r.Dir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			_, _ = r.run("git", "checkout", currentBranch)
			return err
		}
		if err := os.WriteFile(fullPath, data, 0o644); err != nil {
			_, _ = r.run("git", "checkout", currentBranch)
			return err
		}
	}

	// Stage and commit
	_, _ = r.run("git", "add", "-A")
	if _, err := r.run("git", "commit", "--no-verify", "-m", message); err != nil {
		_, _ = r.run("git", "checkout", currentBranch)
		return fmt.Errorf("failed to commit: %w", err)
	}

	// Return to original branch
	_, _ = r.run("git", "checkout", currentBranch)
	return nil
}

// ReadFileFromBranch reads a file from a specific branch without checkout.
func (r *Repository) ReadFileFromBranch(branch, path string) ([]byte, error) {
	out, err := r.run("git", "show", branch+":"+path)
	if err != nil {
		return nil, err
	}
	return []byte(out), nil
}

// ListFilesOnBranch lists files matching a prefix on a branch.
func (r *Repository) ListFilesOnBranch(branch, prefix string) ([]string, error) {
	out, err := r.run("git", "ls-tree", "-r", "--name-only", branch)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if line != "" && strings.HasPrefix(line, prefix) {
			files = append(files, line)
		}
	}
	return files, nil
}

// HeadCommitHash returns the HEAD commit hash.
func (r *Repository) HeadCommitHash() (string, error) {
	out, err := r.run("git", "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// LastCommitMessage returns the last commit message.
func (r *Repository) LastCommitMessage() (string, error) {
	out, err := r.run("git", "log", "-1", "--format=%B")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// DiffStat returns the diff stat for a commit.
func (r *Repository) DiffStat(commitHash string) (string, error) {
	out, err := r.run("git", "diff", "--stat", commitHash+"^", commitHash)
	if err != nil {
		// Try without parent (initial commit)
		out, err = r.run("git", "diff", "--stat", "--root", commitHash)
		if err != nil {
			return "", err
		}
	}
	return out, nil
}

// Author returns the configured git author name.
func (r *Repository) Author() string {
	out, _ := r.run("git", "config", "user.name")
	return strings.TrimSpace(out)
}

func (r *Repository) run(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = r.Dir
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%s: %s", err, string(ee.Stderr))
		}
		return "", err
	}
	return string(out), nil
}

func (r *Repository) runSilent(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = r.Dir
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
