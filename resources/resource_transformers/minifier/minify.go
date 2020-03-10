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

package minifier

import (
	"github.com/gohugoio/hugo/minifiers"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/resources/internal"
	"github.com/gohugoio/hugo/resources/resource"
)

// Client for minification of Resource objects. Supported minfiers are:
// css, html, js, json, svg and xml.
type Client struct {
	rs *resources.Spec
	m  minifiers.Client
}

// New creates a new Client given a specification. Note that it is the media types
// configured for the site that is used to match files to the correct minifier.
func New(rs *resources.Spec) (*Client, error) {
	m, err := minifiers.New(rs.MediaTypes, rs.OutputFormats, rs.Cfg)
	if err != nil {
		return nil, err
	}
	return &Client{rs: rs, m: m}, nil
}

type minifyTransformation struct {
	rs *resources.Spec
	m  minifiers.Client
}

func (t *minifyTransformation) Key() internal.ResourceTransformationKey {
	return internal.NewResourceTransformationKey("minify")
}

func (t *minifyTransformation) Transform(ctx *resources.ResourceTransformationCtx) error {
	_ = t.m.Minify(ctx.InMediaType, ctx.To, ctx.From)
	ctx.AddOutPathIdentifier(".min")
	return nil
}

func (c *Client) Minify(res resources.ResourceTransformer) (resource.Resource, error) {
	return res.Transform(&minifyTransformation{
		rs: c.rs,
		m:  c.m,
	})

}
