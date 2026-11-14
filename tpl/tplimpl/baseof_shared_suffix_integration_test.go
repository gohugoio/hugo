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

package tplimpl_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// Issue #15590: a baseof template without an explicit output format identifier
// (e.g. layouts/baseof.xml) had its media type resolved to whichever output
// format claims the shared suffix first. When a custom XML format (e.g. Atom)
// and the built-in RSS format both use the "xml" suffix, the RSS output lost
// its base template and rendered the overlay standalone -- which, wrapped in
// {{ define "main" }}, silently produced an empty file.
func TestBaseofXMLSharedSuffix(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
[outputs]
home = ["RSS", "ATOM"]

[mediaTypes."application/atom"]
suffixes = ["xml"]

[outputFormats.ATOM]
mediaType = "application/atom"
baseName = "atom"
-- content/p1.md --
---
title: hello
---
body
-- layouts/baseof.xml --
BASEOF-XML-START{{ block "main" . }}{{ end }}BASEOF-XML-END
-- layouts/list.rss.xml --
{{ define "main" }}RSS-LIST{{ end }}
-- layouts/list.atom.xml --
{{ define "main" }}ATOM-LIST{{ end }}
`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	// Both output formats share the baseof.xml base template.
	b.AssertFileContent("public/index.xml", "BASEOF-XML-START", "RSS-LIST", "BASEOF-XML-END")
	b.AssertFileContent("public/atom.xml", "BASEOF-XML-START", "ATOM-LIST", "BASEOF-XML-END")
}

// An explicit output format identifier (baseof.rss.xml) wins over the shared
// baseof.xml for that output format, while other formats fall back to
// baseof.xml.
func TestBaseofWithExplicitOutputFormat(t *testing.T) {
	files := `
-- hugo.toml --
disableKinds = ["taxonomy", "term"]
[outputs]
home = ["RSS", "ATOM"]

[mediaTypes."application/atom"]
suffixes = ["xml"]

[outputFormats.ATOM]
mediaType = "application/atom"
baseName = "atom"
-- content/p1.md --
---
title: hello
---
body
-- layouts/baseof.rss.xml --
BASEOF-RSS-START{{ block "main" . }}{{ end }}BASEOF-RSS-END
-- layouts/baseof.xml --
BASEOF-XML-START{{ block "main" . }}{{ end }}BASEOF-XML-END
-- layouts/list.rss.xml --
{{ define "main" }}RSS-LIST{{ end }}
-- layouts/list.atom.xml --
{{ define "main" }}ATOM-LIST{{ end }}
`
	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	// RSS has an explicit baseof.rss.xml; it wins over the shared baseof.xml.
	b.AssertFileContent("public/index.xml", "BASEOF-RSS-START", "RSS-LIST", "BASEOF-RSS-END")
	b.AssertFileContent("public/atom.xml", "BASEOF-XML-START", "ATOM-LIST", "BASEOF-XML-END")
}
