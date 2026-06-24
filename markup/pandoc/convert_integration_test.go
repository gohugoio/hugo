// Copyright 2026 The Hugo Authors. All rights reserved.
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

package pandoc_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/markup/pandoc"
)

// See issue 15062.
func TestBibliographySupport(t *testing.T) {
	if !pandoc.Supports() {
		t.Skip("pandoc not installed")
	}
	t.Parallel()

	files := `
-- hugo.toml --
disableKinds = ['home','rss','section','sitemap','taxonomy','term']
[security.exec]
allow = ['pandoc']
-- layouts/page.html --
{{ .Content }}
-- content/p1.pdc --
---
title: home
---

---
bibliography: testdata/foo.bib
citation-style: testdata/ieee.csl
link-citations: true
references:
- type: article-journal
  id: WatsonCrick1953
  author:
  - family: Watson
    given: J. D.
  - family: Crick
    given: F. H. C.
  issued:
    date-parts:
    - - 1953
      - 4
      - 25
  title: 'Molecular structure of nucleic acids: a structure for deoxyribose nucleic acid'
  title-short: Molecular structure of nucleic acids
  container-title: Nature
  volume: 171
  issue: 4356
  page: 737-738
  DOI: 10.1038/171737a0
---

This is a citation: @einstein1905physics

This is another citation: [@WatsonCrick1953, p. 33]
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/p1/index.html",
		`<p>This is a citation: <span class="citation"
data-cites="einstein1905physics"><a href="#ref-einstein1905physics"
role="doc-biblioref">[1]</a></span></p>
<p>This is another citation: <span class="citation"
data-cites="WatsonCrick1953"><a href="#ref-WatsonCrick1953"
role="doc-biblioref">[2, p. 33]</a></span></p>
<div id="ref-einstein1905physics" class="csl-entry" role="listitem">
<div class="csl-left-margin">[1] </div><div class="csl-right-inline">A.
Einstein, <span>“Zur elektrodynamik bewegter
k<span>ö</span>rper,”</span> <em>Annalen der Physik</em>, vol. 322, no.
10, pp. 891–921, 1905.</div>
</div>
<div id="ref-WatsonCrick1953" class="csl-entry" role="listitem">
<div class="csl-left-margin">[2] </div><div class="csl-right-inline">J.
D. Watson and F. H. C. Crick, <span>“Molecular structure of nucleic
acids: a structure for deoxyribose nucleic acid,”</span>
<em>Nature</em>, vol. 171, no. 4356, pp. 737–738, Apr. 1953, doi: <a
href="https://doi.org/10.1038/171737a0">10.1038/171737a0</a>.</div>
</div>
</div>`,
	)
}
