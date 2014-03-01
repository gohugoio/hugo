package hugolib

import (
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"html/template"
	"testing"
)

const SLUG_DOC_1 = "---\ntitle: slug doc 1\nslug: slug-doc-1\naliases:\n - sd1/foo/\n - sd2\n - sd3/\n - sd4.html\n---\nslug doc 1 content\n"

const SLUG_DOC_2 = `---
title: slug doc 2
slug: slug-doc-2
---
slug doc 2 content
`

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

type InMemoryAliasTarget struct {
	target.HTMLRedirectAlias
	files map[string][]byte
}

func (t *InMemoryAliasTarget) Publish(label string, permalink template.HTML) (err error) {
	f, _ := t.Translate(label)
	t.files[f] = []byte("--dummy text--")
	return
}

var urlFakeSource = []source.ByteSource{
	{"content/blue/doc1.md", []byte(SLUG_DOC_1), "blue"},
	{"content/blue/doc2.md", []byte(SLUG_DOC_2), "blue"},
}

func TestPageCount(t *testing.T) {
	files := make(map[string][]byte)
	target := &target.InMemoryTarget{Files: files}
	alias := &InMemoryAliasTarget{files: files}
	s := &Site{
		Target: target,
		Alias:  alias,
		Config: Config{UglyUrls: false},
		Source: &source.InMemorySource{ByteSource: urlFakeSource},
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

	blueIndex := target.Files["blue"]
	if blueIndex == nil {
		t.Errorf("No indexed rendered. %v", target.Files)
	}

	expected := ".."
	if string(blueIndex) != expected {
		t.Errorf("Index template does not match expected: %q, got: %q", expected, string(blueIndex))
	}

	for _, s := range []string{
		"sd1/foo/index.html",
		"sd2/index.html",
		"sd3/index.html",
		"sd4.html",
	} {
		if _, ok := target.Files[s]; !ok {
			t.Errorf("No alias rendered: %s", s)
		}
	}
}
