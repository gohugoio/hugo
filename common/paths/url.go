// Copyright 2021 The Hugo Authors. All rights reserved.
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
	"strings"
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

// MakePermalink combines base URL with content path to create full URL paths.
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
		panic(fmt.Errorf("can't make permalink from absolute link %q", plink))
	}

	base.Path = path.Join(base.Path, p.Path)

	// path.Join will strip off the last /, so put it back if it was there.
	hadTrailingSlash := (plink == "" && strings.HasSuffix(host, "/")) || strings.HasSuffix(p.Path, "/")
	if hadTrailingSlash && !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}

	return base
}

// IsAbsURL determines whether the given path points to an absolute URL.
func IsAbsURL(path string) bool {
	url, err := url.Parse(path)
	if err != nil {
		return false
	}

	return url.IsAbs() || strings.HasPrefix(path, "//")
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

// URLEscape escapes unicode letters.
func URLEscape(uri string) string {
	// escape unicode letters
	u, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}
	return u.String()
}
