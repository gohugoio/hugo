package commands

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

// Issue #1133
func TestNewContentPathSectionWithForwardSlashes(t *testing.T) {
	p, s := newContentPathSection("/post/new.md")
	assert.Equal(t, filepath.FromSlash("/post/new.md"), p)
	assert.Equal(t, "post", s)
}
