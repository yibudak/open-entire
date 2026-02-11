package git

import (
	"fmt"
	"strings"
)

const (
	TrailerCheckpoint  = "Entire-Checkpoint"
	TrailerAttribution = "Entire-Attribution"
)

// AddTrailer appends an Entire trailer to the last commit message.
func (r *Repository) AddTrailer(key, value string) error {
	msg, err := r.LastCommitMessage()
	if err != nil {
		return err
	}

	// Check if trailer already exists
	if strings.Contains(msg, key+":") {
		return nil
	}

	trailer := fmt.Sprintf("%s: %s", key, value)

	// Amend the commit with the trailer
	newMsg := msg + "\n\n" + trailer
	_, err = r.run("git", "commit", "--amend", "-m", newMsg)
	return err
}

// FormatAttributionTrailer formats an attribution trailer value.
func FormatAttributionTrailer(percent float64, agentLines, totalLines int) string {
	return fmt.Sprintf("%.0f%% agent (%d/%d lines)", percent, agentLines, totalLines)
}

// ParseCheckpointTrailer extracts a checkpoint ID from a commit message.
func ParseCheckpointTrailer(message string) string {
	for _, line := range strings.Split(message, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, TrailerCheckpoint+":") {
			return strings.TrimSpace(strings.TrimPrefix(line, TrailerCheckpoint+":"))
		}
	}
	return ""
}
