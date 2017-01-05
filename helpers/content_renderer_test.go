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
)

// Renders a codeblock using Blackfriday
func render(input string) string {
	ctx := newViperProvidedRenderingContext()
	render := getHTMLRenderer(0, ctx)

	buf := &bytes.Buffer{}
	render.BlockCode(buf, []byte(input), "html")
	return buf.String()
}

// Renders a codeblock using Mmark
func renderWithMmark(input string) string {
	ctx := newViperProvidedRenderingContext()
	render := getMmarkHTMLRenderer(0, ctx)

	buf := &bytes.Buffer{}
	render.BlockCode(buf, []byte(input), "html", []byte(""), false, false)
	return buf.String()
}

func TestCodeFence(t *testing.T) {

	if !HasPygments() {
		t.Skip("Skipping Pygments test as Pygments is not installed or available.")
		return
	}

	type test struct {
		enabled         bool
		input, expected string
	}

	// Pygments 2.0 and 2.1 have slightly different outputs so only do partial matching
	data := []test{
		{true, "<html></html>", `(?s)^<div class="highlight"><pre><code class="language-html" data-lang="html">.*?</code></pre></div>\n$`},
		{false, "<html></html>", `(?s)^<pre><code class="language-html">.*?</code></pre>\n$`},
	}

	viper.Reset()
	defer viper.Reset()

	viper.Set("pygmentsStyle", "monokai")
	viper.Set("pygmentsUseClasses", true)

	for i, d := range data {
		viper.Set("pygmentsCodeFences", d.enabled)

		result := render(d.input)

		expectedRe, err := regexp.Compile(d.expected)

		if err != nil {
			t.Fatal("Invalid regexp", err)
		}
		matched := expectedRe.MatchString(result)

		if !matched {
			t.Errorf("Test %d failed. BlackFriday enabled:%t, Expected:\n%q got:\n%q", i, d.enabled, d.expected, result)
		}

		result = renderWithMmark(d.input)
		matched = expectedRe.MatchString(result)
		if !matched {
			t.Errorf("Test %d failed. Mmark enabled:%t, Expected:\n%q got:\n%q", i, d.enabled, d.expected, result)
		}
	}
}

func TestBlackfridayTaskList(t *testing.T) {
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
<li><input type="checkbox" checked disabled class="task-list-item"> On1</li>
<li><input type="checkbox" checked disabled class="task-list-item"> On2</li>
<li><input type="checkbox" disabled class="task-list-item"> Off</li>
</ul>

<p>END</p>
`},
		{`- [x] On1`, false, `<ul>
<li>[x] On1</li>
</ul>
`},
	} {
		blackFridayConfig := NewBlackfriday(viper.GetStringMap("blackfriday"))
		blackFridayConfig.TaskLists = this.taskListEnabled
		ctx := &RenderingContext{Content: []byte(this.markdown), PageFmt: "markdown", Config: blackFridayConfig}

		result := string(RenderBytes(ctx))

		if result != this.expect {
			t.Errorf("[%d] got \n%v but expected \n%v", i, result, this.expect)
		}
	}
}
