// Copyright 2017 The Hugo Authors. All rights reserved.
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

package time_test

import (
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/hugolib"
)

// See issue 15565.
func TestTimeFormatZoneName(t *testing.T) {
	// Don't run in parallel; time.Local is modified.
	// Asia/Seoul has a different offset than New York, the same offset
	// as Tokyo, and a non-zero offset.
	seoul, err := time.LoadLocation("Asia/Seoul")
	qt.Assert(t, err, qt.IsNil)
	local := time.Local
	t.Cleanup(func() { time.Local = local })
	time.Local = seoul

	files := `
-- hugo.toml --
timeZone = 'America/New_York'
disableKinds = ['rss','section','sitemap','taxonomy','term']
-- content/toml.md --
+++
title = 'toml'
date = 2023-01-28T02:44:58
+++
-- content/yaml.md --
---
title: yaml
date: 2023-01-28T02:44:58
---
-- content/json.md --
{"title": "json", "date": "2023-01-28T02:44:58"}
-- content/toml-offset.md --
+++
title = 'toml-offset'
date = 2023-01-28T16:44:58+09:00
+++
-- content/yaml-offset.md --
---
title: yaml-offset
date: 2023-01-28T16:44:58+09:00
---
-- content/json-offset.md --
{"title": "json-offset", "date": "2023-01-28T16:44:58+09:00"}
-- layouts/home.html --
{{ $t := time.AsTime "2023-01-28T07:44:58Z" }}
{{ $newyork := time.In "America/New_York" $t }}
{{ $tokyo := time.In "Asia/Tokyo" $t }}
{{ $london := time.In "Europe/London" $t }}
newyork: {{ time.Format "15:04 MST" $newyork }}|{{ time.Format ":time_full" $newyork }}|{{ (time.AsTime $newyork).Format "MST" }}|
tokyo: {{ time.Format "15:04 MST" $tokyo }}|{{ time.Format ":time_full" $tokyo }}|
london: {{ time.Format "15:04 MST" $london }}|{{ time.Format ":time_full" $london }}|
{{ range .RegularPages }}
{{ .Title }}: {{ time.Format "15:04 MST" .Date }}|
{{- end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		"newyork: 02:44 EST|2:44:58\u202fam Eastern Standard Time|EST|",
		"tokyo: 16:44 JST|4:44:58\u202fpm Japan Standard Time|",
		"london: 07:44 GMT|7:44:58\u202fam Greenwich Mean Time|",
		// Dates without an offset get the timeZone location, so they have a
		// zone name. The fix keeps it.
		"toml: 02:44 EST|",
		"yaml: 02:44 EST|",
		"json: 02:44 EST|",
		// Dates with an offset have no zone name. Go borrows the system time
		// zone's name, Asia/Seoul, because its offset matches. The fix doesn't
		// change this.
		"toml-offset: 16:44 KST|",
		"yaml-offset: 16:44 KST|",
		"json-offset: 16:44 KST|",
	)
}

// See issue 14948.
func TestTimeFormatMonthAbbreviationsEnGB(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
locale = 'en-GB'
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/home.html --
{{ $dates := slice
	"2025-01-17T10:42:00-08:00"
	"2025-02-17T10:42:00-08:00"
	"2025-03-17T10:42:00-07:00"
	"2025-04-17T10:42:00-07:00"
	"2025-05-17T10:42:00-07:00"
	"2025-06-17T10:42:00-07:00"
	"2025-07-17T10:42:00-07:00"
	"2025-08-17T10:42:00-07:00"
	"2025-09-17T10:42:00-07:00"
	"2025-10-17T10:42:00-07:00"
	"2025-11-17T10:42:00-08:00"
	"2025-12-17T10:42:00-08:00"
}}
{{- range $dates }}
	{{- time.Format "Jan" . }}|
{{- end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sept|Oct|Nov|Dec|")
}
