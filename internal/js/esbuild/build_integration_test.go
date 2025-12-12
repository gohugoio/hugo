// Copyright 2025 The Hugo Authors. All rights reserved.
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

// Package js provides functions for building JavaScript resources
package esbuild_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestBuilSameJSWithDifferentOptionsIssue14254(t *testing.T) {
	files := `
-- hugo.toml --
-- assets/js/app1.js --
import * as params from '@params';
console.log(params);
-- layouts/home.html --
{{ $opts1 := dict "targetPath" "app1_1.js" "minify" true "params" (dict "key" "value1") }}
{{ $opts2 := dict "targetPath" "app1_2.js" "minify" true "params" (dict "key" "value2") }}
{{ $js1 := resources.Get "js/app1.js" | js.Build $opts1 }}
js1: {{ $js1.RelPermalink }}|{{ $js1.Content | safeHTML }}$
{{ $js2 := resources.Get "js/app1.js" | js.Build $opts2 }}
js2: {{ $js2.RelPermalink }}|{{ $js2.Content | safeHTML }}$

`
	b := hugolib.Test(t, files)

	b.AssertFileContentRe("public/index.html",
		`js1.*app1_1.*value1`,
		`js2.*app1_2.*value2`,
	)
}
