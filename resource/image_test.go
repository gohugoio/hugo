// Copyright 2017-present The Hugo Authors. All rights reserved.
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

package resource

import (
	"fmt"
	"testing"

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

func TestImageTransform(t *testing.T) {

	assert := require.New(t)

	image := fetchSunset(assert)

	assert.Equal("/a/sunset.jpg", image.RelPermalink())
	assert.Equal("image", image.ResourceType())

	resized, err := image.Resize("300x200")
	assert.NoError(err)
	assert.True(image != resized)
	assert.True(image.genericResource != resized.genericResource)

	resized0x, err := image.Resize("x200")
	assert.NoError(err)
	assert.Equal(320, resized0x.Width())
	assert.Equal(200, resized0x.Height())
	assertFileCache(assert, image.spec.Fs, resized0x.RelPermalink(), 320, 200)

	resizedx0, err := image.Resize("200x")
	assert.NoError(err)
	assert.Equal(200, resizedx0.Width())
	assert.Equal(125, resizedx0.Height())
	assertFileCache(assert, image.spec.Fs, resizedx0.RelPermalink(), 200, 125)

	resizedAndRotated, err := image.Resize("x200 r90")
	assert.NoError(err)
	assert.Equal(125, resizedAndRotated.Width())
	assert.Equal(200, resizedAndRotated.Height())
	assertFileCache(assert, image.spec.Fs, resizedAndRotated.RelPermalink(), 125, 200)

	assert.Equal("/a/sunset_H59e56ffff1bc1d8d122b1403d34e039f_90587_300x200_resize_q75_box_center.jpg", resized.RelPermalink())
	assert.Equal(300, resized.Width())
	assert.Equal(200, resized.Height())

	fitted, err := resized.Fit("50x50")
	assert.NoError(err)
	assert.Equal("/a/sunset_H59e56ffff1bc1d8d122b1403d34e039f_90587_e71d3649737587d41fe50793bf366f6f.jpg", fitted.RelPermalink())
	assert.Equal(50, fitted.Width())
	assert.Equal(31, fitted.Height())

	// Check the MD5 key threshold
	fittedAgain, _ := fitted.Fit("10x20")
	fittedAgain, err = fittedAgain.Fit("10x20")
	assert.NoError(err)
	assert.Equal("/a/sunset_H59e56ffff1bc1d8d122b1403d34e039f_90587_8731035e4934a6e6e09cd10d6f04db93.jpg", fittedAgain.RelPermalink())
	assert.Equal(10, fittedAgain.Width())
	assert.Equal(6, fittedAgain.Height())

	filled, err := image.Fill("200x100 bottomLeft")
	assert.NoError(err)
	assert.Equal("/a/sunset_H59e56ffff1bc1d8d122b1403d34e039f_90587_200x100_fill_q75_box_bottomleft.jpg", filled.RelPermalink())
	assert.Equal(200, filled.Width())
	assert.Equal(100, filled.Height())
	assertFileCache(assert, image.spec.Fs, filled.RelPermalink(), 200, 100)

	// Check cache
	filledAgain, err := image.Fill("200x100 bottomLeft")
	assert.NoError(err)
	assert.True(filled == filledAgain)
	assertFileCache(assert, image.spec.Fs, filledAgain.RelPermalink(), 200, 100)

}

func TestDecodeImaging(t *testing.T) {
	assert := require.New(t)
	m := map[string]interface{}{
		"quality":        42,
		"resampleFilter": "NearestNeighbor",
	}

	imaging, err := decodeImaging(m)

	assert.NoError(err)
	assert.Equal(42, imaging.Quality)
	assert.Equal("nearestneighbor", imaging.ResampleFilter)
}
