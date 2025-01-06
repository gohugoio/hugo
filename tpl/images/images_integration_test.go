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
{{- $text := "https://gohugo.io" }}
{{- $optionMaps := slice
    (dict "text" $text)
    (dict "text" $text "level" "medium")
    (dict "text" $text "level" "medium" "scale" 4)
    (dict "text" $text "level" "low" "scale" 2)
    (dict "text" $text "level" "medium" "scale" 3)
    (dict "text" $text "level" "quartile" "scale" 5)
    (dict "text" $text "level" "high" "scale" 6)
    (dict "text" $text "level" "high" "scale" 6 "targetDir" "foo/bar")
}}
{{- range $k, $opts := $optionMaps }}
    {{- with images.QR $opts }}
<img data-id="{{ $k }}" data-img-hash="{{ .Content | hash.XxHash }}" data-level="{{ $opts.level }}" data-scale="{{ $opts.scale }}" data-targetDir="{{ $opts.targetDir }}" src="{{ .RelPermalink }}">
    {{- end }}
{{- end }}
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html",
		`<img data-id="0" data-img-hash="6ccacf8056c41475" data-level="" data-scale="" data-targetDir="" src="/qr_3601c357f288f47f.png">`,
		`<img data-id="1" data-img-hash="6ccacf8056c41475" data-level="medium" data-scale="" data-targetDir="" src="/qr_3601c357f288f47f.png">`,
		`<img data-id="2" data-img-hash="6ccacf8056c41475" data-level="medium" data-scale="4" data-targetDir="" src="/qr_3601c357f288f47f.png">`,
		`<img data-id="3" data-img-hash="c29338c3d105b156" data-level="low" data-scale="2" data-targetDir="" src="/qr_232594637b3d9ac1.png">`,
		`<img data-id="4" data-img-hash="8f7a639cea917b0e" data-level="medium" data-scale="3" data-targetDir="" src="/qr_5c02e7507f8e86e0.png">`,
		`<img data-id="5" data-img-hash="2d15d6dcb861b5da" data-level="quartile" data-scale="5" data-targetDir="" src="/qr_c49dd961bcc47c06.png">`,
		`<img data-id="6" data-img-hash="113c45f2c091bc4d" data-level="high" data-scale="6" data-targetDir="" src="/qr_17994d3244e3c686.png">`,
		`<img data-id="7" data-img-hash="113c45f2c091bc4d" data-level="high" data-scale="6" data-targetDir="foo/bar" src="/foo/bar/qr_abd2f7b221eee6ea.png">`,
	)

	files = strings.ReplaceAll(files, "low", "foo")

	b, err := hugolib.TestE(t, files)
	b.Assert(err.Error(), qt.Contains, "error correction level must be one of low, medium, quartile, or high")

	files = strings.ReplaceAll(files, "foo", "low")
	files = strings.ReplaceAll(files, "https://gohugo.io", "")

	b, err = hugolib.TestE(t, files)
	b.Assert(err.Error(), qt.Contains, "cannot encode an empty string")
}
