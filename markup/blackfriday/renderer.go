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

package blackfriday

import (
	"bytes"
	"strings"

	"github.com/russross/blackfriday"
)

// hugoHTMLRenderer wraps a blackfriday.Renderer, typically a blackfriday.Html
// adding some custom behaviour.
type hugoHTMLRenderer struct {
	c *blackfridayConverter
	blackfriday.Renderer
}

// BlockCode renders a given text as a block of code.
// Pygments is used if it is setup to handle code fences.
func (r *hugoHTMLRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	if r.c.cfg.MarkupConfig.Highlight.CodeFences {
		str := strings.Trim(string(text), "\n\r")
		highlighted, _ := r.c.cfg.Highlight(str, lang, "")
		out.WriteString(highlighted)
	} else {
		r.Renderer.BlockCode(out, text, lang)
	}
}

// ListItem adds task list support to the Blackfriday renderer.
func (r *hugoHTMLRenderer) ListItem(out *bytes.Buffer, text []byte, flags int) {
	if !r.c.bf.TaskLists {
		r.Renderer.ListItem(out, text, flags)
		return
	}

	switch {
	case bytes.HasPrefix(text, []byte("[ ] ")):
		text = append([]byte(`<label><input type="checkbox" disabled class="task-list-item">`), text[3:]...)
		text = append(text, []byte(`</label>`)...)

	case bytes.HasPrefix(text, []byte("[x] ")) || bytes.HasPrefix(text, []byte("[X] ")):
		text = append([]byte(`<label><input type="checkbox" checked disabled class="task-list-item">`), text[3:]...)
		text = append(text, []byte(`</label>`)...)
	}

	r.Renderer.ListItem(out, text, flags)
}

// List adds task list support to the Blackfriday renderer.
func (r *hugoHTMLRenderer) List(out *bytes.Buffer, text func() bool, flags int) {
	if !r.c.bf.TaskLists {
		r.Renderer.List(out, text, flags)
		return
	}
	marker := out.Len()
	r.Renderer.List(out, text, flags)
	if out.Len() > marker {
		list := out.Bytes()[marker:]
		if bytes.Contains(list, []byte("task-list-item")) {
			// Find the index of the first >, it might be 3 or 4 depending on whether
			// there is a new line at the start, but this is safer than just hardcoding it.
			closingBracketIndex := bytes.Index(list, []byte(">"))
			// Rewrite the buffer from the marker
			out.Truncate(marker)
			// Safely assuming closingBracketIndex won't be -1 since there is a list
			// May be either dl, ul or ol
			list := append(list[:closingBracketIndex], append([]byte(` class="task-list"`), list[closingBracketIndex:]...)...)
			out.Write(list)
		}
	}
}
