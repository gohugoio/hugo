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

type PathBridge struct {
}

func (PathBridge) Base(in string) string {
	return path.Base(in)
}

func (PathBridge) Clean(in string) string {
	return path.Clean(in)
}

func (PathBridge) Dir(in string) string {
	return path.Dir(in)
}

func (PathBridge) Ext(in string) string {
	return path.Ext(in)
}

func (PathBridge) Join(elem ...string) string {
	return path.Join(elem...)
}

func (PathBridge) Separator() string {
	return "/"
}

var pathBridge PathBridge

func sanitizeUrlWithFlags(in string, f purell.NormalizationFlags) string {
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
	// and restore SanitizeUrl() to the way it was.
	//                         -- @anthonyfok, 2015-02-16
	//
	// Begin temporary kludge
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	if !strings.HasPrefix(u.Path, "/") {
		u.Path = "/" + u.Path
	}
	return u.String()
	// End temporary kludge

	//return s

}

// SanitizeUrl sanitizes the input URL string.
func SanitizeUrl(in string) string {
	return sanitizeUrlWithFlags(in, purell.FlagsSafe|purell.FlagRemoveTrailingSlash|purell.FlagRemoveDotSegments|purell.FlagRemoveDuplicateSlashes|purell.FlagRemoveUnnecessaryHostDots|purell.FlagRemoveEmptyPortSeparator)
}

// SanitizeUrlKeepTrailingSlash is the same as SanitizeUrl, but will keep any trailing slash.
func SanitizeUrlKeepTrailingSlash(in string) string {
	return sanitizeUrlWithFlags(in, purell.FlagsSafe|purell.FlagRemoveDotSegments|purell.FlagRemoveDuplicateSlashes|purell.FlagRemoveUnnecessaryHostDots|purell.FlagRemoveEmptyPortSeparator)
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

	// path strips traling slash, ignore root path.
	if newPath != "/" && strings.HasSuffix(relativePath, "/") {
		newPath += "/"
	}
	return newPath
}

func UrlizeAndPrep(in string) string {
	return UrlPrep(viper.GetBool("UglyUrls"), Urlize(in))
}

func UrlPrep(ugly bool, in string) string {
	if ugly {
		x := Uglify(SanitizeUrl(in))
		return x
	}
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
	return PrettiyPath(in, pathBridge)
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
	}

	name, ext := FileAndExt(in, pathBridge)
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
