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
	"github.com/PuerkitoBio/purell"
	"net/url"
	"path"
	"strings"
)

// SanitizeUrl sanitizes the input URL string.
func SanitizeUrl(in string) string {
	url, err := purell.NormalizeURLString(in, purell.FlagsSafe|purell.FlagRemoveTrailingSlash|purell.FlagRemoveDotSegments|purell.FlagRemoveDuplicateSlashes|purell.FlagRemoveUnnecessaryHostDots|purell.FlagRemoveEmptyPortSeparator)
	if err != nil {
		return in
	}
	return url
}

// Similar to MakePath, but with Unicode handling
// Example:
//     uri: Vim (text editor)
//     urlize: vim-text-editor
func Urlize(uri string) string {
	sanitized := MakePathToLower(uri)

	// escape unicode letters
	parsedUri, err := url.Parse(sanitized)
	if err != nil {
		// if net/url can not parse URL it's meaning Sanitize works incorrect
		panic(err)
	}
	x := parsedUri.String()
	return x
}

// Combines base URL with content path to create full URL paths.
// Example
//    base:   http://spf13.com/
//    path:   post/how-i-blog
//    result: http://spf13.com/post/how-i-blog
func MakePermalink(host, plink string) *url.URL {

	base, err := url.Parse(host)
	if err != nil {
		panic(err)
	}

	p, err := url.Parse(plink)
	if err != nil {
		panic(err)
	}

	if p.Host != "" {
		panic(fmt.Errorf("Can't make permalink from absolute link %q", plink))
	}

	base.Path = path.Join(base.Path, p.Path)

	// path.Join will strip off the last /, so put it back if it was there.
	if strings.HasSuffix(p.Path, "/") && !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}

	return base
}

// AddContextRoot adds the context root to an URL if it's not already set.
// For relative URL entries on sites with a base url with a context root set (i.e. http://example.com/mysite),
// relative URLs must not include the context root if canonifyUrls is enabled. But if it's disabled, it must be set.
func AddContextRoot(baseUrl, relativePath string) string {

	url, err := url.Parse(baseUrl)
	if err != nil {
		panic(err)
	}

	newPath := path.Join(url.Path, relativePath)

	// path strips traling slash
	if strings.HasSuffix(relativePath, "/") {
		newPath += "/"
	}
	return newPath
}

func UrlPrep(ugly bool, in string) string {
	if ugly {
		x := Uglify(SanitizeUrl(in))
		return x
	} else {
		x := PrettifyUrl(SanitizeUrl(in))
		if path.Ext(x) == ".xml" {
			return x
		}
		url, err := purell.NormalizeURLString(x, purell.FlagAddTrailingSlash)
		if err != nil {
			fmt.Printf("ERROR returned by NormalizeURLString. Returning in = %q\n", in)
			return in
		}
		return url
	}
}

// PrettifyUrl takes a URL string and returns a semantic, clean URL.
func PrettifyUrl(in string) string {
	x := PrettifyUrlPath(in)

	if path.Base(x) == "index.html" {
		return path.Dir(x)
	}

	if in == "" {
		return "/"
	}

	return x
}

// PrettifyUrlPath takes a URL path to a content and converts it
// to enable pretty URLs.
//     /section/name.html       becomes /section/name/index.html
//     /section/name/           becomes /section/name/index.html
//     /section/name/index.html becomes /section/name/index.html
func PrettifyUrlPath(in string) string {
	if path.Ext(in) == "" {
		// /section/name/  -> /section/name/index.html
		if len(in) < 2 {
			return "/"
		}
		return path.Join(path.Clean(in), "index.html")
	} else {
		name, ext := ResourceAndExt(in)
		if name == "index" {
			// /section/name/index.html -> /section/name/index.html
			return path.Clean(in)
		} else {
			// /section/name.html -> /section/name/index.html
			return path.Join(path.Dir(in), name, "index"+ext)
		}
	}
}

// Uglify does the opposite of PrettifyUrlPath().
//     /section/name/index.html becomes /section/name.html
//     /section/name/           becomes /section/name.html
//     /section/name.html       becomes /section/name.html
func Uglify(in string) string {
	if path.Ext(in) == "" {
		if len(in) < 2 {
			return "/"
		}
		// /section/name/  -> /section/name.html
		return path.Clean(in) + ".html"
	} else {
		name, ext := ResourceAndExt(in)
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

// Same as FileAndExt, but for URLs.
func ResourceAndExt(in string) (name string, ext string) {
	ext = path.Ext(in)
	base := path.Base(in)

	return extractFilename(in, ext, base, "/"), ext
}
