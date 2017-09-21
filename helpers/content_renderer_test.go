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
	"regexp"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// Renders a codeblock using Blackfriday
func (c ContentSpec) render(input string) string {
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewBlackfriday()}
	render := c.getHTMLRenderer(0, ctx)

	buf := &bytes.Buffer{}
	render.BlockCode(buf, []byte(input), "html")
	return buf.String()
}

// Renders a codeblock using Mmark
func (c ContentSpec) renderWithMmark(input string) string {
	ctx := &RenderingContext{Cfg: c.cfg, Config: c.NewBlackfriday()}
	render := c.getMmarkHTMLRenderer(0, ctx)

	buf := &bytes.Buffer{}
	render.BlockCode(buf, []byte(input), "html", []byte(""), false, false)
	return buf.String()
}

func TestCodeFence(t *testing.T) {
	assert := require.New(t)

	type test struct {
		enabled         bool
		input, expected string
	}

	// Pygments 2.0 and 2.1 have slightly different outputs so only do partial matching
	data := []test{
		{true, "<html></html>", `(?s)^<div class="highlight">\n?<pre.*><code class="language-html" data-lang="html">.*?</code></pre>\n?</div>\n?$`},
		{false, "<html></html>", `(?s)^<pre.*><code class="language-html">.*?</code></pre>\n$`},
	}

	for _, useClassic := range []bool{false, true} {
		for i, d := range data {
			v := viper.New()
			v.Set("pygmentsStyle", "monokai")
			v.Set("pygmentsUseClasses", true)
			v.Set("pygmentsCodeFences", d.enabled)
			v.Set("pygmentsUseClassic", useClassic)

			c, err := NewContentSpec(v)
			assert.NoError(err)

			result := c.render(d.input)

			expectedRe, err := regexp.Compile(d.expected)

			if err != nil {
				t.Fatal("Invalid regexp", err)
			}
			matched := expectedRe.MatchString(result)

			if !matched {
				t.Errorf("Test %d failed. BlackFriday enabled:%t, Expected:\n%q got:\n%q", i, d.enabled, d.expected, result)
			}

			result = c.renderWithMmark(d.input)
			matched = expectedRe.MatchString(result)
			if !matched {
				t.Errorf("Test %d failed. Mmark enabled:%t, Expected:\n%q got:\n%q", i, d.enabled, d.expected, result)
			}
		}
	}
}

func TestBlackfridayTaskList(t *testing.T) {
	c := newTestContentSpec()

	for i, this := range []struct {
		markdown        string
		taskListEnabled bool
		expect          string
	}{
		{`
TODO:

- [x] On1
- [X] On2
- [ ] Off

END
`, true, `<p>TODO:</p>

<ul class="task-list">
<li><label><input type="checkbox" checked disabled class="task-list-item"> On1</label></li>
<li><label><input type="checkbox" checked disabled class="task-list-item"> On2</label></li>
<li><label><input type="checkbox" disabled class="task-list-item"> Off</label></li>
</ul>

<p>END</p>
`},
		{`- [x] On1`, false, `<ul>
<li>[x] On1</li>
</ul>
`},
		{`* [ ] Off

END`, true, `<ul class="task-list">
<li><label><input type="checkbox" disabled class="task-list-item"> Off</label></li>
</ul>

<p>END</p>
`},
	} {
		blackFridayConfig := c.NewBlackfriday()
		blackFridayConfig.TaskLists = this.taskListEnabled
		ctx := &RenderingContext{Content: []byte(this.markdown), PageFmt: "markdown", Config: blackFridayConfig}

		result := string(c.RenderBytes(ctx))

		if result != this.expect {
			t.Errorf("[%d] got \n%v but expected \n%v", i, result, this.expect)
		}
	}
}
