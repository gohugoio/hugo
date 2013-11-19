package hugolib

import (
	"html/template"
	"testing"
)

func TestPermalink(t *testing.T) {
	tests := []struct {
		file        string
		dir         string
		base        template.URL
		slug        string
		expectedAbs string
		expectedRel string
	}{
		{"x/y/z/boofar.md", "x/y/z", "", "", "/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "x/y/z/", "", "", "/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "x/y/z/", "", "boofar", "/x/y/z/boofar/", "/x/y/z/boofar/"},
		{"x/y/z/boofar.md", "x/y/z", "http://barnew/", "", "http://barnew/x/y/z/boofar", "/x/y/z/boofar"},
		{"x/y/z/boofar.md", "x/y/z/", "http://barnew/", "boofar", "http://barnew/x/y/z/boofar/", "/x/y/z/boofar/"},
	}

	for _, test := range tests {
		p := &Page{
			Node: Node{
				UrlPath: UrlPath{Section: "z"},
				Site:    SiteInfo{BaseUrl: test.base},
			},
			File: File{FileName: test.file, Dir: test.dir},
		}

		if test.slug != "" {
			p.update(map[string]interface{}{
				"slug": test.slug,
			})
		}

		u, err := p.Permalink()
		if err != nil {
			t.Errorf("Unable to process permalink: %s", err)
		}

		expected := test.expectedAbs
		if u != expected {
			t.Errorf("Expected abs url: %s, got: %s", expected, u)
		}

		u, err = p.RelPermalink()
		if err != nil {
			t.Errorf("Unable to process permalink: %s", err)
		}

		expected = test.expectedRel
		if u != expected {
			t.Errorf("Expected abs url: %s, got: %s", expected, u)
		}
	}
}
