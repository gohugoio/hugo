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
	"github.com/spf13/viper"
	"testing"
)

// Renders a codeblock using Blackfriday
func render(input string) string {
	ctx := &RenderingContext{}
	render := GetHTMLRenderer(0, ctx)

	buf := &bytes.Buffer{}
	render.BlockCode(buf, []byte(input), "html")
	return buf.String()
}

// Renders a codeblock using Mmark
func renderWithMmark(input string) string {
	ctx := &RenderingContext{}
	render := GetMmarkHtmlRenderer(0, ctx)

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
	data := []test{
		{true, "<html></html>", "<div class=\"highlight\"><pre><code class=\"language-html\" data-lang=\"html\"><span class=\"nt\">&lt;html&gt;&lt;/html&gt;</span>\n</code></pre></div>\n"},
		{false, "<html></html>", "<pre><code class=\"language-html\">&lt;html&gt;&lt;/html&gt;</code></pre>\n"},
	}

	viper.Reset()
	defer viper.Reset()

	viper.Set("PygmentsStyle", "monokai")
	viper.Set("PygmentsUseClasses", true)

	for i, d := range data {
		viper.Set("PygmentsCodeFences", d.enabled)

		result := render(d.input)
		if result != d.expected {
			t.Errorf("Test %d failed. BlackFriday enabled:%t, Expected:\n%q got:\n%q", i, d.enabled, d.expected, result)
		}

		result = renderWithMmark(d.input)
		if result != d.expected {
			t.Errorf("Test %d failed. Mmark enabled:%t, Expected:\n%q got:\n%q", i, d.enabled, d.expected, result)
		}
	}
}
