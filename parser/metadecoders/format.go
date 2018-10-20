// Copyright 2018 The Hugo Authors. All rights reserved.
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

package metadecoders

import (
	"strings"

	"github.com/gohugoio/hugo/parser/pageparser"
)

type Format string

const (
	// These are the supported metdata  formats in Hugo. Most of these are also
	// supported as /data formats.
	ORG  Format = "org"
	JSON Format = "json"
	TOML Format = "toml"
	YAML Format = "yaml"
)

// FormatFromString turns formatStr, typically a file extension without any ".",
// into a Format. It returns an empty string for unknown formats.
func FormatFromString(formatStr string) Format {
	formatStr = strings.ToLower(formatStr)
	switch formatStr {
	case "yaml", "yml":
		return YAML
	case "json":
		return JSON
	case "toml":
		return TOML
	case "org":
		return ORG
	}

	return ""

}

// FormatFromFrontMatterType will return empty if not supported.
func FormatFromFrontMatterType(typ pageparser.ItemType) Format {
	switch typ {
	case pageparser.TypeFrontMatterJSON:
		return JSON
	case pageparser.TypeFrontMatterORG:
		return ORG
	case pageparser.TypeFrontMatterTOML:
		return TOML
	case pageparser.TypeFrontMatterYAML:
		return YAML
	default:
		return ""
	}
}
