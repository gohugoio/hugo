package hugolib

import (
	"path/filepath"
	"testing"
)

func TestDegenerateMissingFolderInPageFilename(t *testing.T) {
	p := NewPage(filepath.Join("foobar"))
	if p.Section != "" {
		t.Fatalf("No section should be set for a file path: foobar")
	}
}

func TestCreateNewPage(t *testing.T) {
	toCheck := []map[string]string{
		{"input": filepath.Join("sub", "foobar.html"), "expect": "sub"},
		{"input": filepath.Join("content", "sub", "foobar.html"), "expect": "sub"},
		{"input": filepath.Join("content", "dub", "sub", "foobar.html"), "expect": "sub"},
	}

	for _, el := range toCheck {
		p := NewPage(el["input"])
		if p.Section != el["expect"] {
			t.Fatalf("Section not set to %s for page %s. Got: %s", el["expect"], el["input"], p.Section)
		}
	}
}


func TestSettingOutFileOnPageContainsCorrectSlashes(t *testing.T) {
	s := NewSite(&Config{})
	p := NewPage(filepath.Join("sub", "foobar"))
	s.setOutFile(p)
}
