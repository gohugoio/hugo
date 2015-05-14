package transform

import (
	"bytes"
	"github.com/spf13/hugo/helpers"
	"path/filepath"
	"strings"
	"testing"
)

const H5_JS_CONTENT_DOUBLE_QUOTE = "<!DOCTYPE html><html><head><script src=\"foobar.js\"></script><script src=\"/barfoo.js\"></script></head><body><nav><h1>title</h1></nav><article>content <a href=\"foobar\">foobar</a>. <a href=\"/foobar\">Follow up</a></article></body></html>"
const H5_JS_CONTENT_SINGLE_QUOTE = "<!DOCTYPE html><html><head><script src='foobar.js'></script><script src='/barfoo.js'></script></head><body><nav><h1>title</h1></nav><article>content <a href='foobar'>foobar</a>. <a href='/foobar'>Follow up</a></article></body></html>"
const H5_JS_CONTENT_ABS_URL = "<!DOCTYPE html><html><head><script src=\"http://user@host:10234/foobar.js\"></script></head><body><nav><h1>title</h1></nav><article>content <a href=\"https://host/foobar\">foobar</a>. Follow up</article></body></html>"
const H5_JS_CONTENT_ABS_URL_SCHEMALESS = "<!DOCTYPE html><html><head><script src=\"//host/foobar.js\"></script><script src='//host2/barfoo.js'></head><body><nav><h1>title</h1></nav><article>content <a href=\"//host/foobar\">foobar</a>. <a href='//host2/foobar'>Follow up</a></article></body></html>"
const CORRECT_OUTPUT_SRC_HREF_DQ = "<!DOCTYPE html><html><head><script src=\"foobar.js\"></script><script src=\"http://base/barfoo.js\"></script></head><body><nav><h1>title</h1></nav><article>content <a href=\"foobar\">foobar</a>. <a href=\"http://base/foobar\">Follow up</a></article></body></html>"
const CORRECT_OUTPUT_SRC_HREF_SQ = "<!DOCTYPE html><html><head><script src='foobar.js'></script><script src='http://base/barfoo.js'></script></head><body><nav><h1>title</h1></nav><article>content <a href='foobar'>foobar</a>. <a href='http://base/foobar'>Follow up</a></article></body></html>"

const H5_XML_CONTENT_ABS_URL = "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?><feed xmlns=\"http://www.w3.org/2005/Atom\"><entry><content type=\"html\">&lt;p&gt;&lt;a href=&#34;/foobar&#34;&gt;foobar&lt;/a&gt;&lt;/p&gt; &lt;p&gt;A video: &lt;iframe src=&#39;/foo&#39;&gt;&lt;/iframe&gt;&lt;/p&gt;</content></entry></feed>"
const CORRECT_OUTPUT_SRC_HREF_IN_XML = "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?><feed xmlns=\"http://www.w3.org/2005/Atom\"><entry><content type=\"html\">&lt;p&gt;&lt;a href=&#34;http://base/foobar&#34;&gt;foobar&lt;/a&gt;&lt;/p&gt; &lt;p&gt;A video: &lt;iframe src=&#39;http://base/foo&#39;&gt;&lt;/iframe&gt;&lt;/p&gt;</content></entry></feed>"
const H5_XML_CONTENT_GUARDED = "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?><feed xmlns=\"http://www.w3.org/2005/Atom\"><entry><content type=\"html\">&lt;p&gt;&lt;a href=&#34;//foobar&#34;&gt;foobar&lt;/a&gt;&lt;/p&gt; &lt;p&gt;A video: &lt;iframe src=&#39;//foo&#39;&gt;&lt;/iframe&gt;&lt;/p&gt;</content></entry></feed>"

// additional sanity tests for replacements testing
const REPLACE_1 = "No replacements."
const REPLACE_2 = "ᚠᛇᚻ ᛒᛦᚦ ᚠᚱᚩᚠᚢᚱ\nᚠᛁᚱᚪ ᚷᛖᚻᚹᛦᛚᚳᚢᛗ"

// Issue: 816, schemaless links combined with others
const REPLACE_SCHEMALESS_HTML = `Pre. src='//schemaless' src='/normal'  <a href="//schemaless">Schemaless</a>. <a href="/normal">normal</a>. Post.`
const REPLACE_SCHEMALESS_HTML_CORRECT = `Pre. src='//schemaless' src='http://base/normal'  <a href="//schemaless">Schemaless</a>. <a href="http://base/normal">normal</a>. Post.`
const REPLACE_SCHEMALESS_XML = `Pre. src=&#39;//schemaless&#39; src=&#39;/normal&#39;  <a href=&#39;//schemaless&#39;>Schemaless</a>. <a href=&#39;/normal&#39;>normal</a>. Post.`
const REPLACE_SCHEMALESS_XML_CORRECT = `Pre. src=&#39;//schemaless&#39; src=&#39;http://base/normal&#39;  <a href=&#39;//schemaless&#39;>Schemaless</a>. <a href=&#39;http://base/normal&#39;>normal</a>. Post.`

// srcset=
const SRCSET_BASIC = `Pre. <img srcset="/img/small.jpg 200w, /img/medium.jpg 300w, /img/big.jpg 700w" alt="text" src="/img/foo.jpg">`
const SRCSET_BASIC_CORRECT = `Pre. <img srcset="http://base/img/small.jpg 200w, http://base/img/medium.jpg 300w, http://base/img/big.jpg 700w" alt="text" src="http://base/img/foo.jpg">`
const SRCSET_SINGLE_QUOTE = `Pre. <img srcset='/img/small.jpg 200w, /img/big.jpg 700w' alt="text" src="/img/foo.jpg"> POST.`
const SRCSET_SINGLE_QUOTE_CORRECT = `Pre. <img srcset='http://base/img/small.jpg 200w, http://base/img/big.jpg 700w' alt="text" src="http://base/img/foo.jpg"> POST.`
const SRCSET_XML_BASIC = `Pre. <img srcset=&#34;/img/small.jpg 200w, /img/big.jpg 700w&#34; alt=&#34;text&#34; src=&#34;/img/foo.jpg&#34;>`
const SRCSET_XML_BASIC_CORRECT = `Pre. <img srcset=&#34;http://base/img/small.jpg 200w, http://base/img/big.jpg 700w&#34; alt=&#34;text&#34; src=&#34;http://base/img/foo.jpg&#34;>`
const SRCSET_XML_SINGLE_QUOTE = `Pre. <img srcset=&#34;/img/small.jpg 200w, /img/big.jpg 700w&#34; alt=&#34;text&#34; src=&#34;/img/foo.jpg&#34;>`
const SRCSET_XML_SINGLE_QUOTE_CORRECT = `Pre. <img srcset=&#34;http://base/img/small.jpg 200w, http://base/img/big.jpg 700w&#34; alt=&#34;text&#34; src=&#34;http://base/img/foo.jpg&#34;>`
const SRCSET_VARIATIONS = `Pre. 
Missing start quote: <img srcset=/img/small.jpg 200w, /img/big.jpg 700w" alt="text"> src='/img/foo.jpg'> FOO. 
<img srcset='/img.jpg'> 
schemaless: <img srcset='//img.jpg' src='//basic.jpg'>
schemaless2: <img srcset="//img.jpg" src="//basic.jpg2> POST
`
const SRCSET_VARIATIONS_CORRECT = `Pre. 
Missing start quote: <img srcset=/img/small.jpg 200w, /img/big.jpg 700w" alt="text"> src='http://base/img/foo.jpg'> FOO. 
<img srcset='http://base/img.jpg'> 
schemaless: <img srcset='//img.jpg' src='//basic.jpg'>
schemaless2: <img srcset="//img.jpg" src="//basic.jpg2> POST
`
const SRCSET_XML_VARIATIONS = `Pre. 
Missing start quote: &lt;img srcset=/img/small.jpg 200w /img/big.jpg 700w&quot; alt=&quot;text&quot;&gt; src=&#39;/img/foo.jpg&#39;&gt; FOO. 
&lt;img srcset=&#39;/img.jpg&#39;&gt; 
schemaless: &lt;img srcset=&#39;//img.jpg&#39; src=&#39;//basic.jpg&#39;&gt;
schemaless2: &lt;img srcset=&quot;//img.jpg&quot; src=&quot;//basic.jpg2&gt; POST
`
const SRCSET_XML_VARIATIONS_CORRECT = `Pre. 
Missing start quote: &lt;img srcset=/img/small.jpg 200w /img/big.jpg 700w&quot; alt=&quot;text&quot;&gt; src=&#39;http://base/img/foo.jpg&#39;&gt; FOO. 
&lt;img srcset=&#39;http://base/img.jpg&#39;&gt; 
schemaless: &lt;img srcset=&#39;//img.jpg&#39; src=&#39;//basic.jpg&#39;&gt;
schemaless2: &lt;img srcset=&quot;//img.jpg&quot; src=&quot;//basic.jpg2&gt; POST
`

const REL_PATH_VARIATIONS = `PRE. a href="/img/small.jpg" POST.`
const REL_PATH_VARIATIONS_CORRECT = `PRE. a href="../../img/small.jpg" POST.`

const testBaseURL = "http://base/"

var abs_url_bench_tests = []test{
	{H5_JS_CONTENT_DOUBLE_QUOTE, CORRECT_OUTPUT_SRC_HREF_DQ},
	{H5_JS_CONTENT_SINGLE_QUOTE, CORRECT_OUTPUT_SRC_HREF_SQ},
	{H5_JS_CONTENT_ABS_URL, H5_JS_CONTENT_ABS_URL},
	{H5_JS_CONTENT_ABS_URL_SCHEMALESS, H5_JS_CONTENT_ABS_URL_SCHEMALESS},
}

var xml_abs_url_bench_tests = []test{
	{H5_XML_CONTENT_ABS_URL, CORRECT_OUTPUT_SRC_HREF_IN_XML},
	{H5_XML_CONTENT_GUARDED, H5_XML_CONTENT_GUARDED},
}

var sanity_tests = []test{{REPLACE_1, REPLACE_1}, {REPLACE_2, REPLACE_2}}
var extra_tests_html = []test{{REPLACE_SCHEMALESS_HTML, REPLACE_SCHEMALESS_HTML_CORRECT}}
var abs_url_tests = append(abs_url_bench_tests, append(sanity_tests, extra_tests_html...)...)
var extra_tests_xml = []test{{REPLACE_SCHEMALESS_XML, REPLACE_SCHEMALESS_XML_CORRECT}}
var xml_abs_url_tests = append(xml_abs_url_bench_tests, append(sanity_tests, extra_tests_xml...)...)
var srcset_tests = []test{{SRCSET_BASIC, SRCSET_BASIC_CORRECT}, {SRCSET_SINGLE_QUOTE, SRCSET_SINGLE_QUOTE_CORRECT}, {SRCSET_VARIATIONS, SRCSET_VARIATIONS_CORRECT}}
var srcset_xml_tests = []test{
	{SRCSET_XML_BASIC, SRCSET_XML_BASIC_CORRECT},
	{SRCSET_XML_SINGLE_QUOTE, SRCSET_XML_SINGLE_QUOTE_CORRECT},
	{SRCSET_XML_VARIATIONS, SRCSET_XML_VARIATIONS_CORRECT}}

var relurl_tests = []test{{REL_PATH_VARIATIONS, REL_PATH_VARIATIONS_CORRECT}}

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
	if err := tr.Apply(out, helpers.StringToReader("Test: f4 f3 f1 f2 f1 The End."), []byte("")); err != nil {
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
		apply(b.Errorf, tr, abs_url_bench_tests)
	}
}

func BenchmarkAbsURLSrcset(b *testing.B) {
	tr := NewChain(AbsURL)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, srcset_tests)
	}
}

func BenchmarkXMLAbsURLSrcset(b *testing.B) {
	tr := NewChain(AbsURLInXML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, srcset_xml_tests)
	}
}

func TestAbsURL(t *testing.T) {
	tr := NewChain(AbsURL)

	apply(t.Errorf, tr, abs_url_tests)

}

func TestRelativeURL(t *testing.T) {
	tr := NewChain(AbsURL)

	applyWithPath(t.Errorf, tr, relurl_tests, helpers.GetDottedRelativePath(filepath.FromSlash("/post/sub/")))

}

func TestAbsURLSrcSet(t *testing.T) {
	tr := NewChain(AbsURL)

	apply(t.Errorf, tr, srcset_tests)
}

func TestAbsXMLURLSrcSet(t *testing.T) {
	tr := NewChain(AbsURLInXML)

	apply(t.Errorf, tr, srcset_xml_tests)
}

func BenchmarkXMLAbsURL(b *testing.B) {
	tr := NewChain(AbsURLInXML)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, xml_abs_url_bench_tests)
	}
}

func TestXMLAbsURL(t *testing.T) {
	tr := NewChain(AbsURLInXML)
	apply(t.Errorf, tr, xml_abs_url_tests)
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
