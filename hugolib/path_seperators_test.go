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
	toCheck := []map[string]string{
		{"input": filepath.Join("sub", "foobar.html"), "expect": "sub"},
		{"input": filepath.Join("content", "sub", "foobar.html"), "expect": "sub"},
		{"input": filepath.Join("content", "dub", "sub", "foobar.html"), "expect": "sub"},
	}

	for _, el := range toCheck {
		p, err := ReadFrom(strings.NewReader(SIMPLE_PAGE_YAML), el["input"])
		if err != nil {
			t.Fatalf("Reading from SIMPLE_PAGE_YAML resulted in an error: %s", err)
		}
		if p.Section != el["expect"] {
			t.Fatalf("Section not set to %s for page %s. Got: %s", el["expect"], el["input"], p.Section)
		}
	}
}

func TestSettingOutFileOnPageContainsCorrectSlashes(t *testing.T) {
	s := &Site{Config: Config{}}
	p := NewPage(filepath.Join("sub", "foobar"))
	s.setOutFile(p)
}
