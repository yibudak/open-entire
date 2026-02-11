package claude

import (
	"os"
	"path/filepath"
	"strings"
)

// EncodePath converts a filesystem path to Claude's encoded format.
// e.g., /Users/foo/myrepo -> -Users-foo-myrepo
func EncodePath(path string) string {
	return strings.ReplaceAll(path, "/", "-")
}

// ProjectDir returns the Claude Code project directory for a repo.
func ProjectDir(repoDir string) string {
	home, _ := os.UserHomeDir()
	encoded := EncodePath(repoDir)
	return filepath.Join(home, ".claude", "projects", encoded)
}

// SessionFiles returns all .jsonl session files for a project.
func SessionFiles(repoDir string) ([]string, error) {
	dir := ProjectDir(repoDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".jsonl") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	return files, nil
}

// SubagentFiles returns subagent JSONL files for a session.
func SubagentFiles(repoDir, sessionID string) ([]string, error) {
	dir := filepath.Join(ProjectDir(repoDir), sessionID, "subagents")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".jsonl") {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	return files, nil
}
