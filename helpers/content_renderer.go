// Copyright 2015 The Hugo Authors. All rights reserved.
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

package helpers

import (
	"bytes"
	"html"

	"github.com/miekg/mmark"
	"github.com/russross/blackfriday"
	"github.com/spf13/viper"
)

// Wraps a blackfriday.Renderer, typically a blackfriday.Html
// Enabling Hugo to customise the rendering experience
type HugoHtmlRenderer struct {
	blackfriday.Renderer
}

func (renderer *HugoHtmlRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	if viper.GetBool("PygmentsCodeFences") {
		opts := viper.GetString("PygmentsOptions")
		str := html.UnescapeString(string(text))
		out.WriteString(Highlight(str, lang, opts))
	} else {
		renderer.Renderer.BlockCode(out, text, lang)
	}
}

// Wraps a mmark.Renderer, typically a mmark.html
// Enabling Hugo to customise the rendering experience
type HugoMmarkHtmlRenderer struct {
	mmark.Renderer
}

func (renderer *HugoMmarkHtmlRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string, caption []byte, subfigure bool, callouts bool) {
	if viper.GetBool("PygmentsCodeFences") {
		str := html.UnescapeString(string(text))
		out.WriteString(Highlight(str, lang, ""))
	} else {
		renderer.Renderer.BlockCode(out, text, lang, caption, subfigure, callouts)
	}
}
