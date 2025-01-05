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

// Package resources provides template functions for working with resources.
package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/tpl/css"
	"github.com/gohugoio/hugo/tpl/js"

	"github.com/gohugoio/hugo/resources/postpub"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"

	"github.com/gohugoio/hugo/resources/resource_factories/bundler"
	"github.com/gohugoio/hugo/resources/resource_factories/create"
	"github.com/gohugoio/hugo/resources/resource_transformers/integrity"
	"github.com/gohugoio/hugo/resources/resource_transformers/minifier"
	"github.com/gohugoio/hugo/resources/resource_transformers/templates"

	"github.com/spf13/cast"
)

// New returns a new instance of the resources-namespaced template functions.
func New(deps *deps.Deps) (*Namespace, error) {
	if deps.ResourceSpec == nil {
		return &Namespace{}, nil
	}

	minifyClient, err := minifier.New(deps.ResourceSpec)
	if err != nil {
		return nil, err
	}

	return &Namespace{
		deps:            deps,
		createClient:    create.New(deps.ResourceSpec),
		bundlerClient:   bundler.New(deps.ResourceSpec),
		integrityClient: integrity.New(deps.ResourceSpec),
		minifyClient:    minifyClient,
		templatesClient: templates.New(deps.ResourceSpec, deps),
	}, nil
}

var _ resource.ResourceFinder = (*Namespace)(nil)

// Namespace provides template functions for the "resources" namespace.
type Namespace struct {
	deps *deps.Deps

	createClient    *create.Client
	bundlerClient   *bundler.Client
	integrityClient *integrity.Client
	minifyClient    *minifier.Client
	templatesClient *templates.Client

	// We moved some CSS and JS related functions to the css and js package in Hugo 0.128.0.
	// Keep this here until the deprecation period is over.
	cssNs *css.Namespace
	jsNs  *js.Namespace
}

// Copy copies r to the new targetPath in s.
func (ns *Namespace) Copy(s any, r resource.Resource) (resource.Resource, error) {
	targetPath, err := cast.ToStringE(s)
	if err != nil {
		panic(err)
	}
	return ns.createClient.Copy(r, targetPath)
}

// Get locates the filename given in Hugo's assets filesystem
// and creates a Resource object that can be used for further transformations.
func (ns *Namespace) Get(filename any) resource.Resource {
	filenamestr, err := cast.ToStringE(filename)
	if err != nil {
		panic(err)
	}

	if filenamestr == "" {
		return nil
	}

	r, err := ns.createClient.Get(filenamestr)
	if err != nil {
		panic(err)
	}

	return r
}

// GetRemote gets the URL (via HTTP(s)) in the first argument in args and creates Resource object that can be used for
// further transformations.
//
// A second argument may be provided with an option map.
//
// Note: This method does not return any error as a second return value,
// for any error situations the error can be checked in .Err.
func (ns *Namespace) GetRemote(args ...any) (resource.Resource, error) {
	get := func(args ...any) (resource.Resource, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, errors.New("must provide an URL and optionally an options map")
		}

		urlstr, err := cast.ToStringE(args[0])
		if err != nil {
			return nil, err
		}

		var options map[string]any

		if len(args) > 1 {
			options, err = maps.ToStringMapE(args[1])
			if err != nil {
				return nil, err
			}
		}

		return ns.createClient.FromRemote(urlstr, options)
	}

	r, err := get(args...)
	if err != nil {
		switch v := err.(type) {
		case *create.HTTPError:
			return nil, resource.NewResourceError(v, v.Data)
		default:
			return nil, resource.NewResourceError(err, nil)
		}
	}
	return r, nil
}

// GetMatch finds the first Resource matching the given pattern, or nil if none found.
//
// It looks for files in the assets file system.
//
// See Match for a more complete explanation about the rules used.
func (ns *Namespace) GetMatch(pattern any) resource.Resource {
	patternStr, err := cast.ToStringE(pattern)
	if err != nil {
		panic(err)
	}

	r, err := ns.createClient.GetMatch(patternStr)
	if err != nil {
		panic(err)
	}

	return r
}

// ByType returns resources of a given resource type (e.g. "image").
func (ns *Namespace) ByType(typ any) resource.Resources {
	return ns.createClient.ByType(cast.ToString(typ))
}

// Match gets all resources matching the given base path prefix, e.g
// "*.png" will match all png files. The "*" does not match path delimiters (/),
// so if you organize your resources in sub-folders, you need to be explicit about it, e.g.:
// "images/*.png". To match any PNG image anywhere in the bundle you can do "**.png", and
// to match all PNG images below the images folder, use "images/**.jpg".
//
// The matching is case insensitive.
//
// Match matches by using the files name with path relative to the file system root
// with Unix style slashes (/) and no leading slash, e.g. "images/logo.png".
//
// See https://github.com/gobwas/glob for the full rules set.
//
// It looks for files in the assets file system.
//
// See Match for a more complete explanation about the rules used.
func (ns *Namespace) Match(pattern any) resource.Resources {
	patternStr, err := cast.ToStringE(pattern)
	if err != nil {
		panic(err)
	}

	r, err := ns.createClient.Match(patternStr)
	if err != nil {
		panic(err)
	}

	return r
}

// Concat concatenates a slice of Resource objects. These resources must
// (currently) be of the same Media Type.
func (ns *Namespace) Concat(targetPathIn any, r any) (resource.Resource, error) {
	targetPath, err := cast.ToStringE(targetPathIn)
	if err != nil {
		return nil, err
	}

	var rr resource.Resources

	switch v := r.(type) {
	case resource.Resources:
		rr = v
	case resource.ResourcesConverter:
		rr = v.ToResources()
	default:
		return nil, fmt.Errorf("expected slice of Resource objects, received %T instead", r)
	}

	if len(rr) == 0 {
		return nil, errors.New("must provide one or more Resource objects to concat")
	}

	return ns.bundlerClient.Concat(targetPath, rr)
}

// FromString creates a Resource from a string published to the relative target path.
func (ns *Namespace) FromString(targetPathIn, contentIn any) (resource.Resource, error) {
	targetPath, err := cast.ToStringE(targetPathIn)
	if err != nil {
		return nil, err
	}
	content, err := cast.ToStringE(contentIn)
	if err != nil {
		return nil, err
	}

	return ns.createClient.FromString(targetPath, content)
}

// ExecuteAsTemplate creates a Resource from a Go template, parsed and executed with
// the given data, and published to the relative target path.
func (ns *Namespace) ExecuteAsTemplate(ctx context.Context, args ...any) (resource.Resource, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("must provide targetPath, the template data context and a Resource object")
	}
	targetPath, err := cast.ToStringE(args[0])
	if err != nil {
		return nil, err
	}
	data := args[1]

	r, ok := args[2].(resources.ResourceTransformer)
	if !ok {
		return nil, fmt.Errorf("type %T not supported in Resource transformations", args[2])
	}

	return ns.templatesClient.ExecuteAsTemplate(ctx, r, targetPath, data)
}

// Fingerprint transforms the given Resource with a MD5 hash of the content in
// the RelPermalink and Permalink.
func (ns *Namespace) Fingerprint(args ...any) (resource.Resource, error) {
	if len(args) < 1 {
		return nil, errors.New("must provide a Resource object")
	}

	if len(args) > 2 {
		return nil, errors.New("must not provide more arguments than Resource and hash algorithm")
	}

	var algo string
	resIdx := 0

	if len(args) == 2 {
		resIdx = 1
		var err error
		algo, err = cast.ToStringE(args[0])
		if err != nil {
			return nil, err
		}
	}

	r, ok := args[resIdx].(resources.ResourceTransformer)
	if !ok {
		return nil, fmt.Errorf("%T can not be transformed", args[resIdx])
	}

	return ns.integrityClient.Fingerprint(r, algo)
}

// Minify minifies the given Resource using the MediaType to pick the correct
// minifier.
func (ns *Namespace) Minify(r resources.ResourceTransformer) (resource.Resource, error) {
	return ns.minifyClient.Minify(r)
}

// ToCSS converts the given Resource to CSS. You can optional provide an Options object
// as second argument. As an option, you can e.g. specify e.g. the target path (string)
// for the converted CSS resource.
// Deprecated: Moved to the css namespace in Hugo 0.128.0.
func (ns *Namespace) ToCSS(args ...any) (resource.Resource, error) {
	hugo.Deprecate("resources.ToCSS", "Use css.Sass instead.", "v0.128.0")
	return ns.cssNs.Sass(args...)
}

// PostCSS processes the given Resource with PostCSS.
// Deprecated: Moved to the css namespace in Hugo 0.128.0.
func (ns *Namespace) PostCSS(args ...any) (resource.Resource, error) {
	hugo.Deprecate("resources.PostCSS", "Use css.PostCSS instead.", "v0.128.0")
	return ns.cssNs.PostCSS(args...)
}

// PostProcess processes r after the build.
func (ns *Namespace) PostProcess(r resource.Resource) (postpub.PostPublishedResource, error) {
	return ns.deps.ResourceSpec.PostProcess(r)
}

// Babel processes the given Resource with Babel.
// Deprecated: Moved to the js namespace in Hugo 0.128.0.
func (ns *Namespace) Babel(args ...any) (resource.Resource, error) {
	hugo.Deprecate("resources.Babel", "Use js.Babel.", "v0.128.0")
	return ns.jsNs.Babel(args...)
}
