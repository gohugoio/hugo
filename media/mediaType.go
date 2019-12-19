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

package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/mitchellh/mapstructure"
)

const (
	defaultDelimiter = "."
)

// Type (also known as MIME type and content type) is a two-part identifier for
// file formats and format contents transmitted on the Internet.
// For Hugo's use case, we use the top-level type name / subtype name + suffix.
// One example would be application/svg+xml
// If suffix is not provided, the sub type will be used.
// See // https://en.wikipedia.org/wiki/Media_type
type Type struct {
	MainType string `json:"mainType"` // i.e. text
	SubType  string `json:"subType"`  // i.e. html

	// This is the optional suffix after the "+" in the MIME type,
	//  e.g. "xml" in "applicatiion/rss+xml".
	mimeSuffix string

	Delimiter string `json:"delimiter"` // e.g. "."

	// TODO(bep) make this a string to make it hashable + method
	Suffixes []string `json:"suffixes"`

	// Set when doing lookup by suffix.
	fileSuffix string
}

// FromStringAndExt is same as FromString, but adds the file extension to the type.
func FromStringAndExt(t, ext string) (Type, error) {
	tp, err := fromString(t)
	if err != nil {
		return tp, err
	}
	tp.Suffixes = []string{strings.TrimPrefix(ext, ".")}
	return tp, nil
}

// FromString creates a new Type given a type string on the form MainType/SubType and
// an optional suffix, e.g. "text/html" or "text/html+html".
func fromString(t string) (Type, error) {
	t = strings.ToLower(t)
	parts := strings.Split(t, "/")
	if len(parts) != 2 {
		return Type{}, fmt.Errorf("cannot parse %q as a media type", t)
	}
	mainType := parts[0]
	subParts := strings.Split(parts[1], "+")

	subType := strings.Split(subParts[0], ";")[0]

	var suffix string

	if len(subParts) > 1 {
		suffix = subParts[1]
	}

	return Type{MainType: mainType, SubType: subType, mimeSuffix: suffix}, nil
}

// Type returns a string representing the main- and sub-type of a media type, e.g. "text/css".
// A suffix identifier will be appended after a "+" if set, e.g. "image/svg+xml".
// Hugo will register a set of default media types.
// These can be overridden by the user in the configuration,
// by defining a media type with the same Type.
func (m Type) Type() string {
	// Examples are
	// image/svg+xml
	// text/css
	if m.mimeSuffix != "" {
		return m.MainType + "/" + m.SubType + "+" + m.mimeSuffix
	}
	return m.MainType + "/" + m.SubType
}

func (m Type) String() string {
	return m.Type()
}

// FullSuffix returns the file suffix with any delimiter prepended.
func (m Type) FullSuffix() string {
	return m.Delimiter + m.Suffix()
}

// Suffix returns the file suffix without any delmiter prepended.
func (m Type) Suffix() string {
	if m.fileSuffix != "" {
		return m.fileSuffix
	}
	if len(m.Suffixes) > 0 {
		return m.Suffixes[0]
	}
	// There are MIME types without file suffixes.
	return ""
}

// Definitions from https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types etc.
// Note that from Hugo 0.44 we only set Suffix if it is part of the MIME type.
var (
	CalendarType   = Type{MainType: "text", SubType: "calendar", Suffixes: []string{"ics"}, Delimiter: defaultDelimiter}
	CSSType        = Type{MainType: "text", SubType: "css", Suffixes: []string{"css"}, Delimiter: defaultDelimiter}
	SCSSType       = Type{MainType: "text", SubType: "x-scss", Suffixes: []string{"scss"}, Delimiter: defaultDelimiter}
	SASSType       = Type{MainType: "text", SubType: "x-sass", Suffixes: []string{"sass"}, Delimiter: defaultDelimiter}
	CSVType        = Type{MainType: "text", SubType: "csv", Suffixes: []string{"csv"}, Delimiter: defaultDelimiter}
	HTMLType       = Type{MainType: "text", SubType: "html", Suffixes: []string{"html"}, Delimiter: defaultDelimiter}
	JavascriptType = Type{MainType: "application", SubType: "javascript", Suffixes: []string{"js"}, Delimiter: defaultDelimiter}
	JSONType       = Type{MainType: "application", SubType: "json", Suffixes: []string{"json"}, Delimiter: defaultDelimiter}
	RSSType        = Type{MainType: "application", SubType: "rss", mimeSuffix: "xml", Suffixes: []string{"xml"}, Delimiter: defaultDelimiter}
	XMLType        = Type{MainType: "application", SubType: "xml", Suffixes: []string{"xml"}, Delimiter: defaultDelimiter}
	SVGType        = Type{MainType: "image", SubType: "svg", mimeSuffix: "xml", Suffixes: []string{"svg"}, Delimiter: defaultDelimiter}
	TextType       = Type{MainType: "text", SubType: "plain", Suffixes: []string{"txt"}, Delimiter: defaultDelimiter}
	TOMLType       = Type{MainType: "application", SubType: "toml", Suffixes: []string{"toml"}, Delimiter: defaultDelimiter}
	YAMLType       = Type{MainType: "application", SubType: "yaml", Suffixes: []string{"yaml", "yml"}, Delimiter: defaultDelimiter}

	// Common image types
	PNGType  = Type{MainType: "image", SubType: "png", Suffixes: []string{"png"}, Delimiter: defaultDelimiter}
	JPEGType = Type{MainType: "image", SubType: "jpeg", Suffixes: []string{"jpg", "jpeg"}, Delimiter: defaultDelimiter}
	GIFType  = Type{MainType: "image", SubType: "gif", Suffixes: []string{"gif"}, Delimiter: defaultDelimiter}
	TIFFType = Type{MainType: "image", SubType: "tiff", Suffixes: []string{"tif", "tiff"}, Delimiter: defaultDelimiter}
	BMPType  = Type{MainType: "image", SubType: "bmp", Suffixes: []string{"bmp"}, Delimiter: defaultDelimiter}

	// Common video types
	AVIType  = Type{MainType: "video", SubType: "x-msvideo", Suffixes: []string{"avi"}, Delimiter: defaultDelimiter}
	MPEGType = Type{MainType: "video", SubType: "mpeg", Suffixes: []string{"mpg", "mpeg"}, Delimiter: defaultDelimiter}
	MP4Type  = Type{MainType: "video", SubType: "mp4", Suffixes: []string{"mp4"}, Delimiter: defaultDelimiter}
	OGGType  = Type{MainType: "video", SubType: "ogg", Suffixes: []string{"ogv"}, Delimiter: defaultDelimiter}
	WEBMType = Type{MainType: "video", SubType: "webm", Suffixes: []string{"webm"}, Delimiter: defaultDelimiter}
	GPPType  = Type{MainType: "video", SubType: "3gpp", Suffixes: []string{"3gpp", "3gp"}, Delimiter: defaultDelimiter}

	OctetType = Type{MainType: "application", SubType: "octet-stream"}
)

// DefaultTypes is the default media types supported by Hugo.
var DefaultTypes = Types{
	CalendarType,
	CSSType,
	CSVType,
	SCSSType,
	SASSType,
	HTMLType,
	JavascriptType,
	JSONType,
	RSSType,
	XMLType,
	SVGType,
	TextType,
	OctetType,
	YAMLType,
	TOMLType,
	PNGType,
	JPEGType,
	AVIType,
	MPEGType,
	MP4Type,
	OGGType,
	WEBMType,
	GPPType,
}

func init() {
	sort.Sort(DefaultTypes)
}

// Types is a slice of media types.
type Types []Type

func (t Types) Len() int           { return len(t) }
func (t Types) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Types) Less(i, j int) bool { return t[i].Type() < t[j].Type() }

// GetByType returns a media type for tp.
func (t Types) GetByType(tp string) (Type, bool) {
	for _, tt := range t {
		if strings.EqualFold(tt.Type(), tp) {
			return tt, true
		}
	}

	if !strings.Contains(tp, "+") {
		// Try with the main and sub type
		parts := strings.Split(tp, "/")
		if len(parts) == 2 {
			return t.GetByMainSubType(parts[0], parts[1])
		}
	}

	return Type{}, false
}

// BySuffix will return all media types matching a suffix.
func (t Types) BySuffix(suffix string) []Type {
	var types []Type
	for _, tt := range t {
		if match := tt.matchSuffix(suffix); match != "" {
			types = append(types, tt)
		}
	}
	return types
}

// GetFirstBySuffix will return the first media type matching the given suffix.
func (t Types) GetFirstBySuffix(suffix string) (Type, bool) {
	for _, tt := range t {
		if match := tt.matchSuffix(suffix); match != "" {
			tt.fileSuffix = match
			return tt, true
		}
	}
	return Type{}, false
}

// GetBySuffix gets a media type given as suffix, e.g. "html".
// It will return false if no format could be found, or if the suffix given
// is ambiguous.
// The lookup is case insensitive.
func (t Types) GetBySuffix(suffix string) (tp Type, found bool) {
	for _, tt := range t {
		if match := tt.matchSuffix(suffix); match != "" {
			if found {
				// ambiguous
				found = false
				return
			}
			tp = tt
			tp.fileSuffix = match
			found = true
		}
	}
	return
}

func (m Type) matchSuffix(suffix string) string {
	for _, s := range m.Suffixes {
		if strings.EqualFold(suffix, s) {
			return s
		}
	}

	return ""
}

// GetByMainSubType gets a media type given a main and a sub type e.g. "text" and "plain".
// It will return false if no format could be found, or if the combination given
// is ambiguous.
// The lookup is case insensitive.
func (t Types) GetByMainSubType(mainType, subType string) (tp Type, found bool) {
	for _, tt := range t {
		if strings.EqualFold(mainType, tt.MainType) && strings.EqualFold(subType, tt.SubType) {
			if found {
				// ambiguous
				found = false
				return
			}

			tp = tt
			found = true
		}
	}
	return
}

func suffixIsRemoved() error {
	return errors.New(`MediaType.Suffix is removed. Before Hugo 0.44 this was used both to set a custom file suffix and as way
to augment the mediatype definition (what you see after the "+", e.g. "image/svg+xml").

This had its limitations. For one, it was only possible with one file extension per MIME type.

Now you can specify multiple file suffixes using "suffixes", but you need to specify the full MIME type
identifier:

[mediaTypes]
[mediaTypes."image/svg+xml"]
suffixes = ["svg", "abc" ]

In most cases, it will be enough to just change:

[mediaTypes]
[mediaTypes."my/custom-mediatype"]
suffix = "txt"

To:

[mediaTypes]
[mediaTypes."my/custom-mediatype"]
suffixes = ["txt"]

Note that you can still get the Media Type's suffix from a template: {{ $mediaType.Suffix }}. But this will now map to the MIME type filename.
`)
}

// DecodeTypes takes a list of media type configurations and merges those,
// in the order given, with the Hugo defaults as the last resort.
func DecodeTypes(mms ...map[string]interface{}) (Types, error) {
	var m Types

	// Maps type string to Type. Type string is the full application/svg+xml.
	mmm := make(map[string]Type)
	for _, dt := range DefaultTypes {
		suffixes := make([]string, len(dt.Suffixes))
		copy(suffixes, dt.Suffixes)
		dt.Suffixes = suffixes
		mmm[dt.Type()] = dt
	}

	for _, mm := range mms {
		for k, v := range mm {
			var mediaType Type

			mediaType, found := mmm[k]
			if !found {
				var err error
				mediaType, err = fromString(k)
				if err != nil {
					return m, err
				}
			}

			if err := mapstructure.WeakDecode(v, &mediaType); err != nil {
				return m, err
			}

			vm := v.(map[string]interface{})
			maps.ToLower(vm)
			_, delimiterSet := vm["delimiter"]
			_, suffixSet := vm["suffix"]

			if suffixSet {
				return Types{}, suffixIsRemoved()
			}

			// The user may set the delimiter as an empty string.
			if !delimiterSet && len(mediaType.Suffixes) != 0 {
				mediaType.Delimiter = defaultDelimiter
			}

			mmm[k] = mediaType

		}
	}

	for _, v := range mmm {
		m = append(m, v)
	}
	sort.Sort(m)

	return m, nil
}

// MarshalJSON returns the JSON encoding of m.
func (m Type) MarshalJSON() ([]byte, error) {
	type Alias Type
	return json.Marshal(&struct {
		Type   string `json:"type"`
		String string `json:"string"`
		Alias
	}{
		Type:   m.Type(),
		String: m.String(),
		Alias:  (Alias)(m),
	})
}
