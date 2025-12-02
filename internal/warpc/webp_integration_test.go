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

package warpc_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestWebPMisc(t *testing.T) {
	files := `
-- assets/sunrise.webp --
sourcefilename: ../../resources/testdata/sunrise.webp
-- layouts/home.html --
{{ $image := resources.Get "sunrise.webp" }}
{{ $resized := $image.Resize "123x456" }}
Original Width/Height: {{ $image.Width }}/{{ $image.Height }}|
Resized Width:/Height {{ $resized.Width }}/{{ $resized.Height }}|
Resized RelPermalink: {{ $resized.RelPermalink }}|
`

	for range 1 {
		b := hugolib.Test(t, files)

		b.AssertPublishDir("! 000000 sunrise_hu_94448e81dce95acf.webp")

		b.AssertFileContent("public/index.html",
			"Original Width/Height: 1024/640|",
			"Resized Width:/Height 123/456|",
			"Resized RelPermalink: /sunrise_hu_6095509b5348ba46.webp|",
		)
	}
}
