package hugolib

import (
	"bytes"
	"github.com/spf13/hugo/target"
	"html/template"
	"io"
	"testing"
)

const SLUG_DOC_1 = "---\ntitle: slug doc 1\nslug: slug-doc-1\naliases:\n - sd1/foo/\n - sd2\n - sd3/\n - sd4.html\n---\nslug doc 1 content"

//const SLUG_DOC_1 = "---\ntitle: slug doc 1\nslug: slug-doc-1\n---\nslug doc 1 content"
const SLUG_DOC_2 = "---\ntitle: slug doc 2\nslug: slug-doc-2\n---\nslug doc 2 content"

const INDEX_TEMPLATE = "{{ range .Data.Pages }}.{{ end }}"

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustReturn(ret *Page, err error) *Page {
	if err != nil {
		panic(err)
	}
	return ret
}

type InMemoryTarget struct {
	files map[string][]byte
}

func (t *InMemoryTarget) Publish(label string, reader io.Reader) (err error) {
	if t.files == nil {
		t.files = make(map[string][]byte)
	}
	bytes := new(bytes.Buffer)
	bytes.ReadFrom(reader)
	t.files[label] = bytes.Bytes()
	return
}

func (t *InMemoryTarget) Translate(label string) (dest string, err error) {
	return label, nil
}

type InMemoryAliasTarget struct {
	target.HTMLRedirectAlias
	files map[string][]byte
}

func (t *InMemoryAliasTarget) Publish(label string, permalink template.HTML) (err error) {
	f, _ := t.Translate(label)
	t.files[f] = []byte("--dummy text--")
	return
}

var urlFakeSource = []byteSource{
	{"content/blue/doc1.md", []byte(SLUG_DOC_1)},
	{"content/blue/doc2.md", []byte(SLUG_DOC_2)},
}

func TestPageCount(t *testing.T) {
	files := make(map[string][]byte)
	target := &InMemoryTarget{files: files}
	alias := &InMemoryAliasTarget{files: files}
	s := &Site{
		Target: target,
		Alias:  alias,
		Config: Config{UglyUrls: false},
		Source: &inMemorySource{urlFakeSource},
	}
	s.initializeSiteInfo()
	s.prepTemplates()
	must(s.addTemplate("indexes/blue.html", INDEX_TEMPLATE))

	if err := s.CreatePages(); err != nil {
		t.Errorf("Unable to create pages: %s", err)
	}
	if err := s.BuildSiteMeta(); err != nil {
		t.Errorf("Unable to build site metadata: %s", err)
	}

	if err := s.RenderLists(); err != nil {
		t.Errorf("Unable to render site lists: %s", err)
	}

	if err := s.RenderAliases(); err != nil {
		t.Errorf("Unable to render site lists: %s", err)
	}

	blueIndex := target.files["blue"]
	if blueIndex == nil {
		t.Errorf("No indexed rendered. %v", target.files)
	}

	expected := "<html><head></head><body>..</body></html>"
	if string(blueIndex) != expected {
		t.Errorf("Index template does not match expected: %q, got: %q", expected, string(blueIndex))
	}

	for _, s := range []string{
		"sd1/foo/index.html",
		"sd2/index.html",
		"sd3/index.html",
		"sd4.html",
	} {
		if _, ok := target.files[s]; !ok {
			t.Errorf("No alias rendered: %s", s)
		}
	}
}
