package hooks

import (
	"fmt"
	"os"
	"path/filepath"
)

const entireMarker = "# managed by entire"

// Install installs Entire git hooks into the repository.
func Install(repoDir string, force bool) error {
	hooksDir := filepath.Join(repoDir, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return err
	}

	hookFiles := map[string]string{
		"post-commit": postCommitScript,
		"pre-push":    prePushScript,
	}

	for name, script := range hookFiles {
		path := filepath.Join(hooksDir, name)

		// Check for existing hook
		if _, err := os.Stat(path); err == nil && !force {
			data, _ := os.ReadFile(path)
			if len(data) > 0 && !isEntireHook(string(data)) {
				return fmt.Errorf("hook %s already exists (use --force to overwrite)", name)
			}
		}

		if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
			return fmt.Errorf("failed to write %s hook: %w", name, err)
		}
	}

	return nil
}

// Remove removes Entire git hooks from the repository.
func Remove(repoDir string) error {
	hooksDir := filepath.Join(repoDir, ".git", "hooks")

	for _, name := range []string{"post-commit", "pre-push"} {
		path := filepath.Join(hooksDir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if isEntireHook(string(data)) {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove %s hook: %w", name, err)
			}
		}
	}

	return nil
}

// IsInstalled checks if Entire hooks are installed.
func IsInstalled(repoDir string) bool {
	path := filepath.Join(repoDir, ".git", "hooks", "post-commit")
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return isEntireHook(string(data))
}

func isEntireHook(content string) bool {
	return len(content) > 0 && (content[0:min(len(content), 100)] != "" && containsMarker(content))
}

func containsMarker(s string) bool {
	for i := 0; i <= len(s)-len(entireMarker); i++ {
		if s[i:i+len(entireMarker)] == entireMarker {
			return true
		}
	}
	return false
}
