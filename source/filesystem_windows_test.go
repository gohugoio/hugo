package source

//
// NOTE, any changes here need to be reflected in filesystem_linux_test.go
//

// Note the case of the volume drive.  It must be the same in all examples.
var platformBase = "C:\\foo\\"
var platformPaths = []TestPath{
	{"foobar", "foobar", "aaa", "", ""},
	{"b\\1file", "1file", "aaa", "b", "b\\"},
	{"c\\d\\2file", "2file", "aaa", "c", "c\\d\\"},
	{"C:\\foo\\e\\f\\3file", "3file", "aaa", "e", "e\\f\\"}, // note volume case is equal to platformBase
	{"section\\foo.rss", "foo.rss", "aaa", "section", "section\\"},
}
