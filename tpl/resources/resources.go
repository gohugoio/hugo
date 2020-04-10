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
	"errors"
	"fmt"
	"path/filepath"

	"github.com/gohugoio/hugo/resources/postpub"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/resource"
	_errors "github.com/pkg/errors"

	"github.com/gohugoio/hugo/resources/resource_factories/bundler"
	"github.com/gohugoio/hugo/resources/resource_factories/create"
	"github.com/gohugoio/hugo/resources/resource_transformers/integrity"
	"github.com/gohugoio/hugo/resources/resource_transformers/minifier"
	"github.com/gohugoio/hugo/resources/resource_transformers/postcss"
	"github.com/gohugoio/hugo/resources/resource_transformers/templates"
	"github.com/gohugoio/hugo/resources/resource_transformers/tocss/scss"
	"github.com/spf13/cast"
)

// New returns a new instance of the resources-namespaced template functions.
func New(deps *deps.Deps) (*Namespace, error) {
	if deps.ResourceSpec == nil {
		return &Namespace{}, nil
	}

	scssClient, err := scss.New(deps.BaseFs.Assets, deps.ResourceSpec)
	if err != nil {
		return nil, err
	}

	minifyClient, err := minifier.New(deps.ResourceSpec)
	if err != nil {
		return nil, err
	}

	return &Namespace{
		deps:            deps,
		scssClient:      scssClient,
		createClient:    create.New(deps.ResourceSpec),
		bundlerClient:   bundler.New(deps.ResourceSpec),
		integrityClient: integrity.New(deps.ResourceSpec),
		minifyClient:    minifyClient,
		postcssClient:   postcss.New(deps.ResourceSpec),
		templatesClient: templates.New(deps.ResourceSpec, deps),
	}, nil
}

// Namespace provides template functions for the "resources" namespace.
type Namespace struct {
	deps *deps.Deps

	createClient    *create.Client
	bundlerClient   *bundler.Client
	scssClient      *scss.Client
	integrityClient *integrity.Client
	minifyClient    *minifier.Client
	postcssClient   *postcss.Client
	templatesClient *templates.Client
}

// Get locates the filename given in Hugo's assets filesystem
// and creates a Resource object that can be used for further transformations.
func (ns *Namespace) Get(filename interface{}) (resource.Resource, error) {
	filenamestr, err := cast.ToStringE(filename)
	if err != nil {
		return nil, err
	}

	filenamestr = filepath.Clean(filenamestr)

	return ns.createClient.Get(filenamestr)

}

// GetMatch finds the first Resource matching the given pattern, or nil if none found.
//
// It looks for files in the assets file system.
//
// See Match for a more complete explanation about the rules used.
func (ns *Namespace) GetMatch(pattern interface{}) (resource.Resource, error) {
	patternStr, err := cast.ToStringE(pattern)
	if err != nil {
		return nil, err
	}

	return ns.createClient.GetMatch(patternStr)

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
func (ns *Namespace) Match(pattern interface{}) (resource.Resources, error) {
	patternStr, err := cast.ToStringE(pattern)
	if err != nil {
		return nil, err
	}

	return ns.createClient.Match(patternStr)
}

// Concat concatenates a slice of Resource objects. These resources must
// (currently) be of the same Media Type.
func (ns *Namespace) Concat(targetPathIn interface{}, r interface{}) (resource.Resource, error) {
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
		return nil, fmt.Errorf("slice %T not supported in concat", r)
	}

	if len(rr) == 0 {
		return nil, errors.New("must provide one or more Resource objects to concat")
	}

	return ns.bundlerClient.Concat(targetPath, rr)
}

// FromString creates a Resource from a string published to the relative target path.
func (ns *Namespace) FromString(targetPathIn, contentIn interface{}) (resource.Resource, error) {
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
func (ns *Namespace) ExecuteAsTemplate(args ...interface{}) (resource.Resource, error) {
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

	return ns.templatesClient.ExecuteAsTemplate(r, targetPath, data)
}

// Fingerprint transforms the given Resource with a MD5 hash of the content in
// the RelPermalink and Permalink.
func (ns *Namespace) Fingerprint(args ...interface{}) (resource.Resource, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("must provide a Resource and (optional) crypto algo")
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

// ToCSS converts the given Resource to CSS. You can optional provide an Options
// object or a target path (string) as first argument.
func (ns *Namespace) ToCSS(args ...interface{}) (resource.Resource, error) {
	var (
		r          resources.ResourceTransformer
		m          map[string]interface{}
		targetPath string
		err        error
		ok         bool
	)

	r, targetPath, ok = ns.resolveIfFirstArgIsString(args)

	if !ok {
		r, m, err = ns.resolveArgs(args)
		if err != nil {
			return nil, err
		}
	}

	var options scss.Options
	if targetPath != "" {
		options.TargetPath = targetPath
	} else if m != nil {
		options, err = scss.DecodeOptions(m)
		if err != nil {
			return nil, err
		}
	}

	return ns.scssClient.ToCSS(r, options)
}

// PostCSS processes the given Resource with PostCSS
func (ns *Namespace) PostCSS(args ...interface{}) (resource.Resource, error) {
	r, m, err := ns.resolveArgs(args)
	if err != nil {
		return nil, err
	}
	var options postcss.Options
	if m != nil {
		options, err = postcss.DecodeOptions(m)
		if err != nil {
			return nil, err
		}
	}

	return ns.postcssClient.Process(r, options)
}

func (ns *Namespace) PostProcess(r resource.Resource) (postpub.PostPublishedResource, error) {
	return ns.deps.ResourceSpec.PostProcess(r)
}

// We allow string or a map as the first argument in some cases.
func (ns *Namespace) resolveIfFirstArgIsString(args []interface{}) (resources.ResourceTransformer, string, bool) {
	if len(args) != 2 {
		return nil, "", false
	}

	v1, ok1 := args[0].(string)
	if !ok1 {
		return nil, "", false
	}
	v2, ok2 := args[1].(resources.ResourceTransformer)

	return v2, v1, ok2
}

// This roundabout way of doing it is needed to get both pipeline behaviour and options as arguments.
func (ns *Namespace) resolveArgs(args []interface{}) (resources.ResourceTransformer, map[string]interface{}, error) {
	if len(args) == 0 {
		return nil, nil, errors.New("no Resource provided in transformation")
	}

	if len(args) == 1 {
		r, ok := args[0].(resources.ResourceTransformer)
		if !ok {
			return nil, nil, fmt.Errorf("type %T not supported in Resource transformations", args[0])
		}
		return r, nil, nil
	}

	r, ok := args[1].(resources.ResourceTransformer)
	if !ok {
		if _, ok := args[1].(map[string]interface{}); !ok {
			return nil, nil, fmt.Errorf("no Resource provided in transformation")
		}
		return nil, nil, fmt.Errorf("type %T not supported in Resource transformations", args[0])
	}

	m, err := maps.ToStringMapE(args[0])
	if err != nil {
		return nil, nil, _errors.Wrap(err, "invalid options type")
	}

	return r, m, nil
}
