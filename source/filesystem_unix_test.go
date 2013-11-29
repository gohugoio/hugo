// +build linux darwin !windows

package source

//
// NOTE, any changes here need to be reflected in filesystem_windows_test.go
//
var platformBase = "/base/"
var platformPaths = []TestPath{
	{"foobar", "foobar", "aaa", "", ""},
	{"b/1file", "1file", "aaa", "b", "b/"},
	{"c/d/2file", "2file", "aaa", "c", "c/d/"},
	{"/base/e/f/3file", "3file", "aaa", "e", "e/f/"},
	{"section/foo.rss", "foo.rss", "aaa", "section", "section/"},
}
