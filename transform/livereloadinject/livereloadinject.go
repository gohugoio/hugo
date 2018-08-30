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
	"bytes"
	"fmt"

	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/transform"
)

// New creates a function that can be used
// to inject a script tag for the livereload JavaScript in a HTML document.
func New(port int) transform.Transformer {
	return func(ft transform.FromTo) error {
		b := ft.From().Bytes()
		endBodyTag := "</body>"
		match := []byte(endBodyTag)
		replaceTemplate := `<script data-no-instant>document.write('<script src="/livereload.js?port=%d&mindelay=10"></' + 'script>')</script>%s`
		replace := []byte(fmt.Sprintf(replaceTemplate, port, endBodyTag))

		newcontent := bytes.Replace(b, match, replace, 1)
		if len(newcontent) == len(b) {
			endBodyTag = "</BODY>"
			replace := []byte(fmt.Sprintf(replaceTemplate, port, endBodyTag))
			match := []byte(endBodyTag)
			newcontent = bytes.Replace(b, match, replace, 1)
		}

		if _, err := ft.To().Write(newcontent); err != nil {
			helpers.DistinctWarnLog.Println("Failed to inject LiveReload script:", err)
		}
		return nil
	}
}
