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

package helpers

import (
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/paths"
)

// URLize is similar to MakePath, but with Unicode handling
// Example:
//
//	uri: Vim (text editor)
//	urlize: vim-text-editor
func (p *PathSpec) URLize(uri string) string {
	return p.URLEscape(p.MakePathSanitized(uri))
}

// URLizeFilename creates an URL from a filename by escaping unicode letters
// and turn any filepath separator into forward slashes.
func (p *PathSpec) URLizeFilename(filename string) string {
	return p.URLEscape(filepath.ToSlash(filename))
}

// URLEscape escapes unicode letters.
func (p *PathSpec) URLEscape(uri string) string {
	// escape unicode letters
	parsedURI, err := url.Parse(uri)
	if err != nil {
		// if net/url can not parse URL it means Sanitize works incorrectly
		panic(err)
	}
	x := parsedURI.String()
	return x
}

// AbsURL creates an absolute URL from the relative path given and the BaseURL set in config.
func (p *PathSpec) AbsURL(in string, addLanguage bool) string {
	isAbs, err := p.IsAbsURL(in)
	if err != nil {
		return in
	}
	if isAbs || strings.HasPrefix(in, "//") {
		// It  is already  absolute, return it as is.
		return in
	}

	baseURL := p.getBaseURLRoot(in)

	if addLanguage {
		prefix := p.GetLanguagePrefix()
		if prefix != "" {
			hasPrefix := false
			// avoid adding language prefix if already present
			in2 := in
			if strings.HasPrefix(in, "/") {
				in2 = in[1:]
			}
			if in2 == prefix {
				hasPrefix = true
			} else {
				hasPrefix = strings.HasPrefix(in2, prefix+"/")
			}

			if !hasPrefix {
				addSlash := in == "" || strings.HasSuffix(in, "/")
				in = path.Join(prefix, in)

				if addSlash {
					in += "/"
				}
			}
		}
	}

	return paths.MakePermalink(baseURL, in).String()
}

func (p *PathSpec) getBaseURLRoot(path string) string {
	if strings.HasPrefix(path, "/") {
		// Treat it as relative to the server root.
		return p.Cfg.BaseURL().WithoutPath
	} else {
		// Treat it as relative to the baseURL.
		return p.Cfg.BaseURL().WithPath
	}
}

func (p *PathSpec) IsAbsURL(in string) (bool, error) {
	// Fast path.
	if strings.HasPrefix(in, "http://") || strings.HasPrefix(in, "https://") {
		return true, nil
	}
	u, err := url.Parse(in)
	if err != nil {
		return false, err
	}
	return u.IsAbs(), nil
}

func (p *PathSpec) RelURL(in string, addLanguage bool) string {
	isAbs, err := p.IsAbsURL(in)
	if err != nil {
		return in
	}
	baseURL := p.getBaseURLRoot(in)
	canonifyURLs := p.Cfg.CanonifyURLs()

	if (!strings.HasPrefix(in, baseURL) && isAbs) || strings.HasPrefix(in, "//") {
		return in
	}

	u := in

	if strings.HasPrefix(in, baseURL) {
		u = strings.TrimPrefix(u, baseURL)
	}

	if addLanguage {
		prefix := p.GetLanguagePrefix()
		if prefix != "" {
			hasPrefix := false
			// avoid adding language prefix if already present
			in2 := in
			if strings.HasPrefix(in, "/") {
				in2 = in[1:]
			}
			if in2 == prefix {
				hasPrefix = true
			} else {
				hasPrefix = strings.HasPrefix(in2, prefix+"/")
			}

			if !hasPrefix {
				hadSlash := strings.HasSuffix(u, "/")

				u = path.Join(prefix, u)

				if hadSlash {
					u += "/"
				}
			}
		}
	}

	if !canonifyURLs {
		u = paths.AddContextRoot(baseURL, u)
	}

	if in == "" && !strings.HasSuffix(u, "/") && strings.HasSuffix(baseURL, "/") {
		u += "/"
	}

	if !strings.HasPrefix(u, "/") {
		u = "/" + u
	}

	return u
}

// PrependBasePath prepends any baseURL sub-folder to the given resource
func (p *PathSpec) PrependBasePath(rel string, isAbs bool) string {
	basePath := p.GetBasePath(!isAbs)
	if basePath != "" {
		rel = filepath.ToSlash(rel)
		// Need to prepend any path from the baseURL
		hadSlash := strings.HasSuffix(rel, "/")
		rel = path.Join(basePath, rel)
		if hadSlash {
			rel += "/"
		}
	}
	return rel
}
