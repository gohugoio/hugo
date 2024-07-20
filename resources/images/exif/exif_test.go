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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bep/imagemeta"
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
	x, err := d.Decode("", imagemeta.JPEG, f)
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

	v, found = x.Tags["ModifyDate"]
	c.Assert(found, qt.Equals, true)
	c.Assert(v, qt.Equals, "2017:11:23 09:56:54")

	// Verify that it survives a round-trip to JSON and back.
	data, err := json.Marshal(x)
	c.Assert(err, qt.IsNil)
	x2 := &ExifInfo{}
	err = json.Unmarshal(data, x2)
	c.Assert(err, qt.IsNil)

	c.Assert(x2, eq, x)
}

func TestExifPNG(t *testing.T) {
	c := qt.New(t)

	f, err := os.Open(filepath.FromSlash("../../testdata/gohugoio.png"))
	c.Assert(err, qt.IsNil)
	defer f.Close()

	d, err := NewDecoder()
	c.Assert(err, qt.IsNil)
	_, err = d.Decode("", imagemeta.PNG, f)
	c.Assert(err, qt.IsNil)
}

func TestIssue8079(t *testing.T) {
	c := qt.New(t)

	f, err := os.Open(filepath.FromSlash("../../testdata/iss8079.jpg"))
	c.Assert(err, qt.IsNil)
	defer f.Close()

	d, err := NewDecoder()
	c.Assert(err, qt.IsNil)
	x, err := d.Decode("", imagemeta.JPEG, f)
	c.Assert(err, qt.IsNil)
	c.Assert(x.Tags["ImageDescription"], qt.Equals, "Città del Vaticano #nanoblock #vatican #vaticancity")
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
		_, err = d.Decode("", imagemeta.JPEG, f)
		c.Assert(err, qt.IsNil)
		f.Seek(0, 0)
	}
}

var eq = qt.CmpEquals(
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

func TestIssue10738(t *testing.T) {
	c := qt.New(t)

	testFunc := func(c *qt.C, path, include string) any {
		c.Helper()
		f, err := os.Open(filepath.FromSlash(path))
		c.Assert(err, qt.IsNil)
		defer f.Close()

		d, err := NewDecoder(IncludeFields(include))
		c.Assert(err, qt.IsNil)
		x, err := d.Decode("", imagemeta.JPEG, f)
		c.Assert(err, qt.IsNil)

		// Verify that it survives a round-trip to JSON and back.
		data, err := json.Marshal(x)
		c.Assert(err, qt.IsNil)
		x2 := &ExifInfo{}
		err = json.Unmarshal(data, x2)
		c.Assert(err, qt.IsNil)

		c.Assert(x2, eq, x)

		v, found := x.Tags["ExposureTime"]
		c.Assert(found, qt.Equals, true)
		return v
	}

	type args struct {
		path    string // imagePath
		include string // includeFields
	}

	type want struct {
		vN int64 // numerator
		vD int64 // denominator
	}

	type testCase struct {
		name string
		args args
		want want
	}

	tests := []testCase{
		{
			"canon_cr2_fraction", args{
				path:    "../../testdata/issue10738/canon_cr2_fraction.jpg",
				include: "Lens|Date|ExposureTime",
			}, want{
				1,
				500,
			},
		},
		{
			"canon_cr2_integer", args{
				path:    "../../testdata/issue10738/canon_cr2_integer.jpg",
				include: "Lens|Date|ExposureTime",
			}, want{
				10,
				1,
			},
		},
		{
			"dji_dng_fraction", args{
				path:    "../../testdata/issue10738/dji_dng_fraction.jpg",
				include: "Lens|Date|ExposureTime",
			}, want{
				1,
				4000,
			},
		},
		{
			"fuji_raf_fraction", args{
				path:    "../../testdata/issue10738/fuji_raf_fraction.jpg",
				include: "Lens|Date|ExposureTime",
			}, want{
				1,
				250,
			},
		},
		{
			"fuji_raf_integer", args{
				path:    "../../testdata/issue10738/fuji_raf_integer.jpg",
				include: "Lens|Date|ExposureTime",
			}, want{
				1,
				1,
			},
		},
		{
			"leica_dng_fraction", args{
				path:    "../../testdata/issue10738/leica_dng_fraction.jpg",
				include: "Lens|Date|ExposureTime",
			}, want{
				1,
				100,
			},
		},
		{
			"lumix_rw2_fraction", args{
				path:    "../../testdata/issue10738/lumix_rw2_fraction.jpg",
				include: "Lens|Date|ExposureTime",
			}, want{
				1,
				400,
			},
		},
		{
			"nikon_nef_d5600", args{
				path:    "../../testdata/issue10738/nikon_nef_d5600.jpg",
				include: "Lens|Date|ExposureTime",
			}, want{
				1,
				1000,
			},
		},
		{
			"nikon_nef_fraction", args{
				path:    "../../testdata/issue10738/nikon_nef_fraction.jpg",
				include: "Lens|Date|ExposureTime",
			}, want{
				1,
				640,
			},
		},
		{
			"nikon_nef_integer", args{
				path:    "../../testdata/issue10738/nikon_nef_integer.jpg",
				include: "Lens|Date|ExposureTime",
			}, want{
				30,
				1,
			},
		},
		{
			"nikon_nef_fraction_2", args{
				path:    "../../testdata/issue10738/nikon_nef_fraction_2.jpg",
				include: "Lens|Date|ExposureTime",
			}, want{
				1,
				6400,
			},
		},
		{
			"sony_arw_fraction", args{
				path:    "../../testdata/issue10738/sony_arw_fraction.jpg",
				include: "Lens|Date|ExposureTime",
			}, want{
				1,
				160,
			},
		},
		{
			"sony_arw_integer", args{
				path:    "../../testdata/issue10738/sony_arw_integer.jpg",
				include: "Lens|Date|ExposureTime",
			}, want{
				4,
				1,
			},
		},
	}

	for _, tt := range tests {
		c.Run(tt.name, func(c *qt.C) {
			got := testFunc(c, tt.args.path, tt.args.include)
			switch v := got.(type) {
			case float64:
				c.Assert(v, qt.Equals, float64(tt.want.vN))
			case imagemeta.Rat[uint32]:
				r, err := imagemeta.NewRat[uint32](uint32(tt.want.vN), uint32(tt.want.vD))
				c.Assert(err, qt.IsNil)
				c.Assert(v, eq, r)
			default:
				c.Fatalf("unexpected type: %T", got)
			}
		})
	}
}
