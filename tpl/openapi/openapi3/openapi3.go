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
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	kopenapi3 "github.com/getkin/kin-openapi/openapi3"
	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	resourcestpl "github.com/gohugoio/hugo/tpl/resources"
	"github.com/mitchellh/mapstructure"
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
	cache       *dynacache.Partition[string, *OpenAPIDocument]
	deps        *deps.Deps
	resourcesNs *resourcestpl.Namespace
}

// OpenAPIDocument represents an OpenAPI 3 document.
type OpenAPIDocument struct {
	*kopenapi3.T
	identityGroup identity.Identity
}

func (o *OpenAPIDocument) GetIdentityGroup() identity.Identity {
	return o.identityGroup
}

type unmarshalOptions struct {
	// Options passed to resources.GetRemote when resolving remote $ref.
	GetRemote map[string]any
}

// Unmarshal unmarshals the given resource into an OpenAPI 3 document.
func (ns *Namespace) Unmarshal(ctx context.Context, args ...any) (*OpenAPIDocument, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("must provide a Resource and optionally an options map")
	}

	r := args[0].(resource.UnmarshableResource)
	key := r.Key()
	if key == "" {
		return nil, errors.New("no Key set in Resource")
	}

	var opts unmarshalOptions
	if len(args) > 1 {
		optsm, err := maps.ToStringMapE(args[1])
		if err != nil {
			return nil, err
		}
		if err := mapstructure.WeakDecode(optsm, &opts); err != nil {
			return nil, err
		}
		key += "_" + hashing.HashString(optsm)
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
			err = metadecoders.UnmarshalYaml(b, s)
		default:
			err = metadecoders.Default.UnmarshalTo(b, f, s)
		}
		if err != nil {
			return nil, err
		}

		var resourcePath string
		if res, ok := r.(resource.Resource); ok {
			resourcePath = resources.InternalResourceSourcePath(res)
		}
		var relDir string
		if resourcePath != "" {
			if rel, ok := ns.deps.Assets.MakePathRelative(filepath.FromSlash(resourcePath), true); ok {
				relDir = filepath.Dir(rel)
			}
		}

		var idm identity.Manager = identity.NopManager
		if v := identity.GetDependencyManager(r); v != nil {
			idm = v
		}
		idg := identity.FirstIdentity(r)

		resolver := &refResolver{
			ctx:     ctx,
			idm:     idm,
			opts:    opts,
			relBase: filepath.ToSlash(relDir),
			ns:      ns,
		}

		loader := kopenapi3.NewLoader()
		loader.IsExternalRefsAllowed = true
		loader.ReadFromURIFunc = resolver.resolveExternalRef
		err = loader.ResolveRefsIn(s, nil)

		return &OpenAPIDocument{T: s, identityGroup: idg}, err
	})
	if err != nil {
		return nil, err
	}

	return v, nil
}

type refResolver struct {
	ctx     context.Context
	idm     identity.Manager
	opts    unmarshalOptions
	relBase string
	ns      *Namespace
}

// resolveExternalRef resolves external references in OpenAPI documents by either fetching
// remote URLs or loading local files from the assets directory, depending on the reference location.
func (r *refResolver) resolveExternalRef(loader *kopenapi3.Loader, loc *url.URL) ([]byte, error) {
	if loc.Scheme != "" && loc.Host != "" {
		res, err := r.ns.resourcesNs.GetRemote(loc.String(), r.opts.GetRemote)
		if err != nil {
			return nil, fmt.Errorf("failed to get remote ref %q: %w", loc.String(), err)
		}
		content, err := resources.InternalResourceSourceContent(r.ctx, res)
		if err != nil {
			return nil, fmt.Errorf("failed to read remote ref %q: %w", loc.String(), err)
		}
		r.idm.AddIdentity(identity.FirstIdentity(res))
		return []byte(content), nil
	}

	var filename string
	if strings.HasPrefix(loc.Path, "/") {
		filename = loc.Path
	} else {
		filename = path.Join(r.relBase, loc.Path)
	}

	res := r.ns.resourcesNs.Get(filename)
	if res == nil {
		return nil, fmt.Errorf("local ref %q not found", loc.String())
	}
	content, err := resources.InternalResourceSourceContent(r.ctx, res)
	if err != nil {
		return nil, fmt.Errorf("failed to read local ref %q: %w", loc.String(), err)
	}
	r.idm.AddIdentity(identity.FirstIdentity(res))
	return []byte(content), nil
}
