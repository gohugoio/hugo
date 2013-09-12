package hugolib

import (
	"bytes"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"testing"
)

const ALIAS_DOC_1 = "---\ntitle: alias doc\naliases:\n  - \"alias1/\"\n  - \"alias-2/\"\n---\naliases\n"

type byteSource struct {
	name    string
	content []byte
}

var fakeSource = []byteSource{
	{"foo/bar/file.md", []byte(SIMPLE_PAGE)},
	{"alias/test/file1.md", []byte(ALIAS_DOC_1)},
	//{"slug/test/file1.md", []byte(SLUG_DOC_1)},
}

type inMemorySource struct {
	byteSource []byteSource
}

func (i *inMemorySource) Files() (files []*source.File) {
	files = make([]*source.File, len(i.byteSource))
	for i, fake := range i.byteSource {
		files[i] = &source.File{
			Name:     fake.name,
			Contents: bytes.NewReader(fake.content),
		}
	}
	return
}

func checkShowPlanExpected(t *testing.T, s *Site, expected string) {
	out := new(bytes.Buffer)
	if err := s.ShowPlan(out); err != nil {
		t.Fatalf("ShowPlan unexpectedly returned an error: %s", err)
	}
	got := out.String()
	if got != expected {
		t.Errorf("ShowPlan expected:\n%q\ngot\n%q", expected, got)
	}
}

func TestDegenerateNoFiles(t *testing.T) {
	checkShowPlanExpected(t, new(Site), "No source files provided.\n")
}

func TestDegenerateNoTarget(t *testing.T) {
	s := &Site{
		Source: &inMemorySource{fakeSource},
	}
	must(s.CreatePages())
	expected := "foo/bar/file.md\n canonical => !no target specified!\n" +
		"alias/test/file1.md\n canonical => !no target specified!\n"
	checkShowPlanExpected(t, s, expected)
}

func TestFileTarget(t *testing.T) {
	s := &Site{
		Source: &inMemorySource{fakeSource},
		Target: new(target.Filesystem),
		Alias:  new(target.HTMLRedirectAlias),
	}
	must(s.CreatePages())
	expected := "foo/bar/file.md\n canonical => foo/bar/file/index.html\n" +
		"alias/test/file1.md\n" +
		" canonical => alias/test/file1/index.html\n" +
		" alias1/ => alias1/index.html\n" +
		" alias-2/ => alias-2/index.html\n"
	checkShowPlanExpected(t, s, expected)
}

func TestFileTargetUgly(t *testing.T) {
	s := &Site{
		Target: &target.Filesystem{UglyUrls: true},
		Source: &inMemorySource{fakeSource},
		Alias:  new(target.HTMLRedirectAlias),
	}
	s.CreatePages()
	expected := "foo/bar/file.md\n canonical => foo/bar/file.html\n" +
		"alias/test/file1.md\n" +
		" canonical => alias/test/file1.html\n" +
		" alias1/ => alias1/index.html\n" +
		" alias-2/ => alias-2/index.html\n"
	checkShowPlanExpected(t, s, expected)
}

func TestFileTargetPublishDir(t *testing.T) {
	s := &Site{
		Target: &target.Filesystem{PublishDir: "../public"},
		Source: &inMemorySource{fakeSource},
		Alias:  &target.HTMLRedirectAlias{PublishDir: "../public"},
	}

	must(s.CreatePages())
	expected := "foo/bar/file.md\n canonical => ../public/foo/bar/file/index.html\n" +
		"alias/test/file1.md\n" +
		" canonical => ../public/alias/test/file1/index.html\n" +
		" alias1/ => ../public/alias1/index.html\n" +
		" alias-2/ => ../public/alias-2/index.html\n"
	checkShowPlanExpected(t, s, expected)
}
