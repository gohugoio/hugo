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

package hstring

import (
	"html/template"

	"github.com/gohugoio/hugo/common/types"
)

var _ types.PrintableValueProvider = HTML("")

// HTML is a string that represents rendered HTML.
// When printed in templates it will be rendered as template.HTML and considered safe so no need to pipe it into `safeHTML`.
// This type was introduced as a wasy to prevent a common case of inifinite recursion in the template rendering
// when the `linkify` option is enabled with a common (wrong) construct like `{{ .Text | .Page.RenderString }}` in a hook template.
type HTML string

func (s HTML) String() string {
	return string(s)
}

func (s HTML) PrintableValue() any {
	return template.HTML(s)
}
