package hugolib

import (
	"testing"
)

// TestTOCPanicMinimal reproduces the index out of range panic in the TOC transformer
func TestTOCPanicMinimal(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = "http://example.com/"
[markup.goldmark.extensions.passthrough]
enable = true
[markup.goldmark.extensions.passthrough.delimiters]
inline = [['$', '$']]
-- content/p1.md --
---
title: p1
---
# **$a$**
-- layouts/_default/single.html --
{{ .Summary }}
`

	// This should not panic.
	// Before the fix, this would panic with "index out of range [41] with length 41"
	// (or a similar ID depending on the Goldmark version/build).
	Test(t, files)
}
