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
	p, err := NewPageFrom(strings.NewReader(SIMPLE_PAGE_YAML), path.Join("foobar"))
	if err != nil {
		t.Fatalf("Error in NewPageFrom")
	}
	if p.Section != "" {
		t.Fatalf("No section should be set for a file path: foobar")
	}
}

func TestNewPageWithFilePath(t *testing.T) {
	toCheck := []struct {
		input   string
		section string
		layout  []string
	}{
		{path.Join("sub", "foobar.html"), "sub", L("sub/single.html", "_default/single.html")},
		{path.Join("content", "foobar.html"), "", L("page/single.html", "_default/single.html")},
		{path.Join("content", "sub", "foobar.html"), "sub", L("sub/single.html", "_default/single.html")},
		{path.Join("content", "dub", "sub", "foobar.html"), "dub/sub", L("dub/sub/single.html", "dub/single.html", "_default/single.html")},
	}

	for _, el := range toCheck {
		p, err := NewPageFrom(strings.NewReader(SIMPLE_PAGE_YAML), el.input)
		p.guessSection()
		if err != nil {
			t.Errorf("Reading from SIMPLE_PAGE_YAML resulted in an error: %s", err)
		}
		if p.Section != el.section {
			t.Errorf("Section not set to %s for page %s. Got: %s", el.section, el.input, p.Section)
		}

		for _, y := range el.layout {
			el.layout = append(el.layout, "theme/"+y)
		}

		if !listEqual(p.Layout(), el.layout) {
			t.Errorf("Layout incorrect. Expected: '%s', Got: '%s'", el.layout, p.Layout())
		}
	}
}
