package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatAttributionTrailer(t *testing.T) {
	result := FormatAttributionTrailer(73, 146, 200)
	assert.Equal(t, "73% agent (146/200 lines)", result)
}

func TestParseCheckpointTrailer(t *testing.T) {
	msg := `feat: Add user authentication

Some description here.

Entire-Checkpoint: a3b2c4d5e6f7
Entire-Attribution: 73% agent (146/200 lines)`

	id := ParseCheckpointTrailer(msg)
	assert.Equal(t, "a3b2c4d5e6f7", id)
}

func TestParseCheckpointTrailerNotFound(t *testing.T) {
	msg := "feat: Normal commit without trailers"
	id := ParseCheckpointTrailer(msg)
	assert.Equal(t, "", id)
}
