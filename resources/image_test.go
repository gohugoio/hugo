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
	"image"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/spf13/afero"

	"github.com/disintegration/gift"

	"github.com/gohugoio/hugo/helpers"

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/images"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/google/go-cmp/cmp"

	"github.com/gohugoio/hugo/htesting/hqt"

	qt "github.com/frankban/quicktest"
)

var eq = qt.CmpEquals(
	cmp.Comparer(func(p1, p2 *resourceAdapter) bool {
		return p1.resourceAdapterInner == p2.resourceAdapterInner
	}),
	cmp.Comparer(func(p1, p2 os.FileInfo) bool {
		return p1.Name() == p2.Name() && p1.Size() == p2.Size() && p1.IsDir() == p2.IsDir()
	}),
	cmp.Comparer(func(p1, p2 *genericResource) bool { return p1 == p2 }),
	cmp.Comparer(func(m1, m2 media.Type) bool {
		return m1.Type() == m2.Type()
	}),
	cmp.Comparer(
		func(v1, v2 *big.Rat) bool {
			return v1.RatString() == v2.RatString()
		},
	),
	cmp.Comparer(func(v1, v2 time.Time) bool {
		return v1.Unix() == v2.Unix()
	}),
)

func TestImageTransformBasic(t *testing.T) {
	c := qt.New(t)

	image := fetchSunset(c)

	fileCache := image.(specProvider).getSpec().FileCaches.ImageCache().Fs

	assertWidthHeight := func(img resource.Image, w, h int) {
		c.Helper()
		c.Assert(img, qt.Not(qt.IsNil))
		c.Assert(img.Width(), qt.Equals, w)
		c.Assert(img.Height(), qt.Equals, h)
	}

	c.Assert(image.RelPermalink(), qt.Equals, "/a/sunset.jpg")
	c.Assert(image.ResourceType(), qt.Equals, "image")
	assertWidthHeight(image, 900, 562)

	resized, err := image.Resize("300x200")
	c.Assert(err, qt.IsNil)
	c.Assert(image != resized, qt.Equals, true)
	c.Assert(image, qt.Not(eq), resized)
	assertWidthHeight(resized, 300, 200)
	assertWidthHeight(image, 900, 562)

	resized0x, err := image.Resize("x200")
	c.Assert(err, qt.IsNil)
	assertWidthHeight(resized0x, 320, 200)
	assertFileCache(c, fileCache, path.Base(resized0x.RelPermalink()), 320, 200)

	resizedx0, err := image.Resize("200x")
	c.Assert(err, qt.IsNil)
	assertWidthHeight(resizedx0, 200, 125)
	assertFileCache(c, fileCache, path.Base(resizedx0.RelPermalink()), 200, 125)

	resizedAndRotated, err := image.Resize("x200 r90")
	c.Assert(err, qt.IsNil)
	assertWidthHeight(resizedAndRotated, 125, 200)

	assertWidthHeight(resized, 300, 200)
	c.Assert(resized.RelPermalink(), qt.Equals, "/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_300x200_resize_q68_linear.jpg")

	fitted, err := resized.Fit("50x50")
	c.Assert(err, qt.IsNil)
	c.Assert(fitted.RelPermalink(), qt.Equals, "/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_625708021e2bb281c9f1002f88e4753f.jpg")
	assertWidthHeight(fitted, 50, 33)

	// Check the MD5 key threshold
	fittedAgain, _ := fitted.Fit("10x20")
	fittedAgain, err = fittedAgain.Fit("10x20")
	c.Assert(err, qt.IsNil)
	c.Assert(fittedAgain.RelPermalink(), qt.Equals, "/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_3f65ba24dc2b7fba0f56d7f104519157.jpg")
	assertWidthHeight(fittedAgain, 10, 7)

	filled, err := image.Fill("200x100 bottomLeft")
	c.Assert(err, qt.IsNil)
	c.Assert(filled.RelPermalink(), qt.Equals, "/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x100_fill_q68_linear_bottomleft.jpg")
	assertWidthHeight(filled, 200, 100)

	smart, err := image.Fill("200x100 smart")
	c.Assert(err, qt.IsNil)
	c.Assert(smart.RelPermalink(), qt.Equals, fmt.Sprintf("/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x100_fill_q68_linear_smart%d.jpg", 1))
	assertWidthHeight(smart, 200, 100)

	// Check cache
	filledAgain, err := image.Fill("200x100 bottomLeft")
	c.Assert(err, qt.IsNil)
	c.Assert(filled, eq, filledAgain)
}

func TestImageTransformFormat(t *testing.T) {
	c := qt.New(t)

	image := fetchSunset(c)

	fileCache := image.(specProvider).getSpec().FileCaches.ImageCache().Fs

	assertExtWidthHeight := func(img resource.Image, ext string, w, h int) {
		c.Helper()
		c.Assert(img, qt.Not(qt.IsNil))
		c.Assert(helpers.Ext(img.RelPermalink()), qt.Equals, ext)
		c.Assert(img.Width(), qt.Equals, w)
		c.Assert(img.Height(), qt.Equals, h)
	}

	c.Assert(image.RelPermalink(), qt.Equals, "/a/sunset.jpg")
	c.Assert(image.ResourceType(), qt.Equals, "image")
	assertExtWidthHeight(image, ".jpg", 900, 562)

	imagePng, err := image.Resize("450x png")
	c.Assert(err, qt.IsNil)
	c.Assert(imagePng.RelPermalink(), qt.Equals, "/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_450x0_resize_linear.png")
	c.Assert(imagePng.ResourceType(), qt.Equals, "image")
	assertExtWidthHeight(imagePng, ".png", 450, 281)
	c.Assert(imagePng.Name(), qt.Equals, "sunset.jpg")
	c.Assert(imagePng.MediaType().String(), qt.Equals, "image/png")

	assertFileCache(c, fileCache, path.Base(imagePng.RelPermalink()), 450, 281)

	imageGif, err := image.Resize("225x gif")
	c.Assert(err, qt.IsNil)
	c.Assert(imageGif.RelPermalink(), qt.Equals, "/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_225x0_resize_linear.gif")
	c.Assert(imageGif.ResourceType(), qt.Equals, "image")
	assertExtWidthHeight(imageGif, ".gif", 225, 141)
	c.Assert(imageGif.Name(), qt.Equals, "sunset.jpg")
	c.Assert(imageGif.MediaType().String(), qt.Equals, "image/gif")

	assertFileCache(c, fileCache, path.Base(imageGif.RelPermalink()), 225, 141)
}

// https://github.com/gohugoio/hugo/issues/4261
func TestImageTransformLongFilename(t *testing.T) {
	c := qt.New(t)

	image := fetchImage(c, "1234567890qwertyuiopasdfghjklzxcvbnm5to6eeeeee7via8eleph.jpg")
	c.Assert(image, qt.Not(qt.IsNil))

	resized, err := image.Resize("200x")
	c.Assert(err, qt.IsNil)
	c.Assert(resized, qt.Not(qt.IsNil))
	c.Assert(resized.Width(), qt.Equals, 200)
	c.Assert(resized.RelPermalink(), qt.Equals, "/a/_hu59e56ffff1bc1d8d122b1403d34e039f_90587_65b757a6e14debeae720fe8831f0a9bc.jpg")
	resized, err = resized.Resize("100x")
	c.Assert(err, qt.IsNil)
	c.Assert(resized, qt.Not(qt.IsNil))
	c.Assert(resized.Width(), qt.Equals, 100)
	c.Assert(resized.RelPermalink(), qt.Equals, "/a/_hu59e56ffff1bc1d8d122b1403d34e039f_90587_c876768085288f41211f768147ba2647.jpg")
}

// Issue 6137
func TestImageTransformUppercaseExt(t *testing.T) {
	c := qt.New(t)
	image := fetchImage(c, "sunrise.JPG")

	resized, err := image.Resize("200x")
	c.Assert(err, qt.IsNil)
	c.Assert(resized, qt.Not(qt.IsNil))
	c.Assert(resized.Width(), qt.Equals, 200)
}

// https://github.com/gohugoio/hugo/issues/5730
func TestImagePermalinkPublishOrder(t *testing.T) {
	for _, checkOriginalFirst := range []bool{true, false} {
		name := "OriginalFirst"
		if !checkOriginalFirst {
			name = "ResizedFirst"
		}

		t.Run(name, func(t *testing.T) {
			c := qt.New(t)
			spec, workDir := newTestResourceOsFs(c)
			defer func() {
				os.Remove(workDir)
			}()

			check1 := func(img resource.Image) {
				resizedLink := "/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_100x50_resize_q75_box.jpg"
				c.Assert(img.RelPermalink(), qt.Equals, resizedLink)
				assertImageFile(c, spec.PublishFs, resizedLink, 100, 50)
			}

			check2 := func(img resource.Image) {
				c.Assert(img.RelPermalink(), qt.Equals, "/a/sunset.jpg")
				assertImageFile(c, spec.PublishFs, "a/sunset.jpg", 900, 562)
			}

			orignal := fetchImageForSpec(spec, c, "sunset.jpg")
			c.Assert(orignal, qt.Not(qt.IsNil))

			if checkOriginalFirst {
				check2(orignal)
			}

			resized, err := orignal.Resize("100x50")
			c.Assert(err, qt.IsNil)

			check1(resized.(resource.Image))

			if !checkOriginalFirst {
				check2(orignal)
			}
		})
	}
}

func TestImageTransformConcurrent(t *testing.T) {
	var wg sync.WaitGroup

	c := qt.New(t)

	spec, workDir := newTestResourceOsFs(c)
	defer func() {
		os.Remove(workDir)
	}()

	image := fetchImageForSpec(spec, c, "sunset.jpg")

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				img := image
				for k := 0; k < 2; k++ {
					r1, err := img.Resize(fmt.Sprintf("%dx", id-k))
					if err != nil {
						t.Error(err)
					}

					if r1.Width() != id-k {
						t.Errorf("Width: %d:%d", r1.Width(), j)
					}

					r2, err := r1.Resize(fmt.Sprintf("%dx", id-k-1))
					if err != nil {
						t.Error(err)
					}

					img = r2
				}
			}
		}(i + 20)
	}

	wg.Wait()
}

func TestImageWithMetadata(t *testing.T) {
	c := qt.New(t)

	image := fetchSunset(c)

	meta := []map[string]interface{}{
		{
			"title": "My Sunset",
			"name":  "Sunset #:counter",
			"src":   "*.jpg",
		},
	}

	c.Assert(AssignMetadata(meta, image), qt.IsNil)
	c.Assert(image.Name(), qt.Equals, "Sunset #1")

	resized, err := image.Resize("200x")
	c.Assert(err, qt.IsNil)
	c.Assert(resized.Name(), qt.Equals, "Sunset #1")
}

func TestImageResize8BitPNG(t *testing.T) {
	c := qt.New(t)

	image := fetchImage(c, "gohugoio.png")

	c.Assert(image.MediaType().Type(), qt.Equals, "image/png")
	c.Assert(image.RelPermalink(), qt.Equals, "/a/gohugoio.png")
	c.Assert(image.ResourceType(), qt.Equals, "image")

	resized, err := image.Resize("800x")
	c.Assert(err, qt.IsNil)
	c.Assert(resized.MediaType().Type(), qt.Equals, "image/png")
	c.Assert(resized.RelPermalink(), qt.Equals, "/a/gohugoio_hu0e1b9e4a4be4d6f86c7b37b9ccce3fbc_73886_800x0_resize_linear_2.png")
	c.Assert(resized.Width(), qt.Equals, 800)
}

func TestImageResizeInSubPath(t *testing.T) {
	c := qt.New(t)

	image := fetchImage(c, "sub/gohugoio2.png")

	c.Assert(image.MediaType(), eq, media.PNGType)
	c.Assert(image.RelPermalink(), qt.Equals, "/a/sub/gohugoio2.png")
	c.Assert(image.ResourceType(), qt.Equals, "image")

	resized, err := image.Resize("101x101")
	c.Assert(err, qt.IsNil)
	c.Assert(resized.MediaType().Type(), qt.Equals, "image/png")
	c.Assert(resized.RelPermalink(), qt.Equals, "/a/sub/gohugoio2_hu0e1b9e4a4be4d6f86c7b37b9ccce3fbc_73886_101x101_resize_linear_2.png")
	c.Assert(resized.Width(), qt.Equals, 101)

	publishedImageFilename := filepath.Clean(resized.RelPermalink())

	spec := image.(specProvider).getSpec()

	assertImageFile(c, spec.BaseFs.PublishFs, publishedImageFilename, 101, 101)
	c.Assert(spec.BaseFs.PublishFs.Remove(publishedImageFilename), qt.IsNil)

	// Cleare mem cache to simulate reading from the file cache.
	spec.imageCache.clear()

	resizedAgain, err := image.Resize("101x101")
	c.Assert(err, qt.IsNil)
	c.Assert(resizedAgain.RelPermalink(), qt.Equals, "/a/sub/gohugoio2_hu0e1b9e4a4be4d6f86c7b37b9ccce3fbc_73886_101x101_resize_linear_2.png")
	c.Assert(resizedAgain.Width(), qt.Equals, 101)
	assertImageFile(c, image.(specProvider).getSpec().BaseFs.PublishFs, publishedImageFilename, 101, 101)
}

func TestSVGImage(t *testing.T) {
	c := qt.New(t)
	spec := newTestResourceSpec(specDescriptor{c: c})
	svg := fetchResourceForSpec(spec, c, "circle.svg")
	c.Assert(svg, qt.Not(qt.IsNil))
}

func TestSVGImageContent(t *testing.T) {
	c := qt.New(t)
	spec := newTestResourceSpec(specDescriptor{c: c})
	svg := fetchResourceForSpec(spec, c, "circle.svg")
	c.Assert(svg, qt.Not(qt.IsNil))

	content, err := svg.Content()
	c.Assert(err, qt.IsNil)
	c.Assert(content, hqt.IsSameType, "")
	c.Assert(content.(string), qt.Contains, `<svg height="100" width="100">`)
}

func TestImageExif(t *testing.T) {
	c := qt.New(t)
	fs := afero.NewMemMapFs()
	spec := newTestResourceSpec(specDescriptor{fs: fs, c: c})
	image := fetchResourceForSpec(spec, c, "sunset.jpg").(resource.Image)

	getAndCheckExif := func(c *qt.C, image resource.Image) {
		x, err := image.Exif()
		c.Assert(err, qt.IsNil)
		c.Assert(x, qt.Not(qt.IsNil))

		c.Assert(x.Date.Format("2006-01-02"), qt.Equals, "2017-10-27")

		// Malaga: https://goo.gl/taazZy
		c.Assert(x.Lat, qt.Equals, float64(36.59744166666667))
		c.Assert(x.Long, qt.Equals, float64(-4.50846))

		v, found := x.Tags["LensModel"]
		c.Assert(found, qt.Equals, true)
		lensModel, ok := v.(string)
		c.Assert(ok, qt.Equals, true)
		c.Assert(lensModel, qt.Equals, "smc PENTAX-DA* 16-50mm F2.8 ED AL [IF] SDM")
		resized, _ := image.Resize("300x200")
		x2, _ := resized.Exif()
		c.Assert(x2, eq, x)
	}

	getAndCheckExif(c, image)
	image = fetchResourceForSpec(spec, c, "sunset.jpg").(resource.Image)
	// This will read from file cache.
	getAndCheckExif(c, image)

}

func BenchmarkImageExif(b *testing.B) {

	getImages := func(c *qt.C, b *testing.B, fs afero.Fs) []resource.Image {
		spec := newTestResourceSpec(specDescriptor{fs: fs, c: c})
		images := make([]resource.Image, b.N)
		for i := 0; i < b.N; i++ {
			images[i] = fetchResourceForSpec(spec, c, "sunset.jpg", strconv.Itoa(i)).(resource.Image)
		}
		return images
	}

	getAndCheckExif := func(c *qt.C, image resource.Image) {
		x, err := image.Exif()
		c.Assert(err, qt.IsNil)
		c.Assert(x, qt.Not(qt.IsNil))
		c.Assert(x.Long, qt.Equals, float64(-4.50846))

	}

	b.Run("Cold cache", func(b *testing.B) {
		b.StopTimer()
		c := qt.New(b)
		images := getImages(c, b, afero.NewMemMapFs())

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			getAndCheckExif(c, images[i])
		}

	})

	b.Run("Cold cache, 10", func(b *testing.B) {
		b.StopTimer()
		c := qt.New(b)
		images := getImages(c, b, afero.NewMemMapFs())

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			for j := 0; j < 10; j++ {
				getAndCheckExif(c, images[i])
			}
		}

	})

	b.Run("Warm cache", func(b *testing.B) {
		b.StopTimer()
		c := qt.New(b)
		fs := afero.NewMemMapFs()
		images := getImages(c, b, fs)
		for i := 0; i < b.N; i++ {
			getAndCheckExif(c, images[i])
		}

		images = getImages(c, b, fs)

		b.StartTimer()
		for i := 0; i < b.N; i++ {
			getAndCheckExif(c, images[i])
		}

	})

}

// usesFMA indicates whether "fused multiply and add" (FMA) instruction is
// used.  The command "grep FMADD go/test/codegen/floats.go" can help keep
// the FMA-using architecture list updated.
var usesFMA = runtime.GOARCH == "s390x" ||
	runtime.GOARCH == "ppc64" ||
	runtime.GOARCH == "ppc64le" ||
	runtime.GOARCH == "arm64"

// goldenEqual compares two NRGBA images.  It is used in golden tests only.
// A small tolerance is allowed on architectures using "fused multiply and add"
// (FMA) instruction to accommodate for floating-point rounding differences
// with control golden images that were generated on amd64 architecture.
// See https://golang.org/ref/spec#Floating_point_operators
// and https://github.com/gohugoio/hugo/issues/6387 for more information.
//
// Borrowed from https://github.com/disintegration/gift/blob/a999ff8d5226e5ab14b64a94fca07c4ac3f357cf/gift_test.go#L598-L625
// Copyright (c) 2014-2019 Grigory Dryapak
// Licensed under the MIT License.
func goldenEqual(img1, img2 *image.NRGBA) bool {
	maxDiff := 0
	if usesFMA {
		maxDiff = 1
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

func TestImageOperationsGolden(t *testing.T) {
	c := qt.New(t)
	c.Parallel()

	devMode := false

	testImages := []string{"sunset.jpg", "gohugoio8.png", "gohugoio24.png"}

	spec, workDir := newTestResourceOsFs(c)
	defer func() {
		if !devMode {
			os.Remove(workDir)
		}
	}()

	if devMode {
		fmt.Println(workDir)
	}

	// Test PNGs with alpha channel.
	for _, img := range []string{"gopher-hero8.png", "gradient-circle.png"} {
		orig := fetchImageForSpec(spec, c, img)
		for _, resizeSpec := range []string{"200x #e3e615", "200x jpg #e3e615"} {
			resized, err := orig.Resize(resizeSpec)
			c.Assert(err, qt.IsNil)
			rel := resized.RelPermalink()
			c.Log("resize", rel)
			c.Assert(rel, qt.Not(qt.Equals), "")
		}
	}

	for _, img := range testImages {

		orig := fetchImageForSpec(spec, c, img)
		for _, resizeSpec := range []string{"200x100", "600x", "200x r90 q50 Box"} {
			resized, err := orig.Resize(resizeSpec)
			c.Assert(err, qt.IsNil)
			rel := resized.RelPermalink()
			c.Log("resize", rel)
			c.Assert(rel, qt.Not(qt.Equals), "")
		}

		for _, fillSpec := range []string{"300x200 Gaussian Smart", "100x100 Center", "300x100 TopLeft NearestNeighbor", "400x200 BottomLeft"} {
			resized, err := orig.Fill(fillSpec)
			c.Assert(err, qt.IsNil)
			rel := resized.RelPermalink()
			c.Log("fill", rel)
			c.Assert(rel, qt.Not(qt.Equals), "")
		}

		for _, fitSpec := range []string{"300x200 Linear"} {
			resized, err := orig.Fit(fitSpec)
			c.Assert(err, qt.IsNil)
			rel := resized.RelPermalink()
			c.Log("fit", rel)
			c.Assert(rel, qt.Not(qt.Equals), "")
		}

		f := &images.Filters{}

		filters := []gift.Filter{
			f.Grayscale(),
			f.GaussianBlur(6),
			f.Saturation(50),
			f.Sepia(100),
			f.Brightness(30),
			f.ColorBalance(10, -10, -10),
			f.Colorize(240, 50, 100),
			f.Gamma(1.5),
			f.UnsharpMask(1, 1, 0),
			f.Sigmoid(0.5, 7),
			f.Pixelate(5),
			f.Invert(),
			f.Hue(22),
			f.Contrast(32.5),
		}

		resized, err := orig.Fill("400x200 center")
		c.Assert(err, qt.IsNil)

		for _, filter := range filters {
			resized, err := resized.Filter(filter)
			c.Assert(err, qt.IsNil)
			rel := resized.RelPermalink()
			c.Logf("filter: %v %s", filter, rel)
			c.Assert(rel, qt.Not(qt.Equals), "")
		}

		resized, err = resized.Filter(filters[0:4])
		c.Assert(err, qt.IsNil)
		rel := resized.RelPermalink()
		c.Log("filter all", rel)
		c.Assert(rel, qt.Not(qt.Equals), "")
	}

	if devMode {
		return
	}

	dir1 := filepath.Join(workDir, "resources/_gen/images")
	dir2 := filepath.FromSlash("testdata/golden")

	// The two dirs above should now be the same.
	dirinfos1, err := ioutil.ReadDir(dir1)
	c.Assert(err, qt.IsNil)
	dirinfos2, err := ioutil.ReadDir(dir2)
	c.Assert(err, qt.IsNil)
	c.Assert(len(dirinfos1), qt.Equals, len(dirinfos2))

	for i, fi1 := range dirinfos1 {
		fi2 := dirinfos2[i]
		c.Assert(fi1.Name(), qt.Equals, fi2.Name())

		f1, err := os.Open(filepath.Join(dir1, fi1.Name()))
		c.Assert(err, qt.IsNil)
		f2, err := os.Open(filepath.Join(dir2, fi2.Name()))
		c.Assert(err, qt.IsNil)

		img1, _, err := image.Decode(f1)
		c.Assert(err, qt.IsNil)
		img2, _, err := image.Decode(f2)
		c.Assert(err, qt.IsNil)

		nrgba1 := image.NewNRGBA(img1.Bounds())
		gift.New().Draw(nrgba1, img1)
		nrgba2 := image.NewNRGBA(img2.Bounds())
		gift.New().Draw(nrgba2, img2)

		if !goldenEqual(nrgba1, nrgba2) {
			switch fi1.Name() {
			case "gohugoio8_hu7f72c00afdf7634587afaa5eff2a25b2_73538_4c320010919da2d8b63ed24818b4d8e1.png",
				"gohugoio8_hu7f72c00afdf7634587afaa5eff2a25b2_73538_9d4c2220235b3c2d9fa6506be571560f.png",
				"gohugoio8_hu7f72c00afdf7634587afaa5eff2a25b2_73538_c74bb417b961e09cf1aac2130b7b9b85.png",
				"gohugoio8_hu7f72c00afdf7634587afaa5eff2a25b2_73538_300x200_fill_gaussian_smart1_2.png":
				c.Log("expectedly differs from golden due to dithering:", fi1.Name())
			default:
				t.Errorf("resulting image differs from golden: %s", fi1.Name())
			}
		}

		if !usesFMA {
			c.Assert(fi1, eq, fi2)

			_, err = f1.Seek(0, 0)
			c.Assert(err, qt.IsNil)
			_, err = f2.Seek(0, 0)
			c.Assert(err, qt.IsNil)

			hash1, err := helpers.MD5FromReader(f1)
			c.Assert(err, qt.IsNil)
			hash2, err := helpers.MD5FromReader(f2)
			c.Assert(err, qt.IsNil)

			c.Assert(hash1, qt.Equals, hash2)
		}

		f1.Close()
		f2.Close()
	}

}

func BenchmarkResizeParallel(b *testing.B) {
	c := qt.New(b)
	img := fetchSunset(c)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := rand.Intn(10) + 10
			resized, err := img.Resize(strconv.Itoa(w) + "x")
			if err != nil {
				b.Fatal(err)
			}
			_, err = resized.Resize(strconv.Itoa(w-1) + "x")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
