// Copyright © 2013 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
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
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

type FrontmatterType struct {
	markstart, markend []byte
	Parse              func([]byte) (interface{}, error)
	includeMark        bool
}

func InterfaceToConfig(in interface{}, mark rune) ([]byte, error) {
	if in == nil {
		return []byte{}, fmt.Errorf("input was nil")
	}

	b := new(bytes.Buffer)

	switch mark {
	case rune(YAML_LEAD[0]):
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
	case rune(TOML_LEAD[0]):
		err := toml.NewEncoder(b).Encode(in)
		if err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	case rune(JSON_LEAD[0]):
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
		return nil, fmt.Errorf("Unsupported Format provided")
	}
}

func InterfaceToFrontMatter(in interface{}, mark rune) ([]byte, error) {
	if in == nil {
		return []byte{}, fmt.Errorf("input was nil")
	}

	b := new(bytes.Buffer)

	switch mark {
	case rune(YAML_LEAD[0]):
		_, err := b.Write([]byte(YAML_DELIM_UNIX))
		if err != nil {
			return nil, err
		}
		by, err := yaml.Marshal(in)
		if err != nil {
			return nil, err
		}
		b.Write(by)
		_, err = b.Write([]byte(YAML_DELIM_UNIX))
		if err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	case rune(TOML_LEAD[0]):
		_, err := b.Write([]byte(TOML_DELIM_UNIX))
		if err != nil {
			return nil, err
		}

		err = toml.NewEncoder(b).Encode(in)
		if err != nil {
			return nil, err
		}
		_, err = b.Write([]byte("\n" + TOML_DELIM_UNIX))
		if err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	case rune(JSON_LEAD[0]):
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
		return nil, fmt.Errorf("Unsupported Format provided")
	}
}

func FormatToLeadRune(kind string) rune {
	switch FormatSanitize(kind) {
	case "yaml":
		return rune([]byte(YAML_LEAD)[0])
	case "json":
		return rune([]byte(JSON_LEAD)[0])
	default:
		return rune([]byte(TOML_LEAD)[0])
	}
}

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

func DetectFrontMatter(mark rune) (f *FrontmatterType) {
	switch mark {
	case '-':
		return &FrontmatterType{[]byte(YAML_DELIM), []byte(YAML_DELIM), HandleYAMLMetaData, false}
	case '+':
		return &FrontmatterType{[]byte(TOML_DELIM), []byte(TOML_DELIM), HandleTOMLMetaData, false}
	case '{':
		return &FrontmatterType{[]byte{'{'}, []byte{'}'}, HandleJSONMetaData, true}
	default:
		return nil
	}
}

func HandleTOMLMetaData(datum []byte) (interface{}, error) {
	m := map[string]interface{}{}
	datum = removeTOMLIdentifier(datum)
	if _, err := toml.Decode(string(datum), &m); err != nil {
		return m, err
	}
	return m, nil
}

func removeTOMLIdentifier(datum []byte) []byte {
	return bytes.Replace(datum, []byte(TOML_DELIM), []byte(""), -1)
}

func HandleYAMLMetaData(datum []byte) (interface{}, error) {
	m := map[string]interface{}{}
	if err := yaml.Unmarshal(datum, &m); err != nil {
		return m, err
	}
	return m, nil
}

func HandleJSONMetaData(datum []byte) (interface{}, error) {
	var f interface{}
	if err := json.Unmarshal(datum, &f); err != nil {
		return f, err
	}
	return f, nil
}
