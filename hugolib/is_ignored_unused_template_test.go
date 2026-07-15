package hugolib

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestIsIgnoredUnusedTemplate(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		path     string
		want     bool
	}{
		{"exact match", []string{"/_partials/unusedpartial.html"}, "/_partials/unusedpartial.html", true},
		{"no match", []string{"/_partials/unusedpartial.html"}, "/_shortcodes/unusedshortcode.html", false},
		{"glob match", []string{"/_partials/*"}, "/_partials/unusedpartial.html", true},
		{"glob no match", []string{"/_partials/*"}, "/_shortcodes/unused.html", false},
		{"multiple patterns match", []string{"/_default/*", "/_partials/*"}, "/_partials/unused.html", true},
		{"empty patterns", []string{}, "/_partials/unused.html", false},
		{"path with leading slash, pattern without", []string{"_partials/unusedpartial.html"}, "/_partials/unusedpartial.html", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isIgnoredUnusedTemplate(tt.patterns, tt.path)
			qt.Assert(t, got, qt.Equals, tt.want)
		})
	}
}
