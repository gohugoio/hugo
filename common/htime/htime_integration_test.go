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

package htime_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// Issue #11267
func TestApplyWithContext(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
defaultContentLanguage = 'it'
-- layouts/index.html --
{{ $dates := slice
  "2022-01-03"
  "2022-02-01"
  "2022-03-02"
  "2022-04-07"
  "2022-05-06"
  "2022-06-04"
  "2022-07-03"
  "2022-08-01"
  "2022-09-06"
  "2022-10-05"
  "2022-11-03"
  "2022-12-02"
}}
{{ range $dates }}
	{{ . | time.Format "month: _January_ weekday: _Monday_" }}
	{{ . | time.Format "month: _Jan_ weekday: _Mon_" }}
{{ end }}
  `

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", `
month: _gennaio_ weekday: _lunedì_
month: _gen_ weekday: _lun_
month: _febbraio_ weekday: _martedì_
month: _feb_ weekday: _mar_
month: _marzo_ weekday: _mercoledì_
month: _mar_ weekday: _mer_
month: _aprile_ weekday: _giovedì_
month: _apr_ weekday: _gio_
month: _maggio_ weekday: _venerdì_
month: _mag_ weekday: _ven_
month: _giugno_ weekday: _sabato_
month: _giu_ weekday: _sab_
month: _luglio_ weekday: _domenica_
month: _lug_ weekday: _dom_
month: _agosto_ weekday: _lunedì_
month: _ago_ weekday: _lun_
month: _settembre_ weekday: _martedì_
month: _set_ weekday: _mar_
month: _ottobre_ weekday: _mercoledì_
month: _ott_ weekday: _mer_
month: _novembre_ weekday: _giovedì_
month: _nov_ weekday: _gio_
month: _dicembre_ weekday: _venerdì_
month: _dic_ weekday: _ven_
`)
}
