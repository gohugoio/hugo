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

// +build extended

package resources

import (
	"testing"

	"github.com/gohugoio/hugo/media"

	qt "github.com/frankban/quicktest"
)

func TestImageResizeWebP(t *testing.T) {
	c := qt.New(t)

	image := fetchImage(c, "sunset.webp")

	c.Assert(image.MediaType(), qt.Equals, media.WEBPType)
	c.Assert(image.RelPermalink(), qt.Equals, "/a/sunset.webp")
	c.Assert(image.ResourceType(), qt.Equals, "image")
	c.Assert(image.Exif(), qt.IsNil)

	resized, err := image.Resize("123x")
	c.Assert(err, qt.IsNil)
	c.Assert(image.MediaType(), qt.Equals, media.WEBPType)
	c.Assert(resized.RelPermalink(), qt.Equals, "/a/sunset_hu36ee0b61ba924719ad36da960c273f96_59826_123x0_resize_q68_h2_linear_2.webp")
	c.Assert(resized.Width(), qt.Equals, 123)
}
