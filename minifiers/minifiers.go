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

// Package minifiers contains minifiers mapped to MIME types. This package is used
// in both the resource transformation, i.e. resources.Minify, and in the publishing
// chain.
package minifiers

import (
	"io"
	"regexp"

	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/transform"

	"github.com/gohugoio/hugo/media"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
)

// Client wraps a minifier.
type Client struct {
	m *minify.M
}

// Transformer returns a func that can be used in the transformer publishing chain.
// TODO(bep) minify config etc
func (m Client) Transformer(mediatype media.Type) transform.Transformer {
	_, params, min := m.m.Match(mediatype.Type())
	if min == nil {
		// No minifier for this MIME type
		return nil
	}

	return func(ft transform.FromTo) error {
		// Note that the source io.Reader will already be buffered, but it implements
		// the Bytes() method, which is recognized by the Minify library.
		return min.Minify(m.m, ft.To(), ft.From(), params)
	}
}

// Minify tries to minify the src into dst given a MIME type.
func (m Client) Minify(mediatype media.Type, dst io.Writer, src io.Reader) error {
	return m.m.Minify(mediatype.Type(), dst, src)
}

// New creates a new Client with the provided MIME types as the mapping foundation.
// The HTML minifier is also registered for additional HTML types (AMP etc.) in the
// provided list of output formats.
func New(mediaTypes media.Types, outputFormats output.Formats) Client {
	m := minify.New()
	htmlMin := &html.Minifier{
		KeepDocumentTags:        true,
		KeepConditionalComments: true,
	}

	// We use the Type definition of the media types defined in the site if found.
	addMinifierFunc(m, mediaTypes, "text/css", "css", css.Minify)
	addMinifierFunc(m, mediaTypes, "application/javascript", "js", js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	addMinifierFunc(m, mediaTypes, "application/json", "json", json.Minify)
	addMinifierFunc(m, mediaTypes, "image/svg+xml", "svg", svg.Minify)
	addMinifierFunc(m, mediaTypes, "application/xml", "xml", xml.Minify)
	addMinifierFunc(m, mediaTypes, "application/rss", "xml", xml.Minify)

	// HTML
	addMinifier(m, mediaTypes, "text/html", "html", htmlMin)
	for _, of := range outputFormats {
		if of.IsHTML {
			addMinifier(m, mediaTypes, of.MediaType.Type(), "html", htmlMin)
		}
	}
	return Client{m: m}

}

func addMinifier(m *minify.M, mt media.Types, typeString, suffix string, min minify.Minifier) {
	resolvedTypeStr := resolveMediaTypeString(mt, typeString, suffix)
	m.Add(resolvedTypeStr, min)
	if resolvedTypeStr != typeString {
		m.Add(typeString, min)
	}
}

func addMinifierFunc(m *minify.M, mt media.Types, typeString, suffix string, fn minify.MinifierFunc) {
	resolvedTypeStr := resolveMediaTypeString(mt, typeString, suffix)
	m.AddFunc(resolvedTypeStr, fn)
	if resolvedTypeStr != typeString {
		m.AddFunc(typeString, fn)
	}
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
