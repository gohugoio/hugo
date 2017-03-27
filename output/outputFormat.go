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

package output

import (
	"fmt"
	"strings"

	"github.com/spf13/hugo/media"
)

var (
	// An ordered list of built-in output formats
	//
	// See https://www.ampproject.org/learn/overview/
	AMPFormat = Format{
		Name:      "AMP",
		MediaType: media.HTMLType,
		BaseName:  "index",
		Path:      "amp",
		Rel:       "amphtml",
		IsHTML:    true,
	}

	CalendarFormat = Format{
		Name:        "Calendar",
		MediaType:   media.CalendarType,
		IsPlainText: true,
		Protocol:    "webcal://",
		BaseName:    "index",
		Rel:         "alternate",
	}

	CSSFormat = Format{
		Name:      "CSS",
		MediaType: media.CSSType,
		BaseName:  "styles",
		Rel:       "stylesheet",
	}

	HTMLFormat = Format{
		Name:      "HTML",
		MediaType: media.HTMLType,
		BaseName:  "index",
		Rel:       "canonical",
		IsHTML:    true,
	}

	JSONFormat = Format{
		Name:        "JSON",
		MediaType:   media.JSONType,
		BaseName:    "index",
		IsPlainText: true,
		Rel:         "alternate",
	}

	RSSFormat = Format{
		Name:      "RSS",
		MediaType: media.RSSType,
		BaseName:  "index",
		NoUgly:    true,
		Rel:       "alternate",
	}
)

var builtInTypes = map[string]Format{
	strings.ToLower(AMPFormat.Name):      AMPFormat,
	strings.ToLower(CalendarFormat.Name): CalendarFormat,
	strings.ToLower(CSSFormat.Name):      CSSFormat,
	strings.ToLower(HTMLFormat.Name):     HTMLFormat,
	strings.ToLower(JSONFormat.Name):     JSONFormat,
	strings.ToLower(RSSFormat.Name):      RSSFormat,
}

type Formats []Format

func (formats Formats) GetByName(name string) (f Format, found bool) {
	for _, ff := range formats {
		if name == ff.Name {
			f = ff
			found = true
			return
		}
	}
	return
}

// Format represents an output representation, usually to a file on disk.
type Format struct {
	// The Name is used as an identifier. Internal output formats (i.e. HTML and RSS)
	// can be overridden by providing a new definition for those types.
	Name string

	MediaType media.Type

	// Must be set to a value when there are two or more conflicting mediatype for the same resource.
	Path string

	// The base output file name used when not using "ugly URLs", defaults to "index".
	BaseName string

	// The value to use for rel links
	//
	// See https://www.w3schools.com/tags/att_link_rel.asp
	//
	// AMP has a special requirement in this department, see:
	// https://www.ampproject.org/docs/guides/deploy/discovery
	// I.e.:
	// <link rel="amphtml" href="https://www.example.com/url/to/amp/document.html">
	Rel string

	// The protocol to use, i.e. "webcal://". Defaults to the protocol of the baseURL.
	Protocol string

	// IsPlainText decides whether to use text/template or html/template
	// as template parser.
	IsPlainText bool

	// IsHTML returns whether this format is int the HTML family. This includes
	// HTML, AMP etc. This is used to decide when to create alias redirects etc.
	IsHTML bool

	// Enable to ignore the global uglyURLs setting.
	NoUgly bool
}

func GetFormat(key string) (Format, bool) {
	found, ok := builtInTypes[key]
	if !ok {
		found, ok = builtInTypes[strings.ToLower(key)]
	}
	return found, ok
}

// TODO(bep) outputs rewamp on global config?
func GetFormats(keys ...string) (Formats, error) {
	var types []Format

	for _, key := range keys {
		tpe, ok := GetFormat(key)
		if !ok {
			return types, fmt.Errorf("OutputFormat with key %q not found", key)
		}
		types = append(types, tpe)
	}

	return types, nil
}

func (t Format) BaseFilename() string {
	return t.BaseName + "." + t.MediaType.Suffix
}
