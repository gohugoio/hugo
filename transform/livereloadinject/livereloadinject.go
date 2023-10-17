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

type insertPoint struct {
	regexp  *regexp.Regexp
	prepend bool
}

var insertPoints = []insertPoint{
	{regexp: regexp.MustCompile(`(?is)<head(?:\s[^>]*)?>`)},
	{regexp: regexp.MustCompile(`(?is)<html(?:\s[^>]*)?>`)},
	{regexp: regexp.MustCompile(`(?is)<!doctype(?:\s[^>]*)?>`)},
	{regexp: regexp.MustCompile(`(?is)<[a-z]`), prepend: true},
}

// New creates a function that can be used
// to inject a script tag for the livereload JavaScript in a HTML document.
func New(baseURL url.URL) transform.Transformer {
	return func(ft transform.FromTo) error {
		b := ft.From().Bytes()
		idx := -1

		for _, p := range insertPoints {
			if match := p.regexp.FindIndex(b); match != nil {
				if p.prepend {
					idx = match[0]
				} else {
					idx = match[1]
				}
				break
			}
		}

		if idx == -1 {
			loggers.Log().Warnf("Failed to find location to inject LiveReload script.")
			return nil
		}

		path := strings.TrimSuffix(baseURL.Path, "/")

		src := path + "/livereload.js?mindelay=10&v=2"
		src += "&port=" + baseURL.Port()
		src += "&path=" + strings.TrimPrefix(path+"/livereload", "/")

		c := make([]byte, len(b))
		copy(c, b)

		script := []byte(fmt.Sprintf(`<script src="%s" data-no-instant defer></script>`, html.EscapeString(src)))

		c = append(c[:idx], append(script, c[idx:]...)...)

		if _, err := ft.To().Write(c); err != nil {
			loggers.Log().Warnf("Failed to inject LiveReload script:", err)
		}
		return nil
	}
}
