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
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/niklasfasching/go-org/org"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
	yaml "gopkg.in/yaml.v2"
)

// Decoder provides some configuration options for the decoders.
type Decoder struct {
	// Delimiter is the field delimiter used in the CSV decoder. It defaults to ','.
	Delimiter rune

	// Comment, if not 0, is the comment character ued in the CSV decoder. Lines beginning with the
	// Comment character without preceding whitespace are ignored.
	Comment rune
}

// OptionsKey is used in cache keys.
func (d Decoder) OptionsKey() string {
	var sb strings.Builder
	sb.WriteRune(d.Delimiter)
	sb.WriteRune(d.Comment)
	return sb.String()
}

// Default is a Decoder in its default configuration.
var Default = Decoder{
	Delimiter: ',',
}

// UnmarshalToMap will unmarshall data in format f into a new map. This is
// what's needed for Hugo's front matter decoding.
func (d Decoder) UnmarshalToMap(data []byte, f Format) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if data == nil {
		return m, nil
	}

	err := d.UnmarshalTo(data, f, &m)

	return m, err
}

// UnmarshalFileToMap is the same as UnmarshalToMap, but reads the data from
// the given filename.
func (d Decoder) UnmarshalFileToMap(fs afero.Fs, filename string) (map[string]interface{}, error) {
	format := FormatFromString(filename)
	if format == "" {
		return nil, errors.Errorf("%q is not a valid configuration format", filename)
	}

	data, err := afero.ReadFile(fs, filename)
	if err != nil {
		return nil, err
	}
	return d.UnmarshalToMap(data, format)
}

// UnmarshalStringTo tries to unmarshal data to a new instance of type typ.
func (d Decoder) UnmarshalStringTo(data string, typ interface{}) (interface{}, error) {
	data = strings.TrimSpace(data)
	// We only check for the possible types in YAML, JSON and TOML.
	switch typ.(type) {
	case string:
		return data, nil
	case map[string]interface{}:
		format := d.FormatFromContentString(data)
		return d.UnmarshalToMap([]byte(data), format)
	case []interface{}:
		// A standalone slice. Let YAML handle it.
		return d.Unmarshal([]byte(data), YAML)
	case bool:
		return cast.ToBoolE(data)
	case int:
		return cast.ToIntE(data)
	case int64:
		return cast.ToInt64E(data)
	case float64:
		return cast.ToFloat64E(data)
	default:
		return nil, errors.Errorf("unmarshal: %T not supportedd", typ)
	}
}

// Unmarshal will unmarshall data in format f into an interface{}.
// This is what's needed for Hugo's /data handling.
func (d Decoder) Unmarshal(data []byte, f Format) (interface{}, error) {
	if data == nil {
		switch f {
		case CSV:
			return make([][]string, 0), nil
		default:
			return make(map[string]interface{}), nil
		}

	}
	var v interface{}
	err := d.UnmarshalTo(data, f, &v)

	return v, err
}

// UnmarshalTo unmarshals data in format f into v.
func (d Decoder) UnmarshalTo(data []byte, f Format, v interface{}) error {

	var err error

	switch f {
	case ORG:
		err = d.unmarshalORG(data, v)
	case JSON:
		err = json.Unmarshal(data, v)
	case TOML:
		err = toml.Unmarshal(data, v)
	case YAML:
		err = yaml.Unmarshal(data, v)
		if err != nil {
			return toFileError(f, errors.Wrap(err, "failed to unmarshal YAML"))
		}

		// To support boolean keys, the YAML package unmarshals maps to
		// map[interface{}]interface{}. Here we recurse through the result
		// and change all maps to map[string]interface{} like we would've
		// gotten from `json`.
		var ptr interface{}
		switch v.(type) {
		case *map[string]interface{}:
			ptr = *v.(*map[string]interface{})
		case *interface{}:
			ptr = *v.(*interface{})
		default:
			// Not a map.
		}

		if ptr != nil {
			if mm, changed := stringifyMapKeys(ptr); changed {
				switch v.(type) {
				case *map[string]interface{}:
					*v.(*map[string]interface{}) = mm.(map[string]interface{})
				case *interface{}:
					*v.(*interface{}) = mm
				}
			}
		}
	case CSV:
		return d.unmarshalCSV(data, v)

	default:
		return errors.Errorf("unmarshal of format %q is not supported", f)
	}

	if err == nil {
		return nil
	}

	return toFileError(f, errors.Wrap(err, "unmarshal failed"))

}

func (d Decoder) unmarshalCSV(data []byte, v interface{}) error {
	r := csv.NewReader(bytes.NewReader(data))
	r.Comma = d.Delimiter
	r.Comment = d.Comment

	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	switch v.(type) {
	case *interface{}:
		*v.(*interface{}) = records
	default:
		return errors.Errorf("CSV cannot be unmarshaled into %T", v)

	}

	return nil

}

func parseORGDate(s string) string {
	r := regexp.MustCompile(`[<\[](\d{4}-\d{2}-\d{2}) .*[>\]]`)
	if m := r.FindStringSubmatch(s); m != nil {
		return m[1]
	}
	return s
}

func (d Decoder) unmarshalORG(data []byte, v interface{}) error {
	config := org.New()
	config.Log = jww.WARN
	document := config.Parse(bytes.NewReader(data), "")
	if document.Error != nil {
		return document.Error
	}
	frontMatter := make(map[string]interface{}, len(document.BufferSettings))
	for k, v := range document.BufferSettings {
		k = strings.ToLower(k)
		if strings.HasSuffix(k, "[]") {
			frontMatter[k[:len(k)-2]] = strings.Fields(v)
		} else if k == "tags" || k == "categories" || k == "aliases" {
			jww.WARN.Printf("Please use '#+%s[]:' notation, automatic conversion is deprecated.", k)
			frontMatter[k] = strings.Fields(v)
		} else if k == "date" {
			frontMatter[k] = parseORGDate(v)
		} else {
			frontMatter[k] = v
		}
	}
	switch v.(type) {
	case *map[string]interface{}:
		*v.(*map[string]interface{}) = frontMatter
	default:
		*v.(*interface{}) = frontMatter
	}
	return nil
}

func toFileError(f Format, err error) error {
	return herrors.ToFileError(string(f), err)
}

// stringifyMapKeys recurses into in and changes all instances of
// map[interface{}]interface{} to map[string]interface{}. This is useful to
// work around the impedance mismatch between JSON and YAML unmarshaling that's
// described here: https://github.com/go-yaml/yaml/issues/139
//
// Inspired by https://github.com/stripe/stripe-mock, MIT licensed
func stringifyMapKeys(in interface{}) (interface{}, bool) {

	switch in := in.(type) {
	case []interface{}:
		for i, v := range in {
			if vv, replaced := stringifyMapKeys(v); replaced {
				in[i] = vv
			}
		}
	case map[string]interface{}:
		for k, v := range in {
			if vv, changed := stringifyMapKeys(v); changed {
				in[k] = vv
			}
		}
	case map[interface{}]interface{}:
		res := make(map[string]interface{})
		var (
			ok  bool
			err error
		)
		for k, v := range in {
			var ks string

			if ks, ok = k.(string); !ok {
				ks, err = cast.ToStringE(k)
				if err != nil {
					ks = fmt.Sprintf("%v", k)
				}
			}
			if vv, replaced := stringifyMapKeys(v); replaced {
				res[ks] = vv
			} else {
				res[ks] = v
			}
		}
		return res, true
	}

	return nil, false
}
