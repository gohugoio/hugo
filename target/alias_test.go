package target

import (
	"path/filepath"
	"testing"
)

func TestHTMLRedirectAlias(t *testing.T) {
	var o Translator
	o = new(HTMLRedirectAlias)

	tests := []struct {
		value    string
		expected string
	}{
		{"", ""},
		{"s", filepath.FromSlash("s/index.html")},
		{"/", filepath.FromSlash("/index.html")},
		{"alias 1", filepath.FromSlash("alias-1/index.html")},
		{"alias 2/", filepath.FromSlash("alias-2/index.html")},
		{"alias 3.html", "alias-3.html"},
		{"alias4.html", "alias4.html"},
		{"/alias 5.html", filepath.FromSlash("/alias-5.html")},
		{"/трям.html", filepath.FromSlash("/трям.html")},
	}

	for _, test := range tests {
		path, err := o.Translate(test.value)
		if err != nil {
			t.Fatalf("Translate returned an error: %s", err)
		}

		if path != test.expected {
			t.Errorf("Expected: %s, got: %s", test.expected, path)
		}
	}
}
