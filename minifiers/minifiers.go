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
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
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
		KeepEndTags:             true,
		KeepDefaultAttrVals:     true,
	}

	cssMin := &css.Minifier{
		Decimals: -1,
		KeepCSS2: true,
	}

	// We use the Type definition of the media types defined in the site if found.
	addMinifier(m, mediaTypes, "css", cssMin)
	addMinifierFunc(m, mediaTypes, "js", js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-|ld\\+)?json$"), json.Minify)
	addMinifierFunc(m, mediaTypes, "json", json.Minify)
	addMinifierFunc(m, mediaTypes, "svg", svg.Minify)
	addMinifierFunc(m, mediaTypes, "xml", xml.Minify)

	// HTML
	addMinifier(m, mediaTypes, "html", htmlMin)
	for _, of := range outputFormats {
		if of.IsHTML {
			m.Add(of.MediaType.Type(), htmlMin)
		}
	}

	return Client{m: m}

}

func addMinifier(m *minify.M, mt media.Types, suffix string, min minify.Minifier) {
	types := mt.BySuffix(suffix)
	for _, t := range types {
		m.Add(t.Type(), min)
	}
}

func addMinifierFunc(m *minify.M, mt media.Types, suffix string, min minify.MinifierFunc) {
	types := mt.BySuffix(suffix)
	for _, t := range types {
		m.AddFunc(t.Type(), min)
	}
}
