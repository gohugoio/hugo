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

package minifiers

import (
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/media"

	"github.com/gohugoio/hugo/resource"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
)

// Client for minification of Resource objects. Supported minfiers are:
// css, html, js, json, svg and xml.
type Client struct {
	rs *resource.Spec
	m  *minify.M
}

// New creates a new Client given a specification. Note that it is the media types
// configured for the site that is used to match files to the correct minifier.
func New(rs *resource.Spec) *Client {
	m := minify.New()
	mt := rs.MediaTypes

	// We use the Type definition of the media types defined in the site if found.
	addMinifierFunc(m, mt, "text/css", "css", css.Minify)
	addMinifierFunc(m, mt, "text/html", "html", html.Minify)
	addMinifierFunc(m, mt, "application/javascript", "js", js.Minify)
	addMinifierFunc(m, mt, "application/json", "json", json.Minify)
	addMinifierFunc(m, mt, "image/svg+xml", "svg", svg.Minify)
	addMinifierFunc(m, mt, "application/xml", "xml", xml.Minify)
	addMinifierFunc(m, mt, "application/rss", "xml", xml.Minify)

	return &Client{rs: rs, m: m}
}

func addMinifierFunc(m *minify.M, mt media.Types, typeString, suffix string, fn minify.MinifierFunc) {
	resolvedTypeStr := resolveMediaTypeString(mt, typeString, suffix)
	m.AddFunc(resolvedTypeStr, fn)
	if resolvedTypeStr != typeString {
		m.AddFunc(typeString, fn)
	}
}

type minifyTransformation struct {
	rs *resource.Spec
	m  *minify.M
}

func (t *minifyTransformation) Key() resource.ResourceTransformationKey {
	return resource.NewResourceTransformationKey("minify")
}

func (t *minifyTransformation) Transform(ctx *resource.ResourceTransformationCtx) error {
	mtype := resolveMediaTypeString(
		t.rs.MediaTypes,
		ctx.InMediaType.Type(),
		helpers.ExtNoDelimiter(ctx.InPath),
	)
	if err := t.m.Minify(mtype, ctx.To, ctx.From); err != nil {
		return err
	}
	ctx.AddOutPathIdentifier(".min")
	return nil
}

func (c *Client) Minify(res resource.Resource) (resource.Resource, error) {
	return c.rs.Transform(
		res,
		&minifyTransformation{
			rs: c.rs,
			m:  c.m},
	)
}

func resolveMediaTypeString(types media.Types, typeStr, suffix string) string {
	if m, found := resolveMediaType(types, typeStr, suffix); found {
		return m.Type()
	}
	// Fall back to the default.
	return typeStr
}

// Make sure we match the matching pattern with what the user have actually defined
// in his or hers media types configuration.
func resolveMediaType(types media.Types, typeStr, suffix string) (media.Type, bool) {
	if m, found := types.GetByType(typeStr); found {
		return m, true
	}

	if m, found := types.GetFirstBySuffix(suffix); found {
		return m, true
	}

	return media.Type{}, false

}
