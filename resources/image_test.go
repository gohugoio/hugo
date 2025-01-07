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

package resources_test

import (
	"context"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/bep/imagemeta"

	"github.com/gohugoio/hugo/common/paths"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/images"
	"github.com/google/go-cmp/cmp"

	"github.com/gohugoio/hugo/htesting/hqt"

	qt "github.com/frankban/quicktest"
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
	// cmp.Comparer(func(p1, p2 *genericResource) bool { return p1 == p2 }),
	cmp.Comparer(func(m1, m2 media.Type) bool {
		return m1.Type == m2.Type
	}),
	cmp.Comparer(
		func(v1, v2 imagemeta.Rat[uint32]) bool {
			return v1.String() == v2.String()
		},
	),
	cmp.Comparer(
		func(v1, v2 imagemeta.Rat[int32]) bool {
			return v1.String() == v2.String()
		},
	),
	cmp.Comparer(func(v1, v2 time.Time) bool {
		return v1.Unix() == v2.Unix()
	}),
)

func TestImageTransformBasic(t *testing.T) {
	c := qt.New(t)

	_, image := fetchSunset(c)

	assertWidthHeight := func(img images.ImageResource, w, h int) {
		assertWidthHeight(c, img, w, h)
	}

	gotColors, err := image.Colors()
	c.Assert(err, qt.IsNil)
	expectedColors := images.HexStringsToColors("#2d2f33", "#a49e93", "#d39e59", "#a76936", "#737a84", "#7c838b")
	c.Assert(len(gotColors), qt.Equals, len(expectedColors))
	for i := range gotColors {
		c1, c2 := gotColors[i], expectedColors[i]
		c.Assert(c1.ColorHex(), qt.Equals, c2.ColorHex())
		c.Assert(c1.ColorGo(), qt.DeepEquals, c2.ColorGo())
		c.Assert(c1.Luminance(), qt.Equals, c2.Luminance())
	}

	c.Assert(image.RelPermalink(), qt.Equals, "/a/sunset.jpg")
	c.Assert(image.ResourceType(), qt.Equals, "image")
	assertWidthHeight(image, 900, 562)

	resized, err := image.Resize("300x200")
	c.Assert(err, qt.IsNil)
	c.Assert(image != resized, qt.Equals, true)
	assertWidthHeight(resized, 300, 200)
	assertWidthHeight(image, 900, 562)

	resized0x, err := image.Resize("x200")
	c.Assert(err, qt.IsNil)
	assertWidthHeight(resized0x, 320, 200)

	resizedx0, err := image.Resize("200x")
	c.Assert(err, qt.IsNil)
	assertWidthHeight(resizedx0, 200, 125)

	resizedAndRotated, err := image.Resize("x200 r90")
	c.Assert(err, qt.IsNil)
	assertWidthHeight(resizedAndRotated, 125, 200)

	assertWidthHeight(resized, 300, 200)
	c.Assert(resized.RelPermalink(), qt.Equals, "/a/sunset_hu2082030801149749592.jpg")

	fitted, err := resized.Fit("50x50")
	c.Assert(err, qt.IsNil)
	c.Assert(fitted.RelPermalink(), qt.Equals, "/a/sunset_hu16263619592447877226.jpg")
	assertWidthHeight(fitted, 50, 33)

	// Check the MD5 key threshold
	fittedAgain, _ := fitted.Fit("10x20")
	fittedAgain, err = fittedAgain.Fit("10x20")
	c.Assert(err, qt.IsNil)
	c.Assert(fittedAgain.RelPermalink(), qt.Equals, "/a/sunset_hu847809310637164306.jpg")
	assertWidthHeight(fittedAgain, 10, 7)

	filled, err := image.Fill("200x100 bottomLeft")
	c.Assert(err, qt.IsNil)
	c.Assert(filled.RelPermalink(), qt.Equals, "/a/sunset_hu18289448341423092707.jpg")
	assertWidthHeight(filled, 200, 100)

	smart, err := image.Fill("200x100 smart")
	c.Assert(err, qt.IsNil)
	c.Assert(smart.RelPermalink(), qt.Equals, "/a/sunset_hu11649371610839769766.jpg")
	assertWidthHeight(smart, 200, 100)

	// Check cache
	filledAgain, err := image.Fill("200x100 bottomLeft")
	c.Assert(err, qt.IsNil)
	c.Assert(filled, qt.Equals, filledAgain)

	cropped, err := image.Crop("300x300 topRight")
	c.Assert(err, qt.IsNil)
	c.Assert(cropped.RelPermalink(), qt.Equals, "/a/sunset_hu2242042514052853140.jpg")
	assertWidthHeight(cropped, 300, 300)

	smartcropped, err := image.Crop("200x200 smart")
	c.Assert(err, qt.IsNil)
	c.Assert(smartcropped.RelPermalink(), qt.Equals, "/a/sunset_hu12983255101170993571.jpg")
	assertWidthHeight(smartcropped, 200, 200)

	// Check cache
	croppedAgain, err := image.Crop("300x300 topRight")
	c.Assert(err, qt.IsNil)
	c.Assert(cropped, qt.Equals, croppedAgain)
}

func TestImageProcess(t *testing.T) {
	c := qt.New(t)
	_, img := fetchSunset(c)
	resized, err := img.Process("resiZe 300x200")
	c.Assert(err, qt.IsNil)
	assertWidthHeight(c, resized, 300, 200)
	rotated, err := resized.Process("R90")
	c.Assert(err, qt.IsNil)
	assertWidthHeight(c, rotated, 200, 300)
	converted, err := img.Process("png")
	c.Assert(err, qt.IsNil)
	c.Assert(converted.MediaType().Type, qt.Equals, "image/png")

	checkProcessVsMethod := func(action, spec string) {
		var expect images.ImageResource
		var err error
		switch action {
		case images.ActionCrop:
			expect, err = img.Crop(spec)
		case images.ActionFill:
			expect, err = img.Fill(spec)
		case images.ActionFit:
			expect, err = img.Fit(spec)
		case images.ActionResize:
			expect, err = img.Resize(spec)
		}
		c.Assert(err, qt.IsNil)
		got, err := img.Process(spec + " " + action)
		c.Assert(err, qt.IsNil)
		assertWidthHeight(c, got, expect.Width(), expect.Height())
		c.Assert(got.MediaType(), qt.Equals, expect.MediaType())
	}

	checkProcessVsMethod(images.ActionCrop, "300x200 topleFt")
	checkProcessVsMethod(images.ActionFill, "300x200 topleft")
	checkProcessVsMethod(images.ActionFit, "300x200 png")
	checkProcessVsMethod(images.ActionResize, "300x R90")
}

func TestImageTransformFormat(t *testing.T) {
	c := qt.New(t)

	_, image := fetchSunset(c)

	assertExtWidthHeight := func(img images.ImageResource, ext string, w, h int) {
		c.Helper()
		c.Assert(img, qt.Not(qt.IsNil))
		c.Assert(paths.Ext(img.RelPermalink()), qt.Equals, ext)
		c.Assert(img.Width(), qt.Equals, w)
		c.Assert(img.Height(), qt.Equals, h)
	}

	c.Assert(image.RelPermalink(), qt.Equals, "/a/sunset.jpg")
	c.Assert(image.ResourceType(), qt.Equals, "image")
	assertExtWidthHeight(image, ".jpg", 900, 562)

	imagePng, err := image.Resize("450x png")
	c.Assert(err, qt.IsNil)
	c.Assert(imagePng.RelPermalink(), qt.Equals, "/a/sunset_hu11737890885216583918.png")
	c.Assert(imagePng.ResourceType(), qt.Equals, "image")
	assertExtWidthHeight(imagePng, ".png", 450, 281)
	c.Assert(imagePng.Name(), qt.Equals, "sunset.jpg")
	c.Assert(imagePng.MediaType().String(), qt.Equals, "image/png")

	imageGif, err := image.Resize("225x gif")
	c.Assert(err, qt.IsNil)
	c.Assert(imageGif.RelPermalink(), qt.Equals, "/a/sunset_hu1431827106749674475.gif")
	c.Assert(imageGif.ResourceType(), qt.Equals, "image")
	assertExtWidthHeight(imageGif, ".gif", 225, 141)
	c.Assert(imageGif.Name(), qt.Equals, "sunset.jpg")
	c.Assert(imageGif.MediaType().String(), qt.Equals, "image/gif")
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

			check1 := func(img images.ImageResource) {
				resizedLink := "/a/sunset_hu7919355342577096259.jpg"
				c.Assert(img.RelPermalink(), qt.Equals, resizedLink)
				assertImageFile(c, spec.PublishFs, resizedLink, 100, 50)
			}

			check2 := func(img images.ImageResource) {
				c.Assert(img.RelPermalink(), qt.Equals, "/a/sunset.jpg")
				assertImageFile(c, spec.PublishFs, "a/sunset.jpg", 900, 562)
			}

			original := fetchImageForSpec(spec, c, "sunset.jpg")
			c.Assert(original, qt.Not(qt.IsNil))

			if checkOriginalFirst {
				check2(original)
			}

			resized, err := original.Resize("100x50")
			c.Assert(err, qt.IsNil)

			check1(resized)

			if !checkOriginalFirst {
				check2(original)
			}
		})
	}
}

func TestImageBugs(t *testing.T) {
	c := qt.New(t)

	// Issue #4261
	c.Run("Transform long filename", func(c *qt.C) {
		_, image := fetchImage(c, "1234567890qwertyuiopasdfghjklzxcvbnm5to6eeeeee7via8eleph.jpg")
		c.Assert(image, qt.Not(qt.IsNil))

		resized, err := image.Resize("200x")
		c.Assert(err, qt.IsNil)
		c.Assert(resized, qt.Not(qt.IsNil))
		c.Assert(resized.Width(), qt.Equals, 200)
		c.Assert(resized.RelPermalink(), qt.Equals, "/a/1234567890qwertyuiopasdfghjklzxcvbnm5to6eeeeee7via8eleph_hu9514381480012510326.jpg")
		resized, err = resized.Resize("100x")
		c.Assert(err, qt.IsNil)
		c.Assert(resized, qt.Not(qt.IsNil))
		c.Assert(resized.Width(), qt.Equals, 100)
		c.Assert(resized.RelPermalink(), qt.Equals, "/a/1234567890qwertyuiopasdfghjklzxcvbnm5to6eeeeee7via8eleph_hu1776700126481066216.jpg")
	})

	// Issue #6137
	c.Run("Transform upper case extension", func(c *qt.C) {
		_, image := fetchImage(c, "sunrise.JPG")

		resized, err := image.Resize("200x")
		c.Assert(err, qt.IsNil)
		c.Assert(resized, qt.Not(qt.IsNil))
		c.Assert(resized.Width(), qt.Equals, 200)
	})

	// Issue #7955
	c.Run("Fill with smartcrop", func(c *qt.C) {
		_, sunset := fetchImage(c, "sunset.jpg")

		for _, test := range []struct {
			originalDimensions string
			targetWH           int
		}{
			{"408x403", 400},
			{"425x403", 400},
			{"459x429", 400},
			{"476x442", 400},
			{"544x403", 400},
			{"476x468", 400},
			{"578x585", 550},
			{"578x598", 550},
		} {
			c.Run(test.originalDimensions, func(c *qt.C) {
				image, err := sunset.Resize(test.originalDimensions)
				c.Assert(err, qt.IsNil)
				resized, err := image.Fill(fmt.Sprintf("%dx%d smart", test.targetWH, test.targetWH))
				c.Assert(err, qt.IsNil)
				c.Assert(resized, qt.Not(qt.IsNil))
				c.Assert(resized.Width(), qt.Equals, test.targetWH)
				c.Assert(resized.Height(), qt.Equals, test.targetWH)
			})
		}
	})
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

func TestImageResize8BitPNG(t *testing.T) {
	c := qt.New(t)

	_, image := fetchImage(c, "gohugoio.png")

	c.Assert(image.MediaType().Type, qt.Equals, "image/png")
	c.Assert(image.RelPermalink(), qt.Equals, "/a/gohugoio.png")
	c.Assert(image.ResourceType(), qt.Equals, "image")
	c.Assert(image.Exif(), qt.IsNotNil)

	resized, err := image.Resize("800x")
	c.Assert(err, qt.IsNil)
	c.Assert(resized.MediaType().Type, qt.Equals, "image/png")
	c.Assert(resized.RelPermalink(), qt.Equals, "/a/gohugoio_hu8582372628235034388.png")
	c.Assert(resized.Width(), qt.Equals, 800)
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

	content, err := svg.Content(context.Background())
	c.Assert(err, qt.IsNil)
	c.Assert(content, hqt.IsSameType, "")
	c.Assert(content.(string), qt.Contains, `<svg height="100" width="100">`)
}

func TestImageExif(t *testing.T) {
	c := qt.New(t)
	fs := afero.NewMemMapFs()
	spec := newTestResourceSpec(specDescriptor{fs: fs, c: c})
	image := fetchResourceForSpec(spec, c, "sunset.jpg").(images.ImageResource)

	getAndCheckExif := func(c *qt.C, image images.ImageResource) {
		x := image.Exif()
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
		x2 := resized.Exif()

		c.Assert(x2, eq, x)
	}

	getAndCheckExif(c, image)
	image = fetchResourceForSpec(spec, c, "sunset.jpg").(images.ImageResource)
	// This will read from file cache.
	getAndCheckExif(c, image)
}

func TestImageColorsLuminance(t *testing.T) {
	c := qt.New(t)

	_, image := fetchSunset(c)
	c.Assert(image, qt.Not(qt.IsNil))
	colors, err := image.Colors()
	c.Assert(err, qt.IsNil)
	c.Assert(len(colors), qt.Equals, 6)
	var prevLuminance float64
	for i, color := range colors {
		luminance := color.Luminance()
		c.Assert(err, qt.IsNil)
		c.Assert(luminance > 0, qt.IsTrue)
		c.Assert(luminance, qt.Not(qt.Equals), prevLuminance, qt.Commentf("i=%d", i))
		prevLuminance = luminance
	}
}

func BenchmarkImageExif(b *testing.B) {
	getImages := func(c *qt.C, b *testing.B, fs afero.Fs) []images.ImageResource {
		spec := newTestResourceSpec(specDescriptor{fs: fs, c: c})
		imgs := make([]images.ImageResource, b.N)
		for i := 0; i < b.N; i++ {
			imgs[i] = fetchResourceForSpec(spec, c, "sunset.jpg", strconv.Itoa(i)).(images.ImageResource)
		}
		return imgs
	}

	getAndCheckExif := func(c *qt.C, image images.ImageResource) {
		x := image.Exif()
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

func BenchmarkResizeParallel(b *testing.B) {
	c := qt.New(b)
	_, img := fetchSunset(c)

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

func assertWidthHeight(c *qt.C, img images.ImageResource, w, h int) {
	c.Helper()
	c.Assert(img, qt.Not(qt.IsNil))
	c.Assert(img.Width(), qt.Equals, w)
	c.Assert(img.Height(), qt.Equals, h)
}
