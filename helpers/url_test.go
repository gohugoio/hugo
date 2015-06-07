package helpers

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestURLize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  foo bar  ", "foo-bar"},
		{"foo.bar/foo_bar-foo", "foo.bar/foo_bar-foo"},
		{"foo,bar:foo%bar", "foobarfoobar"},
		{"foo/bar.html", "foo/bar.html"},
		{"трям/трям", "%D1%82%D1%80%D1%8F%D0%BC/%D1%82%D1%80%D1%8F%D0%BC"},
	}

	for _, test := range tests {
		output := URLize(test.input)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}

func TestAbsURL(t *testing.T) {
	defer viper.Reset()
	tests := []struct {
		input    string
		baseURL  string
		expected string
	}{
		{"/test/foo", "http://base/", "http://base/test/foo"},
		{"", "http://base/ace/", "http://base/ace/"},
		{"/test/2/foo/", "http://base", "http://base/test/2/foo/"},
		{"http://abs", "http://base/", "http://abs"},
		{"//schemaless", "http://base/", "//schemaless"},
	}

	for _, test := range tests {
		viper.Reset()
		viper.Set("BaseURL", test.baseURL)
		output := AbsURL(test.input)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}

func TestRelURL(t *testing.T) {
	defer viper.Reset()
	//defer viper.Set("canonifyURLs", viper.GetBool("canonifyURLs"))
	tests := []struct {
		input    string
		baseURL  string
		canonify bool
		expected string
	}{
		{"/test/foo", "http://base/", false, "/test/foo"},
		{"test.css", "http://base/sub", false, "/sub/test.css"},
		{"test.css", "http://base/sub", true, "/test.css"},
		{"/test/", "http://base/", false, "/test/"},
		{"/test/", "http://base/sub/", false, "/sub/test/"},
		{"/test/", "http://base/sub/", true, "/test/"},
		{"", "http://base/ace/", false, "/ace/"},
		{"", "http://base/ace", false, "/ace"},
		{"http://abs", "http://base/", false, "http://abs"},
		{"//schemaless", "http://base/", false, "//schemaless"},
	}

	for i, test := range tests {
		viper.Reset()
		viper.Set("BaseURL", test.baseURL)
		viper.Set("canonifyURLs", test.canonify)

		output := RelURL(test.input)
		if output != test.expected {
			t.Errorf("[%d][%t] Expected %#v, got %#v\n", i, test.canonify, test.expected, output)
		}
	}
}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://foo.bar/", "http://foo.bar"},
		{"http://foo.bar", "http://foo.bar"},          // issue #1105
		{"http://foo.bar/zoo/", "http://foo.bar/zoo"}, // issue #931
	}

	for i, test := range tests {
		o1 := SanitizeURL(test.input)
		o2 := SanitizeURLKeepTrailingSlash(test.input)

		expected2 := test.expected

		if strings.HasSuffix(test.input, "/") && !strings.HasSuffix(expected2, "/") {
			expected2 += "/"
		}

		if o1 != test.expected {
			t.Errorf("[%d] 1: Expected %#v, got %#v\n", i, test.expected, o1)
		}
		if o2 != expected2 {
			t.Errorf("[%d] 2: Expected %#v, got %#v\n", i, expected2, o2)
		}
	}
}

func TestMakePermalink(t *testing.T) {
	type test struct {
		host, link, output string
	}

	data := []test{
		{"http://abc.com/foo", "post/bar", "http://abc.com/foo/post/bar"},
		{"http://abc.com/foo/", "post/bar", "http://abc.com/foo/post/bar"},
		{"http://abc.com", "post/bar", "http://abc.com/post/bar"},
		{"http://abc.com", "bar", "http://abc.com/bar"},
		{"http://abc.com/foo/bar", "post/bar", "http://abc.com/foo/bar/post/bar"},
		{"http://abc.com/foo/bar", "post/bar/", "http://abc.com/foo/bar/post/bar/"},
	}

	for i, d := range data {
		output := MakePermalink(d.host, d.link).String()
		if d.output != output {
			t.Errorf("Test #%d failed. Expected %q got %q", i, d.output, output)
		}
	}
}

func TestURLPrep(t *testing.T) {
	type test struct {
		ugly   bool
		input  string
		output string
	}

	data := []test{
		{false, "/section/name.html", "/section/name/"},
		{true, "/section/name/index.html", "/section/name.html"},
	}
	for i, d := range data {
		output := URLPrep(d.ugly, d.input)
		if d.output != output {
			t.Errorf("Test #%d failed. Expected %q got %q", i, d.output, output)
		}
	}

}

func TestAddContextRoot(t *testing.T) {
	tests := []struct {
		baseURL  string
		url      string
		expected string
	}{
		{"http://example.com/sub/", "/foo", "/sub/foo"},
		{"http://example.com/sub/", "/foo/index.html", "/sub/foo/index.html"},
		{"http://example.com/sub1/sub2", "/foo", "/sub1/sub2/foo"},
		{"http://example.com", "/foo", "/foo"},
		// cannot guess that the context root is already added int the example below
		{"http://example.com/sub/", "/sub/foo", "/sub/sub/foo"},
		{"http://example.com/тря", "/трям/", "/тря/трям/"},
		{"http://example.com", "/", "/"},
		{"http://example.com/bar", "//", "/bar/"},
	}

	for _, test := range tests {
		output := AddContextRoot(test.baseURL, test.url)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}

func TestPretty(t *testing.T) {
	assert.Equal(t, PrettifyURLPath("/section/name.html"), "/section/name/index.html")
	assert.Equal(t, PrettifyURLPath("/section/sub/name.html"), "/section/sub/name/index.html")
	assert.Equal(t, PrettifyURLPath("/section/name/"), "/section/name/index.html")
	assert.Equal(t, PrettifyURLPath("/section/name/index.html"), "/section/name/index.html")
	assert.Equal(t, PrettifyURLPath("/index.html"), "/index.html")
	assert.Equal(t, PrettifyURLPath("/name.xml"), "/name/index.xml")
	assert.Equal(t, PrettifyURLPath("/"), "/")
	assert.Equal(t, PrettifyURLPath(""), "/")
	assert.Equal(t, PrettifyURL("/section/name.html"), "/section/name")
	assert.Equal(t, PrettifyURL("/section/sub/name.html"), "/section/sub/name")
	assert.Equal(t, PrettifyURL("/section/name/"), "/section/name")
	assert.Equal(t, PrettifyURL("/section/name/index.html"), "/section/name")
	assert.Equal(t, PrettifyURL("/index.html"), "/")
	assert.Equal(t, PrettifyURL("/name.xml"), "/name/index.xml")
	assert.Equal(t, PrettifyURL("/"), "/")
	assert.Equal(t, PrettifyURL(""), "/")
}

func TestUgly(t *testing.T) {
	assert.Equal(t, Uglify("/section/name.html"), "/section/name.html")
	assert.Equal(t, Uglify("/section/sub/name.html"), "/section/sub/name.html")
	assert.Equal(t, Uglify("/section/name/"), "/section/name.html")
	assert.Equal(t, Uglify("/section/name/index.html"), "/section/name.html")
	assert.Equal(t, Uglify("/index.html"), "/index.html")
	assert.Equal(t, Uglify("/name.xml"), "/name.xml")
	assert.Equal(t, Uglify("/"), "/")
	assert.Equal(t, Uglify(""), "/")
}
