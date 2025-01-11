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
	"github.com/gohugoio/hugo/resources/images/imagetesting"
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
    (dict)
    (dict "level" "medium")
    (dict "level" "medium" "scale" 4)
    (dict "level" "low" "scale" 2)
    (dict "level" "medium" "scale" 3)
    (dict "level" "quartile" "scale" 5)
    (dict "level" "high" "scale" 6)
    (dict "level" "high" "scale" 6 "targetDir" "foo/bar")
}}
{{- range $k, $opts := $optionMaps }}
    {{- with images.QR $text $opts }}
<img data-id="{{ $k }}" data-img-hash="{{ .Content | hash.XxHash }}" data-level="{{ $opts.level }}" data-scale="{{ $opts.scale }}" data-targetDir="{{ $opts.targetDir }}" src="{{ .RelPermalink }}">
    {{- end }}
{{- end }}
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html",
		`<img data-id="0" data-img-hash="6ccacf8056c41475" data-level="" data-scale="" data-targetDir="" src="/qr_924bf7d80a564b23.png">`,
		`<img data-id="1" data-img-hash="6ccacf8056c41475" data-level="medium" data-scale="" data-targetDir="" src="/qr_924bf7d80a564b23.png">`,
		`<img data-id="2" data-img-hash="6ccacf8056c41475" data-level="medium" data-scale="4" data-targetDir="" src="/qr_924bf7d80a564b23.png">`,
		`<img data-id="3" data-img-hash="c29338c3d105b156" data-level="low" data-scale="2" data-targetDir="" src="/qr_9bf1ce25c5f2c058.png">`,
		`<img data-id="4" data-img-hash="8f7a639cea917b0e" data-level="medium" data-scale="3" data-targetDir="" src="/qr_7af14b329dd10af7.png">`,
		`<img data-id="5" data-img-hash="2d15d6dcb861b5da" data-level="quartile" data-scale="5" data-targetDir="" src="/qr_9600ecb2010c2185.png">`,
		`<img data-id="6" data-img-hash="113c45f2c091bc4d" data-level="high" data-scale="6" data-targetDir="" src="/qr_bdc74ee7f5c11cc6.png">`,
		`<img data-id="7" data-img-hash="113c45f2c091bc4d" data-level="high" data-scale="6" data-targetDir="foo/bar" src="/foo/bar/qr_14162f02f2b83fff.png">`,
	)

	files = strings.ReplaceAll(files, "low", "foo")

	b, err := hugolib.TestE(t, files)
	b.Assert(err.Error(), qt.Contains, "error correction level must be one of low, medium, quartile, or high")

	files = strings.ReplaceAll(files, "foo", "low")
	files = strings.ReplaceAll(files, "https://gohugo.io", "")

	b, err = hugolib.TestE(t, files)
	b.Assert(err.Error(), qt.Contains, "cannot encode an empty string")
}

func TestImagesGoldenFuncs(t *testing.T) {
	t.Parallel()

	if imagetesting.SkipGoldenTests {
		t.Skip("Skip golden test on this architecture")
	}

	// Will be used as the base folder for generated images.
	name := "funcs"

	files := `
-- hugo.toml --
-- assets/sunset.jpg --
sourcefilename: ../../resources/testdata/sunset.jpg

-- layouts/index.html --
Home.

{{ template "copy" (dict "name" "qr-default.png" "img" (images.QR "https://gohugo.io"))  }}
{{ template "copy" (dict "name" "qr-level-high_scale-6.png" "img" (images.QR "https://gohugo.io" (dict "level" "high" "scale" 6)))  }}

{{ define "copy"}}
{{ if lt (len (path.Ext .name)) 4 }}
	{{ errorf "No extension in %q" .name }}
{{ end }}
{{ $img := .img }}
{{ $name := printf "images/%s" .name  }}
{{ with $img | resources.Copy $name }}
{{ .Publish }}
{{ end }}
{{ end }}
`

	opts := imagetesting.DefaultGoldenOpts
	opts.T = t
	opts.Name = name
	opts.Files = files

	imagetesting.RunGolden(opts)
}
