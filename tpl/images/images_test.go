// Copyright 2017 The Hugo Authors. All rights reserved.
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

package images

import (
	"bytes"
	"image"
	"image/color"
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/gohugoio/hugo/resources/images"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
)

type tstNoStringer struct{}

type widthHeight struct {
	Width  int
	Height int
}

var configTests = []struct {
	path   any
	input  widthHeight
	expect any
}{
	{
		path:  "a.png",
		input: widthHeight{10, 10},
		expect: image.Config{
			Width:      10,
			Height:     10,
			ColorModel: color.NRGBAModel,
		},
	},
	{
		path:  "a.png",
		input: widthHeight{10, 10},
		expect: image.Config{
			Width:      10,
			Height:     10,
			ColorModel: color.NRGBAModel,
		},
	},
	{
		path:  "b.png",
		input: widthHeight{20, 15},
		expect: image.Config{
			Width:      20,
			Height:     15,
			ColorModel: color.NRGBAModel,
		},
	},
	{
		path:  "a.png",
		input: widthHeight{20, 15},
		expect: image.Config{
			Width:      10,
			Height:     10,
			ColorModel: color.NRGBAModel,
		},
	},
	// errors
	{path: tstNoStringer{}, expect: false},
	{path: "non-existent.png", expect: false},
	{path: "", expect: false},
}

func TestNSConfig(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	afs := afero.NewMemMapFs()
	v := config.New()
	v.Set("workingDir", "/a/b")
	d := testconfig.GetTestDeps(afs, v)
	bcfg := d.Conf

	ns := New(d)

	for _, test := range configTests {

		// check for expected errors early to avoid writing files
		if b, ok := test.expect.(bool); ok && !b {
			_, err := ns.Config(test.path)
			c.Assert(err, qt.Not(qt.IsNil))
			continue
		}

		// cast path to string for afero.WriteFile
		sp, err := cast.ToStringE(test.path)
		c.Assert(err, qt.IsNil)
		img := blankImage(d.ResourceSpec.Imaging.Codec, test.input.Width, test.input.Height)
		afero.WriteFile(ns.deps.Fs.Source, filepath.Join(bcfg.WorkingDir(), sp), img, 0o755)

		result, err := ns.Config(test.path)

		c.Assert(err, qt.IsNil)
		c.Assert(result, qt.Equals, test.expect)
		c.Assert(len(ns.cache), qt.Not(qt.Equals), 0)
	}
}

func blankImage(codec *images.Codec, width, height int) []byte {
	var buf bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	cfg := images.ImageConfig{
		TargetFormat: images.PNG,
	}
	if err := codec.EncodeTo(cfg, &buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
