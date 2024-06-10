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
	"fmt"
	"net/http"
	"strings"
)

var zero Type

const (
	DefaultDelimiter = "."
)

// MediaType (also known as MIME type and content type) is a two-part identifier for
// file formats and format contents transmitted on the Internet.
// For Hugo's use case, we use the top-level type name / subtype name + suffix.
// One example would be application/svg+xml
// If suffix is not provided, the sub type will be used.
// <docsmeta>{ "name": "MediaType" }</docsmeta>
type Type struct {
	// The full MIME type string, e.g. "application/rss+xml".
	Type string `json:"-"`

	// The top-level type name, e.g. "application".
	MainType string `json:"mainType"`
	// The subtype name, e.g. "rss".
	SubType string `json:"subType"`
	// The delimiter before the suffix, e.g. ".".
	Delimiter string `json:"delimiter"`

	// FirstSuffix holds the first suffix defined for this MediaType.
	FirstSuffix SuffixInfo `json:"-"`

	// This is the optional suffix after the "+" in the MIME type,
	//  e.g. "xml" in "application/rss+xml".
	mimeSuffix string

	// E.g. "jpg,jpeg"
	// Stored as a string to make Type comparable.
	// For internal use only.
	SuffixesCSV string `json:"-"`
}

// SuffixInfo holds information about a Media Type's suffix.
type SuffixInfo struct {
	// Suffix is the suffix without the delimiter, e.g. "xml".
	Suffix string `json:"suffix"`

	// FullSuffix is the suffix with the delimiter, e.g. ".xml".
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

// FromStringAndExt creates a Type from a MIME string and a given extensions
func FromStringAndExt(t string, ext ...string) (Type, error) {
	tp, err := FromString(t)
	if err != nil {
		return tp, err
	}
	for i, e := range ext {
		ext[i] = strings.TrimPrefix(e, ".")
	}
	tp.SuffixesCSV = strings.Join(ext, ",")
	tp.Delimiter = DefaultDelimiter
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

	var typ string
	if suffix != "" {
		typ = mainType + "/" + subType + "+" + suffix
	} else {
		typ = mainType + "/" + subType
	}

	return Type{Type: typ, MainType: mainType, SubType: subType, mimeSuffix: suffix}, nil
}

// For internal use.
func (m Type) String() string {
	return m.Type
}

// Suffixes returns all valid file suffixes for this type.
func (m Type) Suffixes() []string {
	if m.SuffixesCSV == "" {
		return nil
	}

	return strings.Split(m.SuffixesCSV, ",")
}

// IsText returns whether this Type is a text format.
// Note that this may currently return false negatives.
// TODO(bep) improve
// For internal use.
func (m Type) IsText() bool {
	if m.MainType == "text" {
		return true
	}
	switch m.SubType {
	case "javascript", "json", "rss", "xml", "svg", "toml", "yml", "yaml":
		return true
	}
	return false
}

// For internal use.
func (m Type) IsHTML() bool {
	return m.SubType == Builtin.HTMLType.SubType
}

// For internal use.
func (m Type) IsMarkdown() bool {
	return m.SubType == Builtin.MarkdownType.SubType
}

func InitMediaType(m *Type) {
	m.init()
}

func (m *Type) init() {
	m.FirstSuffix.FullSuffix = ""
	m.FirstSuffix.Suffix = ""
	if suffixes := m.Suffixes(); suffixes != nil {
		m.FirstSuffix.Suffix = suffixes[0]
		m.FirstSuffix.FullSuffix = m.Delimiter + m.FirstSuffix.Suffix
	}
}

func newMediaType(main, sub string, suffixes []string) Type {
	t := Type{MainType: main, SubType: sub, SuffixesCSV: strings.Join(suffixes, ","), Delimiter: DefaultDelimiter}
	t.init()
	return t
}

func newMediaTypeWithMimeSuffix(main, sub, mimeSuffix string, suffixes []string) Type {
	mt := newMediaType(main, sub, suffixes)
	mt.mimeSuffix = mimeSuffix
	mt.init()
	return mt
}

// Types is a slice of media types.
// <docsmeta>{ "name": "MediaTypes" }</docsmeta>
type Types []Type

func (t Types) Len() int           { return len(t) }
func (t Types) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Types) Less(i, j int) bool { return t[i].Type < t[j].Type }

// GetBestMatch returns the best match for the given media type string.
func (t Types) GetBestMatch(s string) (Type, bool) {
	// First try an exact match.
	if mt, found := t.GetByType(s); found {
		return mt, true
	}

	// Try main type.
	if mt, found := t.GetBySubType(s); found {
		return mt, true
	}

	// Try extension.
	if mt, _, found := t.GetFirstBySuffix(s); found {
		return mt, true
	}

	return Type{}, false
}

// GetByType returns a media type for tp.
func (t Types) GetByType(tp string) (Type, bool) {
	for _, tt := range t {
		if strings.EqualFold(tt.Type, tp) {
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

func (t Types) IsTextSuffix(suffix string) bool {
	suffix = strings.ToLower(suffix)
	for _, tt := range t {
		if tt.hasSuffix(suffix) {
			return tt.IsText()
		}
	}
	return false
}

func (m Type) hasSuffix(suffix string) bool {
	return strings.Contains(","+m.SuffixesCSV+",", ","+suffix+",")
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

// GetBySubType gets a media type given a sub type e.g. "plain".
func (t Types) GetBySubType(subType string) (tp Type, found bool) {
	for _, tt := range t {
		if strings.EqualFold(subType, tt.SubType) {
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
		Type:     m.Type,
		String:   m.String(),
		Suffixes: strings.Split(m.SuffixesCSV, ","),
	})
}
