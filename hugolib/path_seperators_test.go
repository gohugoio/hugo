package hugolib

import (
	"path"
	"strings"
	"testing"
)

var SIMPLE_PAGE_YAML = `---
contenttype: ""
---
Sample Text
`

func TestDegenerateMissingFolderInPageFilename(t *testing.T) {
	p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE_YAML), path.Join("foobar"))
	if err != nil {
		t.Fatalf("Error in ReadFrom")
	}
	if p.Section != "" {
		t.Fatalf("No section should be set for a file path: foobar")
	}
}

func TestNewPageWithFilePath(t *testing.T) {
	toCheck := []struct {
		input   string
		section string
		layout  string
	}{
		{path.Join("sub", "foobar.html"), "sub", "sub/single.html"},
		{path.Join("content", "sub", "foobar.html"), "sub", "sub/single.html"},
		{path.Join("content", "dub", "sub", "foobar.html"), "sub", "sub/single.html"},
	}

	for _, el := range toCheck {
		p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE_YAML), el.input)
		p.guessSection()
		if err != nil {
			t.Fatalf("Reading from SIMPLE_PAGE_YAML resulted in an error: %s", err)
		}
		if p.Section != el.section {
			t.Fatalf("Section not set to %s for page %s. Got: %s", el.section, el.input, p.Section)
		}

		if p.Layout() != el.layout {
			t.Fatalf("Layout incorrect. Expected: '%s', Got: '%s'", el.layout, p.Layout())
		}
	}
}
