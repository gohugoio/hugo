// Copyright 2019 The Hugo Authors. All rights reserved.
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

package transform

import (
	"io/ioutil"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/resources/resource"
	"github.com/pkg/errors"

	"github.com/spf13/cast"
)

// Unmarshal unmarshals the data given, which can be either a string
// or a Resource. Supported formats are JSON, TOML, YAML, and CSV.
// You can optionally provide an options map as the first argument.
func (ns *Namespace) Unmarshal(args ...interface{}) (interface{}, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("unmarshal takes 1 or 2 arguments")
	}

	var data interface{}
	var decoder = metadecoders.Default

	if len(args) == 1 {
		data = args[0]
	} else {
		m, ok := args[0].(map[string]interface{})
		if !ok {
			return nil, errors.New("first argument must be a map")
		}

		var err error

		data = args[1]
		decoder, err = decodeDecoder(m)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to decode options")
		}
	}

	if r, ok := data.(unmarshableResource); ok {
		key := r.Key()

		if key == "" {
			return nil, errors.New("no Key set in Resource")
		}

		if decoder != metadecoders.Default {
			key += decoder.OptionsKey()
		}

		return ns.cache.GetOrCreate(key, func() (interface{}, error) {
			f := metadecoders.FormatFromMediaType(r.MediaType())
			if f == "" {
				return nil, errors.Errorf("MIME %q not supported", r.MediaType())
			}

			reader, err := r.ReadSeekCloser()
			if err != nil {
				return nil, err
			}
			defer reader.Close()

			b, err := ioutil.ReadAll(reader)
			if err != nil {
				return nil, err
			}

			return decoder.Unmarshal(b, f)
		})
	}

	dataStr, err := cast.ToStringE(data)
	if err != nil {
		return nil, errors.Errorf("type %T not supported", data)
	}

	key := helpers.MD5String(dataStr)

	return ns.cache.GetOrCreate(key, func() (interface{}, error) {
		f := decoder.FormatFromContentString(dataStr)
		if f == "" {
			return nil, errors.New("unknown format")
		}

		return decoder.Unmarshal([]byte(dataStr), f)
	})
}

// All the relevant resources implements this interface.
type unmarshableResource interface {
	resource.ReadSeekCloserResource
	resource.Identifier
}

func decodeDecoder(m map[string]interface{}) (metadecoders.Decoder, error) {
	opts := metadecoders.Default

	if m == nil {
		return opts, nil
	}

	// mapstructure does not support string to rune conversion, so do that manually.
	// See https://github.com/mitchellh/mapstructure/issues/151
	for k, v := range m {
		if strings.EqualFold(k, "Delimiter") {
			r, err := stringToRune(v)
			if err != nil {
				return opts, err
			}
			opts.Delimiter = r
			delete(m, k)

		} else if strings.EqualFold(k, "Comment") {
			r, err := stringToRune(v)
			if err != nil {
				return opts, err
			}
			opts.Comment = r
			delete(m, k)
		}
	}

	err := mapstructure.WeakDecode(m, &opts)

	return opts, err
}

func stringToRune(v interface{}) (rune, error) {
	s, err := cast.ToStringE(v)
	if err != nil {
		return 0, err
	}

	if len(s) == 0 {
		return 0, nil
	}

	var r rune

	for i, rr := range s {
		if i == 0 {
			r = rr
		} else {
			return 0, errors.Errorf("invalid character: %q", v)
		}
	}

	return r, nil
}
