// Copyright © 2013-2015 Steve Francia <spf@spf13.com>.
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
	"strings"

	"github.com/PuerkitoBio/purell"
	"github.com/spf13/viper"
)

type pathBridge struct {
}

func (pathBridge) Base(in string) string {
	return path.Base(in)
}

func (pathBridge) Clean(in string) string {
	return path.Clean(in)
}

func (pathBridge) Dir(in string) string {
	return path.Dir(in)
}

func (pathBridge) Ext(in string) string {
	return path.Ext(in)
}

func (pathBridge) Join(elem ...string) string {
	return path.Join(elem...)
}

func (pathBridge) Separator() string {
	return "/"
}

var pb pathBridge

func sanitizeURLWithFlags(in string, f purell.NormalizationFlags) string {
	s, err := purell.NormalizeURLString(in, f)
	if err != nil {
		return in
	}

	// Temporary workaround for the bug fix and resulting
	// behavioral change in purell.NormalizeURLString():
	// a leading '/' was inadvertently to relative links,
	// but no longer, see #878.
	//
	// I think the real solution is to allow Hugo to
	// make relative URL with relative path,
	// e.g. "../../post/hello-again/", as wished by users
	// in issues #157, #622, etc., without forcing
	// relative URLs to begin with '/'.
	// Once the fixes are in, let's remove this kludge
	// and restore SanitizeURL() to the way it was.
	//                         -- @anthonyfok, 2015-02-16
	//
	// Begin temporary kludge
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	if len(u.Path) > 0 && !strings.HasPrefix(u.Path, "/") {
		u.Path = "/" + u.Path
	}
	return u.String()
	// End temporary kludge

	//return s

}

// SanitizeURL sanitizes the input URL string.
func SanitizeURL(in string) string {
	return sanitizeURLWithFlags(in, purell.FlagsSafe|purell.FlagRemoveTrailingSlash|purell.FlagRemoveDotSegments|purell.FlagRemoveDuplicateSlashes|purell.FlagRemoveUnnecessaryHostDots|purell.FlagRemoveEmptyPortSeparator)
}

// SanitizeURLKeepTrailingSlash is the same as SanitizeURL, but will keep any trailing slash.
func SanitizeURLKeepTrailingSlash(in string) string {
	return sanitizeURLWithFlags(in, purell.FlagsSafe|purell.FlagRemoveDotSegments|purell.FlagRemoveDuplicateSlashes|purell.FlagRemoveUnnecessaryHostDots|purell.FlagRemoveEmptyPortSeparator)
}

// Similar to MakePath, but with Unicode handling
// Example:
//     uri: Vim (text editor)
//     urlize: vim-text-editor
func URLize(uri string) string {
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
	hadTrailingSlash := (plink == "" && strings.HasSuffix(host, "/")) || strings.HasSuffix(p.Path, "/")
	if hadTrailingSlash && !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}

	return base
}

// AbsURL creates a absolute URL from the relative path given and the BaseURL set in config.
func AbsURL(path string) string {
	if strings.HasPrefix(path, "http") || strings.HasPrefix(path, "//") {
		return path
	}
	return MakePermalink(viper.GetString("BaseURL"), path).String()
}

// RelURL creates a URL relative to the BaseURL root.
// Note: The result URL will not include the context root if canonifyURLs is enabled.
func RelURL(path string) string {
	baseURL := viper.GetString("BaseURL")
	canonifyURLs := viper.GetBool("canonifyURLs")
	if (!strings.HasPrefix(path, baseURL) && strings.HasPrefix(path, "http")) || strings.HasPrefix(path, "//") {
		return path
	}

	u := path

	if strings.HasPrefix(path, baseURL) {
		u = strings.TrimPrefix(u, baseURL)
	}

	if !canonifyURLs {
		u = AddContextRoot(baseURL, u)
	}
	if path == "" && !strings.HasSuffix(u, "/") && strings.HasSuffix(baseURL, "/") {
		u += "/"
	}

	if !strings.HasPrefix(u, "/") {
		u = "/" + u
	}

	return u
}

// AddContextRoot adds the context root to an URL if it's not already set.
// For relative URL entries on sites with a base url with a context root set (i.e. http://example.com/mysite),
// relative URLs must not include the context root if canonifyURLs is enabled. But if it's disabled, it must be set.
func AddContextRoot(baseURL, relativePath string) string {

	url, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	newPath := path.Join(url.Path, relativePath)

	// path strips traling slash, ignore root path.
	if newPath != "/" && strings.HasSuffix(relativePath, "/") {
		newPath += "/"
	}
	return newPath
}

func URLizeAndPrep(in string) string {
	return URLPrep(viper.GetBool("UglyURLs"), URLize(in))
}

func URLPrep(ugly bool, in string) string {
	if ugly {
		x := Uglify(SanitizeURL(in))
		return x
	}
	x := PrettifyURL(SanitizeURL(in))
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

// PrettifyURL takes a URL string and returns a semantic, clean URL.
func PrettifyURL(in string) string {
	x := PrettifyURLPath(in)

	if path.Base(x) == "index.html" {
		return path.Dir(x)
	}

	if in == "" {
		return "/"
	}

	return x
}

// PrettifyURLPath takes a URL path to a content and converts it
// to enable pretty URLs.
//     /section/name.html       becomes /section/name/index.html
//     /section/name/           becomes /section/name/index.html
//     /section/name/index.html becomes /section/name/index.html
func PrettifyURLPath(in string) string {
	return PrettiyPath(in, pb)
}

// Uglify does the opposite of PrettifyURLPath().
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
	}

	name, ext := FileAndExt(in, pb)
	if name == "index" {
		// /section/name/index.html -> /section/name.html
		d := path.Dir(in)
		if len(d) > 1 {
			return d + ext
		}
		return in
	}
	// /section/name.html -> /section/name.html
	return path.Clean(in)
}
