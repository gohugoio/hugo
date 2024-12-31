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
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
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

func TestQR(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/index.html --
{{ $text := "https://gohugo.io" }}
{{ $optionMaps := slice
	(dict "text" $text)
	(dict "text" $text "level" "low")
	(dict "text" $text "level" "medium")
	(dict "text" $text "level" "quartile")
	(dict "text" $text "level" "high")
	(dict "text" $text "targetDir" "foo")
	(dict "text" $text "level" "high" "targetDir" "foo/bar")
}}
{{ range $k, $opts := $optionMaps }}
	{{ with images.QR $opts }}
		<img data-id="{{ $k }}" data-hash="{{ .Content | hash.XxHash }}" data-level="{{ $opts.level }}" data-targetDir="{{ $opts.targetDir }}" src="{{ .RelPermalink }}">
	{{ end }}
{{ end }}
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html",
		`<img data-id="0" data-hash="5c8f15c6a5da74b1" data-level="" data-targetDir="" src="/qr_30596359e57a22d8.png">`,
		`<img data-id="1" data-hash="863b6fca7913f6ec" data-level="low" data-targetDir="" src="/qr_9fa47d3ceaf0993c.png">`,
		`<img data-id="2" data-hash="5c8f15c6a5da74b1" data-level="medium" data-targetDir="" src="/qr_30596359e57a22d8.png">`,
		`<img data-id="3" data-hash="2e6b70fbd4a4442d" data-level="quartile" data-targetDir="" src="/qr_0a5f8db4478f4066.png">`,
		`<img data-id="4" data-hash="b2d62c862af4e3f6" data-level="high" data-targetDir="" src="/qr_a6f66d8f08c8af75.png">`,
		`<img data-id="5" data-hash="5c8f15c6a5da74b1" data-level="" data-targetDir="foo" src="/foo/qr_30596359e57a22d8.png">`,
		`<img data-id="6" data-hash="b2d62c862af4e3f6" data-level="high" data-targetDir="foo/bar" src="/foo/bar/qr_a6f66d8f08c8af75.png">`,
	)

	files = strings.ReplaceAll(files, "low", "foo")

	b, err := hugolib.TestE(t, files)
	b.Assert(err.Error(), qt.Contains, "error correction level must be one of low, medium, quartile, or high")

	files = strings.ReplaceAll(files, "foo", "low")
	files = strings.ReplaceAll(files, "https://gohugo.io", "")

	b, err = hugolib.TestE(t, files)
	b.Assert(err.Error(), qt.Contains, "cannot encode an empty string")
}
