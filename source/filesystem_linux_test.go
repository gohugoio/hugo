package source

//
// NOTE, any changes here need to be reflected in filesystem_windows_test.go
//
var platformBase = "foo/bar/boo/"
var platformPaths = []TestPath{
	{"foobar", "foobar", "aaa", "", ""},
	{"b/1file", "1file", "aaa", "b", "b/"},
	{"c/d/2file", "2file", "aaa", "d", "c/d/"},
	{"/e/f/3file", "3file", "aaa", "f", "e/f/"},
	{"section\\foo.rss", "foo.rss", "aaa", "section", "section/"},
}
