// Copyright 2020 The Hugo Authors. All rights reserved.
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

package publisher

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/minifiers"
	"github.com/gohugoio/hugo/output"

	qt "github.com/frankban/quicktest"
)

func TestClassCollector(t *testing.T) {
	c := qt.New((t))
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	f := func(tags, classes, ids string) HTMLElements {
		var tagss, classess, idss []string
		if tags != "" {
			tagss = strings.Split(tags, " ")
		}
		if classes != "" {
			classess = strings.Split(classes, " ")
		}
		if ids != "" {
			idss = strings.Split(ids, " ")
		}
		return HTMLElements{
			Tags:    tagss,
			Classes: classess,
			IDs:     idss,
		}
	}

	skipMinifyTest := map[string]bool{
		"Script tags content should be skipped": true, // https://github.com/tdewolff/minify/issues/396
	}

	for _, test := range []struct {
		name   string
		html   string
		expect HTMLElements
	}{
		{"basic", `<body class="b a"></body>`, f("body", "a b", "")},
		{"duplicates", `<div class="b a b"></div><div class="b a b"></div>x'`, f("div", "a b", "")},
		{"single quote", `<body class='b a'></body>`, f("body", "a b", "")},
		{"no quote", `<body class=b id=myelement></body>`, f("body", "b", "myelement")},
		{"short", `<i>`, f("i", "", "")},
		{"invalid", `< body class="b a"></body><div></div>`, f("div", "", "")},
		// https://github.com/gohugoio/hugo/issues/7318
		{"thead", `<table class="cl1">
    <thead class="cl2"><tr class="cl3"><td class="cl4"></td></tr></thead>
    <tbody class="cl5"><tr class="cl6"><td class="cl7"></td></tr></tbody>
</table>`, f("table tbody td thead tr", "cl1 cl2 cl3 cl4 cl5 cl6 cl7", "")},
		{"thead uppercase", `<TABLE class="CL1">
    <THEAD class="CL2"><TR class="CL3"><TD class="CL4"></TD></TR></THEAD>
    <TBODY class="CL5"><TR class="CL6"><TD class="CL7"></TD></TR></TBODY>
</TABLE>`, f("table tbody td thead tr", "CL1 CL2 CL3 CL4 CL5 CL6 CL7", "")},
		// https://github.com/gohugoio/hugo/issues/7161
		{"minified a href", `<a class="b a" href=/></a>`, f("a", "a b", "")},
		{"AlpineJS bind 1", `<body>
    <div x-bind:class="{
        'class1': data.open,
        'class2 class3': data.foo == 'bar'
         }">
    </div>
</body>`, f("body div", "class1 class2 class3", "")},
		{
			"AlpineJS bind 2", `<div x-bind:class="{ 'bg-black':  filter.checked }" class="inline-block mr-1 mb-2 rounded  bg-gray-300 px-2 py-2">FOO</div>`,
			f("div", "bg-black bg-gray-300 inline-block mb-2 mr-1 px-2 py-2 rounded", ""),
		},
		{"AlpineJS bind 3", `<div x-bind:class="{ 'text-gray-800':  !checked, 'text-white': checked }"></div>`, f("div", "text-gray-800 text-white", "")},
		{"AlpineJS bind 4", `<div x-bind:class="{ 'text-gray-800':  !checked, 
					 'text-white': checked }"></div>`, f("div", "text-gray-800 text-white", "")},
		{"AlpineJS bind 5", `<a x-bind:class="{
                'text-a': a && b,
                'text-b': !a && b || c,
                'pl-3': a === 1,
                 pl-2: b == 3,
                'text-gray-600': (a > 1)
                }" class="block w-36 cursor-pointer pr-3 no-underline capitalize"></a>`, f("a", "block capitalize cursor-pointer no-underline pl-2 pl-3 pr-3 text-a text-b text-gray-600 w-36", "")},
		{"AlpineJS bind 6", `<button :class="isActive(32) ? 'border-gray-500 bg-white pt border-t-2' : 'border-transparent hover:bg-gray-100'"></button>`, f("button", "bg-white border-gray-500 border-t-2 border-transparent hover:bg-gray-100 pt", "")},
		{"AlpineJS bind 7", `<button :class="{ 'border-gray-500 bg-white pt border-t-2': isActive(32), 'border-transparent hover:bg-gray-100': !isActive(32) }"></button>`, f("button", "bg-white border-gray-500 border-t-2 border-transparent hover:bg-gray-100 pt", "")},
		{"AlpineJS transition 1", `<div x-transition:enter-start="opacity-0 transform mobile:-translate-x-8 sm:-translate-y-8">`, f("div", "mobile:-translate-x-8 opacity-0 sm:-translate-y-8 transform", "")},
		{"Vue bind", `<div v-bind:class="{ active: isActive }"></div>`, f("div", "active", "")},
		// Issue #7746
		{"Apostrophe inside attribute value", `<a class="missingclass" title="Plus d'information">my text</a><div></div>`, f("a div", "missingclass", "")},
		// Issue #7567
		{"Script tags content should be skipped", `<script><span>foo</span><span>bar</span></script><div class="foo"></div>`, f("div script", "foo", "")},
		{"Style tags content should be skipped", `<style>p{color: red;font-size: 20px;}</style><div class="foo"></div>`, f("div style", "foo", "")},
		{"Pre tags content should be skipped", `<pre class="preclass"><span>foo</span><span>bar</span></pre><div class="foo"></div>`, f("div pre", "foo preclass", "")},
		{"Textarea tags content should be skipped", `<textarea class="textareaclass"><span>foo</span><span>bar</span></textarea><div class="foo"></div>`, f("div textarea", "foo textareaclass", "")},
		{"DOCTYPE should beskipped", `<!DOCTYPE html>`, f("", "", "")},
		{"Comments should be skipped", `<!-- example comment -->`, f("", "", "")},
		{"Comments with elements before and after", `<div></div><!-- example comment --><span><span>`, f("div span", "", "")},
		{"Self closing tag", `<div><hr/></div>`, f("div hr", "", "")},
		// svg with self closing style tag.
		{"SVG with self closing style tag", `<svg><style/><g><path class="foo"/></g></svg>`, f("g path style svg", "foo", "")},
		// Issue #8530
		{"Comment with single quote", `<!-- Hero Area Image d'accueil --><i class="foo">`, f("i", "foo", "")},
		{"Uppercase tags", `<DIV></DIV>`, f("div", "", "")},
		{"Predefined tags with distinct casing", `<script>if (a < b) { nothing(); }</SCRIPT><div></div>`, f("div script", "", "")},
		// Issue #8417
		{"Tabs inline", `<hr	id="a" class="foo"><div class="bar">d</div>`, f("div hr", "bar foo", "a")},
		{"Tabs on multiple rows", `<form
			id="a"
			action="www.example.com"
			method="post"
></form>
<div id="b" class="foo">d</div>`, f("div form", "foo", "a b")},
		{"Big input, multibyte runes", strings.Repeat(`神真美好 `, rnd.Intn(500)+1) + "<div id=\"神真美好\" class=\"foo\">" + strings.Repeat(`神真美好 `, rnd.Intn(100)+1) + "   <span>神真美好</span>", f("div span", "foo", "神真美好")},
	} {
		for _, variant := range []struct {
			minify bool
		}{
			{minify: false},
			{minify: true},
		} {

			name := fmt.Sprintf("%s--minify-%t", test.name, variant.minify)

			c.Run(name, func(c *qt.C) {
				w := newHTMLElementsCollectorWriter(newHTMLElementsCollector(
					config.BuildStats{Enable: true},
				))
				if variant.minify {
					if skipMinifyTest[test.name] {
						c.Skip("skip minify test")
					}
					m, _ := minifiers.New(media.DefaultTypes, output.DefaultFormats, testconfig.GetTestConfig(nil, nil))
					m.Minify(media.Builtin.HTMLType, w, strings.NewReader(test.html))

				} else {
					var buff bytes.Buffer
					buff.WriteString(test.html)
					io.Copy(w, &buff)
				}
				got := w.collector.getHTMLElements()
				c.Assert(got, qt.DeepEquals, test.expect)
			})
		}
	}
}

func TestEndsWithTag(t *testing.T) {
	c := qt.New((t))

	for _, test := range []struct {
		name    string
		s       string
		tagName string
		expect  bool
	}{
		{"empty", "", "div", false},
		{"no match", "foo", "div", false},
		{"no close", "foo<div>", "div", false},
		{"no close 2", "foo/div>", "div", false},
		{"no close 2", "foo//div>", "div", false},
		{"no tag", "foo</>", "div", false},
		{"match", "foo</div>", "div", true},
		{"match space", "foo<  / div>", "div", true},
		{"match space 2", "foo<  / div   \n>", "div", true},
		{"match case", "foo</DIV>", "div", true},
		{"self closing", `</defs><g><g><path fill="#010101" d=asdf"/>`, "div", false},
	} {
		c.Run(test.name, func(c *qt.C) {
			got := isClosedByTag([]byte(test.s), []byte(test.tagName))
			c.Assert(got, qt.Equals, test.expect)
		})
	}
}

func BenchmarkElementsCollectorWriter(b *testing.B) {
	const benchHTML = `
<!DOCTYPE html>
<html>
<head>
<title>title</title>
<style>
	a {color: red;}
	.c {color: blue;}
</style>
</head>
<body id="i1" class="a b c d">
<a class="c d e"></a>
<hr>
<a class="c d e"></a>
<a class="c d e"></a>
<hr>
<a id="i2" class="c d e f"></a>
<a id="i3" class="c d e"></a>
<a class="c d e"></a>
<p>To force<br> line breaks<br> in a text,<br> use the br<br> element.</p>
<hr>
<a class="c d e"></a>
<a class="c d e"></a>
<a class="c d e"></a>
<a class="c d e"></a>
<table>
  <thead class="ch">
  <tr>
    <th>Month</th>
    <th>Savings</th>
  </tr>
  </thead>
  <tbody class="cb">
  <tr>
    <td>January</td>
    <td>$100</td>
  </tr>
  <tr>
    <td>February</td>
    <td>$200</td>
  </tr>
  </tbody>
  <tfoot class="cf">
  <tr>
    <td></td>
    <td>$300</td>
  </tr>
  </tfoot>
</table>
</body>
</html>
`
	for i := 0; i < b.N; i++ {
		w := newHTMLElementsCollectorWriter(newHTMLElementsCollector(
			config.BuildStats{Enable: true},
		))
		fmt.Fprint(w, benchHTML)

	}
}

func BenchmarkElementsCollectorWriterPre(b *testing.B) {
	const benchHTML = `
<pre class="preclass">
<span>foo</span><span>bar</span>
<!-- many more span elements -->
<span class="foo">foo</span>
<span class="bar">bar</span>
<span class="baz">baz</span>
<span class="qux">qux</span>
<span class="quux">quux</span>
<span class="quuz">quuz</span>
<span class="corge">corge</span>
</pre>
<div class="foo"></div>

`
	w := newHTMLElementsCollectorWriter(newHTMLElementsCollector(
		config.BuildStats{Enable: true},
	))
	for i := 0; i < b.N; i++ {
		fmt.Fprint(w, benchHTML)
	}
}
