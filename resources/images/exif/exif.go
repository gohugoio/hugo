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
	"io"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/bep/tmc"

	_exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

const (
	exifTimeLayout = "2006:01:02 15:04:05"
	// IFD ID used for efficiency. ID obtained from https://exiftool.org/TagNames/EXIF.html
	dateTimeTagId = 0x9003
)

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

	idfm *exifcommon.IfdMapping
	ti   *_exif.TagIndex
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
	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		return nil, err
	}
	ti := _exif.NewTagIndex()
	d := &Decoder{idfm: im, ti: ti}
	for _, opt := range options {
		if err := opt(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (d *Decoder) Decode(r io.Reader) (*Exif, error) {
	rawExif, err := _exif.SearchAndExtractExifWithReader(r)

	if err != nil {
		// No Exif data
		return nil, nil
	}

	_, index, err := _exif.Collect(d.idfm, d.ti, rawExif)
	if err != nil {
		return nil, err
	}

	var tm time.Time
	var lat, long float64

	rootIfd := index.RootIfd

	ifd, err := rootIfd.ChildWithIfdPath(exifcommon.IfdGpsInfoStandardIfdIdentity)
	if !d.noLatLong && err == nil {
		if gpsInfo, err := ifd.GpsInfo(); err == nil {
			lat = gpsInfo.Latitude.Decimal()
			long = gpsInfo.Longitude.Decimal()
		}
	}

	exifIfd, _ := rootIfd.ChildWithIfdPath(exifcommon.IfdExifStandardIfdIdentity)

	if results, err := exifIfd.FindTagWithId(dateTimeTagId); !d.noDate && err == nil && len(results) == 1 {
		dateTimeTagEntry := results[0]
		dateTimeRaw, _ := dateTimeTagEntry.Value()
		dateTimeString := dateTimeRaw.(string)

		tm, _ = time.Parse(exifTimeLayout, dateTimeString)
	}

	vals, err := extractTags(&index, d.includeFieldsRe, d.excludeFieldsrRe)
	if err != nil {
		// No Exif metadata
		return nil, nil
	}

	ex := &Exif{Lat: lat, Long: long, Date: tm, Tags: vals}

	return ex, nil
}

// Code borrowed from exif.DateTime and adjusted.
func tryParseDate(s string) (time.Time, error) {
	dateStr := strings.TrimRight(s, "\x00")
	// TODO(bep): look for timezone offset, GPS time, etc.
	timeZone := time.Local
	return time.ParseInLocation(exifTimeLayout, dateStr, timeZone)
}

func processRational(n int64, d int64) interface{} {
	rat := big.NewRat(n, d)
	if n != 1 {
		f, _ := rat.Float64()
		return f
	}
	return rat
}

func extractTagValue(ite *_exif.IfdTagEntry) (val interface{}, err error) {
	tagName := ite.TagName()
	tagType := ite.TagType()
	unitCount := ite.UnitCount()
	valueRaw, err := ite.Value()
	if err != nil {
		return nil, err
	}

	switch tagType {
	case exifcommon.TypeAscii:
		s, _ := valueRaw.(string)
		if strings.Contains(tagName, "DateTime") {
			if d, err := tryParseDate(s); err == nil {
				return d, nil
			}
		}
		return s, nil

	case exifcommon.TypeUndefined:
		return valueRaw, nil
	}

	var rv []interface{}

	for i := 0; i < int(unitCount); i++ {
		var val interface{}
		switch tagType {

		case exifcommon.TypeRational:
			vals, _ := valueRaw.([]exifcommon.Rational)
			r := vals[i]
			val = processRational(int64(r.Numerator), int64(r.Denominator))

		case exifcommon.TypeSignedRational:
			vals, _ := valueRaw.([]exifcommon.SignedRational)
			r := vals[i]
			val = processRational(int64(r.Numerator), int64(r.Denominator))

		case exifcommon.TypeShort:
			vals, _ := valueRaw.([]uint16)
			i := vals[i]
			val = int(i)

		case exifcommon.TypeByte:
			vals, _ := valueRaw.([]uint8)
			val = vals[i]
		}

		if val != nil {
			rv = append(rv, val)
		}
	}

	if unitCount == 1 && len(rv) == 1 {
		return rv[0], nil
	}

	return rv, nil
}

func extractTags(ifdIndex *_exif.IfdIndex, includeMatcher *regexp.Regexp, excludeMatcher *regexp.Regexp) (map[string]interface{}, error) {
	vals := make(map[string]interface{})
	err := ifdIndex.RootIfd.EnumerateTagsRecursively(func(ifd *_exif.Ifd, ite *_exif.IfdTagEntry) error {
		if ite != nil {
			name := ite.TagName()

			if excludeMatcher != nil && excludeMatcher.MatchString(name) {
				return nil
			}
			if includeMatcher != nil && !includeMatcher.MatchString(name) {
				return nil
			}
			val, err := extractTagValue(ite)
			if err != nil {
				return err
			}
			vals[name] = val
		}
		return nil
	})

	return vals, err
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
