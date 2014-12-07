package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlize(t *testing.T) {
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
		output := Urlize(test.input)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
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

func TestUrlPrep(t *testing.T) {
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
		output := UrlPrep(d.ugly, d.input)
		if d.output != output {
			t.Errorf("Test #%d failed. Expected %q got %q", i, d.output, output)
		}
	}

}

func TestPretty(t *testing.T) {
	assert.Equal(t, PrettifyUrlPath("/section/name.html"), "/section/name/index.html")
	assert.Equal(t, PrettifyUrlPath("/section/sub/name.html"), "/section/sub/name/index.html")
	assert.Equal(t, PrettifyUrlPath("/section/name/"), "/section/name/index.html")
	assert.Equal(t, PrettifyUrlPath("/section/name/index.html"), "/section/name/index.html")
	assert.Equal(t, PrettifyUrlPath("/index.html"), "/index.html")
	assert.Equal(t, PrettifyUrlPath("/name.xml"), "/name/index.xml")
	assert.Equal(t, PrettifyUrlPath("/"), "/")
	assert.Equal(t, PrettifyUrlPath(""), "/")
	assert.Equal(t, PrettifyUrl("/section/name.html"), "/section/name")
	assert.Equal(t, PrettifyUrl("/section/sub/name.html"), "/section/sub/name")
	assert.Equal(t, PrettifyUrl("/section/name/"), "/section/name")
	assert.Equal(t, PrettifyUrl("/section/name/index.html"), "/section/name")
	assert.Equal(t, PrettifyUrl("/index.html"), "/")
	assert.Equal(t, PrettifyUrl("/name.xml"), "/name/index.xml")
	assert.Equal(t, PrettifyUrl("/"), "/")
	assert.Equal(t, PrettifyUrl(""), "/")
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
