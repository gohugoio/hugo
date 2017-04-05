// Copyright 2017-present The Hugo Authors. All rights reserved.
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
	"fmt"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// A media type (also known as MIME type and content type) is a two-part identifier for
// file formats and format contents transmitted on the Internet.
// For Hugo's use case, we use the top-level type name / subtype name + suffix.
// One example would be image/jpeg+jpg
// If suffix is not provided, the sub type will be used.
// See // https://en.wikipedia.org/wiki/Media_type
type Type struct {
	MainType string // i.e. text
	SubType  string // i.e. html
	Suffix   string // i.e html
}

// FromTypeString creates a new Type given a type sring on the form MainType/SubType and
// an optional suffix, e.g. "text/html" or "text/html+html".
func FromString(t string) (Type, error) {
	t = strings.ToLower(t)
	parts := strings.Split(t, "/")
	if len(parts) != 2 {
		return Type{}, fmt.Errorf("cannot parse %q as a media type", t)
	}
	mainType := parts[0]
	subParts := strings.Split(parts[1], "+")

	subType := subParts[0]
	var suffix string

	if len(subParts) == 1 {
		suffix = subType
	} else {
		suffix = subParts[1]
	}

	return Type{MainType: mainType, SubType: subType, Suffix: suffix}, nil
}

// Type returns a string representing the main- and sub-type of a media type, i.e. "text/css".
// Hugo will register a set of default media types.
// These can be overridden by the user in the configuration,
// by defining a media type with the same Type.
func (m Type) Type() string {
	return fmt.Sprintf("%s/%s", m.MainType, m.SubType)
}

func (m Type) String() string {
	if m.Suffix != "" {
		return fmt.Sprintf("%s/%s+%s", m.MainType, m.SubType, m.Suffix)
	}
	return fmt.Sprintf("%s/%s", m.MainType, m.SubType)
}

var (
	CalendarType   = Type{"text", "calendar", "ics"}
	CSSType        = Type{"text", "css", "css"}
	CSVType        = Type{"text", "csv", "csv"}
	HTMLType       = Type{"text", "html", "html"}
	JavascriptType = Type{"application", "javascript", "js"}
	JSONType       = Type{"application", "json", "json"}
	RSSType        = Type{"application", "rss", "xml"}
	XMLType        = Type{"application", "xml", "xml"}
	TextType       = Type{"text", "plain", "txt"}
)

var DefaultTypes = Types{
	CalendarType,
	CSSType,
	CSVType,
	HTMLType,
	JavascriptType,
	JSONType,
	RSSType,
	XMLType,
	TextType,
}

func init() {
	sort.Sort(DefaultTypes)
}

type Types []Type

func (t Types) Len() int           { return len(t) }
func (t Types) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Types) Less(i, j int) bool { return t[i].Type() < t[j].Type() }

func (t Types) GetByType(tp string) (Type, bool) {
	for _, tt := range t {
		if strings.EqualFold(tt.Type(), tp) {
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
		if strings.EqualFold(suffix, tt.Suffix) {
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

// DecodeTypes takes a list of media type configurations and merges those,
// in the order given, with the Hugo defaults as the last resort.
func DecodeTypes(maps ...map[string]interface{}) (Types, error) {
	m := make(Types, len(DefaultTypes))
	copy(m, DefaultTypes)

	for _, mm := range maps {
		for k, v := range mm {
			// It may be tempting to put the full media type in the key, e.g.
			//  "text/css+css", but that will break the logic below.
			if strings.Contains(k, "+") {
				return Types{}, fmt.Errorf("media type keys cannot contain any '+' chars. Valid example is %q", "text/css")
			}

			found := false
			for i, vv := range m {
				// Match by type, i.e. "text/css"
				if strings.EqualFold(k, vv.Type()) {
					// Merge it with the existing
					if err := mapstructure.WeakDecode(v, &m[i]); err != nil {
						return m, err
					}
					found = true
				}
			}
			if !found {
				mediaType, err := FromString(k)
				if err != nil {
					return m, err
				}

				if err := mapstructure.WeakDecode(v, &mediaType); err != nil {
					return m, err
				}

				m = append(m, mediaType)
			}
		}
	}

	sort.Sort(m)

	return m, nil
}

func (t Type) MarshalJSON() ([]byte, error) {
	type Alias Type
	return json.Marshal(&struct {
		Type   string
		String string
		Alias
	}{
		Type:   t.Type(),
		String: t.String(),
		Alias:  (Alias)(t),
	})
}
