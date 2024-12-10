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

package resource_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestResourcesMount(t *testing.T) {
	files := `
-- hugo.toml --
-- assets/text/txt1.txt --
Text 1.
-- assets/text/txt2.txt --
Text 2.
-- assets/text/sub/txt3.txt --
Text 3.
-- assets/text/sub/txt4.txt --
Text 4.
-- content/mybundle/index.md --
---
title: "My Bundle"
---
-- content/mybundle/txt1.txt --
Text 1.
-- content/mybundle/sub/txt2.txt --
Text 1.
-- layouts/index.html --
{{ $mybundle := site.GetPage "mybundle" }}
{{ $subResources := resources.Match "/text/sub/*.*"  }}
{{ $subResourcesMount :=  $subResources.Mount "/text/sub" "/newroot" }}
resources:text/txt1.txt:{{ with resources.Get "text/txt1.txt" }}{{ .Name }}{{ end }}|
resources:text/txt2.txt:{{ with resources.Get "text/txt2.txt" }}{{ .Name }}{{ end }}|
resources:text/sub/txt3.txt:{{ with resources.Get "text/sub/txt3.txt" }}{{ .Name }}{{ end }}|
subResources.range:{{ range $subResources }}{{ .Name }}|{{ end }}|
subResources:"text/sub/txt3.txt:{{ with $subResources.Get "text/sub/txt3.txt" }}{{ .Name }}{{ end }}|
subResourcesMount:/newroot/txt3.txt:{{ with $subResourcesMount.Get "/newroot/txt3.txt" }}{{ .Name }}{{ end }}|
page:txt1.txt:{{ with $mybundle.Resources.Get "txt1.txt" }}{{ .Name }}{{ end }}|
page:./txt1.txt:{{ with $mybundle.Resources.Get "./txt1.txt" }}{{ .Name }}{{ end }}|
page:sub/txt2.txt:{{ with $mybundle.Resources.Get "sub/txt2.txt" }}{{ .Name }}{{ end }}|
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
resources:text/txt1.txt:/text/txt1.txt|
resources:text/txt2.txt:/text/txt2.txt|
resources:text/sub/txt3.txt:/text/sub/txt3.txt|
subResources:"text/sub/txt3.txt:/text/sub/txt3.txt|
subResourcesMount:/newroot/txt3.txt:/text/sub/txt3.txt|
page:txt1.txt:txt1.txt|
page:./txt1.txt:txt1.txt|
page:sub/txt2.txt:sub/txt2.txt|
`)
}

func TestResourcesMountOnRename(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term", "home", "sitemap"]
-- content/mybundle/index.md --
---
title: "My Bundle"
resources:
- name: /foo/bars.txt
  src: foo/txt1.txt
- name: foo/bars2.txt
  src: foo/txt2.txt
---
-- content/mybundle/foo/txt1.txt --
Text 1.
-- content/mybundle/foo/txt2.txt --
Text 2.
-- layouts/_default/single.html --
Single.
{{ $mybundle := site.GetPage "mybundle" }}
Resources:{{ range $mybundle.Resources }}Name: {{ .Name }}|{{ end }}$
{{ $subResourcesMount :=  $mybundle.Resources.Mount "/foo" "/newroot" }}
 {{ $subResourcesMount2 :=  $mybundle.Resources.Mount "foo" "/newroot" }}
{{ $subResourcesMount3 :=  $mybundle.Resources.Mount "foo" "." }}
subResourcesMount:/newroot/bars.txt:{{ with $subResourcesMount.Get "/newroot/bars.txt" }}{{ .Name }}{{ end }}|
subResourcesMount:/newroot/bars2.txt:{{ with $subResourcesMount.Get "/newroot/bars2.txt" }}{{ .Name }}{{ end }}|
subResourcesMount2:/newroot/bars2.txt:{{ with $subResourcesMount2.Get "/newroot/bars2.txt" }}{{ .Name }}{{ end }}|
subResourcesMount3:bars2.txt:{{ with $subResourcesMount3.Get "bars2.txt" }}{{ .Name }}{{ end }}|
`
	b := hugolib.Test(t, files)
	b.AssertFileContent("public/mybundle/index.html",
		"Resources:Name: foo/bars.txt|Name: foo/bars2.txt|$",
		"subResourcesMount:/newroot/bars.txt:|\nsubResourcesMount:/newroot/bars2.txt:|",
		"subResourcesMount2:/newroot/bars2.txt:foo/bars2.txt|",
		"subResourcesMount3:bars2.txt:foo/bars2.txt|",
	)
}
