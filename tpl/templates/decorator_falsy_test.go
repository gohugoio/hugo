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

package templates_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// TestDecoratorPartialFalsyReturn tests that partial decorators work correctly
// when a partial returns a falsy value (false, nil, ""). This was a bug in
// v0.154.0-v0.154.5 where the decorator stack would become unbalanced.
// See https://github.com/gohugoio/hugo/issues/14419
func TestDecoratorPartialFalsyReturn(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ["section", "taxonomy", "term", "sitemap", "RSS"]
-- content/p1.md --
---
title: "Page 1"
---
-- layouts/_partials/a.html --
{{ $result := dict }}
{{ with partialCached "b.html" . .RelPermalink }}
  {{ $result = . }}
{{ end }}
{{ return $result }}
-- layouts/_partials/b.html --
{{ $result := dict }}
{{ with partialCached "c.html" . "key1" }}
  {{ $result = merge $result (dict "c" .) }}
{{ end }}
{{ with partialCached "d.html" . "key2" }}
  {{ $result = merge $result (dict "d" .) }}
{{ end }}
{{ return $result }}
-- layouts/_partials/c.html --
{{ return false }}
-- layouts/_partials/d.html --
{{ return "truthy" }}
-- layouts/home.html --
{{ range site.RegularPages }}
{{ with partialCached "a.html" (dict "Page" . "RelPermalink" .RelPermalink) .RelPermalink }}d:{{ .d }}{{ end }}
{{ end }}$
`
	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "d:truthy", "$")
}
