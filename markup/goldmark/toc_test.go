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

// Package goldmark converts Markdown to HTML using Goldmark.
package goldmark

import (
	"testing"

	"github.com/gohugoio/hugo/markup/markup_config"

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/markup/converter"

	qt "github.com/frankban/quicktest"
)

func TestToc(t *testing.T) {
	c := qt.New(t)

	content := `
# Header 1

## First h2

Some text.

### H3

Some more text.

## Second h2

And then some.

### Second H3

#### First H4

`
	p, err := Provider.New(
		converter.ProviderConfig{
			MarkupConfig: markup_config.Default,
			Logger:       loggers.NewErrorLogger()})
	c.Assert(err, qt.IsNil)
	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)
	b, err := conv.Convert(converter.RenderContext{Src: []byte(content), RenderTOC: true})
	c.Assert(err, qt.IsNil)
	got := b.(converter.TableOfContentsProvider).TableOfContents().ToHTML(2, 3)
	c.Assert(got, qt.Equals, `<nav id="TableOfContents">
  <ul>
    <li><a href="#first-h2">First h2</a>
      <ul>
        <li><a href="#h3">H3</a></li>
      </ul>
    </li>
    <li><a href="#second-h2">Second h2</a>
      <ul>
        <li><a href="#second-h3">Second H3</a></li>
      </ul>
    </li>
  </ul>
</nav>`, qt.Commentf(got))
}
