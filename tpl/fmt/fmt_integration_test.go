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

package fmt_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

// Issue #11506
func TestErroridf(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
ignoreErrors = ['error-b']
-- layouts/index.html --
{{ erroridf "error-a" "%s" "a"}}
{{ erroridf "error-b" "%s" "b"}}
  `

	b, err := hugolib.TestE(t, files)

	b.Assert(err, qt.IsNotNil)
	b.AssertLogMatches(`^ERROR a\nYou can suppress this error by adding the following to your site configuration:\nignoreLogs = \['error-a'\]\n$`)
}

func TestWarnidf(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
ignoreLogs = ['warning-b']
-- layouts/index.html --
{{ warnidf "warning-a" "%s" "a"}}
{{ warnidf "warning-b" "%s" "b"}}
  `

	b := hugolib.Test(t, files, hugolib.TestOptWarn())
	b.AssertLogContains("WARN  a", "You can suppress this warning", "ignoreLogs", "['warning-a']")
	b.AssertLogNotContains("['warning-b']")
}
