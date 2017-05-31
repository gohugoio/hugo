// Copyright 2015 The Hugo Authors. All rights reserved.
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

package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/chaseadamsio/goorgeous"

	"gopkg.in/yaml.v2"
)

// FrontmatterType represents a type of frontmatter.
type FrontmatterType struct {
	// Parse decodes content into a Go interface.
	Parse func([]byte) (interface{}, error)

	markstart, markend []byte // starting and ending delimiters
	includeMark        bool   // include start and end mark in output
}

// InterfaceToConfig encodes a given input based upon the mark and writes to w.
func InterfaceToConfig(in interface{}, mark rune, w io.Writer) error {
	if in == nil {
		return errors.New("input was nil")
	}

	switch mark {
	case rune(YAMLLead[0]):
		b, err := yaml.Marshal(in)
		if err != nil {
			return err
		}

		_, err = w.Write(b)
		return err

	case rune(TOMLLead[0]):
		return toml.NewEncoder(w).Encode(in)
	case rune(JSONLead[0]):
		b, err := json.MarshalIndent(in, "", "   ")
		if err != nil {
			return err
		}

		_, err = w.Write(b)
		if err != nil {
			return err
		}

		_, err = w.Write([]byte{'\n'})
		return err

	default:
		return errors.New("Unsupported Format provided")
	}
}

// InterfaceToFrontMatter encodes a given input into a frontmatter
// representation based upon the mark with the appropriate front matter delimiters
// surrounding the output, which is written to w.
func InterfaceToFrontMatter(in interface{}, mark rune, w io.Writer) error {
	if in == nil {
		return errors.New("input was nil")
	}

	switch mark {
	case rune(YAMLLead[0]):
		_, err := w.Write([]byte(YAMLDelimUnix))
		if err != nil {
			return err
		}

		err = InterfaceToConfig(in, mark, w)
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(YAMLDelimUnix))
		return err

	case rune(TOMLLead[0]):
		_, err := w.Write([]byte(TOMLDelimUnix))
		if err != nil {
			return err
		}

		err = InterfaceToConfig(in, mark, w)

		if err != nil {
			return err
		}

		_, err = w.Write([]byte("\n" + TOMLDelimUnix))
		return err

	default:
		return InterfaceToConfig(in, mark, w)
	}
}

// FormatToLeadRune takes a given format kind and return the leading front
// matter delimiter.
func FormatToLeadRune(kind string) rune {
	switch FormatSanitize(kind) {
	case "yaml":
		return rune([]byte(YAMLLead)[0])
	case "json":
		return rune([]byte(JSONLead)[0])
	case "org":
		return '#'
	default:
		return rune([]byte(TOMLLead)[0])
	}
}

// FormatSanitize returns the canonical format name for a given kind.
//
// TODO(bep) move to helpers
func FormatSanitize(kind string) string {
	switch strings.ToLower(kind) {
	case "yaml", "yml":
		return "yaml"
	case "toml", "tml":
		return "toml"
	case "json", "js":
		return "json"
	case "org":
		return kind
	default:
		return "toml"
	}
}

// DetectFrontMatter detects the type of frontmatter analysing its first character.
func DetectFrontMatter(mark rune) (f *FrontmatterType) {
	switch mark {
	case '-':
		return &FrontmatterType{HandleYAMLMetaData, []byte(YAMLDelim), []byte(YAMLDelim), false}
	case '+':
		return &FrontmatterType{HandleTOMLMetaData, []byte(TOMLDelim), []byte(TOMLDelim), false}
	case '{':
		return &FrontmatterType{HandleJSONMetaData, []byte{'{'}, []byte{'}'}, true}
	case '#':
		return &FrontmatterType{HandleOrgMetaData, []byte("#+"), []byte("\n"), false}
	default:
		return nil
	}
}

// HandleTOMLMetaData unmarshals TOML-encoded datum and returns a Go interface
// representing the encoded data structure.
func HandleTOMLMetaData(datum []byte) (interface{}, error) {
	m := map[string]interface{}{}
	datum = removeTOMLIdentifier(datum)

	_, err := toml.Decode(string(datum), &m)

	return m, err

}

// removeTOMLIdentifier removes, if necessary, beginning and ending TOML
// frontmatter delimiters from a byte slice.
func removeTOMLIdentifier(datum []byte) []byte {
	ld := len(datum)
	if ld < 8 {
		return datum
	}

	b := bytes.TrimPrefix(datum, []byte(TOMLDelim))
	if ld-len(b) != 3 {
		// No TOML prefix trimmed, so bail out
		return datum
	}

	b = bytes.Trim(b, "\r\n")
	return bytes.TrimSuffix(b, []byte(TOMLDelim))
}

// HandleYAMLMetaData unmarshals YAML-encoded datum and returns a Go interface
// representing the encoded data structure.
func HandleYAMLMetaData(datum []byte) (interface{}, error) {
	m := map[string]interface{}{}
	err := yaml.Unmarshal(datum, &m)
	return m, err
}

// HandleJSONMetaData unmarshals JSON-encoded datum and returns a Go interface
// representing the encoded data structure.
func HandleJSONMetaData(datum []byte) (interface{}, error) {
	if datum == nil {
		// Package json returns on error on nil input.
		// Return an empty map to be consistent with our other supported
		// formats.
		return make(map[string]interface{}), nil
	}

	var f interface{}
	err := json.Unmarshal(datum, &f)
	return f, err
}

// HandleOrgMetaData unmarshals org-mode encoded datum and returns a Go
// interface representing the encoded data structure.
func HandleOrgMetaData(datum []byte) (interface{}, error) {
	return goorgeous.OrgHeaders(datum)
}
