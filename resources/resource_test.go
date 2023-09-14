// Copyright 2023 The Hugo Authors. All rights reserved.
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

package resources_test

import (
	"testing"

	"github.com/gohugoio/hugo/resources"

	"github.com/gohugoio/hugo/media"

	qt "github.com/frankban/quicktest"
)

func TestNewResourceFromFilename(t *testing.T) {
	c := qt.New(t)
	spec := newTestResourceSpec(specDescriptor{c: c})

	writeSource(t, spec.Fs, "assets/a/b/logo.png", "image")
	writeSource(t, spec.Fs, "assets/a/b/data.json", "json")

	r, err := spec.New(resources.ResourceSourceDescriptor{Fs: spec.BaseFs.Assets.Fs, SourceFilename: "a/b/logo.png"})

	c.Assert(err, qt.IsNil)
	c.Assert(r, qt.IsNotNil)
	c.Assert(r.ResourceType(), qt.Equals, "image")
	c.Assert(r.RelPermalink(), qt.Equals, "/a/b/logo.png")
	c.Assert(r.Permalink(), qt.Equals, "https://example.com/a/b/logo.png")

	r, err = spec.New(resources.ResourceSourceDescriptor{Fs: spec.BaseFs.Assets.Fs, SourceFilename: "a/b/data.json"})

	c.Assert(err, qt.IsNil)
	c.Assert(r, qt.IsNotNil)
	c.Assert(r.ResourceType(), qt.Equals, "application")
}

var pngType, _ = media.FromStringAndExt("image/png", "png")
