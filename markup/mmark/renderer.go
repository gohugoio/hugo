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

package mmark

import (
	"bytes"
	"strings"

	"github.com/miekg/mmark"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/markup/internal"
)

// hugoHTMLRenderer wraps a blackfriday.Renderer, typically a blackfriday.Html
// adding some custom behaviour.
type mmarkRenderer struct {
	Cfg       config.Provider
	Config    *internal.BlackFriday
	highlight func(code, lang, optsStr string) (string, error)
	mmark.Renderer
}

// BlockCode renders a given text as a block of code.
func (r *mmarkRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string, caption []byte, subfigure bool, callouts bool) {
	if r.Cfg.GetBool("pygmentsCodeFences") && (lang != "" || r.Cfg.GetBool("pygmentsCodeFencesGuessSyntax")) {
		str := strings.Trim(string(text), "\n\r")
		highlighted, _ := r.highlight(str, lang, "")
		out.WriteString(highlighted)
	} else {
		r.Renderer.BlockCode(out, text, lang, caption, subfigure, callouts)
	}
}
