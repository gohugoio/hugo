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

package transform

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/helpers"
	"github.com/stretchr/testify/assert"
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
	replace4 = `End of file: srcset="/`
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

	relPathVariations        = `PRE. a href="/img/small.jpg" POST.`
	relPathVariationsCorrect = `PRE. a href="../../img/small.jpg" POST.`

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

func TestChainZeroTransformers(t *testing.T) {
	tr := NewChain()
	in := new(bytes.Buffer)
	out := new(bytes.Buffer)
	if err := tr.Apply(in, out, []byte("")); err != nil {
		t.Errorf("A zero transformer chain returned an error.")
	}
}

func TestChaingMultipleTransformers(t *testing.T) {
	f1 := func(ct contentTransformer) {
		ct.Write(bytes.Replace(ct.Content(), []byte("f1"), []byte("f1r"), -1))
	}
	f2 := func(ct contentTransformer) {
		ct.Write(bytes.Replace(ct.Content(), []byte("f2"), []byte("f2r"), -1))
	}
	f3 := func(ct contentTransformer) {
		ct.Write(bytes.Replace(ct.Content(), []byte("f3"), []byte("f3r"), -1))
	}

	f4 := func(ct contentTransformer) {
		ct.Write(bytes.Replace(ct.Content(), []byte("f4"), []byte("f4r"), -1))
	}

	tr := NewChain(f1, f2, f3, f4)

	out := new(bytes.Buffer)
	if err := tr.Apply(out, strings.NewReader("Test: f4 f3 f1 f2 f1 The End."), []byte("")); err != nil {
		t.Errorf("Multi transformer chain returned an error: %s", err)
	}

	expected := "Test: f4r f3r f1r f2r f1r The End."

	if string(out.Bytes()) != expected {
		t.Errorf("Expected %s got %s", expected, string(out.Bytes()))
	}
}

func BenchmarkAbsURL(b *testing.B) {
	tr := NewChain(AbsURL)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, absURLlBenchTests)
	}
}

func BenchmarkAbsURLSrcset(b *testing.B) {
	tr := NewChain(AbsURL)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, srcsetTests)
	}
}

func BenchmarkXMLAbsURLSrcset(b *testing.B) {
	tr := NewChain(AbsURLInXML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, srcsetXMLTests)
	}
}

func TestAbsURL(t *testing.T) {
	tr := NewChain(AbsURL)

	apply(t.Errorf, tr, absURLTests)

}

func TestRelativeURL(t *testing.T) {
	tr := NewChain(AbsURL)

	applyWithPath(t.Errorf, tr, relurlTests, helpers.GetDottedRelativePath(filepath.FromSlash("/post/sub/")))

}

func TestAbsURLSrcSet(t *testing.T) {
	tr := NewChain(AbsURL)

	apply(t.Errorf, tr, srcsetTests)
}

func TestAbsXMLURLSrcSet(t *testing.T) {
	tr := NewChain(AbsURLInXML)

	apply(t.Errorf, tr, srcsetXMLTests)
}

func BenchmarkXMLAbsURL(b *testing.B) {
	tr := NewChain(AbsURLInXML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, xmlAbsURLBenchTests)
	}
}

func TestXMLAbsURL(t *testing.T) {
	tr := NewChain(AbsURLInXML)
	apply(t.Errorf, tr, xmlAbsURLTests)
}

func TestNewEmptyTransforms(t *testing.T) {
	transforms := NewEmptyTransforms()
	assert.Equal(t, 20, cap(transforms))
}

type errorf func(string, ...interface{})

func applyWithPath(ef errorf, tr chain, tests []test, path string) {
	for _, test := range tests {
		out := new(bytes.Buffer)
		var err error
		err = tr.Apply(out, strings.NewReader(test.content), []byte(path))
		if err != nil {
			ef("Unexpected error: %s", err)
		}
		if test.expected != string(out.Bytes()) {
			ef("Expected:\n%s\nGot:\n%s", test.expected, string(out.Bytes()))
		}
	}
}

func apply(ef errorf, tr chain, tests []test) {
	applyWithPath(ef, tr, tests, testBaseURL)
}

type test struct {
	content  string
	expected string
}
