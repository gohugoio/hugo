package hugolib

import (
	"testing"
	"path/filepath"
)

func TestDegenerateMissingFolderInPageFilename(t *testing.T) {
	p := NewPage(filepath.Join("foobar"))
	if p != nil {
		t.Fatalf("Creating a new Page without a subdirectory should result in nil page")
	}
}

func TestSettingOutFileOnPageContainsCorrectSlashes(t *testing.T) {
	s := NewSite(&Config{})
	p := NewPage(filepath.Join("sub", "foobar"))
	s.setOutFile(p)
}
