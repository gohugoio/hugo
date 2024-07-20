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

package i18n_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestI18nFromTheme(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
[module]
[[module.imports]]
path = "mytheme"
-- i18n/en.toml --
[l1]
other = 'l1main'
[l2]
other = 'l2main'
-- themes/mytheme/i18n/en.toml --
[l1]
other = 'l1theme'
[l2]
other = 'l2theme'
[l3]
other = 'l3theme'
-- layouts/index.html --
l1: {{ i18n "l1"  }}|l2: {{ i18n "l2"  }}|l3: {{ i18n "l3"  }}

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
l1: l1main|l2: l2main|l3: l3theme
	`)
}

func TestPassPageToI18n(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
-- content/_index.md --
---
title: "Home"
---
Duis quis irure id nisi sunt minim aliqua occaecat. Aliqua cillum labore consectetur quis culpa tempor quis non officia cupidatat in ad cillum. Velit irure pariatur nisi adipisicing officia reprehenderit commodo esse non.

Ullamco cupidatat nostrud ut reprehenderit. Consequat nisi culpa magna amet tempor velit reprehenderit. Ad minim eiusmod tempor nostrud eu aliquip consectetur commodo ut in aliqua enim. Cupidatat voluptate laborum consequat qui nulla laborum laborum aute ea culpa nulla dolor cillum veniam. Commodo esse tempor qui labore aute aliqua sint nulla do.

Ad deserunt esse nostrud labore. Amet reprehenderit fugiat nostrud eu reprehenderit sit reprehenderit minim deserunt esse id occaecat cillum. Ad qui Lorem cillum laboris ipsum anim in culpa ad dolor consectetur minim culpa.

Lorem cupidatat officia aute in eu commodo anim nulla deserunt occaecat reprehenderit dolore. Eu cupidatat reprehenderit ipsum sit laboris proident. Duis quis nulla tempor adipisicing. Adipisicing amet ad reprehenderit non mollit. Cupidatat proident tempor laborum sit ipsum adipisicing sunt magna labore. Eu irure nostrud cillum exercitation tempor proident. Laborum magna nisi consequat do sint occaecat magna incididunt.

Sit mollit amet esse dolore in labore aliquip eu duis officia incididunt. Esse veniam labore excepteur eiusmod occaecat ullamco magna sunt. Ipsum occaecat exercitation anim fugiat in amet excepteur excepteur aliquip laborum. Aliquip aliqua consequat officia sit sint amet aliqua ipsum eu veniam. Id enim quis ea in eu consequat exercitation occaecat veniam consequat anim nulla adipisicing minim. Ut duis cillum laboris duis non commodo eu aliquip tempor nisi aute do.

Ipsum nulla esse excepteur ut aliqua esse incididunt deserunt veniam dolore est laborum nisi veniam. Magna eiusmod Lorem do tempor incididunt ut aute aliquip ipsum ea laboris culpa. Occaecat do officia velit fugiat culpa eu minim magna sint occaecat sunt. Duis magna proident incididunt est cupidatat proident esse proident ut ipsum non dolor Lorem eiusmod. Officia quis irure id eu aliquip.

Duis anim elit in officia in in aliquip est. Aliquip nisi labore qui elit elit cupidatat ut labore incididunt eiusmod ipsum. Sit irure nulla non cupidatat exercitation sit culpa nisi ex dolore. Culpa nisi duis duis eiusmod commodo nulla.

Et magna aliqua amet qui mollit. Eiusmod aute ut anim ea est fugiat non nisi in laborum ullamco. Proident mollit sunt nostrud irure esse sunt eiusmod deserunt dolor. Irure aute ad magna est consequat duis cupidatat consequat. Enim tempor aute cillum quis ea do enim proident incididunt aliquip cillum tempor minim. Nulla minim tempor proident in excepteur consectetur veniam.

Exercitation tempor nulla incididunt deserunt laboris ad incididunt aliqua exercitation. Adipisicing laboris veniam aute eiusmod qui magna fugiat velit. Aute quis officia anim commodo id fugiat nostrud est. Quis ipsum amet velit adipisicing eu anim minim eu est in culpa aute. Esse in commodo irure enim proident reprehenderit ullamco in dolore aute cillum.

Irure excepteur ex occaecat ipsum laboris fugiat exercitation. Exercitation adipisicing velit excepteur eu culpa consequat exercitation dolore. In laboris aute quis qui mollit minim culpa. Magna velit ea aliquip veniam fugiat mollit veniam.
-- i18n/en.toml --
[a]
other = 'Reading time: {{ .ReadingTime }}'
-- layouts/index.html --
i18n: {{ i18n "a" . }}|

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
	i18n: Reading time: 3|
	`)
}

// Issue 9216
func TestI18nDefaultContentLanguage(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
disableKinds = ['RSS','sitemap','taxonomy','term','page','section']
defaultContentLanguage = 'es'
defaultContentLanguageInSubdir = true
[languages.es]
[languages.fr]
-- i18n/es.toml --
cat = 'gato'
-- i18n/fr.toml --
# this file intentionally empty
-- layouts/index.html --
{{ .Title }}_{{ T "cat" }}
-- content/_index.fr.md --
---
title: home_fr
---
-- content/_index.md --
---
title: home_es
---
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/es/index.html", `home_es_gato`)
	b.AssertFileContent("public/fr/index.html", `home_fr_gato`)
}
