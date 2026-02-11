package checkpoint

// Shard extracts the shard prefix (first 2 chars) from a checkpoint ID.
func Shard(id string) string {
	if len(id) < 2 {
		return id
	}
	return id[:2]
}

// ShardRemainder returns the remaining chars after the shard prefix.
func ShardRemainder(id string) string {
	if len(id) <= 2 {
		return ""
	}
	return id[2:]
}
