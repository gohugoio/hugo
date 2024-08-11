// Copyright 2019 The Hugo Authors. All rights reserved.
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

package template

import (
	"github.com/gohugoio/hugo/common/types"
	template "github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate"
)

/*

This files contains the Hugo related addons. All the other files in this
package is auto generated.

*/

// Export it so we can populate Hugo's func map with it, which makes it faster.
var GoFuncs = funcMap

// Prepare returns a template ready for execution.
func (t *Template) Prepare() (*template.Template, error) {
	if err := t.escape(); err != nil {
		return nil, err
	}
	return t.text, nil
}

// See https://github.com/golang/go/issues/5884
func StripTags(html string) string {
	return stripTags(html)
}

func indirect(a any) any {
	in := doIndirect(a)

	// We have a special Result type that we want to unwrap when printed.
	if pp, ok := in.(types.PrintableValueProvider); ok {
		return pp.PrintableValue()
	}

	return in
}
