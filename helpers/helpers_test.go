package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
		{"  foo bar  ", "foo-bar"},
		{"foo.bar/foo_bar-foo", "foo.bar/foo_bar-foo"},
		{"foo,bar:foo%bar", "foobarfoobar"},
		{"foo/bar.html", "foo/bar.html"},
		{"трям/трям", "трям/трям"},
	}

	for _, test := range tests {
		output := MakePath(test.input)
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
