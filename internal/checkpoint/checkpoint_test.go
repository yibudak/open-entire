package checkpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateID(t *testing.T) {
	id, err := GenerateID()
	require.NoError(t, err)
	assert.Len(t, id, 12)

	// IDs should be unique
	id2, err := GenerateID()
	require.NoError(t, err)
	assert.NotEqual(t, id, id2)
}

func TestShardPath(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"a3b2c4d5e6f7", "a3/b2c4d5e6f7/"},
		{"0f45ffa1b752", "0f/45ffa1b752/"},
		{"ab", "ab/"},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			assert.Equal(t, tt.want, ShardPath(tt.id))
		})
	}
}

func TestMetadataPath(t *testing.T) {
	assert.Equal(t, "a3/b2c4d5e6f7/metadata.json", MetadataPath("a3b2c4d5e6f7"))
}

func TestSessionPath(t *testing.T) {
	assert.Equal(t, "a3/b2c4d5e6f7/0/", SessionPath("a3b2c4d5e6f7", 0))
	assert.Equal(t, "a3/b2c4d5e6f7/1/", SessionPath("a3b2c4d5e6f7", 1))
}

func TestSessionFiles(t *testing.T) {
	files := SessionFiles("a3b2c4d5e6f7", 0)
	assert.Equal(t, "a3/b2c4d5e6f7/0/content_hash.txt", files["content_hash"])
	assert.Equal(t, "a3/b2c4d5e6f7/0/context.md", files["context"])
	assert.Equal(t, "a3/b2c4d5e6f7/0/full.jsonl", files["full"])
	assert.Equal(t, "a3/b2c4d5e6f7/0/metadata.json", files["metadata"])
	assert.Equal(t, "a3/b2c4d5e6f7/0/prompt.txt", files["prompt"])
}

func TestShard(t *testing.T) {
	assert.Equal(t, "a3", Shard("a3b2c4d5e6f7"))
	assert.Equal(t, "0f", Shard("0f45ffa1b752"))
}

func TestShardRemainder(t *testing.T) {
	assert.Equal(t, "b2c4d5e6f7", ShardRemainder("a3b2c4d5e6f7"))
}
