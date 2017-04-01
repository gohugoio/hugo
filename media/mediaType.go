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
	"fmt"
)

type Types []Type

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

// TODO(bep) output mime.AddExtensionType
