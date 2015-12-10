// Copyright 2015 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
