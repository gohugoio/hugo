// Copyright 2018 The Hugo Authors. All rights reserved.
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

package exif

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExif(t *testing.T) {

	assert := require.New(t)
	f, err := os.Open(filepath.FromSlash("testdata/sunset.jpg"))
	assert.NoError(err)
	defer f.Close()

	d, err := NewDecoder(IncludeFields("Lens|Date"))
	assert.NoError(err)
	x, err := d.Decode(f)
	assert.NoError(err)
	assert.Equal("2017-10-27", x.Date.Format("2006-01-02"))

	// Malaga: https://goo.gl/taazZy
	assert.Equal(float64(36.59744166666667), x.Lat)
	assert.Equal(float64(-4.50846), x.Long)

	v, found := x.Values["LensModel"]
	assert.True(found)
	lensModel, ok := v.(string)
	assert.True(ok)
	assert.Equal("smc PENTAX-DA* 16-50mm F2.8 ED AL [IF] SDM", lensModel)

	v, found = x.Values["DateTime"]
	assert.True(found)
	assert.IsType(time.Time{}, v)

}

func TestExifPNG(t *testing.T) {

	assert := require.New(t)
	f, err := os.Open(filepath.FromSlash("testdata/sunsetpng.png"))
	assert.NoError(err)
	defer f.Close()

	d, err := NewDecoder()
	assert.NoError(err)
	_, err = d.Decode(f)
	assert.Error(err)
}

func BenchmarkDecodeExif(b *testing.B) {
	assert := require.New(b)
	f, err := os.Open(filepath.FromSlash("testdata/sunset.jpg"))
	assert.NoError(err)
	defer f.Close()

	d, err := NewDecoder()
	assert.NoError(err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = d.Decode(f)
		assert.NoError(err)
		f.Seek(0, 0)
	}
}
