// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"path/filepath"
	"testing"

	"html/template"

	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"github.com/spf13/viper"
)

const slugDoc1 = "---\ntitle: slug doc 1\nslug: slug-doc-1\naliases:\n - sd1/foo/\n - sd2\n - sd3/\n - sd4.html\n---\nslug doc 1 content\n"

const slugDoc2 = `---
title: slug doc 2
slug: slug-doc-2
---
slug doc 2 content
`

const indexTemplate = "{{ range .Data.Pages }}.{{ end }}"

func must(err error) {
	if err != nil {
		panic(err)
	}
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
	{filepath.FromSlash("content/blue/doc1.md"), []byte(slugDoc1)},
	{filepath.FromSlash("content/blue/doc2.md"), []byte(slugDoc2)},
}

// Issue #1105
func TestShouldNotAddTrailingSlashToBaseURL(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	for i, this := range []struct {
		in       string
		expected string
	}{
		{"http://base.com/", "http://base.com/"},
		{"http://base.com/sub/", "http://base.com/sub/"},
		{"http://base.com/sub", "http://base.com/sub"},
		{"http://base.com", "http://base.com"}} {

		viper.Set("BaseURL", this.in)
		s := &Site{}
		s.initializeSiteInfo()

		if s.Info.BaseURL != template.URL(this.expected) {
			t.Errorf("[%d] got %s expected %s", i, s.Info.BaseURL, this.expected)
		}
	}

}

func TestPageCount(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	hugofs.InitMemFs()

	viper.Set("uglyurls", false)
	viper.Set("paginate", 10)
	s := &Site{
		Source: &source.InMemorySource{ByteSource: urlFakeSource},
	}
	s.initializeSiteInfo()
	s.prepTemplates("indexes/blue.html", indexTemplate)

	if err := s.createPages(); err != nil {
		t.Errorf("Unable to create pages: %s", err)
	}
	if err := s.buildSiteMeta(); err != nil {
		t.Errorf("Unable to build site metadata: %s", err)
	}

	if err := s.renderSectionLists(); err != nil {
		t.Errorf("Unable to render section lists: %s", err)
	}

	if err := s.renderAliases(); err != nil {
		t.Errorf("Unable to render site lists: %s", err)
	}

	_, err := hugofs.Destination().Open("blue")
	if err != nil {
		t.Errorf("No indexed rendered.")
	}

	//expected := ".."
	//if string(blueIndex) != expected {
	//t.Errorf("Index template does not match expected: %q, got: %q", expected, string(blueIndex))
	//}

	for _, s := range []string{
		"sd1/foo/index.html",
		"sd2/index.html",
		"sd3/index.html",
		"sd4.html",
	} {
		if _, err := hugofs.Destination().Open(filepath.FromSlash(s)); err != nil {
			t.Errorf("No alias rendered: %s", s)
		}
	}
}
