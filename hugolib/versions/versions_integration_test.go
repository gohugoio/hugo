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

package versions_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

func TestDefaultContentVersionDoesNotExist(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["rss", "sitemap", "section", "taxonomy", "term"]
defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true
defaultContentVersion = "doesnotexist"
[versions]
[versions."v1.0.0"]
`
	b, err := hugolib.TestE(t, files)
	b.Assert(err, qt.ErrorMatches, `.*failed to decode "versions": the configured defaultContentVersion "doesnotexist" does not exist`)
}
