package hugolib

import (
	"html/template"
	"testing"
)

func TestPermalink(t *testing.T) {
	tests := []struct {
		base     template.URL
		expectedAbs string
		expectedRel string
	}{
		{"", "/x/y/z/boofar", "/x/y/z/boofar"},
		{"http://barnew/", "http://barnew/x/y/z/boofar", "/x/y/z/boofar"},
	}

	for _, test := range tests {
		p := &Page{
			Node: Node{
				UrlPath: UrlPath{Section: "x/y/z"},
				Site:    SiteInfo{BaseUrl: test.base},
			},
			File: File{FileName: "x/y/z/boofar.md"},
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
