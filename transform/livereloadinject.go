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

	"github.com/spf13/viper"
)

func LiveReloadInject(ct contentTransformer) {
	endTags := []string{"</body>", "</BODY>", "</html>", "</HTML>"}
	port := viper.Get("port")
	replaceTemplate := `<script data-no-instant>document.write('<script src="/livereload.js?port=%d&mindelay=10"></' + 'script>')</script>%s`
	var newcontent []byte

	for _, endTag := range endTags {
		replacement := []byte(fmt.Sprintf(replaceTemplate, port, endTag))
		match := []byte(endTag)
		newcontent = bytes.Replace(ct.Content(), match, replacement, 1)
		if len(newcontent) != len(ct.Content()) {
			break
		}
	}
	if len(newcontent) == len(ct.Content()) {
		doctype := "<!DOCTYPE html>"
		if 0 == bytes.Index(newcontent, []byte(doctype)) {
			// this is an HTML document with body and html tag omited
			replace := fmt.Sprintf(replaceTemplate, port, "")
			newcontent = append(ct.Content(), replace...)
			// append livereload to end of document
		}
	}

	ct.Write(newcontent)
}
