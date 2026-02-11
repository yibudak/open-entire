package claude

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	agentName = "claude-code"
	// Consider sessions active if modified within this window
	activeWindow = 5 * time.Minute
)

// Detect checks if Claude Code is active for the given repo.
// Returns the session ID if found.
func Detect(repoDir string) (string, error) {
	projDir := ProjectDir(repoDir)

	entries, err := os.ReadDir(projDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no Claude Code project directory found")
		}
		return "", err
	}

	// Find recently modified .jsonl files
	type sessionFile struct {
		name    string
		modTime time.Time
	}

	var recent []sessionFile
	threshold := time.Now().Add(-activeWindow)

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(threshold) {
			recent = append(recent, sessionFile{
				name:    strings.TrimSuffix(e.Name(), ".jsonl"),
				modTime: info.ModTime(),
			})
		}
	}

	if len(recent) == 0 {
		return "", fmt.Errorf("no active Claude Code session found")
	}

	// Return the most recently modified session
	sort.Slice(recent, func(i, j int) bool {
		return recent[i].modTime.After(recent[j].modTime)
	})

	return recent[0].name, nil
}

// DetectProcess checks if a Claude Code process is running.
func DetectProcess() bool {
	// Check for claude process in /proc or using ps
	entries, err := filepath.Glob("/proc/*/comm")
	if err == nil {
		for _, entry := range entries {
			data, err := os.ReadFile(entry)
			if err == nil && strings.TrimSpace(string(data)) == "claude" {
				return true
			}
		}
	}

	// macOS fallback: check for process by looking at common locations
	// The actual process detection will use go-ps in production
	return false
}
