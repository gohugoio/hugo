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

	"github.com/BurntSushi/toml"
	"launchpad.net/goyaml"
)

type FrontmatterType struct {
	markstart, markend []byte
	Parse              func([]byte) (interface{}, error)
	includeMark        bool
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
		by, err := goyaml.Marshal(in)
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
			fmt.Println("toml encoder failed", in)
			fmt.Println(err)
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
			fmt.Println("json encoder failed", in)
			fmt.Println(err)
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

func DetectFrontMatter(mark rune) (f *FrontmatterType) {
	switch mark {
	case '-':
		return &FrontmatterType{[]byte(YAML_DELIM), []byte(YAML_DELIM), HandleYamlMetaData, false}
	case '+':
		return &FrontmatterType{[]byte(TOML_DELIM), []byte(TOML_DELIM), HandleTomlMetaData, false}
	case '{':
		return &FrontmatterType{[]byte{'{'}, []byte{'}'}, HandleJsonMetaData, true}
	default:
		return nil
	}
}

func HandleTomlMetaData(datum []byte) (interface{}, error) {
	m := map[string]interface{}{}
	datum = removeTomlIdentifier(datum)
	if _, err := toml.Decode(string(datum), &m); err != nil {
		return m, err
	}
	return m, nil
}

func removeTomlIdentifier(datum []byte) []byte {
	return bytes.Replace(datum, []byte(TOML_DELIM), []byte(""), -1)
}

func HandleYamlMetaData(datum []byte) (interface{}, error) {
	m := map[string]interface{}{}
	if err := goyaml.Unmarshal(datum, &m); err != nil {
		return m, err
	}
	return m, nil
}

func HandleJsonMetaData(datum []byte) (interface{}, error) {
	var f interface{}
	if err := json.Unmarshal(datum, &f); err != nil {
		return f, err
	}
	return f, nil
}
