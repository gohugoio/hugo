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

// Package highlight provides code highlighting.
package highlight

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestHighlight(t *testing.T) {
	c := qt.New(t)

	lines := `LINE1
LINE2
LINE3
LINE4
LINE5
`
	coalesceNeeded := `GET /foo HTTP/1.1
Content-Type: application/json
User-Agent: foo

{
  "hello": "world"
}`

	c.Run("Basic", func(c *qt.C) {
		cfg := DefaultConfig
		cfg.NoClasses = false
		h := New(cfg)

		result, _ := h.Highlight(`echo "Hugo Rocks!"`, "bash", "")
		c.Assert(result, qt.Equals, `<div class="highlight"><pre class="chroma"><code class="language-bash" data-lang="bash"><span class="nb">echo</span> <span class="s2">&#34;Hugo Rocks!&#34;</span></code></pre></div>`)
		result, _ = h.Highlight(`echo "Hugo Rocks!"`, "unknown", "")
		c.Assert(result, qt.Equals, `<pre><code class="language-unknown" data-lang="unknown">echo &#34;Hugo Rocks!&#34;</code></pre>`)

	})

	c.Run("Highlight lines, default config", func(c *qt.C) {
		cfg := DefaultConfig
		cfg.NoClasses = false
		h := New(cfg)

		result, _ := h.Highlight(lines, "bash", "linenos=table,hl_lines=2 4-5,linenostart=3")
		c.Assert(result, qt.Contains, "<div class=\"highlight\"><div class=\"chroma\">\n<table class=\"lntable\"><tr><td class=\"lntd\">\n<pre class=\"chroma\"><code><span class")
		c.Assert(result, qt.Contains, "<span class=\"hl\"><span class=\"lnt\">4")

		result, _ = h.Highlight(lines, "bash", "linenos=inline,hl_lines=2")
		c.Assert(result, qt.Contains, "<span class=\"ln\">2</span>LINE2\n</span>")
		c.Assert(result, qt.Not(qt.Contains), "<table")

		result, _ = h.Highlight(lines, "bash", "linenos=true,hl_lines=2")
		c.Assert(result, qt.Contains, "<table")
		c.Assert(result, qt.Contains, "<span class=\"hl\"><span class=\"lnt\">2\n</span>")
	})

	c.Run("Highlight lines, linenumbers default on", func(c *qt.C) {
		cfg := DefaultConfig
		cfg.NoClasses = false
		cfg.LineNos = true
		h := New(cfg)

		result, _ := h.Highlight(lines, "bash", "")
		c.Assert(result, qt.Contains, "<span class=\"lnt\">2\n</span>")
		result, _ = h.Highlight(lines, "bash", "linenos=false,hl_lines=2")
		c.Assert(result, qt.Not(qt.Contains), "class=\"lnt\"")
	})

	c.Run("Highlight lines, linenumbers default on, linenumbers in table default off", func(c *qt.C) {
		cfg := DefaultConfig
		cfg.NoClasses = false
		cfg.LineNos = true
		cfg.LineNumbersInTable = false
		h := New(cfg)

		result, _ := h.Highlight(lines, "bash", "")
		c.Assert(result, qt.Contains, "<span class=\"ln\">2</span>LINE2\n<")
		result, _ = h.Highlight(lines, "bash", "linenos=table")
		c.Assert(result, qt.Contains, "<span class=\"lnt\">1\n</span>")
	})

	c.Run("No language", func(c *qt.C) {
		cfg := DefaultConfig
		cfg.NoClasses = false
		cfg.LineNos = true
		h := New(cfg)

		result, _ := h.Highlight(lines, "", "")
		c.Assert(result, qt.Equals, "<pre><code>LINE1\nLINE2\nLINE3\nLINE4\nLINE5\n</code></pre>")
	})

	c.Run("No language, guess syntax", func(c *qt.C) {
		cfg := DefaultConfig
		cfg.NoClasses = false
		cfg.GuessSyntax = true
		cfg.LineNos = true
		cfg.LineNumbersInTable = false
		h := New(cfg)

		result, _ := h.Highlight(lines, "", "")
		c.Assert(result, qt.Contains, "<span class=\"ln\">2</span>LINE2\n<")
	})

	c.Run("No language, Escape HTML string", func(c *qt.C) {
		cfg := DefaultConfig
		cfg.NoClasses = false
		h := New(cfg)

		result, _ := h.Highlight("Escaping less-than in code block? <fail>", "", "")
		c.Assert(result, qt.Contains, "&lt;fail&gt;")
	})

	c.Run("Highlight lines, default config", func(c *qt.C) {
		cfg := DefaultConfig
		cfg.NoClasses = false
		h := New(cfg)

		result, _ := h.Highlight(coalesceNeeded, "http", "linenos=true,hl_lines=2")
		c.Assert(result, qt.Contains, "hello")
		c.Assert(result, qt.Contains, "}")
	})

}
