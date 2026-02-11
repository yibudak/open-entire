package claude

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodePath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"/Users/foo/myrepo", "-Users-foo-myrepo"},
		{"/home/dev/project", "-home-dev-project"},
		{"/", "-"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, EncodePath(tt.input))
		})
	}
}
