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
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenericResource(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)

	r := spec.newGenericResource(nil, nil, "/public", "/a/foo.css", "foo.css", "css")

	assert.Equal("https://example.com/foo.css", r.Permalink())
	assert.Equal("foo.css", r.RelPermalink())
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
	assert.Equal("a/b/logo.png", r.RelPermalink())
	assert.Equal("https://example.com/a/b/logo.png", r.Permalink())

	r, err = spec.NewResourceFromFilename(nil, "/public", "/root/a/b/data.json", "a/b/data.json")

	assert.NoError(err)
	assert.NotNil(r)
	assert.Equal("json", r.ResourceType())

	cloned := r.(Cloner).WithNewBase("aceof")
	assert.Equal(r.ResourceType(), cloned.ResourceType())
	assert.Equal("/aceof/a/b/data.json", cloned.RelPermalink())
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
		spec.newGenericResource(nil, nil, "/public", "/b/logo2.png", "logo2.png", "image"),
		spec.newGenericResource(nil, nil, "/public", "/b/foo2.css", "foo2.css", "css"),
		spec.newGenericResource(nil, nil, "/public", "/b/foo3.css", "foo3.css", "css")}

	assert.Nil(resources.GetByPrefix("asdf"))
	assert.Equal("logo1.png", resources.GetByPrefix("logo").RelPermalink())
	assert.Equal("foo2.css", resources.GetByPrefix("foo2").RelPermalink())
	assert.Equal("foo1.css", resources.GetByPrefix("foo1").RelPermalink())
	assert.Equal("foo1.css", resources.GetByPrefix("foo1").RelPermalink())
	assert.Nil(resources.GetByPrefix("asdfasdf"))

}
