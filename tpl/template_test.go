// Copyright 2025 The Hugo Authors. All rights reserved.
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

package tpl

import (
	"strings"
	"testing"
)

func TestStripHTML(t *testing.T) {
	type test struct {
		input, expected string
	}
	data := []test{
		{"<h1>strip h1 tag</h1>", "strip h1 tag"},
		{"<p> strip p tag </p>", " strip p tag "},
		{"</br> strip br<br>", " strip br\n"},
		{"</br> strip br2<br />", " strip br2\n"},
		{"This <strong>is</strong> a\nnewline", "This is a newline"},
		{"No Tags", "No Tags"},
		{`<p>Summary Next Line.
<figure >

        <img src="/not/real" />


</figure>
.
More text here.</p>

<p>Some more text</p>`, "Summary Next Line. . More text here.\nSome more text\n"},

		// Issue 9199
		{"<div data-action='click->my-controller#doThing'>qwe</div>", "qwe"},
		{"Hello, World!", "Hello, World!"},
		{"foo&amp;bar", "foo&amp;bar"},
		{`Hello <a href="www.example.com/">World</a>!`, "Hello World!"},
		{"Foo <textarea>Bar</textarea> Baz", "Foo Bar Baz"},
		{"Foo <!-- Bar --> Baz", "Foo Baz"},

		// Issue 15462
		{"<textarea><b>x</b></textarea>", "x"},
		{"<script>var a = 1 < 2; document.write(\"<b>\")</script>foo", "foo"},
		{"<STYLE>p > a { color: red }</STYLE>foo", "foo"},
		{"<script src=\"a.js\" />foo", "foo"},
		{"<script>foo", ""},
		{"I <3 Hugo, 1 < 2 and 3 > 2<b>.</b>", "I <3 Hugo, 1 < 2 and 3 > 2."},
		{"<P>a</P><p>b</p >", "a\nb\n"},
		{"a<br/>b<BR>c<br class=\"x\">d", "a\nb\nc\nd"},
		{`<a href="x>y" title='it"s'>t</a>`, "t"},
		{"<!DOCTYPE html><p>x</p>", "x\n"},
		{"<![CDATA[x<b>]]>y", "y"},
		{"<?xml version=\"1.0\"?>y", "y"},
		{"a<!-- <b> -->b", "ab"},
		{"a</ b>c", "ac"},
		{"foo <b", "foo "},
		{"foo <", "foo <"},
		{"a\t<b>\n </b> \u00a0 b", "a\tb"},
		{"a\u00a0<b> </b>b", "a\u00a0b"},
		{"<ul><li>a</li><li>b</li></ul>", "a b"},
		{"<h2>Title</h2><div>Text</div><table><tr><td>1</td><td>2</td></tr></table>", "Title Text 1 2"},
		{"<p>a<b>b</b>c</p><p>d</p>", "abc\nd\n"},
		{"<svg><path d=\"M0 0\"/></svg>x", "x"},
	}
	for i, d := range data {
		output := StripHTML(d.input)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}
	}
}

func BenchmarkStripHTML(b *testing.B) {
	const para = `<p>Lorem <em>ipsum</em> dolor sit amet, <a href="https://example.com/?a=1&amp;b=2" title="A &gt; B">consectetur</a> adipiscing elit.<br>
Sed do eiusmod <code>tempor</code> incididunt ut labore et dolore magna aliqua.</p>
<!-- a comment -->
<ul>
<li>Ut enim ad minim veniam, quis nostrud</li>
<li>exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</li>
</ul>
<pre><code class="language-go">func main() { fmt.Println("hello") }
</code></pre>
<p>Duis aute irure dolor in reprehenderit — in voluptate velit esse cillum dolore eu fugiat nulla pariatur.</p>
`
	b.Run("Small", func(b *testing.B) {
		for b.Loop() {
			StripHTML(`Hello <a href="www.example.com/">World</a>!`)
		}
	})
	b.Run("NoTags", func(b *testing.B) {
		for b.Loop() {
			StripHTML("Lorem ipsum dolor sit amet, consectetur adipiscing elit.")
		}
	})
	b.Run("Page", func(b *testing.B) {
		s := strings.Repeat(para, 20)
		for b.Loop() {
			StripHTML(s)
		}
	})
}
