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

	"github.com/disintegration/imaging"

	"sync"

	"github.com/stretchr/testify/require"
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

	assert := require.New(t)

	image := fetchSunset(assert)
	fileCache := image.spec.FileCaches.ImageCache().Fs

	assert.Equal("/a/sunset.jpg", image.RelPermalink())
	assert.Equal("image", image.ResourceType())

	resized, err := image.Resize("300x200")
	assert.NoError(err)
	assert.True(image != resized)
	assert.True(image.genericResource != resized.genericResource)
	assert.True(image.sourceFilename != resized.sourceFilename)

	resized0x, err := image.Resize("x200")
	assert.NoError(err)
	assert.Equal(320, resized0x.Width())
	assert.Equal(200, resized0x.Height())

	assertFileCache(assert, fileCache, resized0x.RelPermalink(), 320, 200)

	resizedx0, err := image.Resize("200x")
	assert.NoError(err)
	assert.Equal(200, resizedx0.Width())
	assert.Equal(125, resizedx0.Height())
	assertFileCache(assert, fileCache, resizedx0.RelPermalink(), 200, 125)

	resizedAndRotated, err := image.Resize("x200 r90")
	assert.NoError(err)
	assert.Equal(125, resizedAndRotated.Width())
	assert.Equal(200, resizedAndRotated.Height())
	assertFileCache(assert, fileCache, resizedAndRotated.RelPermalink(), 125, 200)

	assert.Equal("/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_300x200_resize_q68_linear.jpg", resized.RelPermalink())
	assert.Equal(300, resized.Width())
	assert.Equal(200, resized.Height())

	fitted, err := resized.Fit("50x50")
	assert.NoError(err)
	assert.Equal("/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_625708021e2bb281c9f1002f88e4753f.jpg", fitted.RelPermalink())
	assert.Equal(50, fitted.Width())
	assert.Equal(33, fitted.Height())

	// Check the MD5 key threshold
	fittedAgain, _ := fitted.Fit("10x20")
	fittedAgain, err = fittedAgain.Fit("10x20")
	assert.NoError(err)
	assert.Equal("/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_3f65ba24dc2b7fba0f56d7f104519157.jpg", fittedAgain.RelPermalink())
	assert.Equal(10, fittedAgain.Width())
	assert.Equal(6, fittedAgain.Height())

	filled, err := image.Fill("200x100 bottomLeft")
	assert.NoError(err)
	assert.Equal("/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x100_fill_q68_linear_bottomleft.jpg", filled.RelPermalink())
	assert.Equal(200, filled.Width())
	assert.Equal(100, filled.Height())
	assertFileCache(assert, fileCache, filled.RelPermalink(), 200, 100)

	smart, err := image.Fill("200x100 smart")
	assert.NoError(err)
	assert.Equal(fmt.Sprintf("/a/sunset_hu59e56ffff1bc1d8d122b1403d34e039f_90587_200x100_fill_q68_linear_smart%d.jpg", smartCropVersionNumber), smart.RelPermalink())
	assert.Equal(200, smart.Width())
	assert.Equal(100, smart.Height())
	assertFileCache(assert, fileCache, smart.RelPermalink(), 200, 100)

	// Check cache
	filledAgain, err := image.Fill("200x100 bottomLeft")
	assert.NoError(err)
	assert.True(filled == filledAgain)
	assert.True(filled.sourceFilename == filledAgain.sourceFilename)
	assertFileCache(assert, fileCache, filledAgain.RelPermalink(), 200, 100)

}

// https://github.com/gohugoio/hugo/issues/4261
func TestImageTransformLongFilename(t *testing.T) {
	assert := require.New(t)

	image := fetchImage(assert, "1234567890qwertyuiopasdfghjklzxcvbnm5to6eeeeee7via8eleph.jpg")
	assert.NotNil(image)

	resized, err := image.Resize("200x")
	assert.NoError(err)
	assert.NotNil(resized)
	assert.Equal(200, resized.Width())
	assert.Equal("/a/_hu59e56ffff1bc1d8d122b1403d34e039f_90587_65b757a6e14debeae720fe8831f0a9bc.jpg", resized.RelPermalink())
	resized, err = resized.Resize("100x")
	assert.NoError(err)
	assert.NotNil(resized)
	assert.Equal(100, resized.Width())
	assert.Equal("/a/_hu59e56ffff1bc1d8d122b1403d34e039f_90587_c876768085288f41211f768147ba2647.jpg", resized.RelPermalink())
}

func TestImageTransformConcurrent(t *testing.T) {

	var wg sync.WaitGroup

	assert := require.New(t)

	spec := newTestResourceOsFs(assert)

	image := fetchImageForSpec(spec, assert, "sunset.jpg")

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				img := image
				for k := 0; k < 2; k++ {
					r1, err := img.Resize(fmt.Sprintf("%dx", id-k))
					if err != nil {
						t.Fatal(err)
					}

					if r1.Width() != id-k {
						t.Fatalf("Width: %d:%d", r1.Width(), j)
					}

					r2, err := r1.Resize(fmt.Sprintf("%dx", id-k-1))
					if err != nil {
						t.Fatal(err)
					}

					_, err = r2.decodeSource()
					if err != nil {
						t.Fatal("Err decode:", err)
					}

					img = r1
				}
			}
		}(i + 20)
	}

	wg.Wait()
}

func TestDecodeImaging(t *testing.T) {
	assert := require.New(t)
	m := map[string]interface{}{
		"quality":        42,
		"resampleFilter": "NearestNeighbor",
		"anchor":         "topLeft",
	}

	imaging, err := decodeImaging(m)

	assert.NoError(err)
	assert.Equal(42, imaging.Quality)
	assert.Equal("nearestneighbor", imaging.ResampleFilter)
	assert.Equal("topleft", imaging.Anchor)

	m = map[string]interface{}{}

	imaging, err = decodeImaging(m)
	assert.NoError(err)
	assert.Equal(defaultJPEGQuality, imaging.Quality)
	assert.Equal("box", imaging.ResampleFilter)
	assert.Equal("smart", imaging.Anchor)

	_, err = decodeImaging(map[string]interface{}{
		"quality": 123,
	})
	assert.Error(err)

	_, err = decodeImaging(map[string]interface{}{
		"resampleFilter": "asdf",
	})
	assert.Error(err)

	_, err = decodeImaging(map[string]interface{}{
		"anchor": "asdf",
	})
	assert.Error(err)

	imaging, err = decodeImaging(map[string]interface{}{
		"anchor": "Smart",
	})
	assert.NoError(err)
	assert.Equal("smart", imaging.Anchor)

}

func TestImageWithMetadata(t *testing.T) {
	assert := require.New(t)

	image := fetchSunset(assert)

	var meta = []map[string]interface{}{
		{
			"title": "My Sunset",
			"name":  "Sunset #:counter",
			"src":   "*.jpg",
		},
	}

	assert.NoError(AssignMetadata(meta, image))
	assert.Equal("Sunset #1", image.Name())

	resized, err := image.Resize("200x")
	assert.NoError(err)
	assert.Equal("Sunset #1", resized.Name())

}

func TestImageResize8BitPNG(t *testing.T) {

	assert := require.New(t)

	image := fetchImage(assert, "gohugoio.png")

	assert.Equal(imaging.PNG, image.format)
	assert.Equal("/a/gohugoio.png", image.RelPermalink())
	assert.Equal("image", image.ResourceType())

	resized, err := image.Resize("800x")
	assert.NoError(err)
	assert.Equal(imaging.PNG, resized.format)
	assert.Equal("/a/gohugoio_hu0e1b9e4a4be4d6f86c7b37b9ccce3fbc_73886_800x0_resize_linear_2.png", resized.RelPermalink())
	assert.Equal(800, resized.Width())

}

func TestImageResizeInSubPath(t *testing.T) {

	assert := require.New(t)

	image := fetchImage(assert, "sub/gohugoio2.png")
	fileCache := image.spec.FileCaches.ImageCache().Fs

	assert.Equal(imaging.PNG, image.format)
	assert.Equal("/a/sub/gohugoio2.png", image.RelPermalink())
	assert.Equal("image", image.ResourceType())

	resized, err := image.Resize("101x101")
	assert.NoError(err)
	assert.Equal(imaging.PNG, resized.format)
	assert.Equal("/a/sub/gohugoio2_hu0e1b9e4a4be4d6f86c7b37b9ccce3fbc_73886_101x101_resize_linear_2.png", resized.RelPermalink())
	assert.Equal(101, resized.Width())

	assertFileCache(assert, fileCache, resized.RelPermalink(), 101, 101)
	publishedImageFilename := filepath.Clean(resized.RelPermalink())
	assertImageFile(assert, image.spec.BaseFs.PublishFs, publishedImageFilename, 101, 101)
	assert.NoError(image.spec.BaseFs.PublishFs.Remove(publishedImageFilename))

	// Cleare mem cache to simulate reading from the file cache.
	resized.spec.imageCache.clear()

	resizedAgain, err := image.Resize("101x101")
	assert.NoError(err)
	assert.Equal("/a/sub/gohugoio2_hu0e1b9e4a4be4d6f86c7b37b9ccce3fbc_73886_101x101_resize_linear_2.png", resizedAgain.RelPermalink())
	assert.Equal(101, resizedAgain.Width())
	assertFileCache(assert, fileCache, resizedAgain.RelPermalink(), 101, 101)
	assertImageFile(assert, image.spec.BaseFs.PublishFs, publishedImageFilename, 101, 101)

}

func TestSVGImage(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)
	svg := fetchResourceForSpec(spec, assert, "circle.svg")
	assert.NotNil(svg)
}

func TestSVGImageContent(t *testing.T) {
	assert := require.New(t)
	spec := newTestResourceSpec(assert)
	svg := fetchResourceForSpec(spec, assert, "circle.svg")
	assert.NotNil(svg)

	content, err := svg.Content()
	assert.NoError(err)
	assert.IsType("", content)
	assert.Contains(content.(string), `<svg height="100" width="100">`)
}

func BenchmarkResizeParallel(b *testing.B) {
	assert := require.New(b)
	img := fetchSunset(assert)

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
