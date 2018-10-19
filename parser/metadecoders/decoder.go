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
	"encoding/json"

	"github.com/BurntSushi/toml"
	"github.com/chaseadamsio/goorgeous"
	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
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

// UnmarshalToMap will unmarshall data in format f into a new map. This is
// what's needed for Hugo's front matter decoding.
func UnmarshalToMap(data []byte, f Format) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	if data == nil {
		return m, nil
	}

	var err error

	switch f {
	case ORG:
		m, err = goorgeous.OrgHeaders(data)
	case JSON:
		err = json.Unmarshal(data, &m)
	case TOML:
		_, err = toml.Decode(string(data), &m)
	case YAML:
		err = yaml.Unmarshal(data, &m)

		// To support boolean keys, the `yaml` package unmarshals maps to
		// map[interface{}]interface{}. Here we recurse through the result
		// and change all maps to map[string]interface{} like we would've
		// gotten from `json`.
		if err == nil {
			for k, v := range m {
				if vv, changed := stringifyMapKeys(v); changed {
					m[k] = vv
				}
			}
		}
	default:
		return nil, errors.Errorf("unmarshal of format %q is not supported", f)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal failed for format %q", f)
	}

	return m, nil

}
