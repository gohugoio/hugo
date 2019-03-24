// Copyright 2018 The Hugo Authors. All rights reserved.
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

package urlreplacers

import (
	"path/filepath"
	"testing"

	bp "github.com/gohugoio/hugo/bufferpool"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/transform"
)

const (
	h5JsContentDoubleQuote      = "<!DOCTYPE html><html><head><script src=\"foobar.js\"></script><script src=\"/barfoo.js\"></script></head><body><nav><h1>title</h1></nav><article>content <a href=\"foobar\">foobar</a>. <a href=\"/foobar\">Follow up</a></article></body></html>"
	h5JsContentSingleQuote      = "<!DOCTYPE html><html><head><script src='foobar.js'></script><script src='/barfoo.js'></script></head><body><nav><h1>title</h1></nav><article>content <a href='foobar'>foobar</a>. <a href='/foobar'>Follow up</a></article></body></html>"
	h5JsContentAbsURL           = "<!DOCTYPE html><html><head><script src=\"http://user@host:10234/foobar.js\"></script></head><body><nav><h1>title</h1></nav><article>content <a href=\"https://host/foobar\">foobar</a>. Follow up</article></body></html>"
	h5JsContentAbsURLSchemaless = "<!DOCTYPE html><html><head><script src=\"//host/foobar.js\"></script><script src='//host2/barfoo.js'></head><body><nav><h1>title</h1></nav><article>content <a href=\"//host/foobar\">foobar</a>. <a href='//host2/foobar'>Follow up</a></article></body></html>"
	corectOutputSrcHrefDq       = "<!DOCTYPE html><html><head><script src=\"foobar.js\"></script><script src=\"http://base/barfoo.js\"></script></head><body><nav><h1>title</h1></nav><article>content <a href=\"foobar\">foobar</a>. <a href=\"http://base/foobar\">Follow up</a></article></body></html>"
	corectOutputSrcHrefSq       = "<!DOCTYPE html><html><head><script src='foobar.js'></script><script src='http://base/barfoo.js'></script></head><body><nav><h1>title</h1></nav><article>content <a href='foobar'>foobar</a>. <a href='http://base/foobar'>Follow up</a></article></body></html>"

	h5XMLXontentAbsURL        = "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?><feed xmlns=\"http://www.w3.org/2005/Atom\"><entry><content type=\"html\">&lt;p&gt;&lt;a href=&#34;/foobar&#34;&gt;foobar&lt;/a&gt;&lt;/p&gt; &lt;p&gt;A video: &lt;iframe src=&#39;/foo&#39;&gt;&lt;/iframe&gt;&lt;/p&gt;</content></entry></feed>"
	correctOutputSrcHrefInXML = "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?><feed xmlns=\"http://www.w3.org/2005/Atom\"><entry><content type=\"html\">&lt;p&gt;&lt;a href=&#34;http://base/foobar&#34;&gt;foobar&lt;/a&gt;&lt;/p&gt; &lt;p&gt;A video: &lt;iframe src=&#39;http://base/foo&#39;&gt;&lt;/iframe&gt;&lt;/p&gt;</content></entry></feed>"
	h5XMLContentGuarded       = "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?><feed xmlns=\"http://www.w3.org/2005/Atom\"><entry><content type=\"html\">&lt;p&gt;&lt;a href=&#34;//foobar&#34;&gt;foobar&lt;/a&gt;&lt;/p&gt; &lt;p&gt;A video: &lt;iframe src=&#39;//foo&#39;&gt;&lt;/iframe&gt;&lt;/p&gt;</content></entry></feed>"
)

const (
	// additional sanity tests for replacements testing
	replace1 = "No replacements."
	replace2 = "ᚠᛇᚻ ᛒᛦᚦ ᚠᚱᚩᚠᚢᚱ\nᚠᛁᚱᚪ ᚷᛖᚻᚹᛦᛚᚳᚢᛗ"
	replace3 = `End of file: src="/`
	replace5 = `Srcsett with no closing quote: srcset="/img/small.jpg do be do be do.`

	// Issue: 816, schemaless links combined with others
	replaceSchemalessHTML        = `Pre. src='//schemaless' src='/normal'  <a href="//schemaless">Schemaless</a>. <a href="/normal">normal</a>. Post.`
	replaceSchemalessHTMLCorrect = `Pre. src='//schemaless' src='http://base/normal'  <a href="//schemaless">Schemaless</a>. <a href="http://base/normal">normal</a>. Post.`
	replaceSchemalessXML         = `Pre. src=&#39;//schemaless&#39; src=&#39;/normal&#39;  <a href=&#39;//schemaless&#39;>Schemaless</a>. <a href=&#39;/normal&#39;>normal</a>. Post.`
	replaceSchemalessXMLCorrect  = `Pre. src=&#39;//schemaless&#39; src=&#39;http://base/normal&#39;  <a href=&#39;//schemaless&#39;>Schemaless</a>. <a href=&#39;http://base/normal&#39;>normal</a>. Post.`
)

const (
	// srcset=
	srcsetBasic                 = `Pre. <img srcset="/img/small.jpg 200w, /img/medium.jpg 300w, /img/big.jpg 700w" alt="text" src="/img/foo.jpg">`
	srcsetBasicCorrect          = `Pre. <img srcset="http://base/img/small.jpg 200w, http://base/img/medium.jpg 300w, http://base/img/big.jpg 700w" alt="text" src="http://base/img/foo.jpg">`
	srcsetSingleQuote           = `Pre. <img srcset='/img/small.jpg 200w, /img/big.jpg 700w' alt="text" src="/img/foo.jpg"> POST.`
	srcsetSingleQuoteCorrect    = `Pre. <img srcset='http://base/img/small.jpg 200w, http://base/img/big.jpg 700w' alt="text" src="http://base/img/foo.jpg"> POST.`
	srcsetXMLBasic              = `Pre. <img srcset=&#34;/img/small.jpg 200w, /img/big.jpg 700w&#34; alt=&#34;text&#34; src=&#34;/img/foo.jpg&#34;>`
	srcsetXMLBasicCorrect       = `Pre. <img srcset=&#34;http://base/img/small.jpg 200w, http://base/img/big.jpg 700w&#34; alt=&#34;text&#34; src=&#34;http://base/img/foo.jpg&#34;>`
	srcsetXMLSingleQuote        = `Pre. <img srcset=&#34;/img/small.jpg 200w, /img/big.jpg 700w&#34; alt=&#34;text&#34; src=&#34;/img/foo.jpg&#34;>`
	srcsetXMLSingleQuoteCorrect = `Pre. <img srcset=&#34;http://base/img/small.jpg 200w, http://base/img/big.jpg 700w&#34; alt=&#34;text&#34; src=&#34;http://base/img/foo.jpg&#34;>`
	srcsetVariations            = `Pre. 
Missing start quote: <img srcset=/img/small.jpg 200w, /img/big.jpg 700w" alt="text"> src='/img/foo.jpg'> FOO. 
<img srcset='/img.jpg'> 
schemaless: <img srcset='//img.jpg' src='//basic.jpg'>
schemaless2: <img srcset="//img.jpg" src="//basic.jpg2> POST
`
)

const (
	srcsetVariationsCorrect = `Pre. 
Missing start quote: <img srcset=/img/small.jpg 200w, /img/big.jpg 700w" alt="text"> src='http://base/img/foo.jpg'> FOO. 
<img srcset='http://base/img.jpg'> 
schemaless: <img srcset='//img.jpg' src='//basic.jpg'>
schemaless2: <img srcset="//img.jpg" src="//basic.jpg2> POST
`
	srcsetXMLVariations = `Pre. 
Missing start quote: &lt;img srcset=/img/small.jpg 200w /img/big.jpg 700w&quot; alt=&quot;text&quot;&gt; src=&#39;/img/foo.jpg&#39;&gt; FOO. 
&lt;img srcset=&#39;/img.jpg&#39;&gt; 
schemaless: &lt;img srcset=&#39;//img.jpg&#39; src=&#39;//basic.jpg&#39;&gt;
schemaless2: &lt;img srcset=&quot;//img.jpg&quot; src=&quot;//basic.jpg2&gt; POST
`
	srcsetXMLVariationsCorrect = `Pre. 
Missing start quote: &lt;img srcset=/img/small.jpg 200w /img/big.jpg 700w&quot; alt=&quot;text&quot;&gt; src=&#39;http://base/img/foo.jpg&#39;&gt; FOO. 
&lt;img srcset=&#39;http://base/img.jpg&#39;&gt; 
schemaless: &lt;img srcset=&#39;//img.jpg&#39; src=&#39;//basic.jpg&#39;&gt;
schemaless2: &lt;img srcset=&quot;//img.jpg&quot; src=&quot;//basic.jpg2&gt; POST
`

	relPathVariations        = `PRE. a href="/img/small.jpg" input action="/foo.html" POST.`
	relPathVariationsCorrect = `PRE. a href="../../img/small.jpg" input action="../../foo.html" POST.`

	testBaseURL = "http://base/"
)

var (
	absURLlBenchTests = []test{
		{h5JsContentDoubleQuote, corectOutputSrcHrefDq},
		{h5JsContentSingleQuote, corectOutputSrcHrefSq},
		{h5JsContentAbsURL, h5JsContentAbsURL},
		{h5JsContentAbsURLSchemaless, h5JsContentAbsURLSchemaless},
	}

	xmlAbsURLBenchTests = []test{
		{h5XMLXontentAbsURL, correctOutputSrcHrefInXML},
		{h5XMLContentGuarded, h5XMLContentGuarded},
	}

	sanityTests    = []test{{replace1, replace1}, {replace2, replace2}, {replace3, replace3}, {replace3, replace3}, {replace5, replace5}}
	extraTestsHTML = []test{{replaceSchemalessHTML, replaceSchemalessHTMLCorrect}}
	absURLTests    = append(absURLlBenchTests, append(sanityTests, extraTestsHTML...)...)
	extraTestsXML  = []test{{replaceSchemalessXML, replaceSchemalessXMLCorrect}}
	xmlAbsURLTests = append(xmlAbsURLBenchTests, append(sanityTests, extraTestsXML...)...)
	srcsetTests    = []test{{srcsetBasic, srcsetBasicCorrect}, {srcsetSingleQuote, srcsetSingleQuoteCorrect}, {srcsetVariations, srcsetVariationsCorrect}}
	srcsetXMLTests = []test{
		{srcsetXMLBasic, srcsetXMLBasicCorrect},
		{srcsetXMLSingleQuote, srcsetXMLSingleQuoteCorrect},
		{srcsetXMLVariations, srcsetXMLVariationsCorrect}}

	relurlTests = []test{{relPathVariations, relPathVariationsCorrect}}
)

func BenchmarkAbsURL(b *testing.B) {
	tr := transform.New(NewAbsURLTransformer(testBaseURL))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, absURLlBenchTests)
	}
}

func BenchmarkAbsURLSrcset(b *testing.B) {
	tr := transform.New(NewAbsURLTransformer(testBaseURL))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, srcsetTests)
	}
}

func BenchmarkXMLAbsURLSrcset(b *testing.B) {
	tr := transform.New(NewAbsURLInXMLTransformer(testBaseURL))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, srcsetXMLTests)
	}
}

func TestAbsURL(t *testing.T) {
	tr := transform.New(NewAbsURLTransformer(testBaseURL))

	apply(t.Errorf, tr, absURLTests)

}

func TestAbsURLUnqoted(t *testing.T) {
	tr := transform.New(NewAbsURLTransformer(testBaseURL))

	apply(t.Errorf, tr, []test{
		{
			content:  `Link: <a href=/asdf>ASDF</a>`,
			expected: `Link: <a href=http://base/asdf>ASDF</a>`,
		},
		{
			content:  `Link: <a href=/asdf   >ASDF</a>`,
			expected: `Link: <a href=http://base/asdf   >ASDF</a>`,
		},
	})
}

func TestRelativeURL(t *testing.T) {
	tr := transform.New(NewAbsURLTransformer(helpers.GetDottedRelativePath(filepath.FromSlash("/post/sub/"))))

	applyWithPath(t.Errorf, tr, relurlTests)

}

func TestAbsURLSrcSet(t *testing.T) {
	tr := transform.New(NewAbsURLTransformer(testBaseURL))

	apply(t.Errorf, tr, srcsetTests)
}

func TestAbsXMLURLSrcSet(t *testing.T) {
	tr := transform.New(NewAbsURLInXMLTransformer(testBaseURL))

	apply(t.Errorf, tr, srcsetXMLTests)
}

func BenchmarkXMLAbsURL(b *testing.B) {
	tr := transform.New(NewAbsURLInXMLTransformer(testBaseURL))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, xmlAbsURLBenchTests)
	}
}

func TestXMLAbsURL(t *testing.T) {
	tr := transform.New(NewAbsURLInXMLTransformer(testBaseURL))
	apply(t.Errorf, tr, xmlAbsURLTests)
}

func apply(ef errorf, tr transform.Chain, tests []test) {
	applyWithPath(ef, tr, tests)
}

func applyWithPath(ef errorf, tr transform.Chain, tests []test) {
	out := bp.GetBuffer()
	defer bp.PutBuffer(out)

	in := bp.GetBuffer()
	defer bp.PutBuffer(in)

	for _, test := range tests {
		var err error
		in.WriteString(test.content)
		err = tr.Apply(out, in)
		if err != nil {
			ef("Unexpected error: %s", err)
		}
		if test.expected != out.String() {
			ef("Expected:\n%s\nGot:\n%s", test.expected, out.String())
		}
		out.Reset()
		in.Reset()
	}
}

type test struct {
	content  string
	expected string
}

type errorf func(string, ...interface{})
