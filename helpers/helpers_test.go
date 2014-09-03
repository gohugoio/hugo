package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPretty(t *testing.T) {
	assert.Equal(t, PrettifyPath("/section/name.html"), "/section/name/index.html")
	assert.Equal(t, PrettifyPath("/section/sub/name.html"), "/section/sub/name/index.html")
	assert.Equal(t, PrettifyPath("/section/name/"), "/section/name/index.html")
	assert.Equal(t, PrettifyPath("/section/name/index.html"), "/section/name/index.html")
	assert.Equal(t, PrettifyPath("/index.html"), "/index.html")
	assert.Equal(t, PrettifyPath("/name.xml"), "/name/index.xml")
	assert.Equal(t, PrettifyPath("/"), "/")
	assert.Equal(t, PrettifyPath(""), "/")
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

func TestMakePath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  Foo bar  ", "Foo-bar"},
		{"Foo.Bar/foo_Bar-Foo", "Foo.Bar/foo_Bar-Foo"},
		{"fOO,bar:foo%bAR", "fOObarfoobAR"},
		{"FOo/BaR.html", "FOo/BaR.html"},
		{"трям/трям", "трям/трям"},
		{"은행", "은행"},
		{"Банковский кассир", "Банковский-кассир"},
	}

	for _, test := range tests {
		output := MakePath(test.input)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}

func TestMakeToLower(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  foo bar  ", "foo-bar"},
		{"  Foo Bar  ", "foo-bar"},
		{"foo.bar/foo_bar-foo", "foo.bar/foo_bar-foo"},
		{"foo,bar:foo%bar", "foobarfoobar"},
		{"foo/bar.html", "foo/bar.html"},
		{"трям/трям", "трям/трям"},
		{"은행", "은행"},
	}

	for _, test := range tests {
		output := MakePathToLower(test.input)
		if output != test.expected {
			t.Errorf("Expected %#v, got %#v\n", test.expected, output)
		}
	}
}

func TestUrlize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  foo bar  ", "foo-bar"},
		{"Foo And BAR", "foo-and-bar"},
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
