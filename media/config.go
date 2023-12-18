// Copyright 2023 The Hugo Authors. All rights reserved.
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

package media

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

// DefaultTypes is the default media types supported by Hugo.
var DefaultTypes Types

func init() {
	ns, err := DecodeTypes(nil)
	if err != nil {
		panic(err)
	}
	DefaultTypes = ns.Config

	// Initialize the Builtin types with values from DefaultTypes.
	v := reflect.ValueOf(&Builtin).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		builtinType := f.Interface().(Type)
		defaultType, found := DefaultTypes.GetByType(builtinType.Type)
		if !found {
			panic(errors.New("missing default type for builtin type: " + builtinType.Type))
		}
		f.Set(reflect.ValueOf(defaultType))
	}
}

// Hold the configuration for a given media type.
type MediaTypeConfig struct {
	// The file suffixes used for this media type.
	Suffixes []string
	// Delimiter used before suffix.
	Delimiter string
}

// DecodeTypes decodes the given map of media types.
func DecodeTypes(in map[string]any) (*config.ConfigNamespace[map[string]MediaTypeConfig, Types], error) {
	buildConfig := func(v any) (Types, any, error) {
		m, err := maps.ToStringMapE(v)
		if err != nil {
			return nil, nil, err
		}
		if m == nil {
			m = map[string]any{}
		}
		m = maps.CleanConfigStringMap(m)
		// Merge with defaults.
		maps.MergeShallow(m, defaultMediaTypesConfig)

		var types Types

		for k, v := range m {
			mediaType, err := FromString(k)
			if err != nil {
				return nil, nil, err
			}
			if err := mapstructure.WeakDecode(v, &mediaType); err != nil {
				return nil, nil, err
			}
			mm := maps.ToStringMap(v)
			suffixes, found := maps.LookupEqualFold(mm, "suffixes")
			if found {
				mediaType.SuffixesCSV = strings.TrimSpace(strings.ToLower(strings.Join(cast.ToStringSlice(suffixes), ",")))
			}
			if mediaType.SuffixesCSV != "" && mediaType.Delimiter == "" {
				mediaType.Delimiter = DefaultDelimiter
			}
			InitMediaType(&mediaType)
			types = append(types, mediaType)
		}

		sort.Sort(types)

		return types, m, nil
	}

	ns, err := config.DecodeNamespace[map[string]MediaTypeConfig](in, buildConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode media types: %w", err)
	}
	return ns, nil
}
