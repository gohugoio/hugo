// Copyright 2026 The Hugo Authors. All rights reserved.
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

package urls_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// See issue 15473.
func TestJoinPathEmptySlice(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = 'https://example.com'
-- layouts/index.html --
|{{ urls.JoinPath (slice) }}|
|{{ urls.JoinPath (slice) (slice) }}|
|{{ urls.JoinPath (slice) "a" (slice) }}|
|{{ urls.JoinPath (slice "a" "b") }}|
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "||\n||\n|a|\n|a/b|")
}
