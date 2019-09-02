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

package hugolib

import (
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gohugoio/hugo/htesting"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/viper"
)

// We have many tests for the different resize operations etc. in the resource package,
// this is an integration test.
func TestImageResize(t *testing.T) {
	c := qt.New(t)
	// Make this a real as possible.
	workDir, clean, err := htesting.CreateTempDir(hugofs.Os, "image-resize")
	c.Assert(err, qt.IsNil)
	defer clean()

	newBuilder := func() *sitesBuilder {

		v := viper.New()
		v.Set("workingDir", workDir)
		v.Set("baseURL", "https://example.org")

		b := newTestSitesBuilder(t).WithWorkingDir(workDir)
		b.Fs = hugofs.NewDefault(v)
		b.WithViper(v)
		b.WithContent("mybundle/index.md", `
---
title: "My bundle"
---

`)

		b.WithTemplatesAdded("index.html", `
{{ $p := .Site.GetPage "mybundle" }}
{{ $img1 := resources.Get "images/sunset.jpg" }}
{{ $img2 := $p.Resources.GetMatch "sunset.jpg" }}
{{ $img3 := resources.GetMatch "images/*.jpg" }}
{{ $r := $img1.Resize "123x234" }}
{{ $r2 := $r.Resize "12x23" }}
{{ $b := $img2.Resize "345x678" }}
{{ $b2 := $b.Resize "34x67" }}
{{ $c := $img3.Resize "456x789" }}
{{ $fingerprinted := $img1.Resize "350x" | fingerprint }}

{{ $images := slice $r $r2 $b $b2 $c $fingerprinted }}

{{ range $i, $r := $images }}
{{ printf "Resized%d:" (add $i  1) }} {{ $r.Name }}|{{ $r.Width }}|{{ $r.Height }}|{{ $r.MediaType }}|{{ $r.RelPermalink }}|
{{ end }}

`)

		return b
	}

	imageDir := filepath.Join(workDir, "assets", "images")
	bundleDir := filepath.Join(workDir, "content", "mybundle")

	c.Assert(os.MkdirAll(imageDir, 0777), qt.IsNil)
	c.Assert(os.MkdirAll(bundleDir, 0777), qt.IsNil)
	src, err := os.Open("testdata/sunset.jpg")
	c.Assert(err, qt.IsNil)
	out, err := os.Create(filepath.Join(imageDir, "sunset.jpg"))
	c.Assert(err, qt.IsNil)
	_, err = io.Copy(out, src)
	c.Assert(err, qt.IsNil)
	out.Close()

	src.Seek(0, 0)

	out, err = os.Create(filepath.Join(bundleDir, "sunset.jpg"))
	c.Assert(err, qt.IsNil)
	_, err = io.Copy(out, src)
	c.Assert(err, qt.IsNil)
	out.Close()
	src.Close()

	b := newBuilder()
	b.Build(BuildCfg{})

	assertImage := func(width, height int, filename string) {
		filename = filepath.Join(workDir, "public", filename)
		f, err := b.Fs.Destination.Open(filename)
		c.Assert(err, qt.IsNil)
		defer f.Close()
		cfg, err := jpeg.DecodeConfig(f)
		c.Assert(cfg.Width, qt.Equals, width)
		c.Assert(cfg.Height, qt.Equals, height)
	}

	imgExpect := `
Resized1: images/sunset.jpg|123|234|image/jpg|/images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_123x234_resize_q75_box.jpg|
Resized2: images/sunset.jpg|12|23|image/jpg|/images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_ada4bb1a57f77a63306e3bd67286248e.jpg|
Resized3: sunset.jpg|345|678|image/jpg|/mybundle/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_345x678_resize_q75_box.jpg|
Resized4: sunset.jpg|34|67|image/jpg|/mybundle/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_44d8c928664d7c5a67377c6ec58425ce.jpg|
Resized5: images/sunset.jpg|456|789|image/jpg|/images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_456x789_resize_q75_box.jpg|
Resized6: images/sunset.jpg|350|219|image/jpg|/images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_350x0_resize_q75_box.a86fe88d894e5db613f6aa8a80538fefc25b20fa24ba0d782c057adcef616f56.jpg|

`

	b.AssertFileContent(filepath.Join(workDir, "public/index.html"), imgExpect)
	assertImage(350, 219, "images/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_350x0_resize_q75_box.a86fe88d894e5db613f6aa8a80538fefc25b20fa24ba0d782c057adcef616f56.jpg")

	// Build it again to make sure we read images from file cache.
	b = newBuilder()
	b.Build(BuildCfg{})

	b.AssertFileContent(filepath.Join(workDir, "public/index.html"), imgExpect)

}
