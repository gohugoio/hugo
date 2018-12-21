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

package transform

import (
	"io/ioutil"

	"github.com/gohugoio/hugo/common/hugio"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/resource"
	"github.com/pkg/errors"

	"github.com/spf13/cast"
)

// Unmarshal unmarshals the data given, which can be either a string
// or a Resource. Supported formats are JSON, TOML and YAML.
func (ns *Namespace) Unmarshal(data interface{}) (interface{}, error) {

	// All the relevant Resource types implements ReadSeekCloserResource,
	// which should be the most effective way to get the content.
	if r, ok := data.(resource.ReadSeekCloserResource); ok {
		var key string
		var reader hugio.ReadSeekCloser

		if k, ok := r.(resource.Identifier); ok {
			key = k.Key()
		}

		if key == "" {
			reader, err := r.ReadSeekCloser()
			if err != nil {
				return nil, err
			}
			defer reader.Close()

			key, err = helpers.MD5FromReader(reader)
			if err != nil {
				return nil, err
			}

			reader.Seek(0, 0)
		}

		return ns.cache.GetOrCreate(key, func() (interface{}, error) {
			f := metadecoders.FormatFromMediaType(r.MediaType())
			if f == "" {
				return nil, errors.Errorf("MIME %q not supported", r.MediaType())
			}

			if reader == nil {
				var err error
				reader, err = r.ReadSeekCloser()
				if err != nil {
					return nil, err
				}
				defer reader.Close()
			}

			b, err := ioutil.ReadAll(reader)
			if err != nil {
				return nil, err
			}

			return metadecoders.Unmarshal(b, f)
		})

	}

	dataStr, err := cast.ToStringE(data)
	if err != nil {
		return nil, errors.Errorf("type %T not supported", data)
	}

	key := helpers.MD5String(dataStr)

	return ns.cache.GetOrCreate(key, func() (interface{}, error) {
		f := metadecoders.FormatFromContentString(dataStr)
		if f == "" {
			return nil, errors.New("unknown format")
		}

		return metadecoders.Unmarshal([]byte(dataStr), f)
	})
}
