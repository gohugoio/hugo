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
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"html"
	"html/template"
	"io"
	"strings"
	"sync/atomic"

	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/internal/warpc"
	"github.com/gohugoio/hugo/markup/converter/hooks"
	"github.com/gohugoio/hugo/markup/highlight"
	"github.com/gohugoio/hugo/markup/highlight/chromalexers"
	"github.com/gohugoio/hugo/resources"
	"github.com/gohugoio/hugo/tpl"
	"github.com/mitchellh/mapstructure"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/cast"
)

// New returns a new instance of the transform-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	if deps.MemCache == nil {
		panic("must provide MemCache")
	}

	return &Namespace{
		deps: deps,
		cacheUnmarshal: dynacache.GetOrCreatePartition[string, *resources.StaleValue[any]](
			deps.MemCache,
			"/tmpl/transform/unmarshal",
			dynacache.OptionsPartition{Weight: 30, ClearWhen: dynacache.ClearOnChange},
		),
		cacheMath: dynacache.GetOrCreatePartition[string, template.HTML](
			deps.MemCache,
			"/tmpl/transform/math",
			dynacache.OptionsPartition{Weight: 30, ClearWhen: dynacache.ClearNever},
		),
	}
}

// Namespace provides template functions for the "transform" namespace.
type Namespace struct {
	cacheUnmarshal *dynacache.Partition[string, *resources.StaleValue[any]]
	cacheMath      *dynacache.Partition[string, template.HTML]

	id   atomic.Uint32
	deps *deps.Deps
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
	highlighted, err := hl.Highlight(ss, lang, optsv)
	if err != nil {
		return "", err
	}
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
	return chromalexers.Get(language) != nil
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

// XMLEscape returns the given string, removing disallowed characters then
// escaping the result to its XML equivalent.
func (ns *Namespace) XMLEscape(s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	// https://www.w3.org/TR/xml/#NT-Char
	cleaned := strings.Map(func(r rune) rune {
		if r == 0x9 || r == 0xA || r == 0xD ||
			(r >= 0x20 && r <= 0xD7FF) ||
			(r >= 0xE000 && r <= 0xFFFD) ||
			(r >= 0x10000 && r <= 0x10FFFF) {
			return r
		}
		return -1
	}, ss)

	var buf bytes.Buffer
	err = xml.EscapeText(&buf, []byte(cleaned))
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Markdownify renders s from Markdown to HTML.
func (ns *Namespace) Markdownify(ctx context.Context, s any) (template.HTML, error) {
	home := ns.deps.Site.Home()
	if home == nil {
		panic("home must not be nil")
	}
	ss, err := home.RenderString(ctx, s)
	if err != nil {
		return "", err
	}

	// Strip if this is a short inline type of text.
	bb := ns.deps.ContentSpec.TrimShortHTML([]byte(ss), "markdown")

	return helpers.BytesToHTML(bb), nil
}

// Plainify returns a copy of s with all HTML tags removed.
func (ns *Namespace) Plainify(s any) (template.HTML, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", err
	}

	return template.HTML(tpl.StripHTML(ss)), nil
}

// ToMath converts a LaTeX string to math in the given format, default MathML.
// This uses KaTeX to render the math, see https://katex.org/.
func (ns *Namespace) ToMath(ctx context.Context, args ...any) (types.Result[template.HTML], error) {
	var res types.Result[template.HTML]

	if len(args) < 1 {
		return res, errors.New("must provide at least one argument")
	}
	expression, err := cast.ToStringE(args[0])
	if err != nil {
		return res, err
	}

	katexInput := warpc.KatexInput{
		Expression: expression,
		Options: warpc.KatexOptions{
			Output:           "mathml",
			MinRuleThickness: 0.04,
			ErrorColor:       "#cc0000",
			ThrowOnError:     true,
		},
	}

	if len(args) > 1 {
		if err := mapstructure.WeakDecode(args[1], &katexInput.Options); err != nil {
			return res, err
		}
	}

	s := hashing.HashString(args...)
	key := "tomath/" + s[:2] + "/" + s[2:]
	fileCache := ns.deps.ResourceSpec.FileCaches.MiscCache()

	v, err := ns.cacheMath.GetOrCreate(key, func(string) (template.HTML, error) {
		_, r, err := fileCache.GetOrCreate(key, func() (io.ReadCloser, error) {
			message := warpc.Message[warpc.KatexInput]{
				Header: warpc.Header{
					Version: 1,
					ID:      ns.id.Add(1),
				},
				Data: katexInput,
			}

			k, err := ns.deps.WasmDispatchers.Katex()
			if err != nil {
				return nil, err
			}
			result, err := k.Execute(ctx, message)
			if err != nil {
				return nil, err
			}
			return hugio.NewReadSeekerNoOpCloserFromString(result.Data.Output), nil
		})
		if err != nil {
			return "", err
		}

		s, err := hugio.ReadString(r)

		return template.HTML(s), err
	})

	res = types.Result[template.HTML]{
		Value: v,
		Err:   err,
	}

	return res, nil
}

// For internal use.
func (ns *Namespace) Reset() {
	ns.cacheUnmarshal.Clear()
}
