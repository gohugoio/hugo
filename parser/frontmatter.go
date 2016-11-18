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
	"strings"

	toml "github.com/pelletier/go-toml"

	"gopkg.in/yaml.v2"
)

type frontmatterType struct {
	markstart, markend []byte
	Parse              func([]byte) (interface{}, error)
	includeMark        bool
}

func InterfaceToConfig(in interface{}, mark rune) ([]byte, error) {
	if in == nil {
		return []byte{}, errors.New("input was nil")
	}

	b := new(bytes.Buffer)

	switch mark {
	case rune(YAMLLead[0]):
		by, err := yaml.Marshal(in)
		if err != nil {
			return nil, err
		}
		b.Write(by)
		_, err = b.Write([]byte("..."))
		if err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	case rune(TOMLLead[0]):
		tree := toml.TreeFromMap(in.(map[string]interface{}))
		return []byte(tree.String()), nil
	case rune(JSONLead[0]):
		by, err := json.MarshalIndent(in, "", "   ")
		if err != nil {
			return nil, err
		}
		b.Write(by)
		_, err = b.Write([]byte("\n"))
		if err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	default:
		return nil, errors.New("Unsupported Format provided")
	}
}

func InterfaceToFrontMatter(in interface{}, mark rune) ([]byte, error) {
	if in == nil {
		return []byte{}, errors.New("input was nil")
	}

	b := new(bytes.Buffer)

	switch mark {
	case rune(YAMLLead[0]):
		_, err := b.Write([]byte(YAMLDelimUnix))
		if err != nil {
			return nil, err
		}
		by, err := yaml.Marshal(in)
		if err != nil {
			return nil, err
		}
		b.Write(by)
		_, err = b.Write([]byte(YAMLDelimUnix))
		if err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	case rune(TOMLLead[0]):
		_, err := b.Write([]byte(TOMLDelimUnix))
		if err != nil {
			return nil, err
		}

		tree := toml.TreeFromMap(in.(map[string]interface{}))
		b.Write([]byte(tree.String()))
		_, err = b.Write([]byte("\n" + TOMLDelimUnix))
		if err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	case rune(JSONLead[0]):
		by, err := json.MarshalIndent(in, "", "   ")
		if err != nil {
			return nil, err
		}
		b.Write(by)
		_, err = b.Write([]byte("\n"))
		if err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	default:
		return nil, errors.New("Unsupported Format provided")
	}
}

func FormatToLeadRune(kind string) rune {
	switch FormatSanitize(kind) {
	case "yaml":
		return rune([]byte(YAMLLead)[0])
	case "json":
		return rune([]byte(JSONLead)[0])
	default:
		return rune([]byte(TOMLLead)[0])
	}
}

// TODO(bep) move to helpers
func FormatSanitize(kind string) string {
	switch strings.ToLower(kind) {
	case "yaml", "yml":
		return "yaml"
	case "toml", "tml":
		return "toml"
	case "json", "js":
		return "json"
	default:
		return "toml"
	}
}

// DetectFrontMatter detects the type of frontmatter analysing its first character.
func DetectFrontMatter(mark rune) (f *frontmatterType) {
	switch mark {
	case '-':
		return &frontmatterType{[]byte(YAMLDelim), []byte(YAMLDelim), HandleYAMLMetaData, false}
	case '+':
		return &frontmatterType{[]byte(TOMLDelim), []byte(TOMLDelim), HandleTOMLMetaData, false}
	case '{':
		return &frontmatterType{[]byte{'{'}, []byte{'}'}, HandleJSONMetaData, true}
	default:
		return nil
	}
}

func HandleTOMLMetaData(datum []byte) (interface{}, error) {
	m := map[string]interface{}{}
	datum = removeTOMLIdentifier(datum)

	tree, err := toml.Load(string(datum))

	if err != nil {
		return m, err
	}

	m = tree.ToMap()

	return m, nil
}

func removeTOMLIdentifier(datum []byte) []byte {
	return bytes.Replace(datum, []byte(TOMLDelim), []byte(""), -1)
}

func HandleYAMLMetaData(datum []byte) (interface{}, error) {
	m := map[string]interface{}{}
	err := yaml.Unmarshal(datum, &m)
	return m, err
}

func HandleJSONMetaData(datum []byte) (interface{}, error) {
	var f interface{}
	err := json.Unmarshal(datum, &f)
	return f, err
}
