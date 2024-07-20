// Copyright 2018 The Hugo Authors. All rights reserved.
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

package livereloadinject

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/transform"
)

var (
	ignoredSyntax  = regexp.MustCompile(`(?s)^(?:\s+|<!--.*?-->|<\?.*?\?>)*`)
	tagsBeforeHead = []*regexp.Regexp{
		regexp.MustCompile(`(?is)^<!doctype\s[^>]*>`),
		regexp.MustCompile(`(?is)^<html(?:\s[^>]*)?>`),
		regexp.MustCompile(`(?is)^<head(?:\s[^>]*)?>`),
	}
)

// New creates a function that can be used to inject a script tag for
// the livereload JavaScript at the start of an HTML document's head.
func New(baseURL *url.URL) transform.Transformer {
	return func(ft transform.FromTo) error {
		b := ft.From().Bytes()

		// We find the start of the head by reading past (in order)
		// the doctype declaration, HTML start tag and head start tag,
		// all of which are optional, and any whitespace, comments, or
		// XML instructions in-between.
		idx := 0
		for _, tag := range tagsBeforeHead {
			idx += len(ignoredSyntax.Find(b[idx:]))
			idx += len(tag.Find(b[idx:]))
		}

		path := strings.TrimSuffix(baseURL.Path, "/")

		src := path + "/livereload.js?mindelay=10&v=2"
		src += "&port=" + baseURL.Port()
		src += "&path=" + strings.TrimPrefix(path+"/livereload", "/")

		script := []byte(fmt.Sprintf(`<script src="%s" data-no-instant defer></script>`, html.EscapeString(src)))

		c := make([]byte, len(b)+len(script))
		copy(c, b[:idx])
		copy(c[idx:], script)
		copy(c[idx+len(script):], b[idx:])

		if _, err := ft.To().Write(c); err != nil {
			loggers.Log().Warnf("Failed to inject LiveReload script:", err)
		}
		return nil
	}
}
