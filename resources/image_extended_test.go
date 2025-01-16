// Copyright 2024 The Hugo Authors. All rights reserved.
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

//go:build extended
// +build extended

package resources_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/htesting/hqt"
	"github.com/gohugoio/hugo/media"
)

func TestImageResizeWebP(t *testing.T) {
	c := qt.New(t)

	_, image := fetchImage(c, "sunrise.webp")

	c.Assert(image.MediaType(), qt.Equals, media.Builtin.WEBPType)
	c.Assert(image.RelPermalink(), qt.Equals, "/a/sunrise.webp")
	c.Assert(image.ResourceType(), qt.Equals, "image")
	exif := image.Exif()
	c.Assert(exif, qt.Not(qt.IsNil))
	c.Assert(exif.Tags["Copyright"], qt.Equals, "Bjørn Erik Pedersen")
	c.Assert(exif.Lat, hqt.IsSameFloat64, 36.59744166666667)
	c.Assert(exif.Long, hqt.IsSameFloat64, -4.50846)
	c.Assert(exif.Date.IsZero(), qt.Equals, false)

	resized, err := image.Resize("123x")
	c.Assert(err, qt.IsNil)
	c.Assert(image.MediaType(), qt.Equals, media.Builtin.WEBPType)
	c.Assert(resized.RelPermalink(), qt.Equals, "/a/sunrise_hu_a1deb893888915d9.webp")
	c.Assert(resized.Width(), qt.Equals, 123)
}
