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
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("exif failed: %v", r)
		}
	}()

	var tagInfos imagemeta.Tags
	handleTag := func(ti imagemeta.TagInfo) error {
		tagInfos.Add(ti)
		return nil
	}

	shouldInclude := func(ti imagemeta.TagInfo) bool {
		if ti.Source == imagemeta.EXIF {
			if !d.noDate {
				// We need the time tags to calculate the date.
				if isTimeTag(ti.Tag) {
					return true
				}
			}
			if !d.noLatLong {
				// We need to GPS tags to calculate the lat/long.
				if isGPSTag(ti.Tag) {
					return true
				}
			}

			if !strings.HasPrefix(ti.Namespace, "IFD0") {
				// Drop thumbnail tags.
				return false
			}
		}

		if d.shouldExclude(ti.Tag) {
			return false
		}

		return d.shouldInclude(ti.Tag)
	}

	var warnf func(string, ...any)
	if d.warnl != nil {
		// There should be very little warnings (fingers crossed!),
		// but this will typically be unrecognized formats.
		// To be able to possibly get rid of these warnings,
		// we need to know what images are causing them.
		warnf = func(format string, args ...any) {
			format = fmt.Sprintf("%q: %s: ", filename, format)
			d.warnl.Logf(format, args...)
		}
	}

	err = imagemeta.Decode(
		imagemeta.Options{
			R:               r.(io.ReadSeeker),
			ImageFormat:     format,
			ShouldHandleTag: shouldInclude,
			HandleTag:       handleTag,
			Sources:         imagemeta.EXIF, // For now. TODO(bep)
			Warnf:           warnf,
		},
	)

	var tm time.Time
	var lat, long float64

	if !d.noDate {
		tm, _ = tagInfos.GetDateTime()
	}

	if !d.noLatLong {
		lat, long, _ = tagInfos.GetLatLong()
	}

	tags := make(map[string]any)
	for k, v := range tagInfos.All() {
		if d.shouldExclude(k) {
			continue
		}
		if !d.shouldInclude(k) {
			continue
		}
		tags[k] = v.Value
	}

	ex = &ExifInfo{Lat: lat, Long: long, Date: tm, Tags: tags}

	return
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
