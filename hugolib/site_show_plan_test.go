package hugolib

import (
	"bytes"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"testing"
)

const ALIAS_DOC_1 = "---\ntitle: alias doc\naliases:\n  - \"alias1/\"\n  - \"alias-2/\"\n---\naliases\n"

var fakeSource = []source.ByteSource{
	{
		Name:    "foo/bar/file.md",
		Content: []byte(SIMPLE_PAGE),
	},
	{
		Name:    "alias/test/file1.md",
		Content: []byte(ALIAS_DOC_1),
	},
	{
		Name:    "section/somecontent.html",
		Content: []byte(RENDER_NO_FRONT_MATTER),
	},
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
		Source: &source.InMemorySource{ByteSource: fakeSource},
	}
	must(s.CreatePages())
	expected := "foo/bar/file.md (renderer: markdown)\n canonical => !no target specified!\n\n" +
		"alias/test/file1.md (renderer: markdown)\n canonical => !no target specified!\n\n" +
		"section/somecontent.html (renderer: n/a)\n canonical => !no target specified!\n\n"
	checkShowPlanExpected(t, s, expected)
}

func TestFileTarget(t *testing.T) {
	s := &Site{
		Source: &source.InMemorySource{ByteSource: fakeSource},
		Target: new(target.Filesystem),
		Alias:  new(target.HTMLRedirectAlias),
	}
	must(s.CreatePages())
	expected := "foo/bar/file.md (renderer: markdown)\n canonical => foo/bar/file/index.html\n\n" +
		"alias/test/file1.md (renderer: markdown)\n" +
		" canonical => alias/test/file1/index.html\n" +
		" alias1/ => alias1/index.html\n" +
		" alias-2/ => alias-2/index.html\n\n" +
		"section/somecontent.html (renderer: n/a)\n canonical => section/somecontent/index.html\n\n"

	checkShowPlanExpected(t, s, expected)
}

func TestFileTargetUgly(t *testing.T) {
	s := &Site{
		Target: &target.Filesystem{UglyUrls: true},
		Source: &source.InMemorySource{ByteSource: fakeSource},
		Alias:  new(target.HTMLRedirectAlias),
	}
	s.CreatePages()
	expected := "foo/bar/file.md (renderer: markdown)\n canonical => foo/bar/file.html\n\n" +
		"alias/test/file1.md (renderer: markdown)\n" +
		" canonical => alias/test/file1.html\n" +
		" alias1/ => alias1/index.html\n" +
		" alias-2/ => alias-2/index.html\n\n" +
		"section/somecontent.html (renderer: n/a)\n canonical => section/somecontent.html\n\n"
	checkShowPlanExpected(t, s, expected)
}

func TestFileTargetPublishDir(t *testing.T) {
	s := &Site{
		Target: &target.Filesystem{PublishDir: "../public"},
		Source: &source.InMemorySource{ByteSource: fakeSource},
		Alias:  &target.HTMLRedirectAlias{PublishDir: "../public"},
	}

	must(s.CreatePages())
	expected := "foo/bar/file.md (renderer: markdown)\n canonical => ../public/foo/bar/file/index.html\n\n" +
		"alias/test/file1.md (renderer: markdown)\n" +
		" canonical => ../public/alias/test/file1/index.html\n" +
		" alias1/ => ../public/alias1/index.html\n" +
		" alias-2/ => ../public/alias-2/index.html\n\n" +
		"section/somecontent.html (renderer: n/a)\n canonical => ../public/section/somecontent/index.html\n\n"
	checkShowPlanExpected(t, s, expected)
}
