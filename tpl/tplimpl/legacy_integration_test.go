// Copyright 2025 The Hugo Authors. All rights reserved.
//
// Portions Copyright The Go Authors.

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

package tplimpl_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestLegacyPartialIssue13599(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- layouts/partials/mypartial.html --
Mypartial.
-- layouts/_default/index.html --
mypartial:   {{ template "partials/mypartial.html" . }}

`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "Mypartial.")
}
