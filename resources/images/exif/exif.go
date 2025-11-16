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
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bep/imagemeta"
	"github.com/bep/logg"
	"github.com/bep/tmc"
)

// ExifInfo holds the decoded Exif data for an Image.
type ExifInfo struct {
	// GPS latitude in degrees.
	Lat float64

	// GPS longitude in degrees.
	Long float64

	// Image creation date/time.
	Date time.Time

	// A collection of the available Exif tags for this Image.
	Tags Tags
}

type Decoder struct {
	includeFieldsRe  *regexp.Regexp
	excludeFieldsrRe *regexp.Regexp
	noDate           bool
	noLatLong        bool
	warnl            logg.LevelLogger
}

func (d *Decoder) shouldInclude(s string) bool {
	return (d.includeFieldsRe == nil || d.includeFieldsRe.MatchString(s))
}

func (d *Decoder) shouldExclude(s string) bool {
	return d.excludeFieldsrRe != nil && d.excludeFieldsrRe.MatchString(s)
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

func WithWarnLogger(warnl logg.LevelLogger) func(*Decoder) error {
	return func(d *Decoder) error {
		d.warnl = warnl
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

var (
	isTimeTag = func(s string) bool {
		return strings.Contains(s, "Time")
	}
	isGPSTag = func(s string) bool {
		return strings.HasPrefix(s, "GPS")
	}
)

// Filename is only used for logging.
func (d *Decoder) Decode(filename string, format imagemeta.ImageFormat, r io.Reader) (ex *ExifInfo, err error) {
	ex = &ExifInfo{}
	ex.Tags = make(Tags)

	metadata, err := imagemeta.Decode(r, format)
	if err != nil {
		return nil, err
	}

	for k, v := range metadata.Tags {
		if d.shouldInclude(k) && !d.shouldExclude(k) {
			ex.Tags[k] = v
		}
	}

	if !d.noLatLong {
		ex.Lat, ex.Long, err = d.extractGPS(metadata)
		if err != nil && d.warnl != nil {
			d.warnl.Warnf("failed to extract GPS: %v", err)
		}
	}

	if !d.noDate {
		ex.Date, err = d.extractDate(metadata)
		if err != nil && d.warnl != nil {
			d.warnl.Warnf("failed to extract date: %v", err)
		}
	}

	return ex, nil
}

func (d *Decoder) extractGPS(metadata *imagemeta.Metadata) (lat, long float64, err error) {
	latRef, ok := metadata.Tags["GPSLatitudeRef"].(string)
	if !ok {
		return 0, 0, fmt.Errorf("missing GPSLatitudeRef")
	}
	latVals, ok := metadata.Tags["GPSLatitude"].([]imagemeta.Rat[uint32])
	if !ok || len(latVals) != 3 {
		return 0, 0, fmt.Errorf("invalid GPSLatitude")
	}
	lat = float64(latVals[0].Num)/float64(latVals[0].Den) + float64(latVals[1].Num)/float64(latVals[1].Den)/60 + float64(latVals[2].Num)/float64(latVals[2].Den)/3600
	if latRef == "S" {
		lat = -lat
	}

	longRef, ok := metadata.Tags["GPSLongitudeRef"].(string)
	if !ok {
		return 0, 0, fmt.Errorf("missing GPSLongitudeRef")
	}
	longVals, ok := metadata.Tags["GPSLongitude"].([]imagemeta.Rat[uint32])
	if !ok || len(longVals) != 3 {
		return 0, 0, fmt.Errorf("invalid GPSLongitude")
	}
	long = float64(longVals[0].Num)/float64(longVals[0].Den) + float64(longVals[1].Num)/float64(longVals[1].Den)/60 + float64(longVals[2].Num)/float64(longVals[2].Den)/3600
	if longRef == "W" {
		long = -long
	}

	return lat, long, nil
}

func (d *Decoder) extractDate(metadata *imagemeta.Metadata) (t time.Time, err error) {
	dateStr, ok := metadata.Tags["DateTime"].(string)
	if !ok {
		return time.Time{}, fmt.Errorf("missing DateTime")
	}
	// Assume format "2006:01:02 15:04:05"
	return time.Parse("2006:01:02 15:04:05", dateStr)
}

var tcodec *tmc.Codec

func init() {
	newIntadapter := func(target any) tmc.Adapter {
		var bitSize int
		var isSigned bool

		switch target.(type) {
		case int:
			bitSize = 0
			isSigned = true
		case int8:
			bitSize = 8
			isSigned = true
		case int16:
			bitSize = 16
			isSigned = true
		case int32:
			bitSize = 32
			isSigned = true
		case int64:
			bitSize = 64
			isSigned = true
		case uint:
			bitSize = 0
		case uint8:
			bitSize = 8
		case uint16:
			bitSize = 16
		case uint32:
			bitSize = 32
		case uint64:
			bitSize = 64
		}

		intFromString := func(s string) (any, error) {
			if bitSize == 0 {
				return strconv.Atoi(s)
			}

			var v any
			var err error

			if isSigned {
				v, err = strconv.ParseInt(s, 10, bitSize)
			} else {
				v, err = strconv.ParseUint(s, 10, bitSize)
			}

			if err != nil {
				return 0, err
			}

			if isSigned {
				i := v.(int64)
				switch target.(type) {
				case int:
					return int(i), nil
				case int8:
					return int8(i), nil
				case int16:
					return int16(i), nil
				case int32:
					return int32(i), nil
				case int64:
					return i, nil
				}
			}

			i := v.(uint64)
			switch target.(type) {
			case uint:
				return uint(i), nil
			case uint8:
				return uint8(i), nil
			case uint16:
				return uint16(i), nil
			case uint32:
				return uint32(i), nil
			case uint64:
				return i, nil

			}

			return 0, fmt.Errorf("unsupported target type %T", target)
		}

		intToString := func(v any) (string, error) {
			return fmt.Sprintf("%d", v), nil
		}

		return tmc.NewAdapter(target, intFromString, intToString)
	}

	ru, _ := imagemeta.NewRat[uint32](1, 2)
	ri, _ := imagemeta.NewRat[int32](1, 2)
	tmcAdapters := []tmc.Adapter{
		tmc.NewAdapter(ru, nil, nil),
		tmc.NewAdapter(ri, nil, nil),
		newIntadapter(int(1)),
		newIntadapter(int8(1)),
		newIntadapter(int16(1)),
		newIntadapter(int32(1)),
		newIntadapter(int64(1)),
		newIntadapter(uint(1)),
		newIntadapter(uint8(1)),
		newIntadapter(uint16(1)),
		newIntadapter(uint32(1)),
		newIntadapter(uint64(1)),
	}

	tmcAdapters = append(tmc.DefaultTypeAdapters, tmcAdapters...)

	var err error
	tcodec, err = tmc.New(tmc.WithTypeAdapters(tmcAdapters))
	if err != nil {
		panic(err)
	}
}

// Tags is a map of EXIF tags.
type Tags map[string]any

// UnmarshalJSON is for internal use only.
func (v *Tags) UnmarshalJSON(b []byte) error {
	vv := make(map[string]any)
	if err := tcodec.Unmarshal(b, &vv); err != nil {
		return err
	}

	*v = vv

	return nil
}

// MarshalJSON is for internal use only.
func (v Tags) MarshalJSON() ([]byte, error) {
	return tcodec.Marshal(v)
}
