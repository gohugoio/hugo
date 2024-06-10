// Copyright 2024 The Hugo Authors. All rights reserved.
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

package output

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/media"
	"github.com/mitchellh/mapstructure"
)

// OutputFormatConfig configures a single output format.
type OutputFormatConfig struct {
	// The MediaType string. This must be a configured media type.
	MediaType string
	Format
}

var defaultOutputFormat = Format{
	BaseName: "index",
	Rel:      "alternate",
}

func DecodeConfig(mediaTypes media.Types, in any) (*config.ConfigNamespace[map[string]OutputFormatConfig, Formats], error) {
	buildConfig := func(in any) (Formats, any, error) {
		f := make(Formats, len(DefaultFormats))
		copy(f, DefaultFormats)
		if in != nil {
			m, err := maps.ToStringMapE(in)
			if err != nil {
				return nil, nil, fmt.Errorf("failed convert config to map: %s", err)
			}
			m = maps.CleanConfigStringMap(m)

			for k, v := range m {
				found := false
				for i, vv := range f {
					// Both are lower case.
					if k == vv.Name {
						// Merge it with the existing
						if err := decode(mediaTypes, v, &f[i]); err != nil {
							return f, nil, err
						}
						found = true
					}
				}
				if found {
					continue
				}

				newOutFormat := defaultOutputFormat
				if err := decode(mediaTypes, v, &newOutFormat); err != nil {
					return f, nil, err
				}
				newOutFormat.Name = k

				f = append(f, newOutFormat)

			}
		}

		// Also format is a map for documentation purposes.
		docm := make(map[string]OutputFormatConfig, len(f))
		for _, ff := range f {
			docm[ff.Name] = OutputFormatConfig{
				MediaType: ff.MediaType.Type,
				Format:    ff,
			}
		}

		sort.Sort(f)
		return f, docm, nil
	}

	return config.DecodeNamespace[map[string]OutputFormatConfig](in, buildConfig)
}

func decode(mediaTypes media.Types, input any, output *Format) error {
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: func(a reflect.Type, b reflect.Type, c any) (any, error) {
			if a.Kind() == reflect.Map {
				dataVal := reflect.Indirect(reflect.ValueOf(c))
				for _, key := range dataVal.MapKeys() {
					keyStr, ok := key.Interface().(string)
					if !ok {
						// Not a string key
						continue
					}
					if strings.EqualFold(keyStr, "mediaType") {
						// If mediaType is a string, look it up and replace it
						// in the map.
						vv := dataVal.MapIndex(key)
						vvi := vv.Interface()

						switch vviv := vvi.(type) {
						case media.Type:
						// OK
						case string:
							mediaType, found := mediaTypes.GetByType(vviv)
							if !found {
								return c, fmt.Errorf("media type %q not found", vviv)
							}
							dataVal.SetMapIndex(key, reflect.ValueOf(mediaType))
						default:
							return nil, fmt.Errorf("invalid output format configuration; wrong type for media type, expected string (e.g. text/html), got %T", vvi)
						}
					}
				}
			}
			return c, nil
		},
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	if err = decoder.Decode(input); err != nil {
		return fmt.Errorf("failed to decode output format configuration: %w", err)
	}

	return nil
}
