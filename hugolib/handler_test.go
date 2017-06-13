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

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
)

func TestDefaultHandler(t *testing.T) {
	t.Parallel()

	var (
		cfg, fs = newTestCfg()
	)

	cfg.Set("verbose", true)
	cfg.Set("uglyURLs", true)

	writeSource(t, fs, filepath.FromSlash("content/sect/doc1.html"), "---\nmarkup: markdown\n---\n# title\nsome *content*")
	writeSource(t, fs, filepath.FromSlash("content/sect/doc2.html"), "<!doctype html><html><body>more content</body></html>")
	writeSource(t, fs, filepath.FromSlash("content/sect/doc3.md"), "# doc3\n*some* content")
	writeSource(t, fs, filepath.FromSlash("content/sect/doc4.md"), "---\ntitle: doc4\n---\n# doc4\n*some content*")
	writeSource(t, fs, filepath.FromSlash("content/sect/doc3/img1.png"), "‰PNG  ��� IHDR����������:~›U��� IDATWcø��ZMoñ����IEND®B`‚")
	writeSource(t, fs, filepath.FromSlash("content/sect/img2.gif"), "GIF89a��€��ÿÿÿ���,�������D�;")
	writeSource(t, fs, filepath.FromSlash("content/sect/img2.spf"), "****FAKE-FILETYPE****")
	writeSource(t, fs, filepath.FromSlash("content/doc7.html"), "<html><body>doc7 content</body></html>")
	writeSource(t, fs, filepath.FromSlash("content/sect/doc8.html"), "---\nmarkup: md\n---\n# title\nsome *content*")

	writeSource(t, fs, filepath.FromSlash("layouts/_default/single.html"), "{{.Content}}")
	writeSource(t, fs, filepath.FromSlash("head"), "<head><script src=\"script.js\"></script></head>")
	writeSource(t, fs, filepath.FromSlash("head_abs"), "<head><script src=\"/script.js\"></script></head")

	buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	tests := []struct {
		doc      string
		expected string
	}{
		{filepath.FromSlash("public/sect/doc1.html"), "\n\n<h1 id=\"title\">title</h1>\n\n<p>some <em>content</em></p>\n"},
		{filepath.FromSlash("public/sect/doc2.html"), "<!doctype html><html><body>more content</body></html>"},
		{filepath.FromSlash("public/sect/doc3.html"), "\n\n<h1 id=\"doc3\">doc3</h1>\n\n<p><em>some</em> content</p>\n"},
		{filepath.FromSlash("public/sect/doc3/img1.png"), string([]byte("‰PNG  ��� IHDR����������:~›U��� IDATWcø��ZMoñ����IEND®B`‚"))},
		{filepath.FromSlash("public/sect/img2.gif"), string([]byte("GIF89a��€��ÿÿÿ���,�������D�;"))},
		{filepath.FromSlash("public/sect/img2.spf"), string([]byte("****FAKE-FILETYPE****"))},
		{filepath.FromSlash("public/doc7.html"), "<html><body>doc7 content</body></html>"},
		{filepath.FromSlash("public/sect/doc8.html"), "\n\n<h1 id=\"title\">title</h1>\n\n<p>some <em>content</em></p>\n"},
	}

	for _, test := range tests {
		file, err := fs.Destination.Open(test.doc)
		if err != nil {
			t.Fatalf("Did not find %s in target.", test.doc)
		}

		content := helpers.ReaderToString(file)

		if content != test.expected {
			t.Errorf("%s content expected:\n%q\ngot:\n%q", test.doc, test.expected, content)
		}
	}

}
