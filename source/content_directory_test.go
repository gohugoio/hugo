package source

import (
	"github.com/spf13/viper"
	"testing"
)

func TestIgnoreDotFilesAndDirectories(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	tests := []struct {
		path                string
		ignore              bool
		ignoreFilesRegexpes interface{}
	}{
		{".foobar/", true, nil},
		{"foobar/.barfoo/", true, nil},
		{"barfoo.md", false, nil},
		{"foobar/barfoo.md", false, nil},
		{"foobar/.barfoo.md", true, nil},
		{".barfoo.md", true, nil},
		{".md", true, nil},
		{"", true, nil},
		{"foobar/barfoo.md~", true, nil},
		{".foobar/barfoo.md~", true, nil},
		{"foobar~/barfoo.md", false, nil},
		{"foobar/bar~foo.md", false, nil},
		{"foobar/foo.md", true, []string{"\\.md$", "\\.boo$"}},
		{"foobar/foo.html", false, []string{"\\.md$", "\\.boo$"}},
		{"foobar/foo.md", true, []string{"^foo"}},
		{"foobar/foo.md", false, []string{"*", "\\.md$", "\\.boo$"}},
	}

	for _, test := range tests {

		viper.Set("ignoreFiles", test.ignoreFilesRegexpes)

		if ignored := isNonProcessablePath(test.path); test.ignore != ignored {
			t.Errorf("File not ignored.  Expected: %t, got: %t", test.ignore, ignored)
		}
	}
}
