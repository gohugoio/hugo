package hugolib

import (
	"path/filepath"
	"strings"
	"testing"
)

var SIMPLE_PAGE_YAML = `---
contenttype: ""
---
Sample Text
`

func TestDegenerateMissingFolderInPageFilename(t *testing.T) {
	p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE_YAML), filepath.Join("foobar"))
	if err != nil {
		t.Fatalf("Error in ReadFrom")
	}
	if p.Section != "" {
		t.Fatalf("No section should be set for a file path: foobar")
	}
}

func TestNewPageWithFilePath(t *testing.T) {
	toCheck := []struct{
		input string
		section string
		layout string
	}{
		{filepath.Join("sub", "foobar.html"), "sub", "sub/single.html"},
		{filepath.Join("content", "sub", "foobar.html"), "sub", "sub/single.html"},
		{filepath.Join("content", "dub", "sub", "foobar.html"), "sub", "sub/single.html"},
	}

	for _, el := range toCheck {
		p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE_YAML), el.input)
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


