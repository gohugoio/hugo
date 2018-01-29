// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package resource

import (
	"fmt"
	"math/rand"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenericResource(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)

	r := spec.newGenericResource(nil, nil, "/public", "/a/foo.css", "foo.css", "css")

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
	r := spec.newGenericResource(factory, nil, "/public", "/a/foo.css", "foo.css", "css")

	assert.Equal("https://example.com/foo/foo.css", r.Permalink())
	assert.Equal("/foo/foo.css", r.RelPermalink())
	assert.Equal("css", r.ResourceType())
}

func TestNewResourceFromFilename(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)

	writeSource(t, spec.Fs, "/project/a/b/logo.png", "image")
	writeSource(t, spec.Fs, "/root/a/b/data.json", "json")

	r, err := spec.NewResourceFromFilename(nil, "/public",
		filepath.FromSlash("/project/a/b/logo.png"), filepath.FromSlash("a/b/logo.png"))

	assert.NoError(err)
	assert.NotNil(r)
	assert.Equal("image", r.ResourceType())
	assert.Equal("/a/b/logo.png", r.RelPermalink())
	assert.Equal("https://example.com/a/b/logo.png", r.Permalink())

	r, err = spec.NewResourceFromFilename(nil, "/public", "/root/a/b/data.json", "a/b/data.json")

	assert.NoError(err)
	assert.NotNil(r)
	assert.Equal("json", r.ResourceType())

	cloned := r.(Cloner).WithNewBase("aceof")
	assert.Equal(r.ResourceType(), cloned.ResourceType())
	assert.Equal("/aceof/a/b/data.json", cloned.RelPermalink())
}

func TestNewResourceFromFilenameSubPathInBaseURL(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpecForBaseURL(assert, "https://example.com/docs")

	writeSource(t, spec.Fs, "/project/a/b/logo.png", "image")

	r, err := spec.NewResourceFromFilename(nil, "/public",
		filepath.FromSlash("/project/a/b/logo.png"), filepath.FromSlash("a/b/logo.png"))

	assert.NoError(err)
	assert.NotNil(r)
	assert.Equal("image", r.ResourceType())
	assert.Equal("/docs/a/b/logo.png", r.RelPermalink())
	assert.Equal("https://example.com/docs/a/b/logo.png", r.Permalink())
	img := r.(*Image)
	assert.Equal("/a/b/logo.png", img.target())

}

func TestResourcesByType(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)
	resources := Resources{
		spec.newGenericResource(nil, nil, "/public", "/a/foo1.css", "foo1.css", "css"),
		spec.newGenericResource(nil, nil, "/public", "/a/logo.png", "logo.css", "image"),
		spec.newGenericResource(nil, nil, "/public", "/a/foo2.css", "foo2.css", "css"),
		spec.newGenericResource(nil, nil, "/public", "/a/foo3.css", "foo3.css", "css")}

	assert.Len(resources.ByType("css"), 3)
	assert.Len(resources.ByType("image"), 1)

}

func TestResourcesGetByPrefix(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)
	resources := Resources{
		spec.newGenericResource(nil, nil, "/public", "/a/foo1.css", "foo1.css", "css"),
		spec.newGenericResource(nil, nil, "/public", "/a/logo1.png", "logo1.png", "image"),
		spec.newGenericResource(nil, nil, "/public", "/b/Logo2.png", "Logo2.png", "image"),
		spec.newGenericResource(nil, nil, "/public", "/b/foo2.css", "foo2.css", "css"),
		spec.newGenericResource(nil, nil, "/public", "/b/foo3.css", "foo3.css", "css")}

	assert.Nil(resources.GetByPrefix("asdf"))
	assert.Equal("/logo1.png", resources.GetByPrefix("logo").RelPermalink())
	assert.Equal("/logo1.png", resources.GetByPrefix("loGo").RelPermalink())
	assert.Equal("/Logo2.png", resources.GetByPrefix("logo2").RelPermalink())
	assert.Equal("/foo2.css", resources.GetByPrefix("foo2").RelPermalink())
	assert.Equal("/foo1.css", resources.GetByPrefix("foo1").RelPermalink())
	assert.Equal("/foo1.css", resources.GetByPrefix("foo1").RelPermalink())
	assert.Nil(resources.GetByPrefix("asdfasdf"))

	assert.Equal(2, len(resources.ByPrefix("logo")))
	assert.Equal(1, len(resources.ByPrefix("logo2")))

	logo := resources.GetByPrefix("logo")
	assert.NotNil(logo.Params())
	assert.Equal("logo1.png", logo.Name())
	assert.Equal("logo1.png", logo.Title())

}

func TestResourcesGetMatch(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)
	resources := Resources{
		spec.newGenericResource(nil, nil, "/public", "/a/foo1.css", "foo1.css", "css"),
		spec.newGenericResource(nil, nil, "/public", "/a/logo1.png", "logo1.png", "image"),
		spec.newGenericResource(nil, nil, "/public", "/b/Logo2.png", "Logo2.png", "image"),
		spec.newGenericResource(nil, nil, "/public", "/b/foo2.css", "foo2.css", "css"),
		spec.newGenericResource(nil, nil, "/public", "/b/foo3.css", "foo3.css", "css"),
		spec.newGenericResource(nil, nil, "/public", "/b/c/foo4.css", "c/foo4.css", "css"),
		spec.newGenericResource(nil, nil, "/public", "/b/c/foo5.css", "c/foo5.css", "css"),
		spec.newGenericResource(nil, nil, "/public", "/b/c/d/foo6.css", "c/d/foo6.css", "css"),
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

func TestAssignMetadata(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)

	var foo1, foo2, foo3, logo1, logo2, logo3 Resource
	var resources Resources

	for _, this := range []struct {
		metaData   []map[string]interface{}
		assertFunc func(err error)
	}{
		{[]map[string]interface{}{
			map[string]interface{}{
				"title": "My Resource",
				"name":  "My Name",
				"src":   "*",
			},
		}, func(err error) {
			assert.Equal("My Resource", logo1.Title())
			assert.Equal("My Name", logo1.Name())
			assert.Equal("My Name", foo2.Name())

		}},
		{[]map[string]interface{}{
			map[string]interface{}{
				"title": "My Logo",
				"src":   "*loGo*",
			},
			map[string]interface{}{
				"title": "My Resource",
				"name":  "My Name",
				"src":   "*",
			},
		}, func(err error) {
			assert.Equal("My Logo", logo1.Title())
			assert.Equal("My Logo", logo2.Title())
			assert.Equal("My Name", logo1.Name())
			assert.Equal("My Name", foo2.Name())
			assert.Equal("My Name", foo3.Name())
			assert.Equal("My Resource", foo3.Title())

		}},
		{[]map[string]interface{}{
			map[string]interface{}{
				"title": "My Logo",
				"src":   "*loGo*",
				"params": map[string]interface{}{
					"Param1": true,
					"icon":   "logo",
				},
			},
			map[string]interface{}{
				"title": "My Resource",
				"src":   "*",
				"params": map[string]interface{}{
					"Param2": true,
					"icon":   "resource",
				},
			},
		}, func(err error) {
			assert.NoError(err)
			assert.Equal("My Logo", logo1.Title())
			assert.Equal("My Resource", foo3.Title())
			_, p1 := logo2.Params()["param1"]
			_, p2 := foo2.Params()["param2"]
			_, p1_2 := foo2.Params()["param1"]
			_, p2_2 := logo2.Params()["param2"]

			icon1, _ := logo2.Params()["icon"]
			icon2, _ := foo2.Params()["icon"]

			assert.True(p1)
			assert.True(p2)

			// Check merge
			assert.True(p2_2)
			assert.False(p1_2)

			assert.Equal("logo", icon1)
			assert.Equal("resource", icon2)

		}},
		{[]map[string]interface{}{
			map[string]interface{}{
				"name": "Logo Name #:counter",
				"src":  "*logo*",
			},
			map[string]interface{}{
				"title": "Resource #:counter",
				"name":  "Name #:counter",
				"src":   "*",
			},
		}, func(err error) {
			assert.NoError(err)
			assert.Equal("Resource #2", logo2.Title())
			assert.Equal("Logo Name #1", logo2.Name())
			assert.Equal("Resource #4", logo1.Title())
			assert.Equal("Logo Name #2", logo1.Name())
			assert.Equal("Resource #1", foo2.Title())
			assert.Equal("Resource #3", foo1.Title())
			assert.Equal("Name #2", foo1.Name())
			assert.Equal("Resource #5", foo3.Title())

			assert.Equal(logo2, resources.GetByPrefix("logo name #1"))

		}},
		{[]map[string]interface{}{
			map[string]interface{}{
				"title": "Third Logo #:counter",
				"src":   "logo3.png",
			},
			map[string]interface{}{
				"title": "Other Logo #:counter",
				"name":  "Name #:counter",
				"src":   "logo*",
			},
		}, func(err error) {
			assert.NoError(err)
			assert.Equal("Third Logo #1", logo3.Title())
			assert.Equal("Name #3", logo3.Name())
			assert.Equal("Other Logo #1", logo2.Title())
			assert.Equal("Name #1", logo2.Name())
			assert.Equal("Other Logo #2", logo1.Title())
			assert.Equal("Name #2", logo1.Name())

		}},
		{[]map[string]interface{}{
			map[string]interface{}{
				"title": "Third Logo",
				"src":   "logo3.png",
			},
			map[string]interface{}{
				"title": "Other Logo #:counter",
				"name":  "Name #:counter",
				"src":   "logo*",
			},
		}, func(err error) {
			assert.NoError(err)
			assert.Equal("Third Logo", logo3.Title())
			assert.Equal("Name #3", logo3.Name())
			assert.Equal("Other Logo #1", logo2.Title())
			assert.Equal("Name #1", logo2.Name())
			assert.Equal("Other Logo #2", logo1.Title())
			assert.Equal("Name #2", logo1.Name())

		}},
		{[]map[string]interface{}{
			map[string]interface{}{
				"name": "third-logo",
				"src":  "logo3.png",
			},
			map[string]interface{}{
				"title": "Logo #:counter",
				"name":  "Name #:counter",
				"src":   "logo*",
			},
		}, func(err error) {
			assert.NoError(err)
			assert.Equal("Logo #3", logo3.Title())
			assert.Equal("third-logo", logo3.Name())
			assert.Equal("Logo #1", logo2.Title())
			assert.Equal("Name #1", logo2.Name())
			assert.Equal("Logo #2", logo1.Title())
			assert.Equal("Name #2", logo1.Name())

		}},
		{[]map[string]interface{}{
			map[string]interface{}{
				"title": "Third Logo #:counter",
			},
		}, func(err error) {
			// Missing src
			assert.Error(err)

		}},
		{[]map[string]interface{}{
			map[string]interface{}{
				"title": "Title",
				"src":   "[]",
			},
		}, func(err error) {
			// Invalid pattern
			assert.Error(err)

		}},
	} {

		foo2 = spec.newGenericResource(nil, nil, "/public", "/b/foo2.css", "foo2.css", "css")
		logo2 = spec.newGenericResource(nil, nil, "/public", "/b/Logo2.png", "Logo2.png", "image")
		foo1 = spec.newGenericResource(nil, nil, "/public", "/a/foo1.css", "foo1.css", "css")
		logo1 = spec.newGenericResource(nil, nil, "/public", "/a/logo1.png", "logo1.png", "image")
		foo3 = spec.newGenericResource(nil, nil, "/public", "/b/foo3.css", "foo3.css", "css")
		logo3 = spec.newGenericResource(nil, nil, "/public", "/b/logo3.png", "logo3.png", "image")

		resources = Resources{
			foo2,
			logo2,
			foo1,
			logo1,
			foo3,
			logo3,
		}

		this.assertFunc(AssignMetadata(this.metaData, resources...))
	}

}

func BenchmarkResourcesByPrefix(b *testing.B) {
	resources := benchResources(b)
	prefixes := []string{"abc", "jkl", "nomatch", "sub/"}
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resources.ByPrefix(prefixes[rnd.Intn(len(prefixes))])
		}
	})
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

	resources := Resources{spec.newGenericResource(nil, nil, "/public", "/a/"+a100, a100, "css")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resources.Match(pattern)
	}

}

func benchResources(b *testing.B) Resources {
	assert := require.New(b)
	spec := newTestResourceSpec(assert)
	var resources Resources

	for i := 0; i < 30; i++ {
		name := fmt.Sprintf("abcde%d_%d.css", i%5, i)
		resources = append(resources, spec.newGenericResource(nil, nil, "/public", "/a/"+name, name, "css"))
	}

	for i := 0; i < 30; i++ {
		name := fmt.Sprintf("efghi%d_%d.css", i%5, i)
		resources = append(resources, spec.newGenericResource(nil, nil, "/public", "/a/"+name, name, "css"))
	}

	for i := 0; i < 30; i++ {
		name := fmt.Sprintf("jklmn%d_%d.css", i%5, i)
		resources = append(resources, spec.newGenericResource(nil, nil, "/public", "/b/sub/"+name, "sub/"+name, "css"))
	}

	return resources

}

func BenchmarkAssignMetadata(b *testing.B) {
	assert := require.New(b)
	spec := newTestResourceSpec(assert)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		var resources Resources
		var meta = []map[string]interface{}{
			map[string]interface{}{
				"title": "Foo #:counter",
				"name":  "Foo Name #:counter",
				"src":   "foo1*",
			},
			map[string]interface{}{
				"title": "Rest #:counter",
				"name":  "Rest Name #:counter",
				"src":   "*",
			},
		}
		for i := 0; i < 20; i++ {
			name := fmt.Sprintf("foo%d_%d.css", i%5, i)
			resources = append(resources, spec.newGenericResource(nil, nil, "/public", "/a/"+name, name, "css"))
		}
		b.StartTimer()

		if err := AssignMetadata(meta, resources...); err != nil {
			b.Fatal(err)
		}

	}
}
