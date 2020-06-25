// Copyright 2020 The Hugo Authors. All rights reserved.
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

package openapi3

import (
	"io/ioutil"

	gyaml "github.com/ghodss/yaml"

	"github.com/pkg/errors"

	kopenapi3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/gohugoio/hugo/cache/memcache"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/resources/resource"
)

// New returns a new instance of the openapi3-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		cache: deps.MemCache.GetOrCreatePartition("tpl/openapi3", memcache.ClearOnChange),
		deps:  deps,
	}
}

// Namespace provides template functions for the "openapi3".
type Namespace struct {
	cache memcache.Getter
	deps  *deps.Deps
}

func (ns *Namespace) Unmarshal(r resource.UnmarshableResource) (*kopenapi3.Swagger, error) {
	key := r.Key()
	if key == "" {
		return nil, errors.New("no Key set in Resource")
	}

	v, err := ns.cache.GetOrCreate(key, func() memcache.Entry {
		f := metadecoders.FormatFromMediaType(r.MediaType())
		if f == "" {
			return memcache.Entry{Err: errors.Errorf("MIME %q not supported", r.MediaType())}
		}

		reader, err := r.ReadSeekCloser()
		if err != nil {
			return memcache.Entry{Err: err}
		}

		defer reader.Close()

		b, err := ioutil.ReadAll(reader)
		if err != nil {
			return memcache.Entry{Err: err}
		}

		s := &kopenapi3.Swagger{}
		switch f {
		case metadecoders.YAML:
			err = gyaml.Unmarshal(b, s)
		default:
			err = metadecoders.Default.UnmarshalTo(b, f, s)
		}
		if err != nil {
			return memcache.Entry{Err: err}
		}

		err = kopenapi3.NewSwaggerLoader().ResolveRefsIn(s, nil)

		return memcache.Entry{
			Value:     s,
			Err:       err,
			ClearWhen: memcache.ClearOnChange,
			StaleFunc: func() bool {
				return resource.IsStaleAny(r)
			},
		}
	})
	if err != nil {
		return nil, err
	}

	return v.(*kopenapi3.Swagger), nil
}
