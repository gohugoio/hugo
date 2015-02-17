package transform

import (
	"bytes"
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
const REPLACE_SCHEMALESS_XML = `Pre. src=&#34;//schemaless&#34; src=&#34;/normal&#34;  <a href=&#39;//schemaless&#39;>Schemaless</a>. <a href=&#39;/normal&#39;>normal</a>. Post.`
const REPLACE_SCHEMALESS_XML_CORRECT = `Pre. src=&#34;//schemaless&#34; src=&#34;http://base/normal&#34;  <a href=&#39;//schemaless&#39;>Schemaless</a>. <a href=&#39;http://base/normal&#39;>normal</a>. Post.`

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

func TestChainZeroTransformers(t *testing.T) {
	tr := NewChain()
	in := new(bytes.Buffer)
	out := new(bytes.Buffer)
	if err := tr.Apply(in, out); err != nil {
		t.Errorf("A zero transformer chain returned an error.")
	}
}

func BenchmarkAbsUrl(b *testing.B) {
	absURL, _ := AbsURL("http://base")
	tr := NewChain(absURL...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, abs_url_bench_tests)
	}
}

func TestAbsUrl(t *testing.T) {
	absURL, _ := AbsURL("http://base")
	tr := NewChain(absURL...)

	apply(t.Errorf, tr, abs_url_tests)

}

func BenchmarkXmlAbsUrl(b *testing.B) {
	absURLInXML, _ := AbsURLInXML("http://base")
	tr := NewChain(absURLInXML...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, xml_abs_url_bench_tests)
	}
}

func TestXMLAbsUrl(t *testing.T) {
	absURLInXML, _ := AbsURLInXML("http://base")
	tr := NewChain(absURLInXML...)
	apply(t.Errorf, tr, xml_abs_url_tests)
}

type errorf func(string, ...interface{})

func apply(ef errorf, tr chain, tests []test) {
	for _, test := range tests {
		out := new(bytes.Buffer)
		err := tr.Apply(out, strings.NewReader(test.content))
		if err != nil {
			ef("Unexpected error: %s", err)
		}
		if test.expected != string(out.Bytes()) {
			ef("Expected:\n%s\nGot:\n%s", test.expected, string(out.Bytes()))
		}
	}
}

type test struct {
	content  string
	expected string
}
