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

package imagetesting

import (
	"image"
	"image/gif"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp"

	"github.com/disintegration/gift"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/htesting"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
)

var eq = qt.CmpEquals(
	cmp.Comparer(func(p1, p2 os.FileInfo) bool {
		return p1.Name() == p2.Name() && p1.Size() == p2.Size() && p1.IsDir() == p2.IsDir()
	}),
	cmp.Comparer(func(d1, d2 fs.DirEntry) bool {
		p1, err1 := d1.Info()
		p2, err2 := d2.Info()
		if err1 != nil || err2 != nil {
			return false
		}
		return p1.Name() == p2.Name() && p1.Size() == p2.Size() && p1.IsDir() == p2.IsDir()
	}),
)

// GoldenImageTestOpts provides options for a golden image test.
type GoldenImageTestOpts struct {
	// The test.
	T testing.TB

	// Name of the test. Will be used as the base folder for generated images.
	// Slashes allowed and encouraged.
	Name string

	// The test site's files in txttar format.
	Files string

	// Set to true to write golden files to disk.
	WriteFiles bool

	// Set to true to skip any assertions. Useful when adding new golden variants to a test.
	DevMode bool
}

// To rebuild all Golden image tests, toggle WriteFiles=true and run:
// GOARCH=amd64 go test -count 1 -timeout 30s -run "^TestImagesGolden" ./...
// TODO(bep) see if we can do this via flags.
var DefaultGoldenOpts = GoldenImageTestOpts{
	WriteFiles: false,
	DevMode:    false,
}

func RunGolden(opts GoldenImageTestOpts) *hugolib.IntegrationTestBuilder {
	opts.T.Helper()

	c := hugolib.Test(opts.T, opts.Files, hugolib.TestOptWithOSFs()) // hugolib.TestOptWithPrintAndKeepTempDir(true))
	c.AssertFileContent("public/index.html", "Home.")

	outputDir := filepath.Join(c.H.Conf.WorkingDir(), "public", "images")
	goldenBaseDir := filepath.Join("testdata", "images_golden")
	goldenDir := filepath.Join(goldenBaseDir, filepath.FromSlash(opts.Name))
	if opts.WriteFiles {
		c.Assert(htesting.IsRealCI(), qt.IsFalse)
		c.Assert(os.MkdirAll(goldenBaseDir, 0o777), qt.IsNil)
		c.Assert(os.RemoveAll(goldenDir), qt.IsNil)
		c.Assert(hugio.CopyDir(hugofs.Os, outputDir, goldenDir, nil), qt.IsNil)
		return c
	}

	if opts.DevMode {
		c.Assert(htesting.IsRealCI(), qt.IsFalse)
		return c
	}

	decodeAll := func(f *os.File) []image.Image {
		c.Helper()

		var images []image.Image

		if strings.HasSuffix(f.Name(), ".gif") {
			gif, err := gif.DecodeAll(f)
			c.Assert(err, qt.IsNil, qt.Commentf(f.Name()))
			images = make([]image.Image, len(gif.Image))
			for i, img := range gif.Image {
				images[i] = img
			}
		} else {
			img, _, err := image.Decode(f)
			c.Assert(err, qt.IsNil, qt.Commentf(f.Name()))
			images = append(images, img)
		}
		return images
	}

	entries1, err := os.ReadDir(outputDir)
	c.Assert(err, qt.IsNil)
	entries2, err := os.ReadDir(goldenDir)
	c.Assert(err, qt.IsNil)
	c.Assert(len(entries1), qt.Equals, len(entries2))
	for i, e1 := range entries1 {
		c.Assert(filepath.Ext(e1.Name()), qt.Not(qt.Equals), "")
		func() {
			e2 := entries2[i]

			f1, err := os.Open(filepath.Join(outputDir, e1.Name()))
			c.Assert(err, qt.IsNil)
			defer f1.Close()

			f2, err := os.Open(filepath.Join(goldenDir, e2.Name()))
			c.Assert(err, qt.IsNil)
			defer f2.Close()

			imgs2 := decodeAll(f2)
			imgs1 := decodeAll(f1)
			c.Assert(len(imgs1), qt.Equals, len(imgs2))

			if !UsesFMA {
				c.Assert(e1, eq, e2)
				_, err = f1.Seek(0, 0)
				c.Assert(err, qt.IsNil)
				_, err = f2.Seek(0, 0)
				c.Assert(err, qt.IsNil)

				hash1, _, err := hashing.XXHashFromReader(f1)
				c.Assert(err, qt.IsNil)
				hash2, _, err := hashing.XXHashFromReader(f2)
				c.Assert(err, qt.IsNil)

				c.Assert(hash1, qt.Equals, hash2)
			}

			for i, img1 := range imgs1 {
				img2 := imgs2[i]
				nrgba1 := image.NewNRGBA(img1.Bounds())
				gift.New().Draw(nrgba1, img1)
				nrgba2 := image.NewNRGBA(img2.Bounds())
				gift.New().Draw(nrgba2, img2)
				c.Assert(goldenEqual(nrgba1, nrgba2), qt.Equals, true, qt.Commentf(e1.Name()))
			}
		}()
	}
	return c
}

// goldenEqual compares two NRGBA images.  It is used in golden tests only.
// A small tolerance is allowed on architectures using "fused multiply and add"
// (FMA) instruction to accommodate for floating-point rounding differences
// with control golden images that were generated on amd64 architecture.
// See https://golang.org/ref/spec#Floating_point_operators
// and https://github.com/gohugoio/hugo/issues/6387 for more information.
//
// Based on https://github.com/disintegration/gift/blob/a999ff8d5226e5ab14b64a94fca07c4ac3f357cf/gift_test.go#L598-L625
// Copyright (c) 2014-2019 Grigory Dryapak
// Licensed under the MIT License.
func goldenEqual(img1, img2 *image.NRGBA) bool {
	maxDiff := 0
	if runtime.GOARCH != "amd64" {
		// The golden files are created using the AMD64 architecture.
		// Be lenient on other platforms due to floaging point and dithering differences.
		maxDiff = 15
	}
	if !img1.Rect.Eq(img2.Rect) {
		return false
	}
	if len(img1.Pix) != len(img2.Pix) {
		return false
	}
	for i := 0; i < len(img1.Pix); i++ {
		diff := int(img1.Pix[i]) - int(img2.Pix[i])
		if diff < 0 {
			diff = -diff
		}
		if diff > maxDiff {
			return false
		}
	}
	return true
}

// We don't have a CI test environment for these, and there are known dithering issues that makes these time consuming to maintain.
var SkipGoldenTests = runtime.GOARCH == "ppc64" || runtime.GOARCH == "ppc64le" || runtime.GOARCH == "s390x"

// UsesFMA indicates whether "fused multiply and add" (FMA) instruction is
// used.  The command "grep FMADD go/test/codegen/floats.go" can help keep
// the FMA-using architecture list updated.
var UsesFMA = runtime.GOARCH == "s390x" ||
	runtime.GOARCH == "ppc64" ||
	runtime.GOARCH == "ppc64le" ||
	runtime.GOARCH == "arm64" ||
	runtime.GOARCH == "riscv64"
