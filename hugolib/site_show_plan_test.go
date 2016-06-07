// Copyright 2015 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"github.com/spf13/viper"
)

const aliasDoc1 = "---\ntitle: alias doc\naliases:\n  - \"alias1/\"\n  - \"alias-2/\"\n---\naliases\n"

var fakeSource = []source.ByteSource{
	{
		Name:    filepath.FromSlash("foo/bar/file.md"),
		Content: []byte(simplePage),
	},
	{
		Name:    filepath.FromSlash("alias/test/file1.md"),
		Content: []byte(aliasDoc1),
	},
	{
		Name:    filepath.FromSlash("section/somecontent.html"),
		Content: []byte(renderNoFrontmatter),
	},
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
	must(s.createPages())
	expected := "foo/bar/file.md (renderer: markdown)\n canonical => !no target specified!\n\n" +
		"alias/test/file1.md (renderer: markdown)\n canonical => !no target specified!\n\n" +
		"section/somecontent.html (renderer: n/a)\n canonical => !no target specified!\n\n"
	checkShowPlanExpected(t, s, expected)
}

func TestFileTarget(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("DefaultExtension", "html")

	s := &Site{
		Source: &source.InMemorySource{ByteSource: fakeSource},
	}
	s.aliasTarget()
	s.pageTarget()
	must(s.createPages())
	expected := "foo/bar/file.md (renderer: markdown)\n canonical => foo/bar/file/index.html\n\n" +
		"alias/test/file1.md (renderer: markdown)\n" +
		" canonical => alias/test/file1/index.html\n" +
		" alias1/ => alias1/index.html\n" +
		" alias-2/ => alias-2/index.html\n\n" +
		"section/somecontent.html (renderer: n/a)\n canonical => section/somecontent/index.html\n\n"

	checkShowPlanExpected(t, s, expected)
}

func TestPageTargetUgly(t *testing.T) {
	viper.Reset()
	defer viper.Reset()
	viper.Set("DefaultExtension", "html")
	viper.Set("UglyURLs", true)

	s := &Site{
		targets: targetList{page: &target.PagePub{UglyURLs: true}},
		Source:  &source.InMemorySource{ByteSource: fakeSource},
	}
	s.aliasTarget()

	s.createPages()
	expected := "foo/bar/file.md (renderer: markdown)\n canonical => foo/bar/file.html\n\n" +
		"alias/test/file1.md (renderer: markdown)\n" +
		" canonical => alias/test/file1.html\n" +
		" alias1/ => alias1/index.html\n" +
		" alias-2/ => alias-2/index.html\n\n" +
		"section/somecontent.html (renderer: n/a)\n canonical => section/somecontent.html\n\n"
	checkShowPlanExpected(t, s, expected)
}

func TestFileTargetPublishDir(t *testing.T) {
	viper.Reset()
	defer viper.Reset()
	viper.Set("DefaultExtension", "html")

	s := &Site{

		targets: targetList{
			page:  &target.PagePub{PublishDir: "../public"},
			alias: &target.HTMLRedirectAlias{PublishDir: "../public"},
		},
		Source: &source.InMemorySource{ByteSource: fakeSource},
	}

	must(s.createPages())
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
