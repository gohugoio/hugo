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

func TestAbsUrlify(t *testing.T) {
	tr, _ := AbsURL("http://base")
	chain := NewChain(tr...)
	apply(t.Errorf, chain, abs_url_tests)
}

type test struct {
	content  string
	expected string
}

var abs_url_tests = []test{
	{H5_JS_CONTENT_DOUBLE_QUOTE, CORRECT_OUTPUT_SRC_HREF_DQ},
	{H5_JS_CONTENT_SINGLE_QUOTE, CORRECT_OUTPUT_SRC_HREF_SQ},
	{H5_JS_CONTENT_ABS_URL, H5_JS_CONTENT_ABS_URL},
	{H5_JS_CONTENT_ABS_URL_SCHEMALESS, H5_JS_CONTENT_ABS_URL_SCHEMALESS},
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
