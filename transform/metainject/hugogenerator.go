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

package metainject

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/transform"
)

var metaTagsCheck = regexp.MustCompile(`(?i)<meta\s+name=['|"]?generator['|"]?`)
var hugoGeneratorTag = fmt.Sprintf(`<meta name="generator" content="Hugo %s" />`, hugo.CurrentVersion)

// HugoGenerator injects a meta generator tag for Hugo if none present.
func HugoGenerator(ft transform.FromTo) error {
	b := ft.From().Bytes()
	if metaTagsCheck.Match(b) {
		if _, err := ft.To().Write(b); err != nil {
			helpers.DistinctWarnLog.Println("Failed to inject Hugo generator tag:", err)
		}
		return nil
	}

	head := "<head>"
	replace := []byte(fmt.Sprintf("%s\n\t%s", head, hugoGeneratorTag))
	newcontent := bytes.Replace(b, []byte(head), replace, 1)

	if len(newcontent) == len(b) {
		head := "<HEAD>"
		replace := []byte(fmt.Sprintf("%s\n\t%s", head, hugoGeneratorTag))
		newcontent = bytes.Replace(b, []byte(head), replace, 1)
	}

	if _, err := ft.To().Write(newcontent); err != nil {
		helpers.DistinctWarnLog.Println("Failed to inject Hugo generator tag:", err)
	}

	return nil

}
