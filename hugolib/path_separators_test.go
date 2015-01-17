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
	p, err := NewPageFrom(strings.NewReader(SIMPLE_PAGE_YAML), filepath.Join("foobar"))
	if err != nil {
		t.Fatalf("Error in NewPageFrom")
	}
	if p.Section() != "" {
		t.Fatalf("No section should be set for a file path: foobar")
	}
}

func TestNewPageWithFilePath(t *testing.T) {
	toCheck := []struct {
		input   string
		section string
		layout  []string
	}{
		{filepath.Join("sub", "foobar.html"), "sub", L("sub/single.html", "_default/single.html")},
		{filepath.Join("content", "foobar.html"), "", L("page/single.html", "_default/single.html")},
		{filepath.Join("content", "sub", "foobar.html"), "sub", L("sub/single.html", "_default/single.html")},
		{filepath.Join("content", "dub", "sub", "foobar.html"), "dub", L("dub/single.html", "_default/single.html")},
	}

	for _, el := range toCheck {
		p, err := NewPageFrom(strings.NewReader(SIMPLE_PAGE_YAML), el.input)
		if err != nil {
			t.Errorf("Reading from SIMPLE_PAGE_YAML resulted in an error: %s", err)
		}
		if p.Section() != el.section {
			t.Errorf("Section not set to %s for page %s. Got: %s", el.section, el.input, p.Section())
		}

		for _, y := range el.layout {
			el.layout = append(el.layout, "theme/"+y)
		}

		if !listEqual(p.Layout(), el.layout) {
			t.Errorf("Layout incorrect. Expected: '%s', Got: '%s'", el.layout, p.Layout())
		}
	}
}
