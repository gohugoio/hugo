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
	"path/filepath"
	"strings"
)

type Format string

const (
	// These are the supported metadata  formats in Hugo. Most of these are also
	// supported as /data formats.
	ORG  Format = "org"
	JSON Format = "json"
	TOML Format = "toml"
	YAML Format = "yaml"
	CSV  Format = "csv"
	XML  Format = "xml"
)

// FormatFromStrings returns the first non-empty Format from the given strings.
func FormatFromStrings(ss ...string) Format {
	for _, s := range ss {
		if f := FormatFromString(s); f != "" {
			return f
		}
	}
	return ""
}

// FormatFromString turns formatStr, typically a file extension without any ".",
// into a Format. It returns an empty string for unknown formats.
func FormatFromString(formatStr string) Format {
	formatStr = strings.ToLower(formatStr)
	if strings.Contains(formatStr, ".") {
		// Assume a filename
		formatStr = strings.TrimPrefix(filepath.Ext(formatStr), ".")
	}
	switch formatStr {
	case "yaml", "yml":
		return YAML
	case "json":
		return JSON
	case "toml":
		return TOML
	case "org":
		return ORG
	case "csv":
		return CSV
	case "xml":
		return XML
	}

	return ""
}

// FormatFromContentString tries to detect the format (JSON, YAML, TOML or XML)
// in the given string.
// It return an empty string if no format could be detected.
func (d Decoder) FormatFromContentString(data string) Format {
	csvIdx := strings.IndexRune(data, d.Delimiter)
	jsonIdx := strings.Index(data, "{")
	yamlIdx := strings.Index(data, ":")
	xmlIdx := strings.Index(data, "<")
	tomlIdx := strings.Index(data, "=")

	if isLowerIndexThan(csvIdx, jsonIdx, yamlIdx, xmlIdx, tomlIdx) {
		return CSV
	}

	if isLowerIndexThan(jsonIdx, yamlIdx, xmlIdx, tomlIdx) {
		return JSON
	}

	if isLowerIndexThan(yamlIdx, xmlIdx, tomlIdx) {
		return YAML
	}

	if isLowerIndexThan(xmlIdx, tomlIdx) {
		return XML
	}

	if tomlIdx != -1 {
		return TOML
	}

	return ""
}

func isLowerIndexThan(first int, others ...int) bool {
	if first == -1 {
		return false
	}
	for _, other := range others {
		if other != -1 && other < first {
			return false
		}
	}

	return true
}
