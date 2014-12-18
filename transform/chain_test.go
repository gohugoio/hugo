package transform

import (
	"bytes"
	"testing"
)

const H5_JS_CONTENT_ABS_URL_WITH_NAV = "<!DOCTYPE html><html><head><script src=\"/foobar.js\"></script></head><body><nav><ul><li hugo-nav=\"section_0\"></li><li hugo-nav=\"section_1\"></li></ul></nav><article>content <a href=\"/foobar\">foobar</a>. Follow up</article></body></html>"

const CORRECT_OUTPUT_SRC_HREF_WITH_NAV = "<!DOCTYPE html><html><head><script src=\"http://two/foobar.js\"></script></head><body><nav><ul><li hugo-nav=\"section_0\"></li><li hugo-nav=\"section_1\"></li></ul></nav><article>content <a href=\"http://two/foobar\">foobar</a>. Follow up</article></body></html>"

const H5_XML_CONTENT_ABS_URL = "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?><feed xmlns=\"http://www.w3.org/2005/Atom\"><entry><content type=\"html\">&lt;p&gt;&lt;a href=&#34;/foobar&#34;&gt;foobar&lt;/a&gt;&lt;/p&gt; &lt;p&gt;A video: &lt;iframe src=&#39;/foo&#39;&gt;&lt;/iframe&gt;&lt;/p&gt;</content></entry></feed>"

const CORRECT_OUTPUT_SRC_HREF_IN_XML = "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?><feed xmlns=\"http://www.w3.org/2005/Atom\"><entry><content type=\"html\">&lt;p&gt;&lt;a href=&#34;http://xml/foobar&#34;&gt;foobar&lt;/a&gt;&lt;/p&gt; &lt;p&gt;A video: &lt;iframe src=&#39;http://xml/foo&#39;&gt;&lt;/iframe&gt;&lt;/p&gt;</content></entry></feed>"

var two_chain_tests = []test{
	{H5_JS_CONTENT_ABS_URL_WITH_NAV, CORRECT_OUTPUT_SRC_HREF_WITH_NAV},
}

var xml_abs_url_tests = []test{
	{H5_XML_CONTENT_ABS_URL, CORRECT_OUTPUT_SRC_HREF_IN_XML},
}

func TestChainZeroTransformers(t *testing.T) {
	tr := NewChain()
	in := new(bytes.Buffer)
	out := new(bytes.Buffer)
	if err := tr.Apply(in, out); err != nil {
		t.Errorf("A zero transformer chain returned an error.")
	}
}

func BenchmarkChain(b *testing.B) {
	absURL, _ := AbsURL("http://two")
	tr := NewChain(absURL...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apply(b.Errorf, tr, two_chain_tests)
	}
}

func TestXMLAbsUrl(t *testing.T) {
	absURLInXML, _ := AbsURLInXML("http://xml")
	tr := NewChain(absURLInXML...)
	apply(t.Errorf, tr, xml_abs_url_tests)
}
