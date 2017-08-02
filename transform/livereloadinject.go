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

package transform

import (
	"bytes"
	"fmt"
)

// LiveReloadInject returns a function that can be used
// to inject a script tag for the livereload JavaScript in a HTML document.
func LiveReloadInject(port int) func(ct contentTransformer) {
	return func(ct contentTransformer) {
		endBodyTag := "</body>"
		match := []byte(endBodyTag)
		replaceTemplate := `<script data-no-instant>document.write('<script src="/livereload.js?port=%d&mindelay=10"></' + 'script>')</script>%s`
		replace := []byte(fmt.Sprintf(replaceTemplate, port, endBodyTag))

		newcontent := bytes.Replace(ct.Content(), match, replace, 1)
		if len(newcontent) == len(ct.Content()) {
			endBodyTag = "</BODY>"
			replace := []byte(fmt.Sprintf(replaceTemplate, port, endBodyTag))
			match := []byte(endBodyTag)
			newcontent = bytes.Replace(ct.Content(), match, replace, 1)
		}

		ct.Write(newcontent)
	}
}
