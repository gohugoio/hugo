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

package images_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestImageConfigFromModule(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
baseURL = 'http://example.com/'
theme = ["mytheme"]
-- static/images/pixel1.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- themes/mytheme/static/images/pixel2.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- layouts/index.html --
{{ $path := "static/images/pixel1.png" }}
fileExists OK: {{ fileExists $path }}|
imageConfig OK: {{ (imageConfig $path).Width }}|
{{ $path2 := "static/images/pixel2.png" }}
fileExists2 OK: {{ fileExists $path2 }}|
imageConfig2 OK: {{ (imageConfig $path2).Width }}|

  `

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
fileExists OK: true|
imageConfig OK: 1|
fileExists2 OK: true|
imageConfig2 OK: 1|
`)
}
