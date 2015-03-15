package hugolib

import (
	"bytes"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"path/filepath"
	"strings"
	"testing"
)

const ALIAS_DOC_1 = "---\ntitle: alias doc\naliases:\n  - \"alias1/\"\n  - \"alias-2/\"\n---\naliases\n"

var fakeSource = []source.ByteSource{
	{
		Name:    filepath.FromSlash("foo/bar/file.md"),
		Content: []byte(SIMPLE_PAGE),
	},
	{
		Name:    filepath.FromSlash("alias/test/file1.md"),
		Content: []byte(ALIAS_DOC_1),
	},
	{
		Name:    filepath.FromSlash("section/somecontent.html"),
		Content: []byte(RENDER_NO_FRONT_MATTER),
	},
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func checkShowPlanExpected(t *testing.T, s *Site, expected string) {

	out := new(bytes.Buffer)
	if err := s.ShowPlan(out); err != nil {
		t.Fatalf("ShowPlan unexpectedly returned an error: %s", err)
	}
	got := out.String()

	expected = filepath.FromSlash(expected)
	// hackety hack: alias is an Url
	expected = strings.Replace(expected, (helpers.FilePathSeparator + " =>"), "/ =>", -1)
	expected = strings.Replace(expected, "n"+(helpers.FilePathSeparator+"a"), "n/a", -1)
	gotList := strings.Split(got, "\n")
	expectedList := strings.Split(expected, "\n")

	diff := DiffStringSlices(gotList, expectedList)

	if len(diff) > 0 {
		t.Errorf("Got diff in show plan: %s", diff)
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
	}
	s.AliasTarget()
	s.PageTarget()
	must(s.CreatePages())
	expected := "foo/bar/file.md (renderer: markdown)\n canonical => foo/bar/file/index.html\n\n" +
		"alias/test/file1.md (renderer: markdown)\n" +
		" canonical => alias/test/file1/index.html\n" +
		" alias1/ => alias1/index.html\n" +
		" alias-2/ => alias-2/index.html\n\n" +
		"section/somecontent.html (renderer: n/a)\n canonical => section/somecontent/index.html\n\n"

	checkShowPlanExpected(t, s, expected)
}

func TestPageTargetUgly(t *testing.T) {
	s := &Site{
		Targets: targetList{Page: &target.PagePub{UglyUrls: true}},
		Source:  &source.InMemorySource{ByteSource: fakeSource},
	}
	s.AliasTarget()

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

		Targets: targetList{
			Page:  &target.PagePub{PublishDir: "../public"},
			Alias: &target.HTMLRedirectAlias{PublishDir: "../public"},
		},
		Source: &source.InMemorySource{ByteSource: fakeSource},
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

// DiffStringSlices returns the difference between two string slices.
// See:
// http://stackoverflow.com/questions/19374219/how-to-find-the-difference-between-two-slices-of-strings-in-golang
func DiffStringSlices(slice1 []string, slice2 []string) []string {
	diffStr := []string{}
	m := map[string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 1
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}

	return diffStr
}
