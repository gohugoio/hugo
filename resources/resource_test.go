// Copyright 2019 The Hugo Authors. All rights reserved.
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

package resources

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/media"

	qt "github.com/frankban/quicktest"
)

func TestGenericResource(t *testing.T) {
	c := qt.New(t)
	spec := newTestResourceSpec(specDescriptor{c: c})

	r := spec.newGenericResource(nil, nil, nil, "/a/foo.css", "foo.css", media.CSSType)

	c.Assert(r.Permalink(), qt.Equals, "https://example.com/foo.css")
	c.Assert(r.RelPermalink(), qt.Equals, "/foo.css")
	c.Assert(r.ResourceType(), qt.Equals, "css")

}

func TestGenericResourceWithLinkFacory(t *testing.T) {
	c := qt.New(t)
	spec := newTestResourceSpec(specDescriptor{c: c})

	factory := newTargetPaths("/foo")

	r := spec.newGenericResource(nil, factory, nil, "/a/foo.css", "foo.css", media.CSSType)

	c.Assert(r.Permalink(), qt.Equals, "https://example.com/foo/foo.css")
	c.Assert(r.RelPermalink(), qt.Equals, "/foo/foo.css")
	c.Assert(r.Key(), qt.Equals, "/foo/foo.css")
	c.Assert(r.ResourceType(), qt.Equals, "css")
}

func TestNewResourceFromFilename(t *testing.T) {
	c := qt.New(t)
	spec := newTestResourceSpec(specDescriptor{c: c})

	writeSource(t, spec.Fs, "content/a/b/logo.png", "image")
	writeSource(t, spec.Fs, "content/a/b/data.json", "json")

	bfs := afero.NewBasePathFs(spec.Fs.Source, "content")

	r, err := spec.New(ResourceSourceDescriptor{Fs: bfs, SourceFilename: "a/b/logo.png"})

	c.Assert(err, qt.IsNil)
	c.Assert(r, qt.Not(qt.IsNil))
	c.Assert(r.ResourceType(), qt.Equals, "image")
	c.Assert(r.RelPermalink(), qt.Equals, "/a/b/logo.png")
	c.Assert(r.Permalink(), qt.Equals, "https://example.com/a/b/logo.png")

	r, err = spec.New(ResourceSourceDescriptor{Fs: bfs, SourceFilename: "a/b/data.json"})

	c.Assert(err, qt.IsNil)
	c.Assert(r, qt.Not(qt.IsNil))
	c.Assert(r.ResourceType(), qt.Equals, "json")

}

func TestNewResourceFromFilenameSubPathInBaseURL(t *testing.T) {
	c := qt.New(t)
	spec := newTestResourceSpec(specDescriptor{c: c, baseURL: "https://example.com/docs"})

	writeSource(t, spec.Fs, "content/a/b/logo.png", "image")
	bfs := afero.NewBasePathFs(spec.Fs.Source, "content")

	fmt.Println()
	r, err := spec.New(ResourceSourceDescriptor{Fs: bfs, SourceFilename: filepath.FromSlash("a/b/logo.png")})

	c.Assert(err, qt.IsNil)
	c.Assert(r, qt.Not(qt.IsNil))
	c.Assert(r.ResourceType(), qt.Equals, "image")
	c.Assert(r.RelPermalink(), qt.Equals, "/docs/a/b/logo.png")
	c.Assert(r.Permalink(), qt.Equals, "https://example.com/docs/a/b/logo.png")

}

var pngType, _ = media.FromStringAndExt("image/png", "png")

func TestResourcesByType(t *testing.T) {
	c := qt.New(t)
	spec := newTestResourceSpec(specDescriptor{c: c})
	resources := resource.Resources{
		spec.newGenericResource(nil, nil, nil, "/a/foo1.css", "foo1.css", media.CSSType),
		spec.newGenericResource(nil, nil, nil, "/a/logo.png", "logo.css", pngType),
		spec.newGenericResource(nil, nil, nil, "/a/foo2.css", "foo2.css", media.CSSType),
		spec.newGenericResource(nil, nil, nil, "/a/foo3.css", "foo3.css", media.CSSType)}

	c.Assert(len(resources.ByType("css")), qt.Equals, 3)
	c.Assert(len(resources.ByType("image")), qt.Equals, 1)

}

func TestResourcesGetByPrefix(t *testing.T) {
	c := qt.New(t)
	spec := newTestResourceSpec(specDescriptor{c: c})
	resources := resource.Resources{
		spec.newGenericResource(nil, nil, nil, "/a/foo1.css", "foo1.css", media.CSSType),
		spec.newGenericResource(nil, nil, nil, "/a/logo1.png", "logo1.png", pngType),
		spec.newGenericResource(nil, nil, nil, "/b/Logo2.png", "Logo2.png", pngType),
		spec.newGenericResource(nil, nil, nil, "/b/foo2.css", "foo2.css", media.CSSType),
		spec.newGenericResource(nil, nil, nil, "/b/foo3.css", "foo3.css", media.CSSType)}

	c.Assert(resources.GetMatch("asdf*"), qt.IsNil)
	c.Assert(resources.GetMatch("logo*").RelPermalink(), qt.Equals, "/logo1.png")
	c.Assert(resources.GetMatch("loGo*").RelPermalink(), qt.Equals, "/logo1.png")
	c.Assert(resources.GetMatch("logo2*").RelPermalink(), qt.Equals, "/Logo2.png")
	c.Assert(resources.GetMatch("foo2*").RelPermalink(), qt.Equals, "/foo2.css")
	c.Assert(resources.GetMatch("foo1*").RelPermalink(), qt.Equals, "/foo1.css")
	c.Assert(resources.GetMatch("foo1*").RelPermalink(), qt.Equals, "/foo1.css")
	c.Assert(resources.GetMatch("asdfasdf*"), qt.IsNil)

	c.Assert(len(resources.Match("logo*")), qt.Equals, 2)
	c.Assert(len(resources.Match("logo2*")), qt.Equals, 1)

	logo := resources.GetMatch("logo*")
	c.Assert(logo.Params(), qt.Not(qt.IsNil))
	c.Assert(logo.Name(), qt.Equals, "logo1.png")
	c.Assert(logo.Title(), qt.Equals, "logo1.png")

}

func TestResourcesGetMatch(t *testing.T) {
	c := qt.New(t)
	spec := newTestResourceSpec(specDescriptor{c: c})
	resources := resource.Resources{
		spec.newGenericResource(nil, nil, nil, "/a/foo1.css", "foo1.css", media.CSSType),
		spec.newGenericResource(nil, nil, nil, "/a/logo1.png", "logo1.png", pngType),
		spec.newGenericResource(nil, nil, nil, "/b/Logo2.png", "Logo2.png", pngType),
		spec.newGenericResource(nil, nil, nil, "/b/foo2.css", "foo2.css", media.CSSType),
		spec.newGenericResource(nil, nil, nil, "/b/foo3.css", "foo3.css", media.CSSType),
		spec.newGenericResource(nil, nil, nil, "/b/c/foo4.css", "c/foo4.css", media.CSSType),
		spec.newGenericResource(nil, nil, nil, "/b/c/foo5.css", "c/foo5.css", media.CSSType),
		spec.newGenericResource(nil, nil, nil, "/b/c/d/foo6.css", "c/d/foo6.css", media.CSSType),
	}

	c.Assert(resources.GetMatch("logo*").RelPermalink(), qt.Equals, "/logo1.png")
	c.Assert(resources.GetMatch("loGo*").RelPermalink(), qt.Equals, "/logo1.png")
	c.Assert(resources.GetMatch("logo2*").RelPermalink(), qt.Equals, "/Logo2.png")
	c.Assert(resources.GetMatch("foo2*").RelPermalink(), qt.Equals, "/foo2.css")
	c.Assert(resources.GetMatch("foo1*").RelPermalink(), qt.Equals, "/foo1.css")
	c.Assert(resources.GetMatch("foo1*").RelPermalink(), qt.Equals, "/foo1.css")
	c.Assert(resources.GetMatch("*/foo*").RelPermalink(), qt.Equals, "/c/foo4.css")

	c.Assert(resources.GetMatch("asdfasdf"), qt.IsNil)

	c.Assert(len(resources.Match("Logo*")), qt.Equals, 2)
	c.Assert(len(resources.Match("logo2*")), qt.Equals, 1)
	c.Assert(len(resources.Match("c/*")), qt.Equals, 2)

	c.Assert(len(resources.Match("**.css")), qt.Equals, 6)
	c.Assert(len(resources.Match("**/*.css")), qt.Equals, 3)
	c.Assert(len(resources.Match("c/**/*.css")), qt.Equals, 1)

	// Matches only CSS files in c/
	c.Assert(len(resources.Match("c/**.css")), qt.Equals, 3)

	// Matches all CSS files below c/ (including in c/d/)
	c.Assert(len(resources.Match("c/**.css")), qt.Equals, 3)

	// Patterns beginning with a slash will not match anything.
	// We could maybe consider trimming that slash, but let's be explicit about this.
	// (it is possible for users to do a rename)
	// This is analogous to standing in a directory and doing "ls *.*".
	c.Assert(len(resources.Match("/c/**.css")), qt.Equals, 0)

}

func BenchmarkResourcesMatch(b *testing.B) {
	resources := benchResources(b)
	prefixes := []string{"abc*", "jkl*", "nomatch*", "sub/*"}
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resources.Match(prefixes[rnd.Intn(len(prefixes))])
		}
	})
}

// This adds a benchmark for the a100 test case as described by Russ Cox here:
// https://research.swtch.com/glob (really interesting article)
// I don't expect Hugo users to "stumble upon" this problem, so this is more to satisfy
// my own curiosity.
func BenchmarkResourcesMatchA100(b *testing.B) {
	c := qt.New(b)
	spec := newTestResourceSpec(specDescriptor{c: c})
	a100 := strings.Repeat("a", 100)
	pattern := "a*a*a*a*a*a*a*a*b"

	resources := resource.Resources{spec.newGenericResource(nil, nil, nil, "/a/"+a100, a100, media.CSSType)}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resources.Match(pattern)
	}

}

func benchResources(b *testing.B) resource.Resources {
	c := qt.New(b)
	spec := newTestResourceSpec(specDescriptor{c: c})
	var resources resource.Resources

	for i := 0; i < 30; i++ {
		name := fmt.Sprintf("abcde%d_%d.css", i%5, i)
		resources = append(resources, spec.newGenericResource(nil, nil, nil, "/a/"+name, name, media.CSSType))
	}

	for i := 0; i < 30; i++ {
		name := fmt.Sprintf("efghi%d_%d.css", i%5, i)
		resources = append(resources, spec.newGenericResource(nil, nil, nil, "/a/"+name, name, media.CSSType))
	}

	for i := 0; i < 30; i++ {
		name := fmt.Sprintf("jklmn%d_%d.css", i%5, i)
		resources = append(resources, spec.newGenericResource(nil, nil, nil, "/b/sub/"+name, "sub/"+name, media.CSSType))
	}

	return resources

}

func BenchmarkAssignMetadata(b *testing.B) {
	c := qt.New(b)
	spec := newTestResourceSpec(specDescriptor{c: c})

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		var resources resource.Resources
		var meta = []map[string]interface{}{
			{
				"title": "Foo #:counter",
				"name":  "Foo Name #:counter",
				"src":   "foo1*",
			},
			{
				"title": "Rest #:counter",
				"name":  "Rest Name #:counter",
				"src":   "*",
			},
		}
		for i := 0; i < 20; i++ {
			name := fmt.Sprintf("foo%d_%d.css", i%5, i)
			resources = append(resources, spec.newGenericResource(nil, nil, nil, "/a/"+name, name, media.CSSType))
		}
		b.StartTimer()

		if err := AssignMetadata(meta, resources...); err != nil {
			b.Fatal(err)
		}

	}
}
