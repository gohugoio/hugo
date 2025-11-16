// Copyright 2024 The Hugo Authors. All rights reserved.
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

package paths

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type pathBridge struct{}

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

// MakePermalink combines base URL with content path to create full URL paths.
// Example
//
//	base:   http://spf13.com/
//	path:   post/how-i-blog
//	result: http://spf13.com/post/how-i-blog
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
		panic(fmt.Errorf("can't make permalink from absolute link %q", plink))
	}

	base.Path = path.Join(base.Path, p.Path)
	base.Fragment = p.Fragment
	base.RawQuery = p.RawQuery

	// path.Join will strip off the last /, so put it back if it was there.
	hadTrailingSlash := (plink == "" && strings.HasSuffix(host, "/")) || strings.HasSuffix(p.Path, "/")
	if hadTrailingSlash && !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}

	return base
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

	// path strips trailing slash, ignore root path.
	if newPath != "/" && strings.HasSuffix(relativePath, "/") {
		newPath += "/"
	}
	return newPath
}

// URLizeAn

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
//
//	/section/name.html       becomes /section/name/index.html
//	/section/name/           becomes /section/name/index.html
//	/section/name/index.html becomes /section/name/index.html
func PrettifyURLPath(in string) string {
	return prettifyPath(in, pb)
}

// Uglify does the opposite of PrettifyURLPath().
//
//	/section/name/index.html becomes /section/name.html
//	/section/name/           becomes /section/name.html
//	/section/name.html       becomes /section/name.html
func Uglify(in string) string {
	if path.Ext(in) == "" {
		if len(in) < 2 {
			return "/"
		}
		// /section/name/  -> /section/name.html
		return path.Clean(in) + ".html"
	}

	name, ext := fileAndExt(in, pb)
	if name == "index" {
		// /section/name/index.html -> /section/name.html
		d := path.Dir(in)
		if len(d) > 1 {
			return d + ext
		}
		return in
	}
	// /.xml -> /index.xml
	if name == "" {
		return path.Dir(in) + "index" + ext
	}
	// /section/name.html -> /section/name.html
	return path.Clean(in)
}

// URLEscape escapes unicode letters.
func URLEscape(uri string) string {
	// escape unicode letters
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	return u.String()
}

// TrimExt trims the extension from a path..
func TrimExt(in string) string {
	return strings.TrimSuffix(in, path.Ext(in))
}

// From https://github.com/golang/go/blob/e0c76d95abfc1621259864adb3d101cf6f1f90fc/src/cmd/go/internal/web/url.go#L45
func UrlFromFilename(filename string) (*url.URL, error) {
	if !filepath.IsAbs(filename) {
		return nil, fmt.Errorf("filepath must be absolute")
	}

	// If filename has a Windows volume name, convert the volume to a host and prefix
	// per https://blogs.msdn.microsoft.com/ie/2006/12/06/file-uris-in-windows/.
	if vol := filepath.VolumeName(filename); vol != "" {
		if strings.HasPrefix(vol, `\\`) {
			filename = filepath.ToSlash(filename[2:])
			i := strings.IndexByte(filename, '/')

			if i < 0 {
				// A degenerate case.
				// \\host.example.com (without a share name)
				// becomes
				// file://host.example.com/
				return &url.URL{
					Scheme: "file",
					Host:   filename,
					Path:   "/",
				}, nil
			}

			// \\host.example.com\Share\path\to\file
			// becomes
			// file://host.example.com/Share/path/to/file
			return &url.URL{
				Scheme: "file",
				Host:   filename[:i],
				Path:   filepath.ToSlash(filename[i:]),
			}, nil
		}

		// C:\path\to\file
		// becomes
		// file:///C:/path/to/file
		return &url.URL{
			Scheme: "file",
			Path:   "/" + filepath.ToSlash(filename),
		}, nil
	}

	// /path/to/file
	// becomes
	// file:///path/to/file
	return &url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(filename),
	}, nil
}

// UrlStringToFilename converts the URL s to a filename.
// If ParseRequestURI fails, the input is just converted to OS specific slashes and returned.
func UrlStringToFilename(s string) (string, bool) {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return filepath.FromSlash(s), false
	}

	p := u.Path

	if p == "" {
		p, _ = url.QueryUnescape(u.Opaque)
		return filepath.FromSlash(p), false
	}

	if runtime.GOOS != "windows" {
		return p, true
	}

	if len(p) == 0 || p[0] != '/' {
		return filepath.FromSlash(p), false
	}

	p = filepath.FromSlash(p)

	if len(u.Host) == 1 {
		// file://c/Users/...
		return strings.ToUpper(u.Host) + ":" + p, true
	}

	if u.Host != "" && u.Host != "localhost" {
		if filepath.VolumeName(u.Host) != "" {
			return "", false
		}
		return `\\` + u.Host + p, true
	}

	if vol := filepath.VolumeName(p[1:]); vol == "" || strings.HasPrefix(vol, `\\`) {
		return "", false
	}

	return p[1:], true
}
