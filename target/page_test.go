package target

import (
	"path/filepath"
	"testing"
)

func TestPageTranslator(t *testing.T) {
	tests := []struct {
		content  string
		expected string
	}{
		{"/", "index.html"},
		{"index.html", "index.html"},
		{"bar/index.html", "bar/index.html"},
		{"foo", "foo/index.html"},
		{"foo.html", "foo/index.html"},
		{"foo.xhtml", "foo/index.xhtml"},
		{"section", "section/index.html"},
		{"section/", "section/index.html"},
		{"section/foo", "section/foo/index.html"},
		{"section/foo.html", "section/foo/index.html"},
		{"section/foo.rss", "section/foo/index.rss"},
	}

	for _, test := range tests {
		f := new(PagePub)
		dest, err := f.Translate(filepath.FromSlash(test.content))
		expected := filepath.FromSlash(test.expected)
		if err != nil {
			t.Fatalf("Translate returned and unexpected err: %s", err)
		}

		if dest != expected {
			t.Errorf("Translate expected return: %s, got: %s", expected, dest)
		}
	}
}

func TestPageTranslatorBase(t *testing.T) {
	tests := []struct {
		content  string
		expected string
	}{
		{"/", "a/base/index.html"},
	}

	for _, test := range tests {
		f := &PagePub{PublishDir: "a/base"}
		fts := &PagePub{PublishDir: "a/base/"}

		for _, fs := range []*PagePub{f, fts} {
			dest, err := fs.Translate(test.content)
			if err != nil {
				t.Fatalf("Translated returned and err: %s", err)
			}

			if dest != filepath.FromSlash(test.expected) {
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
		{"section", "section.html"},
		{"index.html", "index.html"},
	}

	for _, test := range tests {
		f := &PagePub{UglyUrls: true}
		dest, err := f.Translate(filepath.FromSlash(test.content))
		if err != nil {
			t.Fatalf("Translate returned an unexpected err: %s", err)
		}

		if dest != test.expected {
			t.Errorf("Translate expected return: %s, got: %s", test.expected, dest)
		}
	}
}

func TestTranslateDefaultExtension(t *testing.T) {
	f := &PagePub{DefaultExtension: ".foobar"}
	dest, _ := f.Translate("baz")
	if dest != filepath.FromSlash("baz/index.foobar") {
		t.Errorf("Translate expected return: %s, got %s", "baz/index.foobar", dest)
	}
}
