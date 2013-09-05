package hugolib

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"
	"testing"
)

const SLUG_DOC_1 = "---\ntitle: slug doc 1\nslug: slug-doc-1\n---\nslug doc 1 content"
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

func TestPageCount(t *testing.T) {
	target := new(InMemoryTarget)
	s := &Site{Target: target}
	s.prepTemplates()
	must(s.addTemplate("indexes/blue.html", INDEX_TEMPLATE))
	//s.Files = append(s.Files, "blue/doc1.md")
	//s.Files = append(s.Files, "blue/doc2.md")
	s.Pages = append(s.Pages, mustReturn(ReadFrom(strings.NewReader(SLUG_DOC_1), filepath.FromSlash("content/blue/doc1.md"))))
	s.Pages = append(s.Pages, mustReturn(ReadFrom(strings.NewReader(SLUG_DOC_2), filepath.FromSlash("content/blue/doc2.md"))))

	if err := s.BuildSiteMeta(); err != nil {
		t.Errorf("Unable to build site metadata: %s", err)
	}

	if err := s.RenderLists(); err != nil {
		t.Errorf("Unable to render site lists: %s", err)
	}

	blueIndex := target.files["blue/index.html"]
	if blueIndex == nil {
		t.Errorf("No indexed rendered. %v", target.files)
	}

	if len(blueIndex) != 2 {
		t.Errorf("Number of pages does not equal 2, got %d. %q", len(blueIndex), blueIndex)
	}
}
