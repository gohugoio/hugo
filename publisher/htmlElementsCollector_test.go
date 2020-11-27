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
	"fmt"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestClassCollector(t *testing.T) {
	c := qt.New((t))

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

	for _, test := range []struct {
		name   string
		html   string
		expect HTMLElements
	}{
		{"basic", `<body class="b a"></body>`, f("body", "a b", "")},
		{"duplicates", `<div class="b a b"></div>`, f("div", "a b", "")},
		{"single quote", `<body class='b a'></body>`, f("body", "a b", "")},
		{"no quote", `<body class=b id=myelement></body>`, f("body", "b", "myelement")},
		{"thead", `
		https://github.com/gohugoio/hugo/issues/7318
<table class="cl1">
    <thead class="cl2"><tr class="cl3"><td class="cl4"></td></tr></thead>
    <tbody class="cl5"><tr class="cl6"><td class="cl7"></td></tr></tbody>
</table>`, f("table tbody td thead tr", "cl1 cl2 cl3 cl4 cl5 cl6 cl7", "")},
		// https://github.com/gohugoio/hugo/issues/7161
		{"minified a href", `<a class="b a" href=/></a>`, f("a", "a b", "")},

		{"AlpineJS bind 1", `<body>
			<div x-bind:class="{
        'class1': data.open,
        'class2 class3': data.foo == 'bar'
         }">
			</div>
		</body>`, f("body div", "class1 class2 class3", "")},

		{"Alpine bind 2", `<div x-bind:class="{ 'bg-black':  filter.checked }"
                        class="inline-block mr-1 mb-2 rounded  bg-gray-300 px-2 py-2">FOO</div>`,
			f("div", "bg-black bg-gray-300 inline-block mb-2 mr-1 px-2 py-2 rounded", "")},

		{"Alpine bind 3", `<div x-bind:class="{ 'text-gray-800':  !checked, 'text-white': checked }"></div>`, f("div", "text-gray-800 text-white", "")},
		{"Alpine bind 4", `<div x-bind:class="{ 'text-gray-800':  !checked, 
					 'text-white': checked }"></div>`, f("div", "text-gray-800 text-white", "")},

		{"Alpine bind 5", `<a x-bind:class="{
                'text-a': a && b,
                'text-b': !a && b || c,
                'pl-3': a === 1,
                 pl-2: b == 3,
                'text-gray-600': (a > 1)
      
                }" class="block w-36 cursor-pointer pr-3 no-underline capitalize"></a>`, f("a", "block capitalize cursor-pointer no-underline pl-2 pl-3 pr-3 text-a text-b text-gray-600 w-36", "")},

		{"Alpine transition 1", `<div x-transition:enter-start="opacity-0 transform mobile:-translate-x-8 sm:-translate-y-8">`, f("div", "mobile:-translate-x-8 opacity-0 sm:-translate-y-8 transform", "")},
		{"Vue bind", `<div v-bind:class="{ active: isActive }"></div>`, f("div", "active", "")},
		// https://github.com/gohugoio/hugo/issues/7746
		{"Apostrophe inside attribute value", `<a class="missingclass" title="Plus d'information">my text</a><div></div>`, f("a div", "missingclass", "")},
	} {
		c.Run(test.name, func(c *qt.C) {
			w := newHTMLElementsCollectorWriter(newHTMLElementsCollector())
			fmt.Fprint(w, test.html)
			got := w.collector.getHTMLElements()
			c.Assert(got, qt.DeepEquals, test.expect)
		})
	}

}

func BenchmarkClassCollectorWriter(b *testing.B) {
	const benchHTML = `
<html>
<body id="i1" class="a b c d">
<a class="c d e"></a>
<br>
<a class="c d e"></a>
<a class="c d e"></a>
<br>
<a id="i2" class="c d e f"></a>
<a id="i3" class="c d e"></a>
<a class="c d e"></a>
<br>
<a class="c d e"></a>
<a class="c d e"></a>
<a class="c d e"></a>
<a class="c d e"></a>
</body>
</html>
`
	for i := 0; i < b.N; i++ {
		w := newHTMLElementsCollectorWriter(newHTMLElementsCollector())
		fmt.Fprint(w, benchHTML)

	}
}
