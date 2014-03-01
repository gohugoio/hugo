// Copyright © 2013 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helpers

import (
	"fmt"
	"net/url"
	"path"
)

var _ = fmt.Println

// Similar to MakePath, but with Unicode handling
// Example:
//     uri: Vim (text editor)
//     urlize: vim-text-editor
func Urlize(uri string) string {
	sanitized := MakePath(uri)

	// escape unicode letters
	parsedUri, err := url.Parse(sanitized)
	if err != nil {
		// if net/url can not parse URL it's meaning Sanitize works incorrect
		panic(err)
	}
	x := parsedUri.String()
	return x
}

// Combines a base with a path
// Example
//    base:   http://spf13.com/
//    path:   post/how-i-blog
//    result: http://spf13.com/post/how-i-blog
func MakePermalink(host, plink string) *url.URL {

	base, err := url.Parse(host)
	if err != nil {
		panic(err)
	}

	path, err := url.Parse(plink)
	if err != nil {
		panic(err)
	}
	return base.ResolveReference(path)
}

func UrlPrep(ugly bool, in string) string {
	if ugly {
		return Uglify(in)
	} else {
		return PrettifyUrl(in)
	}
}

// Don't Return /index.html portion.
func PrettifyUrl(in string) string {
	x := PrettifyPath(in)

	if path.Base(x) == "index.html" {
		return path.Dir(x)
	}

	if in == "" {
		return "/"
	}

	return x
}

// /section/name/index.html -> /section/name.html
// /section/name/  -> /section/name.html
// /section/name.html -> /section/name.html
func Uglify(in string) string {
	if path.Ext(in) == "" {
		if len(in) < 2 {
			return "/"
		}
		// /section/name/  -> /section/name.html
		return path.Clean(in) + ".html"
	} else {
		name, ext := FileAndExt(in)
		if name == "index" {
			// /section/name/index.html -> /section/name.html
			d := path.Dir(in)
			if len(d) > 1 {
				return d + ext
			} else {
				return in
			}
		} else {
			// /section/name.html -> /section/name.html
			return path.Clean(in)
		}
	}
}
