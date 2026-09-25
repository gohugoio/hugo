// Copyright 2025 The Hugo Authors. All rights reserved.
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

package page

import (
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/resources/kinds"

	qt "github.com/frankban/quicktest"
)

func TestPagePathsBuilder(t *testing.T) {
	c := qt.New(t)

	d := TargetPathDescriptor{}
	b := getPagePathBuilder(d)
	defer putPagePathBuilder(b)
	b.Add("foo", "bar")

	c.Assert(b.Path(0), qt.Equals, "/foo/bar")
}

// See issue 10393.
func TestTargetPathsExplicitURLWithOutputFormats(t *testing.T) {
	t.Parallel()

	pp := &paths.PathParser{IsContentExt: func(ext string) bool { return ext == "md" }}
	p := pp.Parse(files.ComponentFolderContent, "/posts/_index.md")
	json := output.JSONFormat
	json.Name = "custom"
	json.Path = "api"
	json.BaseName = "data"

	for _, test := range []struct {
		name             string
		format           output.Format
		primaryMediaType string
		url              string
		ugly             bool
		want             string
	}{
		{"html", output.HTMLFormat, output.HTMLFormat.MediaType.Type, "/posts/test.html", false, "/posts/test.html"},
		{"rss", output.RSSFormat, output.HTMLFormat.MediaType.Type, "/posts/test.html", false, "/posts/test.xml"},
		{"rss ugly", output.RSSFormat, output.HTMLFormat.MediaType.Type, "/posts/test.html", true, "/posts/test.xml"},
		{"html arbitrary suffix", output.HTMLFormat, output.HTMLFormat.MediaType.Type, "/posts/test.php", false, "/posts/test.php"},
		{"rss arbitrary suffix", output.RSSFormat, output.HTMLFormat.MediaType.Type, "/posts/test.php", false, "/posts/test.xml"},
		{"custom alternate", json, output.HTMLFormat.MediaType.Type, "/posts/test.php", false, "/api/posts/test.json"},
		{"same media type", output.AMPFormat, output.HTMLFormat.MediaType.Type, "/posts/test.php", false, "/amp/posts/test.php"},
		{"rss primary", output.RSSFormat, output.RSSFormat.MediaType.Type, "/posts/test.feed", false, "/posts/test.feed"},
		{"custom primary", json, json.MediaType.Type, "/posts/test.data", false, "/api/posts/test.data"},
		{"no primary media type", output.RSSFormat, "", "/posts/test.feed", false, "/posts/test.feed"},
	} {
		t.Run(test.name, func(t *testing.T) {
			c := qt.New(t)
			tp := CreateTargetPaths(TargetPathDescriptor{
				Type:             test.format,
				PrimaryMediaType: test.primaryMediaType,
				Kind:             kinds.KindSection,
				Path:             p,
				Section:          p,
				URL:              test.url,
				UglyURLs:         test.ugly,
			})
			c.Assert(filepath.ToSlash(tp.TargetFilename), qt.Equals, test.want)
			c.Assert(tp.Link, qt.Equals, test.want)
		})
	}
}

func BenchmarkPagePathsBuilderPath(b *testing.B) {
	d := TargetPathDescriptor{}
	pb := getPagePathBuilder(d)
	defer putPagePathBuilder(pb)
	pb.Add("foo", "bar")

	for b.Loop() {
		_ = pb.Path(0)
	}
}

func BenchmarkPagePathsBuilderPathDir(b *testing.B) {
	d := TargetPathDescriptor{}
	pb := getPagePathBuilder(d)
	defer putPagePathBuilder(pb)
	pb.Add("foo", "bar")
	pb.prefixPath = "foo/"

	for b.Loop() {
		_ = pb.PathDir()
	}
}
