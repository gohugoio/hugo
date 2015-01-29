package source

import (
	"testing"
)

func TestIgnoreDotFilesAndDirectories(t *testing.T) {
	tests := []struct {
		path   string
		ignore bool
	}{
		{".foobar/", true},
		{"foobar/.barfoo/", true},
		{"barfoo.md", false},
		{"foobar/barfoo.md", false},
		{"foobar/.barfoo.md", true},
		{".barfoo.md", true},
		{".md", true},
		{"", true},
		{"foobar/barfoo.md~", true},
		{".foobar/barfoo.md~", true},
		{"foobar~/barfoo.md", false},
		{"foobar/bar~foo.md", false},
	}

	for _, test := range tests {
		if ignored := isNonProcessablePath(test.path); test.ignore != ignored {
			t.Errorf("File not ignored.  Expected: %t, got: %t", test.ignore, ignored)
		}
	}
}
