package hugolib

import (
	"bytes"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"testing"
)

var fakeSource = []struct {
	name    string
	content []byte
}{
	{"foo/bar/file.md", []byte(SIMPLE_PAGE)},
}

type inMemorySource struct {
}

func (i inMemorySource) Files() (files []*source.File) {
	files = make([]*source.File, len(fakeSource))
	for i, fake := range fakeSource {
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
	s := &Site{Source: new(inMemorySource)}
	must(s.CreatePages())
	expected := "foo/bar/file.md\n canonical => !no target specified!\n"
	checkShowPlanExpected(t, s, expected)
}

func TestFileTarget(t *testing.T) {
	s := &Site{
		Source: new(inMemorySource),
		Target: new(target.Filesystem),
	}
	must(s.CreatePages())
	checkShowPlanExpected(t, s, "foo/bar/file.md\n canonical => foo/bar/file/index.html\n")
}

func TestFileTargetUgly(t *testing.T) {
	s := &Site{
		Target: &target.Filesystem{UglyUrls: true},
		Source: new(inMemorySource),
	}
	s.CreatePages()
	expected := "foo/bar/file.md\n canonical => foo/bar/file.html\n"
	checkShowPlanExpected(t, s, expected)
}

func TestFileTargetPublishDir(t *testing.T) {
	s := &Site{
		Target: &target.Filesystem{PublishDir: "../public"},
		Source: new(inMemorySource),
	}

	must(s.CreatePages())
	expected := "foo/bar/file.md\n canonical => ../public/foo/bar/file/index.html\n"
	checkShowPlanExpected(t, s, expected)
}
