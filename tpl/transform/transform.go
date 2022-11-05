// Copyright 2017 The Hugo Authors. All rights reserved.
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

// Package transform provides template functions for transforming content.
package transform

import (
	"html"
	"html/template"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/gohugoio/hugo/cache/namedmemcache"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/highlight"
	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/cast"
)

// New returns a new instance of the transform-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	cache := namedmemcache.New()
	deps.BuildStartListeners.Add(
		func() {
			cache.Clear()
		})

	return &Namespace{
		cache: cache,
		deps:  deps,
	}
}

// Namespace provides template functions for the "transform" namespace.
type Namespace struct {
	cache *namedmemcache.Cache
	deps  *deps.Deps
}

// Emojify returns a copy of s with all emoji codes replaced with actual emojis.
//
// See http://www.emoji-cheat-sheet.com/
func (ns *Namespace) Emojify(s any) (template.HTML, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	return template.HTML(helpers.Emojify([]byte(ss))), nil
}

// Highlight returns a copy of s as an HTML string with syntax
// highlighting applied.
func (ns *Namespace) Highlight(s any, lang string, opts ...any) (template.HTML, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	var optsv any
	if len(opts) > 0 {
		optsv = opts[0]
	}

	hl := ns.deps.ContentSpec.Converters.GetHighlighter()
	highlighted, _ := hl.Highlight(ss, lang, optsv)
	return template.HTML(highlighted), nil
}

// HighlightCodeBlock highlights a code block on the form received in the codeblock render hooks.
func (ns *Namespace) HighlightCodeBlock(ctx hooks.CodeblockContext, opts ...any) (highlight.HighlightResult, error) {
	var optsv any
	if len(opts) > 0 {
		optsv = opts[0]
	}

	hl := ns.deps.ContentSpec.Converters.GetHighlighter()

	return hl.HighlightCodeBlock(ctx, optsv)
}

// CanHighlight returns whether the given code language is supported by the Chroma highlighter.
func (ns *Namespace) CanHighlight(language string) bool {
	return lexers.Get(language) != nil
}

// HTMLEscape returns a copy of s with reserved HTML characters escaped.
func (ns *Namespace) HTMLEscape(s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	return html.EscapeString(ss), nil
}

// HTMLUnescape returns a copy of s with HTML escape requences converted to plain
// text.
func (ns *Namespace) HTMLUnescape(s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	return html.UnescapeString(ss), nil
}

// Markdownify renders s from Markdown to HTML.
func (ns *Namespace) Markdownify(s any) (template.HTML, error) {

	home := ns.deps.Site.Home()
	if home == nil {
		panic("home must not be nil")
	}
	ss, err := home.RenderString(s)
	if err != nil {
		return "", err
	}

	// Strip if this is a short inline type of text.
	bb := ns.deps.ContentSpec.TrimShortHTML([]byte(ss))

	return helpers.BytesToHTML(bb), nil
}

// Plainify returns a copy of s with all HTML tags removed.
func (ns *Namespace) Plainify(s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	return tpl.StripHTML(ss), nil
}

// For internal use.
func (ns *Namespace) Reset() {
	ns.cache.Clear()
}
