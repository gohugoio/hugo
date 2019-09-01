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

package exif

import (
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gohugoio/hugo/htesting/hqt"
	"github.com/google/go-cmp/cmp"

	qt "github.com/frankban/quicktest"
)

func TestExif(t *testing.T) {
	c := qt.New(t)
	f, err := os.Open(filepath.FromSlash("../../testdata/sunset.jpg"))
	c.Assert(err, qt.IsNil)
	defer f.Close()

	d, err := NewDecoder(IncludeFields("Lens|Date"))
	c.Assert(err, qt.IsNil)
	x, err := d.Decode(f)
	c.Assert(err, qt.IsNil)
	c.Assert(x.Date.Format("2006-01-02"), qt.Equals, "2017-10-27")

	// Malaga: https://goo.gl/taazZy
	c.Assert(x.Lat, qt.Equals, float64(36.59744166666667))
	c.Assert(x.Long, qt.Equals, float64(-4.50846))

	v, found := x.Tags["LensModel"]
	c.Assert(found, qt.Equals, true)
	lensModel, ok := v.(string)
	c.Assert(ok, qt.Equals, true)
	c.Assert(lensModel, qt.Equals, "smc PENTAX-DA* 16-50mm F2.8 ED AL [IF] SDM")

	v, found = x.Tags["DateTime"]
	c.Assert(found, qt.Equals, true)
	c.Assert(v, hqt.IsSameType, time.Time{})

	// Verify that it survives a round-trip to JSON and back.
	data, err := json.Marshal(x)
	c.Assert(err, qt.IsNil)
	x2 := &Exif{}
	err = json.Unmarshal(data, x2)

	c.Assert(x2, eq, x)

}

func TestExifPNG(t *testing.T) {
	c := qt.New(t)

	f, err := os.Open(filepath.FromSlash("../../testdata/gohugoio.png"))
	c.Assert(err, qt.IsNil)
	defer f.Close()

	d, err := NewDecoder()
	c.Assert(err, qt.IsNil)
	_, err = d.Decode(f)
	c.Assert(err, qt.Not(qt.IsNil))
}

func BenchmarkDecodeExif(b *testing.B) {
	c := qt.New(b)
	f, err := os.Open(filepath.FromSlash("../../testdata/sunset.jpg"))
	c.Assert(err, qt.IsNil)
	defer f.Close()

	d, err := NewDecoder()
	c.Assert(err, qt.IsNil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = d.Decode(f)
		c.Assert(err, qt.IsNil)
		f.Seek(0, 0)
	}
}

var eq = qt.CmpEquals(
	cmp.Comparer(
		func(v1, v2 *big.Rat) bool {
			return v1.RatString() == v2.RatString()
		},
	),
	cmp.Comparer(func(v1, v2 time.Time) bool {
		return v1.Unix() == v2.Unix()
	}),
)
