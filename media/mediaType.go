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

// Package media contains Media Type (MIME type) related types and functions.
package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/mitchellh/mapstructure"
)

var zero Type

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
	MainType  string `json:"mainType"`  // i.e. text
	SubType   string `json:"subType"`   // i.e. html
	Delimiter string `json:"delimiter"` // e.g. "."

	// FirstSuffix holds the first suffix defined for this Type.
	FirstSuffix SuffixInfo `json:"firstSuffix"`

	// This is the optional suffix after the "+" in the MIME type,
	//  e.g. "xml" in "application/rss+xml".
	mimeSuffix string

	// E.g. "jpg,jpeg"
	// Stored as a string to make Type comparable.
	suffixesCSV string
}

// SuffixInfo holds information about a Type's suffix.
type SuffixInfo struct {
	Suffix     string `json:"suffix"`
	FullSuffix string `json:"fullSuffix"`
}

// FromContent resolve the Type primarily using http.DetectContentType.
// If http.DetectContentType resolves to application/octet-stream, a zero Type is returned.
// If http.DetectContentType  resolves to text/plain or application/xml, we try to get more specific using types and ext.
func FromContent(types Types, extensionHints []string, content []byte) Type {
	t := strings.Split(http.DetectContentType(content), ";")[0]
	if t == "application/octet-stream" {
		return zero
	}

	var found bool
	m, found := types.GetByType(t)
	if !found {
		if t == "text/xml" {
			// This is how it's configured in Hugo by default.
			m, found = types.GetByType("application/xml")
		}
	}

	if !found {
		return zero
	}

	var mm Type

	for _, extension := range extensionHints {
		extension = strings.TrimPrefix(extension, ".")
		mm, _, found = types.GetFirstBySuffix(extension)
		if found {
			break
		}
	}

	if found {
		if m == mm {
			return m
		}

		if m.IsText() && mm.IsText() {
			// http.DetectContentType isn't brilliant when it comes to common text formats, so we need to do better.
			// For now we say that if it's detected to be a text format and the extension/content type in header reports
			// it to be a text format, then we use that.
			return mm
		}

		// E.g. an image with a *.js extension.
		return zero
	}

	return m
}

// FromStringAndExt creates a Type from a MIME string and a given extension.
func FromStringAndExt(t, ext string) (Type, error) {
	tp, err := FromString(t)
	if err != nil {
		return tp, err
	}
	tp.suffixesCSV = strings.TrimPrefix(ext, ".")
	tp.Delimiter = defaultDelimiter
	tp.init()
	return tp, nil
}

// FromString creates a new Type given a type string on the form MainType/SubType and
// an optional suffix, e.g. "text/html" or "text/html+html".
func FromString(t string) (Type, error) {
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

// For internal use.
func (m Type) String() string {
	return m.Type()
}

// Suffixes returns all valid file suffixes for this type.
func (m Type) Suffixes() []string {
	if m.suffixesCSV == "" {
		return nil
	}

	return strings.Split(m.suffixesCSV, ",")
}

// IsText returns whether this Type is a text format.
// Note that this may currently return false negatives.
// TODO(bep) improve
func (m Type) IsText() bool {
	if m.MainType == "text" {
		return true
	}
	switch m.SubType {
	case "javascript", "json", "rss", "xml", "svg", TOMLType.SubType, YAMLType.SubType:
		return true
	}
	return false
}

func (m *Type) init() {
	m.FirstSuffix.FullSuffix = ""
	m.FirstSuffix.Suffix = ""
	if suffixes := m.Suffixes(); suffixes != nil {
		m.FirstSuffix.Suffix = suffixes[0]
		m.FirstSuffix.FullSuffix = m.Delimiter + m.FirstSuffix.Suffix
	}
}

// WithDelimiterAndSuffixes is used in tests.
func WithDelimiterAndSuffixes(t Type, delimiter, suffixesCSV string) Type {
	t.Delimiter = delimiter
	t.suffixesCSV = suffixesCSV
	t.init()
	return t
}

func newMediaType(main, sub string, suffixes []string) Type {
	t := Type{MainType: main, SubType: sub, suffixesCSV: strings.Join(suffixes, ","), Delimiter: defaultDelimiter}
	t.init()
	return t
}

func newMediaTypeWithMimeSuffix(main, sub, mimeSuffix string, suffixes []string) Type {
	mt := newMediaType(main, sub, suffixes)
	mt.mimeSuffix = mimeSuffix
	mt.init()
	return mt
}

// Definitions from https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types etc.
// Note that from Hugo 0.44 we only set Suffix if it is part of the MIME type.
var (
	CalendarType   = newMediaType("text", "calendar", []string{"ics"})
	CSSType        = newMediaType("text", "css", []string{"css"})
	SCSSType       = newMediaType("text", "x-scss", []string{"scss"})
	SASSType       = newMediaType("text", "x-sass", []string{"sass"})
	CSVType        = newMediaType("text", "csv", []string{"csv"})
	HTMLType       = newMediaType("text", "html", []string{"html"})
	JavascriptType = newMediaType("text", "javascript", []string{"js", "jsm", "mjs"})
	TypeScriptType = newMediaType("text", "typescript", []string{"ts"})
	TSXType        = newMediaType("text", "tsx", []string{"tsx"})
	JSXType        = newMediaType("text", "jsx", []string{"jsx"})

	JSONType           = newMediaType("application", "json", []string{"json"})
	WebAppManifestType = newMediaTypeWithMimeSuffix("application", "manifest", "json", []string{"webmanifest"})
	RSSType            = newMediaTypeWithMimeSuffix("application", "rss", "xml", []string{"xml", "rss"})
	XMLType            = newMediaType("application", "xml", []string{"xml"})
	SVGType            = newMediaTypeWithMimeSuffix("image", "svg", "xml", []string{"svg"})
	TextType           = newMediaType("text", "plain", []string{"txt"})
	TOMLType           = newMediaType("application", "toml", []string{"toml"})
	YAMLType           = newMediaType("application", "yaml", []string{"yaml", "yml"})

	// Common image types
	PNGType  = newMediaType("image", "png", []string{"png"})
	JPEGType = newMediaType("image", "jpeg", []string{"jpg", "jpeg", "jpe", "jif", "jfif"})
	GIFType  = newMediaType("image", "gif", []string{"gif"})
	TIFFType = newMediaType("image", "tiff", []string{"tif", "tiff"})
	BMPType  = newMediaType("image", "bmp", []string{"bmp"})
	WEBPType = newMediaType("image", "webp", []string{"webp"})

	// Common font types
	TrueTypeFontType = newMediaType("font", "ttf", []string{"ttf"})
	OpenTypeFontType = newMediaType("font", "otf", []string{"otf"})

	// Common document types
	PDFType      = newMediaType("application", "pdf", []string{"pdf"})
	MarkdownType = newMediaType("text", "markdown", []string{"md", "markdown"})

	// Common video types
	AVIType  = newMediaType("video", "x-msvideo", []string{"avi"})
	MPEGType = newMediaType("video", "mpeg", []string{"mpg", "mpeg"})
	MP4Type  = newMediaType("video", "mp4", []string{"mp4"})
	OGGType  = newMediaType("video", "ogg", []string{"ogv"})
	WEBMType = newMediaType("video", "webm", []string{"webm"})
	GPPType  = newMediaType("video", "3gpp", []string{"3gpp", "3gp"})

	OctetType = newMediaType("application", "octet-stream", nil)
)

// DefaultTypes is the default media types supported by Hugo.
var DefaultTypes = Types{
	CalendarType,
	CSSType,
	CSVType,
	SCSSType,
	SASSType,
	HTMLType,
	MarkdownType,
	JavascriptType,
	TypeScriptType,
	TSXType,
	JSXType,
	JSONType,
	WebAppManifestType,
	RSSType,
	XMLType,
	SVGType,
	TextType,
	OctetType,
	YAMLType,
	TOMLType,
	PNGType,
	GIFType,
	BMPType,
	JPEGType,
	WEBPType,
	AVIType,
	MPEGType,
	MP4Type,
	OGGType,
	WEBMType,
	GPPType,
	OpenTypeFontType,
	TrueTypeFontType,
	PDFType,
}

func init() {
	sort.Sort(DefaultTypes)

	// Sanity check.
	seen := make(map[Type]bool)
	for _, t := range DefaultTypes {
		if seen[t] {
			panic(fmt.Sprintf("MediaType %s duplicated in list", t))
		}
		seen[t] = true
	}
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
	suffix = strings.ToLower(suffix)
	var types []Type
	for _, tt := range t {
		if tt.hasSuffix(suffix) {
			types = append(types, tt)
		}
	}
	return types
}

// GetFirstBySuffix will return the first type matching the given suffix.
func (t Types) GetFirstBySuffix(suffix string) (Type, SuffixInfo, bool) {
	suffix = strings.ToLower(suffix)
	for _, tt := range t {
		if tt.hasSuffix(suffix) {
			return tt, SuffixInfo{
				FullSuffix: tt.Delimiter + suffix,
				Suffix:     suffix,
			}, true
		}
	}
	return Type{}, SuffixInfo{}, false
}

// GetBySuffix gets a media type given as suffix, e.g. "html".
// It will return false if no format could be found, or if the suffix given
// is ambiguous.
// The lookup is case insensitive.
func (t Types) GetBySuffix(suffix string) (tp Type, si SuffixInfo, found bool) {
	suffix = strings.ToLower(suffix)
	for _, tt := range t {
		if tt.hasSuffix(suffix) {
			if found {
				// ambiguous
				found = false
				return
			}
			tp = tt
			si = SuffixInfo{
				FullSuffix: tt.Delimiter + suffix,
				Suffix:     suffix,
			}
			found = true
		}
	}
	return
}

func (m Type) hasSuffix(suffix string) bool {
	return strings.Contains(","+m.suffixesCSV+",", ","+suffix+",")
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
func DecodeTypes(mms ...map[string]any) (Types, error) {
	var m Types

	// Maps type string to Type. Type string is the full application/svg+xml.
	mmm := make(map[string]Type)
	for _, dt := range DefaultTypes {
		mmm[dt.Type()] = dt
	}

	for _, mm := range mms {
		for k, v := range mm {
			var mediaType Type

			mediaType, found := mmm[k]
			if !found {
				var err error
				mediaType, err = FromString(k)
				if err != nil {
					return m, err
				}
			}

			if err := mapstructure.WeakDecode(v, &mediaType); err != nil {
				return m, err
			}

			vm := maps.ToStringMap(v)
			maps.PrepareParams(vm)
			_, delimiterSet := vm["delimiter"]
			_, suffixSet := vm["suffix"]

			if suffixSet {
				return Types{}, suffixIsRemoved()
			}

			if suffixes, found := vm["suffixes"]; found {
				mediaType.suffixesCSV = strings.TrimSpace(strings.ToLower(strings.Join(cast.ToStringSlice(suffixes), ",")))
			}

			// The user may set the delimiter as an empty string.
			if !delimiterSet && mediaType.suffixesCSV != "" {
				mediaType.Delimiter = defaultDelimiter
			}

			mediaType.init()

			mmm[k] = mediaType

		}
	}

	for _, v := range mmm {
		m = append(m, v)
	}
	sort.Sort(m)

	return m, nil
}

// IsZero reports whether this Type represents a zero value.
// For internal use.
func (m Type) IsZero() bool {
	return m.SubType == ""
}

// MarshalJSON returns the JSON encoding of m.
// For internal use.
func (m Type) MarshalJSON() ([]byte, error) {
	type Alias Type
	return json.Marshal(&struct {
		Alias
		Type     string   `json:"type"`
		String   string   `json:"string"`
		Suffixes []string `json:"suffixes"`
	}{
		Alias:    (Alias)(m),
		Type:     m.Type(),
		String:   m.String(),
		Suffixes: strings.Split(m.suffixesCSV, ","),
	})
}
