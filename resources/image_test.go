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
	"testing"

	"github.com/gohugoio/hugo/htesting/hqt"

	"github.com/disintegration/imaging"

	"sync"

	qt "github.com/frankban/quicktest"
)

func TestParseImageConfig(t *testing.T) {
	for i, this := range []struct {
		in     string
		expect interface{}
	}{
		{"300x400", newImageConfig(300, 400, 0, 0, "", "")},
		{"100x200 bottomRight", newImageConfig(100, 200, 0, 0, "", "BottomRight")},
		{"10x20 topleft Lanczos", newImageConfig(10, 20, 0, 0, "Lanczos", "topleft")},
		{"linear left 10x r180", newImageConfig(10, 0, 0, 180, "linear", "left")},
		{"x20 riGht Cosine q95", newImageConfig(0, 20, 95, 0, "cosine", "right")},

		{"", false},
		{"foo", false},
	} {
		result, err := parseImageConfig(this.in)
		if b, ok := this.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] parseImageConfig didn't return an expected error", i)
			}
		} else {
			if err != nil {
				t.Fatalf("[%d] err: %s", i, err)
			}
			if fmt.Sprint(result) != fmt.Sprint(this.expect) {
				t.Fatalf("[%d] got\n%v\n but expected\n%v", i, result, this.expect)
			}
		}
	}
}

func TestImageTransformBasic(t *testing.T) {

	c := qt.New(t)

	image := fetchSunset(c)
	fileCache := image.spec.FileCaches.ImageCache().Fs

	c.Assert(image.RelPermalink(), qt.Equals, "/a/sunset.jpg")
	c.Assert(image.ResourceType(), qt.Equals, "image")

	resized, err := image.Resize("300x200")
	c.Assert(err, qt.IsNil)
	c.Assert(image != resized, qt.Equals, true)
	c.Assert(image.genericResource != resized.genericResource, qt.Equals, true)
	c.Assert(image.sourceFilename != resized.sourceFilename, qt.Equals, true)

	resized0x, err := image.Resize("x200")
	c.Assert(err, qt.IsNil)
	c.Assert(resized0x.Width(), qt.Equals, 320)
	c.Assert(resized0x.Height(), qt.Equals, 200)

	assertFileCache(c, fileCache, resized0x.RelPermalink(), 320, 200)

	resizedx0, err := image.Resize("200x")
	c.Assert(err, qt.IsNil)
	c.Assert(resizedx0.Width(), qt.Equals, 200)
	c.Assert(resizedx0.Height(), qt.Equals, 125)
	assertFileCache(c, fileCache, resizedx0.RelPermalink(), 200, 125)

	resizedAndRotated, err := image.Resize("x200 r90")
	c.Assert(err, qt.IsNil)
	c.Assert(resizedAndRotated.Width(), qt.Equals, 125)
	c.Assert(resizedAndRotated.Height(), qt.Equals, 200)
	assertFileCache(c, fileCache, resizedAndRotated.RelPermalink(), 125, 200)

	c.Assert(resized.RelPermalink(), qt.Equals, "/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_300x200_resize_q68_linear.jpg")
	c.Assert(resized.Width(), qt.Equals, 300)
	c.Assert(resized.Height(), qt.Equals, 200)

	fitted, err := resized.Fit("50x50")
	c.Assert(err, qt.IsNil)
	c.Assert(fitted.RelPermalink(), qt.Equals, "/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_625708021e2bb281c9f1002f88e4753f.jpg")
	c.Assert(fitted.Width(), qt.Equals, 50)
	c.Assert(fitted.Height(), qt.Equals, 33)

	// Check the MD5 key threshold
	fittedAgain, _ := fitted.Fit("10x20")
	fittedAgain, err = fittedAgain.Fit("10x20")
	c.Assert(err, qt.IsNil)
	c.Assert(fittedAgain.RelPermalink(), qt.Equals, "/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_3f65ba24dc2b7fba0f56d7f104519157.jpg")
	c.Assert(fittedAgain.Width(), qt.Equals, 10)
	c.Assert(fittedAgain.Height(), qt.Equals, 6)

	filled, err := image.Fill("200x100 bottomLeft")
	c.Assert(err, qt.IsNil)
	c.Assert(filled.RelPermalink(), qt.Equals, "/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x100_fill_q68_linear_bottomleft.jpg")
	c.Assert(filled.Width(), qt.Equals, 200)
	c.Assert(filled.Height(), qt.Equals, 100)
	assertFileCache(c, fileCache, filled.RelPermalink(), 200, 100)

	smart, err := image.Fill("200x100 smart")
	c.Assert(err, qt.IsNil)
	c.Assert(smart.RelPermalink(), qt.Equals, fmt.Sprintf("/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x100_fill_q68_linear_smart%d.jpg", smartCropVersionNumber))
	c.Assert(smart.Width(), qt.Equals, 200)
	c.Assert(smart.Height(), qt.Equals, 100)
	assertFileCache(c, fileCache, smart.RelPermalink(), 200, 100)

	// Check cache
	filledAgain, err := image.Fill("200x100 bottomLeft")
	c.Assert(err, qt.IsNil)
	c.Assert(filled == filledAgain, qt.Equals, true)
	c.Assert(filled.sourceFilename == filledAgain.sourceFilename, qt.Equals, true)
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

			check1 := func(img *Image) {
				resizedLink := "/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_100x50_resize_q75_box.jpg"
				c.Assert(img.RelPermalink(), qt.Equals, resizedLink)
				assertImageFile(c, spec.PublishFs, resizedLink, 100, 50)
			}

			check2 := func(img *Image) {
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

			check1(resized)

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

					_, err = r2.decodeSource()
					if err != nil {
						t.Error("Err decode:", err)
					}

					img = r1
				}
			}
		}(i + 20)
	}

	wg.Wait()
}

func TestDecodeImaging(t *testing.T) {
	c := qt.New(t)
	m := map[string]interface{}{
		"quality":        42,
		"resampleFilter": "NearestNeighbor",
		"anchor":         "topLeft",
	}

	imaging, err := decodeImaging(m)

	c.Assert(err, qt.IsNil)
	c.Assert(imaging.Quality, qt.Equals, 42)
	c.Assert(imaging.ResampleFilter, qt.Equals, "nearestneighbor")
	c.Assert(imaging.Anchor, qt.Equals, "topleft")

	m = map[string]interface{}{}

	imaging, err = decodeImaging(m)
	c.Assert(err, qt.IsNil)
	c.Assert(imaging.Quality, qt.Equals, defaultJPEGQuality)
	c.Assert(imaging.ResampleFilter, qt.Equals, "box")
	c.Assert(imaging.Anchor, qt.Equals, "smart")

	_, err = decodeImaging(map[string]interface{}{
		"quality": 123,
	})
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = decodeImaging(map[string]interface{}{
		"resampleFilter": "asdf",
	})
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = decodeImaging(map[string]interface{}{
		"anchor": "asdf",
	})
	c.Assert(err, qt.Not(qt.IsNil))

	imaging, err = decodeImaging(map[string]interface{}{
		"anchor": "Smart",
	})
	c.Assert(err, qt.IsNil)
	c.Assert(imaging.Anchor, qt.Equals, "smart")

}

func TestImageWithMetadata(t *testing.T) {
	c := qt.New(t)

	image := fetchSunset(c)

	var meta = []map[string]interface{}{
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

	c.Assert(image.format, qt.Equals, imaging.PNG)
	c.Assert(image.RelPermalink(), qt.Equals, "/a/gohugoio.png")
	c.Assert(image.ResourceType(), qt.Equals, "image")

	resized, err := image.Resize("800x")
	c.Assert(err, qt.IsNil)
	c.Assert(resized.format, qt.Equals, imaging.PNG)
	c.Assert(resized.RelPermalink(), qt.Equals, "/a/gohugoio_hu0e1b9e4a4be4d6f86c7b37b9ccce3fbc_73886_800x0_resize_linear_2.png")
	c.Assert(resized.Width(), qt.Equals, 800)

}

func TestImageResizeInSubPath(t *testing.T) {

	c := qt.New(t)

	image := fetchImage(c, "sub/gohugoio2.png")
	fileCache := image.spec.FileCaches.ImageCache().Fs

	c.Assert(image.format, qt.Equals, imaging.PNG)
	c.Assert(image.RelPermalink(), qt.Equals, "/a/sub/gohugoio2.png")
	c.Assert(image.ResourceType(), qt.Equals, "image")

	resized, err := image.Resize("101x101")
	c.Assert(err, qt.IsNil)
	c.Assert(resized.format, qt.Equals, imaging.PNG)
	c.Assert(resized.RelPermalink(), qt.Equals, "/a/sub/gohugoio2_hu0e1b9e4a4be4d6f86c7b37b9ccce3fbc_73886_101x101_resize_linear_2.png")
	c.Assert(resized.Width(), qt.Equals, 101)

	assertFileCache(c, fileCache, resized.RelPermalink(), 101, 101)
	publishedImageFilename := filepath.Clean(resized.RelPermalink())
	assertImageFile(c, image.spec.BaseFs.PublishFs, publishedImageFilename, 101, 101)
	c.Assert(image.spec.BaseFs.PublishFs.Remove(publishedImageFilename), qt.IsNil)

	// Cleare mem cache to simulate reading from the file cache.
	resized.spec.imageCache.clear()

	resizedAgain, err := image.Resize("101x101")
	c.Assert(err, qt.IsNil)
	c.Assert(resizedAgain.RelPermalink(), qt.Equals, "/a/sub/gohugoio2_hu0e1b9e4a4be4d6f86c7b37b9ccce3fbc_73886_101x101_resize_linear_2.png")
	c.Assert(resizedAgain.Width(), qt.Equals, 101)
	assertFileCache(c, fileCache, resizedAgain.RelPermalink(), 101, 101)
	assertImageFile(c, image.spec.BaseFs.PublishFs, publishedImageFilename, 101, 101)

}

func TestSVGImage(t *testing.T) {
	c := qt.New(t)
	spec := newTestResourceSpec(c)
	svg := fetchResourceForSpec(spec, c, "circle.svg")
	c.Assert(svg, qt.Not(qt.IsNil))
}

func TestSVGImageContent(t *testing.T) {
	c := qt.New(t)
	spec := newTestResourceSpec(c)
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
