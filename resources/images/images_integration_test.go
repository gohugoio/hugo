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

func TestAutoOrient(t *testing.T) {
	files := `
-- hugo.toml --
-- assets/rotate270.jpg --
sourcefilename: ../testdata/exif/orientation6.jpg
-- layouts/index.html --
{{ $img := resources.Get "rotate270.jpg" }}
W/H original: {{ $img.Width }}/{{ $img.Height }}
{{ $rotated := $img.Filter images.AutoOrient }}
W/H rotated: {{ $rotated.Width }}/{{ $rotated.Height }}
`

	b := hugolib.Test(t, files)
	b.AssertFileContent("public/index.html", "W/H original: 80/40\n\nW/H rotated: 40/80")
}
