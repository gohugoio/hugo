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

// Package openapi3 provides functions for generating OpenAPI v3 (Swagger) documentation.
package openapi3

import (
	"errors"
	"fmt"
	"io"

	gyaml "github.com/ghodss/yaml"

	kopenapi3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/resources/resource"
)

// New returns a new instance of the openapi3-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		cache: dynacache.GetOrCreatePartition[string, *OpenAPIDocument](deps.MemCache, "/tmpl/openapi3", dynacache.OptionsPartition{Weight: 30, ClearWhen: dynacache.ClearOnChange}),
		deps:  deps,
	}
}

// Namespace provides template functions for the "openapi3".
type Namespace struct {
	cache *dynacache.Partition[string, *OpenAPIDocument]
	deps  *deps.Deps
}

// OpenAPIDocument represents an OpenAPI 3 document.
type OpenAPIDocument struct {
	*kopenapi3.T
	identityGroup identity.Identity
}

func (o *OpenAPIDocument) GetIdentityGroup() identity.Identity {
	return o.identityGroup
}

// Unmarshal unmarshals the given resource into an OpenAPI 3 document.
func (ns *Namespace) Unmarshal(r resource.UnmarshableResource) (*OpenAPIDocument, error) {
	key := r.Key()
	if key == "" {
		return nil, errors.New("no Key set in Resource")
	}

	v, err := ns.cache.GetOrCreate(key, func(string) (*OpenAPIDocument, error) {
		f := metadecoders.FormatFromStrings(r.MediaType().Suffixes()...)
		if f == "" {
			return nil, fmt.Errorf("MIME %q not supported", r.MediaType())
		}

		reader, err := r.ReadSeekCloser()
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		b, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		s := &kopenapi3.T{}
		switch f {
		case metadecoders.YAML:
			err = gyaml.Unmarshal(b, s)
		default:
			err = metadecoders.Default.UnmarshalTo(b, f, s)
		}
		if err != nil {
			return nil, err
		}

		err = kopenapi3.NewLoader().ResolveRefsIn(s, nil)

		return &OpenAPIDocument{T: s, identityGroup: identity.FirstIdentity(r)}, err
	})
	if err != nil {
		return nil, err
	}

	return v, nil
}
