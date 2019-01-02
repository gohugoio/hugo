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
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/media"

	"github.com/stretchr/testify/require"
)

func TestGenericResource(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)

	r := spec.newGenericResource(nil, nil, nil, "/a/foo.css", "foo.css", media.CSSType)

	assert.Equal("https://example.com/foo.css", r.Permalink())
	assert.Equal("/foo.css", r.RelPermalink())
	assert.Equal("css", r.ResourceType())

}

func TestGenericResourceWithLinkFacory(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)

	factory := func(s string) string {
		return path.Join("/foo", s)
	}
	r := spec.newGenericResource(nil, factory, nil, "/a/foo.css", "foo.css", media.CSSType)

	assert.Equal("https://example.com/foo/foo.css", r.Permalink())
	assert.Equal("/foo/foo.css", r.RelPermalink())
	assert.Equal("foo.css", r.Key())
	assert.Equal("css", r.ResourceType())
}

func TestNewResourceFromFilename(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)

	writeSource(t, spec.Fs, "content/a/b/logo.png", "image")
	writeSource(t, spec.Fs, "content/a/b/data.json", "json")

	r, err := spec.New(ResourceSourceDescriptor{SourceFilename: "a/b/logo.png"})

	assert.NoError(err)
	assert.NotNil(r)
	assert.Equal("image", r.ResourceType())
	assert.Equal("/a/b/logo.png", r.RelPermalink())
	assert.Equal("https://example.com/a/b/logo.png", r.Permalink())

	r, err = spec.New(ResourceSourceDescriptor{SourceFilename: "a/b/data.json"})

	assert.NoError(err)
	assert.NotNil(r)
	assert.Equal("json", r.ResourceType())

	cloned := r.(resource.Cloner).WithNewBase("aceof")
	assert.Equal(r.ResourceType(), cloned.ResourceType())
	assert.Equal("/aceof/a/b/data.json", cloned.RelPermalink())
}

func TestNewResourceFromFilenameSubPathInBaseURL(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpecForBaseURL(assert, "https://example.com/docs")

	writeSource(t, spec.Fs, "content/a/b/logo.png", "image")

	r, err := spec.New(ResourceSourceDescriptor{SourceFilename: filepath.FromSlash("a/b/logo.png")})

	assert.NoError(err)
	assert.NotNil(r)
	assert.Equal("image", r.ResourceType())
	assert.Equal("/docs/a/b/logo.png", r.RelPermalink())
	assert.Equal("https://example.com/docs/a/b/logo.png", r.Permalink())
	img := r.(*Image)
	assert.Equal(filepath.FromSlash("/a/b/logo.png"), img.targetFilenames()[0])

}

var pngType, _ = media.FromStringAndExt("image/png", "png")

func TestResourcesByType(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)
	resources := resource.Resources{
		spec.newGenericResource(nil, nil, nil, "/a/foo1.css", "foo1.css", media.CSSType),
		spec.newGenericResource(nil, nil, nil, "/a/logo.png", "logo.css", pngType),
		spec.newGenericResource(nil, nil, nil, "/a/foo2.css", "foo2.css", media.CSSType),
		spec.newGenericResource(nil, nil, nil, "/a/foo3.css", "foo3.css", media.CSSType)}

	assert.Len(resources.ByType("css"), 3)
	assert.Len(resources.ByType("image"), 1)

}

func TestResourcesGetByPrefix(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)
	resources := resource.Resources{
		spec.newGenericResource(nil, nil, nil, "/a/foo1.css", "foo1.css", media.CSSType),
		spec.newGenericResource(nil, nil, nil, "/a/logo1.png", "logo1.png", pngType),
		spec.newGenericResource(nil, nil, nil, "/b/Logo2.png", "Logo2.png", pngType),
		spec.newGenericResource(nil, nil, nil, "/b/foo2.css", "foo2.css", media.CSSType),
		spec.newGenericResource(nil, nil, nil, "/b/foo3.css", "foo3.css", media.CSSType)}

	assert.Nil(resources.GetMatch("asdf*"))
	assert.Equal("/logo1.png", resources.GetMatch("logo*").RelPermalink())
	assert.Equal("/logo1.png", resources.GetMatch("loGo*").RelPermalink())
	assert.Equal("/Logo2.png", resources.GetMatch("logo2*").RelPermalink())
	assert.Equal("/foo2.css", resources.GetMatch("foo2*").RelPermalink())
	assert.Equal("/foo1.css", resources.GetMatch("foo1*").RelPermalink())
	assert.Equal("/foo1.css", resources.GetMatch("foo1*").RelPermalink())
	assert.Nil(resources.GetMatch("asdfasdf*"))

	assert.Equal(2, len(resources.Match("logo*")))
	assert.Equal(1, len(resources.Match("logo2*")))

	logo := resources.GetMatch("logo*")
	assert.NotNil(logo.Params())
	assert.Equal("logo1.png", logo.Name())
	assert.Equal("logo1.png", logo.Title())

}

func TestResourcesGetMatch(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)
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

	assert.Equal("/logo1.png", resources.GetMatch("logo*").RelPermalink())
	assert.Equal("/logo1.png", resources.GetMatch("loGo*").RelPermalink())
	assert.Equal("/Logo2.png", resources.GetMatch("logo2*").RelPermalink())
	assert.Equal("/foo2.css", resources.GetMatch("foo2*").RelPermalink())
	assert.Equal("/foo1.css", resources.GetMatch("foo1*").RelPermalink())
	assert.Equal("/foo1.css", resources.GetMatch("foo1*").RelPermalink())
	assert.Equal("/c/foo4.css", resources.GetMatch("*/foo*").RelPermalink())

	assert.Nil(resources.GetMatch("asdfasdf"))

	assert.Equal(2, len(resources.Match("Logo*")))
	assert.Equal(1, len(resources.Match("logo2*")))
	assert.Equal(2, len(resources.Match("c/*")))

	assert.Equal(6, len(resources.Match("**.css")))
	assert.Equal(3, len(resources.Match("**/*.css")))
	assert.Equal(1, len(resources.Match("c/**/*.css")))

	// Matches only CSS files in c/
	assert.Equal(3, len(resources.Match("c/**.css")))

	// Matches all CSS files below c/ (including in c/d/)
	assert.Equal(3, len(resources.Match("c/**.css")))

	// Patterns beginning with a slash will not match anything.
	// We could maybe consider trimming that slash, but let's be explicit about this.
	// (it is possible for users to do a rename)
	// This is analogous to standing in a directory and doing "ls *.*".
	assert.Equal(0, len(resources.Match("/c/**.css")))

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
	assert := require.New(b)
	spec := newTestResourceSpec(assert)
	a100 := strings.Repeat("a", 100)
	pattern := "a*a*a*a*a*a*a*a*b"

	resources := resource.Resources{spec.newGenericResource(nil, nil, nil, "/a/"+a100, a100, media.CSSType)}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resources.Match(pattern)
	}

}

func benchResources(b *testing.B) resource.Resources {
	assert := require.New(b)
	spec := newTestResourceSpec(assert)
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
	assert := require.New(b)
	spec := newTestResourceSpec(assert)

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
