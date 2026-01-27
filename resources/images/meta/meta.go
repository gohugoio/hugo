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

package meta

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
	"github.com/gohugoio/hugo/common/predicate"
	"github.com/gohugoio/hugo/hugofs/hglob"
	"github.com/spf13/cast"
)

// MetaInfo holds the decoded metadata for an Image; what you get in $image.Meta.
// Note by default, only EXIF and IPTC data is decoded, unless configured otherwise.
// If you want a consolidated view of the different tag sections, use the merge template func, e.g. {{ $m := merge .Exif .IPTC .XMP }}.
type MetaInfo struct {
	// GPS latitude in degrees.
	Lat float64

	// GPS longitude in degrees.
	Long float64

	// Image creation date/time.
	Date time.Time

	// Orientation tag value.
	Orientation int

	Exif Tags
	IPTC Tags
	XMP  Tags
}

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
	// For ExifConfig (legacy, regexp-based)
	includeFieldsRe  *regexp.Regexp
	excludeFieldsrRe *regexp.Regexp

	// For MetaConfig (glob-based)
	fieldsPredicate predicate.P[string]

	noDate    bool
	noLatLong bool
	sources   imagemeta.Source
	warnl     logg.LevelLogger
}

func (d *Decoder) shouldIncludeField(s string) bool {
	// Glob-based predicate takes precedence (used by MetaConfig)
	if d.fieldsPredicate != nil {
		return d.fieldsPredicate(strings.ToLower(s))
	}
	// Fall back to regexp-based filtering (used by ExifConfig)
	if d.excludeFieldsrRe != nil && d.excludeFieldsrRe.MatchString(s) {
		return false
	}
	return d.includeFieldsRe == nil || d.includeFieldsRe.MatchString(s)
}

// IncludeFields sets a regexp for fields to include (legacy, for ExifConfig).
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

// ExcludeFields sets a regexp for fields to exclude (legacy, for ExifConfig).
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

// WithFields sets glob patterns for field filtering (for MetaConfig).
// Patterns starting with "! " are exclusions.
func WithFields(patterns []string) func(*Decoder) error {
	return func(d *Decoder) error {
		if len(patterns) == 0 {
			return nil
		}
		// Lowercase patterns for case-insensitive matching
		lowered := make([]string, len(patterns))
		for i, p := range patterns {
			if after, found := strings.CutPrefix(p, hglob.NegationPrefix); found {
				lowered[i] = hglob.NegationPrefix + strings.ToLower(after)
			} else {
				lowered[i] = strings.ToLower(p)
			}
		}
		p, err := predicate.NewStringPredicateFromGlobs(lowered, hglob.GetGlobDot)
		if err != nil {
			return err
		}
		d.fieldsPredicate = p
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

func WithSources(sources ...string) func(*Decoder) error {
	return func(d *Decoder) error {
		var s imagemeta.Source
		for _, source := range sources {
			switch source {
			case "exif":
				s |= imagemeta.EXIF
			case "iptc":
				s |= imagemeta.IPTC
			case "xmp":
				s |= imagemeta.XMP
			}
		}
		d.sources = s
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

	shouldHandleTag := func(ti imagemeta.TagInfo) bool {
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

		return d.shouldIncludeField(ti.Tag)
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

	_, err = imagemeta.Decode(
		imagemeta.Options{
			R:               r.(io.ReadSeeker),
			ImageFormat:     format,
			ShouldHandleTag: shouldHandleTag,
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
		if !d.shouldIncludeField(k) {
			continue
		}
		tags[k] = v.Value
	}

	ex = &ExifInfo{Lat: lat, Long: long, Date: tm, Tags: tags}

	return
}

// DecodeMeta decodes metadata from all sources (EXIF, IPTC, XMP).
// Filename is only used for logging.
func (d *Decoder) DecodeMeta(filename string, format imagemeta.ImageFormat, r io.Reader) (m *MetaInfo, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("metadata decoding failed: %v", r)
		}
	}()

	var tagInfos imagemeta.Tags
	handleTag := func(ti imagemeta.TagInfo) error {
		tagInfos.Add(ti)
		return nil
	}

	shouldHandleTag := func(ti imagemeta.TagInfo) bool {
		// For EXIF, only include tags from IFD0 (skip thumbnail data).
		if ti.Source == imagemeta.EXIF {
			if !strings.HasPrefix(ti.Namespace, "IFD0") {
				return false
			}
		}

		return d.shouldIncludeField(ti.Tag)
	}

	var warnf func(string, ...any)
	if d.warnl != nil {
		warnf = func(format string, args ...any) {
			format = fmt.Sprintf("%q: %s: ", filename, format)
			d.warnl.Logf(format, args...)
		}
	}

	sources := d.sources
	if sources.IsZero() {
		sources = imagemeta.EXIF | imagemeta.IPTC | imagemeta.XMP
	}

	_, err = imagemeta.Decode(
		imagemeta.Options{
			R:               r.(io.ReadSeeker),
			ImageFormat:     format,
			ShouldHandleTag: shouldHandleTag,
			HandleTag:       handleTag,
			Sources:         sources,
			Warnf:           warnf,
		},
	)
	if err != nil {
		return nil, err
	}

	tm, _ := tagInfos.GetDateTime()
	lat, long, _ := tagInfos.GetLatLong()

	exifTags := make(map[string]any)
	iptcTags := make(map[string]any)
	xmpTags := make(map[string]any)

	for k, v := range tagInfos.All() {
		if !d.shouldIncludeField(k) {
			continue
		}
		switch v.Source {
		case imagemeta.EXIF:
			exifTags[k] = v.Value
		case imagemeta.IPTC:
			iptcTags[k] = v.Value
		case imagemeta.XMP:
			xmpTags[k] = v.Value
		}
	}

	var orientation int
	if v, ok := exifTags["Orientation"]; ok {
		orientation = cast.ToInt(v)
	}

	m = &MetaInfo{
		Lat:         lat,
		Long:        long,
		Date:        tm,
		Orientation: orientation,
		Exif:        exifTags,
		IPTC:        iptcTags,
		XMP:         xmpTags,
	}

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
