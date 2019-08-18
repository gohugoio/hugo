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
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/google/go-cmp/cmp"

	"github.com/gohugoio/hugo/htesting/hqt"

	qt "github.com/frankban/quicktest"
)

var eq = qt.CmpEquals(
	cmp.Comparer(func(p1, p2 *resourceAdapter) bool {
		return p1.resourceAdapterInner == p2.resourceAdapterInner
	}),
	cmp.Comparer(func(p1, p2 *genericResource) bool { return p1 == p2 }),
	cmp.Comparer(func(m1, m2 media.Type) bool {
		return m1.Type() == m2.Type()
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
	assertFileCache(c, fileCache, resized0x.RelPermalink(), 320, 200)

	resizedx0, err := image.Resize("200x")
	c.Assert(err, qt.IsNil)
	assertWidthHeight(resizedx0, 200, 125)
	assertFileCache(c, fileCache, resizedx0.RelPermalink(), 200, 125)

	resizedAndRotated, err := image.Resize("x200 r90")
	c.Assert(err, qt.IsNil)
	assertWidthHeight(resizedAndRotated, 125, 200)
	assertFileCache(c, fileCache, resizedAndRotated.RelPermalink(), 125, 200)

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
	assertWidthHeight(fittedAgain, 10, 6)

	filled, err := image.Fill("200x100 bottomLeft")
	c.Assert(err, qt.IsNil)
	c.Assert(filled.RelPermalink(), qt.Equals, "/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x100_fill_q68_linear_bottomleft.jpg")
	assertWidthHeight(filled, 200, 100)
	assertFileCache(c, fileCache, filled.RelPermalink(), 200, 100)

	smart, err := image.Fill("200x100 smart")
	c.Assert(err, qt.IsNil)
	c.Assert(smart.RelPermalink(), qt.Equals, fmt.Sprintf("/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x100_fill_q68_linear_smart%d.jpg", 1))
	assertWidthHeight(smart, 200, 100)
	assertFileCache(c, fileCache, smart.RelPermalink(), 200, 100)

	// Check cache
	filledAgain, err := image.Fill("200x100 bottomLeft")
	c.Assert(err, qt.IsNil)
	c.Assert(filled, eq, filledAgain)
	assertFileCache(c, fileCache, filledAgain.RelPermalink(), 200, 100)
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
			spec := newTestResourceOsFs(c)

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

	spec := newTestResourceOsFs(c)

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
	fileCache := image.(specProvider).getSpec().FileCaches.ImageCache().Fs

	c.Assert(image.MediaType(), eq, media.PNGType)
	c.Assert(image.RelPermalink(), qt.Equals, "/a/sub/gohugoio2.png")
	c.Assert(image.ResourceType(), qt.Equals, "image")

	resized, err := image.Resize("101x101")
	c.Assert(err, qt.IsNil)
	c.Assert(resized.MediaType().Type(), qt.Equals, "image/png")
	c.Assert(resized.RelPermalink(), qt.Equals, "/a/sub/gohugoio2_hu0e1b9e4a4be4d6f86c7b37b9ccce3fbc_73886_101x101_resize_linear_2.png")
	c.Assert(resized.Width(), qt.Equals, 101)

	assertFileCache(c, fileCache, resized.RelPermalink(), 101, 101)
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
	assertFileCache(c, fileCache, resizedAgain.RelPermalink(), 101, 101)
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
