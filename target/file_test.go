package target

import (
	"testing"
)

func TestFileTranslator(t *testing.T) {
	tests := []struct {
		content  string
		expected string
	}{
		{"/", "index.html"},
		{"index.html", "index/index.html"},
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

func TestFileTranslatorBase(t *testing.T) {
	tests := []struct {
		content string
		expected string
	}{
		{"/", "a/base/index.html"},
	}

	for _, test := range tests {
		f := &Filesystem{PublishDir: "a/base"}
		fts := &Filesystem{PublishDir: "a/base/"}

		for _, fs := range []*Filesystem{f, fts} {
			dest, err := fs.Translate(test.content)
			if err != nil {
				t.Fatalf("Translated returned and err: %s", err)
			}

			if dest != test.expected {
				t.Errorf("Translate expected: %s, got: %s", test.expected, dest)
			}
		}
	}
}

func TestTranslateUglyUrls(t *testing.T) {
	tests := []struct {
		content  string
		expected string
	}{
		{"foo.html", "foo.html"},
		{"/", "index.html"},
		{"index.html", "index.html"},
	}

	for _, test := range tests {
		f := &Filesystem{UglyUrls: true}
		dest, err := f.Translate(test.content)
		if err != nil {
			t.Fatalf("Translate returned an unexpected err: %s", err)
		}

		if dest != test.expected {
			t.Errorf("Translate expected return: %s, got: %s", test.expected, dest)
		}
	}
}

func TestTranslateDefaultExtension(t *testing.T) {
	f := &Filesystem{DefaultExtension: ".foobar"}
	dest, _ := f.Translate("baz")
	if dest != "baz/index.foobar" {
		t.Errorf("Translate expected return: %s, got %s", "baz/index.foobar", dest)
	}
}
