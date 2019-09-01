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
	"bytes"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/bep/tmc"

	_exif "github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

const exifTimeLayout = "2006:01:02 15:04:05"

type Exif struct {
	Lat  float64
	Long float64
	Date time.Time
	Tags Tags
}

type Decoder struct {
	includeFieldsRe  *regexp.Regexp
	excludeFieldsrRe *regexp.Regexp
	noDate           bool
	noLatLong        bool
}

func IncludeFields(expression string) func(*Decoder) error {
	return func(d *Decoder) error {
		re, err := compileRegexp(expression)
		if err != nil {
			return err
		}
		d.includeFieldsRe = re
		return nil
	}
}

func ExcludeFields(expression string) func(*Decoder) error {
	return func(d *Decoder) error {
		re, err := compileRegexp(expression)
		if err != nil {
			return err
		}
		d.excludeFieldsrRe = re
		return nil
	}
}

func WithLatLongDisabled(disabled bool) func(*Decoder) error {
	return func(d *Decoder) error {
		d.noLatLong = disabled
		return nil
	}
}

func WithDateDisabled(disabled bool) func(*Decoder) error {
	return func(d *Decoder) error {
		d.noDate = disabled
		return nil
	}
}

func compileRegexp(expression string) (*regexp.Regexp, error) {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return nil, nil
	}
	if !strings.HasPrefix(expression, "(") {
		// Make it case insensitive
		expression = "(?i)" + expression
	}

	return regexp.Compile(expression)

}

func NewDecoder(options ...func(*Decoder) error) (*Decoder, error) {
	d := &Decoder{}
	for _, opt := range options {
		if err := opt(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (d *Decoder) Decode(r io.Reader) (ex *Exif, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Exif failed: %v", r)
		}
	}()

	var x *_exif.Exif
	x, err = _exif.Decode(r)
	if err != nil {
		if err.Error() == "EOF" {

			// Found no Exif
			return nil, nil
		}
		return
	}

	var tm time.Time
	var lat, long float64

	if !d.noDate {
		tm, _ = x.DateTime()
	}

	if !d.noLatLong {
		lat, long, _ = x.LatLong()
	}

	walker := &exifWalker{x: x, vals: make(map[string]interface{}), includeMatcher: d.includeFieldsRe, excludeMatcher: d.excludeFieldsrRe}
	if err = x.Walk(walker); err != nil {
		return
	}

	ex = &Exif{Lat: lat, Long: long, Date: tm, Tags: walker.vals}

	return
}

func decodeTag(x *_exif.Exif, f _exif.FieldName, t *tiff.Tag) (interface{}, error) {
	switch t.Format() {
	case tiff.StringVal, tiff.UndefVal:
		s := nullString(t.Val)
		if strings.Contains(string(f), "DateTime") {
			if d, err := tryParseDate(x, s); err == nil {
				return d, nil
			}
		}
		return s, nil
	case tiff.OtherVal:
		return "unknown", nil
	}

	var rv []interface{}

	for i := 0; i < int(t.Count); i++ {
		switch t.Format() {
		case tiff.RatVal:
			n, d, _ := t.Rat2(i)
			rat := big.NewRat(n, d)
			if n == 1 {
				rv = append(rv, rat)
			} else {
				f, _ := rat.Float64()
				rv = append(rv, f)
			}

		case tiff.FloatVal:
			v, _ := t.Float(i)
			rv = append(rv, v)
		case tiff.IntVal:
			v, _ := t.Int(i)
			rv = append(rv, v)
		}
	}

	if t.Count == 1 {
		if len(rv) == 1 {
			return rv[0], nil
		}
	}

	return rv, nil

}

// Code borrowed from exif.DateTime and adjusted.
func tryParseDate(x *_exif.Exif, s string) (time.Time, error) {
	dateStr := strings.TrimRight(s, "\x00")
	// TODO(bep): look for timezone offset, GPS time, etc.
	timeZone := time.Local
	if tz, _ := x.TimeZone(); tz != nil {
		timeZone = tz
	}
	return time.ParseInLocation(exifTimeLayout, dateStr, timeZone)

}

type exifWalker struct {
	x              *_exif.Exif
	vals           map[string]interface{}
	includeMatcher *regexp.Regexp
	excludeMatcher *regexp.Regexp
}

func (e *exifWalker) Walk(f _exif.FieldName, tag *tiff.Tag) error {
	name := string(f)
	if e.excludeMatcher != nil && e.excludeMatcher.MatchString(name) {
		return nil
	}
	if e.includeMatcher != nil && !e.includeMatcher.MatchString(name) {
		return nil
	}
	val, err := decodeTag(e.x, f, tag)
	if err != nil {
		return err
	}
	e.vals[name] = val
	return nil
}

func nullString(in []byte) string {
	var rv bytes.Buffer
	for _, b := range in {
		if unicode.IsPrint(rune(b)) {
			rv.WriteByte(b)
		}
	}
	rvs := rv.String()
	if utf8.ValidString(rvs) {
		return rvs
	}

	return ""
}

var tcodec *tmc.Codec

func init() {
	var err error
	tcodec, err = tmc.New()
	if err != nil {
		panic(err)
	}
}

type Tags map[string]interface{}

func (v *Tags) UnmarshalJSON(b []byte) error {
	vv := make(map[string]interface{})
	if err := tcodec.Unmarshal(b, &vv); err != nil {
		return err
	}

	*v = vv

	return nil
}

func (v Tags) MarshalJSON() ([]byte, error) {
	return tcodec.Marshal(v)
}
