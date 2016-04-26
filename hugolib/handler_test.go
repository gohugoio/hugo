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
	"path/filepath"
	"testing"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/source"
	"github.com/spf13/hugo/target"
	"github.com/spf13/viper"
)

func TestDefaultHandler(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	hugofs.InitMemFs()
	sources := []source.ByteSource{
		{filepath.FromSlash("sect/doc1.html"), []byte("---\nmarkup: markdown\n---\n# title\nsome *content*")},
		{filepath.FromSlash("sect/doc2.html"), []byte("<!doctype html><html><body>more content</body></html>")},
		{filepath.FromSlash("sect/doc3.md"), []byte("# doc3\n*some* content")},
		{filepath.FromSlash("sect/doc4.md"), []byte("---\ntitle: doc4\n---\n# doc4\n*some content*")},
		{filepath.FromSlash("sect/doc3/img1.png"), []byte("‰PNG  ��� IHDR����������:~›U��� IDATWcø��ZMoñ����IEND®B`‚")},
		{filepath.FromSlash("sect/img2.gif"), []byte("GIF89a��€��ÿÿÿ���,�������D�;")},
		{filepath.FromSlash("sect/img2.spf"), []byte("****FAKE-FILETYPE****")},
		{filepath.FromSlash("doc7.html"), []byte("<html><body>doc7 content</body></html>")},
		{filepath.FromSlash("sect/doc8.html"), []byte("---\nmarkup: md\n---\n# title\nsome *content*")},
	}

	viper.Set("DefaultExtension", "html")
	viper.Set("verbose", true)

	s := &Site{
		Source:  &source.InMemorySource{ByteSource: sources},
		targets: targetList{page: &target.PagePub{UglyURLs: true}},
	}

	s.initializeSiteInfo()

	s.prepTemplates(
		"_default/single.html", "{{.Content}}",
		"head", "<head><script src=\"script.js\"></script></head>",
		"head_abs", "<head><script src=\"/script.js\"></script></head>")

	// From site_test.go
	createAndRenderPages(t, s)

	tests := []struct {
		doc      string
		expected string
	}{
		{filepath.FromSlash("sect/doc1.html"), "\n\n<h1 id=\"title\">title</h1>\n\n<p>some <em>content</em></p>\n"},
		{filepath.FromSlash("sect/doc2.html"), "<!doctype html><html><body>more content</body></html>"},
		{filepath.FromSlash("sect/doc3.html"), "\n\n<h1 id=\"doc3\">doc3</h1>\n\n<p><em>some</em> content</p>\n"},
		{filepath.FromSlash("sect/doc3/img1.png"), string([]byte("‰PNG  ��� IHDR����������:~›U��� IDATWcø��ZMoñ����IEND®B`‚"))},
		{filepath.FromSlash("sect/img2.gif"), string([]byte("GIF89a��€��ÿÿÿ���,�������D�;"))},
		{filepath.FromSlash("sect/img2.spf"), string([]byte("****FAKE-FILETYPE****"))},
		{filepath.FromSlash("doc7.html"), "<html><body>doc7 content</body></html>"},
		{filepath.FromSlash("sect/doc8.html"), "\n\n<h1 id=\"title\">title</h1>\n\n<p>some <em>content</em></p>\n"},
	}

	for _, test := range tests {
		file, err := hugofs.Destination().Open(test.doc)
		if err != nil {
			t.Fatalf("Did not find %s in target.", test.doc)
		}

		content := helpers.ReaderToString(file)

		if content != test.expected {
			t.Errorf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, content)
		}
	}

}
