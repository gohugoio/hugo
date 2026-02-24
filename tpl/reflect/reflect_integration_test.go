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

package reflect_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestIs(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
-- assets/a.png --
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==
-- assets/b.svg --
<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
  <circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red" />
</svg>
-- assets/c.txt --
This is a text file.
-- assets/d.avif --
AAAAHGZ0eXBhdmlmAAAAAGF2aWZtaWYxbWlhZgAAAOptZXRhAAAAAAAAACFoZGxyAAAAAAAAAABwaWN0AAAAAAAAAAAAAAAAAAAAAA5waXRtAAAAAAABAAAAImlsb2MAAAAAREAAAQABAAAAAAEOAAEAAAAAAAAAEgAAACNpaW5mAAAAAAABAAAAFWluZmUCAAAAAAEAAGF2MDEAAAAAamlwcnAAAABLaXBjbwAAABNjb2xybmNseAABAA0ABoAAAAAMYXYxQ4EgAgAAAAAUaXNwZQAAAAAAAAABAAAAAQAAABBwaXhpAAAAAAMICAgAAAAXaXBtYQAAAAAAAAABAAEEAYIDBAAAABptZGF0EgAKBzgABhAQ0GkyBRAAAAtA
-- layouts/home.html --
{{ $a := resources.Get "a.png" }}
{{ $a10 := $a.Fit "10x10" }}
{{ $b := resources.Get "b.svg" }}
{{ $c := resources.Get "c.txt" }}
{{ $d := resources.Get "d.avif" }}
PNG.ResourceType: {{ $a.ResourceType }}
SVG.ResourceType: {{ $b.ResourceType }}
Text.ResourceType: {{ $c.ResourceType }}
AVIF.ResourceType: {{ $d.ResourceType }}
IsSite: false: {{ reflect.IsSite . }}|true: {{ reflect.IsSite .Site }}|true: {{ reflect.IsSite site }}
IsPage: true: {{ reflect.IsPage . }}|false: {{ reflect.IsPage .Site }}|false: {{ reflect.IsPage site }}
IsResource: true: {{ reflect.IsResource . }}|true: {{ reflect.IsResource $a }}|true: {{ reflect.IsResource $b }}|true: {{ reflect.IsResource $c }}



`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
PNG.ResourceType: image
SVG.ResourceType: image
Text.ResourceType: text
AVIF.ResourceType: image
IsSite: false: false|true: true|true: true
IsPage: true: true|false: false|false: false
IsResource: true: true|true: true|true: true|true: true
`)
}
