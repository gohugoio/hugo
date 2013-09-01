package target

import (
	"testing"
)

func TestFileTranslator(t *testing.T) {
	tests := []struct {
		content  string
		expected string
	}{
		{"foo", "foo/index.html"},
		{"foo.html", "foo/index.html"},
		{"foo.xhtml", "foo/index.xhtml"},
		{"section/foo", "section/foo/index.html"},
		{"section/foo.html", "section/foo/index.html"},
		{"section/foo.rss", "section/foo/index.rss"},
	}

	for _, test := range tests {
		f := new(Filesystem)
		dest, err := f.Translate(test.content)
		if err != nil {
			t.Fatalf("Translate returned and unexpected err: %s", err)
		}

		if dest != test.expected {
			t.Errorf("Tranlate expected return: %s, got: %s", test.expected, dest)
		}
	}
}

func TestTranslateUglyUrls(t *testing.T) {
	f := &Filesystem{UglyUrls: true}
	dest, err := f.Translate("foo.html")
	if err != nil {
		t.Fatalf("Translate returned an unexpected err: %s", err)
	}

	if dest != "foo.html" {
		t.Errorf("Translate expected return: %s, got: %s", "foo.html", dest)
	}
}

func TestTranslateDefaultExtension(t *testing.T) {
	f := &Filesystem{DefaultExtension: ".foobar"}
	dest, _ := f.Translate("baz")
	if dest != "baz/index.foobar" {
		t.Errorf("Translate expected return: %s, got %s", "baz/index.foobar", dest)
	}
}
