package source

import (
	"testing"
)

func TestIgnoreDotFiles(t *testing.T) {
	tests := []struct {
		path   string
		ignore bool
	}{
		{"barfoo.md", false},
		{"foobar/barfoo.md", false},
		{"foobar/.barfoo.md", true},
		{".barfoo.md", true},
		{".md", true},
		{"", true},
	}

	for _, test := range tests {
		if ignored := ignoreDotFile(test.path); test.ignore != ignored {
			t.Errorf("File not ignored.  Expected: %t, got: %t", test.ignore, ignored)
		}
	}
}
