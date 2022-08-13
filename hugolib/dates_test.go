// Copyright 2021 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"fmt"
	"time"

	qt "github.com/frankban/quicktest"

	"strings"
	"testing"
)

func TestDateFormatMultilingual(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", `
baseURL = "https://example.org"

defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true

[languages]
[languages.en]
weight=10
[languages.nn]
weight=20
	
`)

	pageWithDate := `---
title: Page
date: 2021-07-18
---	
`

	b.WithContent(
		"_index.en.md", pageWithDate,
		"_index.nn.md", pageWithDate,
	)

	b.WithTemplatesAdded("index.html", `
Date: {{ .Date | time.Format ":date_long" }}
	`)

	b.Build(BuildCfg{})

	b.AssertFileContent("public/en/index.html", `Date: July 18, 2021`)
	b.AssertFileContent("public/nn/index.html", `Date: 18. juli 2021`)

}

func TestTimeZones(t *testing.T) {
	b := newTestSitesBuilder(t)
	b.WithConfigFile("toml", `
baseURL = "https://example.org"

defaultContentLanguage = "en"
defaultContentLanguageInSubDir = true

[languages]
[languages.en]
timeZone="UTC"
weight=10
[languages.nn]
timeZone="America/Antigua"
weight=20
	
`)

	const (
		pageTemplYaml = `---
title: Page
date: %s
lastMod: %s
publishDate: %s
expiryDate: %s
---	
`

		pageTemplTOML = `+++
title="Page"
date=%s
lastMod=%s
publishDate=%s
expiryDate=%s
+++
`

		shortDateTempl = `%d-07-%d`
		longDateTempl  = `%d-07-%d 15:28:01`
	)

	createPageContent := func(pageTempl, dateTempl string, quoted bool) string {
		createDate := func(year, i int) string {
			d := fmt.Sprintf(dateTempl, year, i)
			if quoted {
				return fmt.Sprintf("%q", d)
			}

			return d
		}

		return fmt.Sprintf(
			pageTempl,
			createDate(2021, 10),
			createDate(2021, 11),
			createDate(2021, 12),
			createDate(2099, 13), // This test will fail in 2099 :-)
		)
	}

	b.WithContent(
		// YAML
		"short-date-yaml-unqouted.en.md", createPageContent(pageTemplYaml, shortDateTempl, false),
		"short-date-yaml-unqouted.nn.md", createPageContent(pageTemplYaml, shortDateTempl, false),
		"short-date-yaml-qouted.en.md", createPageContent(pageTemplYaml, shortDateTempl, true),
		"short-date-yaml-qouted.nn.md", createPageContent(pageTemplYaml, shortDateTempl, true),
		"long-date-yaml-unqouted.en.md", createPageContent(pageTemplYaml, longDateTempl, false),
		"long-date-yaml-unqouted.nn.md", createPageContent(pageTemplYaml, longDateTempl, false),
		// TOML
		"short-date-toml-unqouted.en.md", createPageContent(pageTemplTOML, shortDateTempl, false),
		"short-date-toml-unqouted.nn.md", createPageContent(pageTemplTOML, shortDateTempl, false),
		"short-date-toml-qouted.en.md", createPageContent(pageTemplTOML, shortDateTempl, true),
		"short-date-toml-qouted.nn.md", createPageContent(pageTemplTOML, shortDateTempl, true),
	)

	const datesTempl = `
Date: {{ .Date | safeHTML  }}
Lastmod: {{ .Lastmod | safeHTML  }}
PublishDate: {{ .PublishDate | safeHTML  }}
ExpiryDate: {{ .ExpiryDate | safeHTML  }}

	`

	b.WithTemplatesAdded(
		"_default/single.html", datesTempl,
	)

	b.Build(BuildCfg{})

	expectShortDateEn := `
Date: 2021-07-10 00:00:00 +0000 UTC
Lastmod: 2021-07-11 00:00:00 +0000 UTC
PublishDate: 2021-07-12 00:00:00 +0000 UTC
ExpiryDate: 2099-07-13 00:00:00 +0000 UTC`

	expectShortDateNn := strings.ReplaceAll(expectShortDateEn, "+0000 UTC", "-0400 AST")

	expectLongDateEn := `
Date: 2021-07-10 15:28:01 +0000 UTC
Lastmod: 2021-07-11 15:28:01 +0000 UTC
PublishDate: 2021-07-12 15:28:01 +0000 UTC
ExpiryDate: 2099-07-13 15:28:01 +0000 UTC`

	expectLongDateNn := strings.ReplaceAll(expectLongDateEn, "+0000 UTC", "-0400 AST")

	// TODO(bep) create a common proposal for go-yaml, go-toml
	// for a custom date parser hook to handle these time zones.
	// JSON is omitted from this test as JSON does no (to my knowledge)
	// have date literals.

	// YAML
	// Note: This is with go-yaml v2, I suspect v3 will fail with the unquouted values.
	b.AssertFileContent("public/en/short-date-yaml-unqouted/index.html", expectShortDateEn)
	b.AssertFileContent("public/nn/short-date-yaml-unqouted/index.html", expectShortDateNn)
	b.AssertFileContent("public/en/short-date-yaml-qouted/index.html", expectShortDateEn)
	b.AssertFileContent("public/nn/short-date-yaml-qouted/index.html", expectShortDateNn)

	b.AssertFileContent("public/en/long-date-yaml-unqouted/index.html", expectLongDateEn)
	b.AssertFileContent("public/nn/long-date-yaml-unqouted/index.html", expectLongDateNn)

	// TOML
	// These fails: TOML (Burnt Sushi) defaults to local timezone.
	// TODO(bep) check go-toml
	b.AssertFileContent("public/en/short-date-toml-unqouted/index.html", expectShortDateEn)
	b.AssertFileContent("public/nn/short-date-toml-unqouted/index.html", expectShortDateNn)
	b.AssertFileContent("public/en/short-date-toml-qouted/index.html", expectShortDateEn)
	b.AssertFileContent("public/nn/short-date-toml-qouted/index.html", expectShortDateNn)

}

// Issue 8832
func TestTimeZoneInvalid(t *testing.T) {
	b := newTestSitesBuilder(t)

	b.WithConfigFile("toml", `
	
timeZone = "America/LosAngeles"   # Should be America/Los_Angeles
`)

	err := b.CreateSitesE()
	b.Assert(err, qt.Not(qt.IsNil))
	b.Assert(err.Error(), qt.Contains, `failed to load config: invalid timeZone for language "en": unknown time zone America/LosAngeles`)
}

// Issue 8835
func TestTimeOnError(t *testing.T) {
	b := newTestSitesBuilder(t)

	b.WithTemplates("index.html", `time: {{ time "2020-10-20" "invalid-timezone" }}`)
	b.WithContent("p1.md", "")

	b.Assert(b.BuildE(BuildCfg{}), qt.Not(qt.IsNil))

}

func TestTOMLDates(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
timeZone = "America/Los_Angeles"
-- content/_index.md --
---
date: "2020-10-20"
---
-- content/p1.md --
+++
title = "TOML Date with UTC offset"
date = 2021-08-16T06:00:00+00:00
+++


## Foo
-- data/mydata.toml --
date = 2020-10-20
talks = [
	{ date = 2017-01-23, name = "Past talk 1" },
	{ date = 2017-01-24, name = "Past talk 2" },
	{ date = 2017-01-26, name = "Past talk 3" },
	{ date = 2050-02-12, name = "Future talk 1" },
	{ date = 2050-02-13, name = "Future talk 2" },
]
-- layouts/index.html --
{{ $futureTalks := where site.Data.mydata.talks "date" ">" now }}
{{ $pastTalks := where site.Data.mydata.talks "date" "<" now }}

{{ $homeDate := site.Home.Date }}
{{ $p1Date := (site.GetPage "p1").Date }}
Future talks: {{ len $futureTalks }}
Past talks: {{ len $pastTalks }}

Home's Date should be greater than past: {{ gt $homeDate (index $pastTalks 0).date }}
Home's Date should be less than future: {{ lt $homeDate (index $futureTalks 0).date }}
Home's Date should be equal mydata date: {{ eq $homeDate site.Data.mydata.date }}
Full time: {{ $p1Date | time.Format ":time_full" }}
`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/index.html", `
Future talks: 2
Past talks: 3
Home's Date should be greater than past: true
Home's Date should be less than future: true
Home's Date should be equal mydata date: true
Full time: 6:00:00 am UTC
`)
}

func TestTOMLTimeZone(t *testing.T) {
	type testCase struct {
		input          string
		configTimeZone string
		localZoneName  string
		expectedFull   string // :time_full
		expectedLong   string // :time_long
		comment        string
	}

	// America/Los_Angeles: -07:00
	// Asia/Tokyo +09:00
	testCases := []testCase{
		// default
		{
			input:          "2021-08-16T06:00:00+00:00",
			configTimeZone: "",
			localZoneName:  "",
			expectedFull:   "6:00:00 am UTC",
			expectedLong:   "6:00:00 am UTC",
			comment:        "zoneinfo `UTC` is displayed if offset is +00:00 as default",
		},
		{
			input:          "2021-08-16T06:00:00-07:00",
			configTimeZone: "",
			localZoneName:  "",
			expectedFull:   "6:00:00 am ",
			expectedLong:   "6:00:00 am ",
			comment:        "config not set; no zoneinfo",
		},
		{
			input:          "2021-08-16T06:00:00+09:00",
			configTimeZone: "",
			localZoneName:  "",
			expectedFull:   "6:00:00 am ",
			expectedLong:   "6:00:00 am ",
			comment:        "config not set; no zoneinfo",
		},

		// environment variable TZ has no effect
		{
			input:          "2021-08-16T06:00:00+00:00",
			configTimeZone: "",
			localZoneName:  "America/Los_Angeles",
			expectedFull:   "6:00:00 am UTC",
			expectedLong:   "6:00:00 am UTC",
			comment:        "zoneinfo `UTC` is displayed if offset is +00:00 as default",
		},
		{
			input:          "2021-08-16T06:00:00-07:00",
			configTimeZone: "",
			localZoneName:  "America/Los_Angeles",
			expectedFull:   "6:00:00 am ",
			expectedLong:   "6:00:00 am ",
			comment:        "environ set but config not set; no zoneinfo",
		},
		{
			input:          "2021-08-16T06:00:00+09:00",
			configTimeZone: "",
			localZoneName:  "America/Los_Angeles",
			expectedFull:   "6:00:00 am ",
			expectedLong:   "6:00:00 am ",
			comment:        "environ set but config not set; no zoneinfo",
		},

		{
			input:          "2021-08-16T06:00:00+00:00",
			configTimeZone: `timeZone = ""`,
			localZoneName:  "America/Los_Angeles",
			expectedFull:   "6:00:00 am UTC",
			expectedLong:   "6:00:00 am UTC",
			comment:        "zoneinfo `UTC` is displayed if offset is +00:00 as default",
		},
		{
			input:          "2021-08-16T06:00:00-07:00",
			configTimeZone: `timeZone = ""`,
			localZoneName:  "America/Los_Angeles",
			expectedFull:   "6:00:00 am ",
			expectedLong:   "6:00:00 am ",
			comment:        "environ set but timezone in config set to empty string; no zoneinfo",
		},
		{
			input:          "2021-08-16T06:00:00+09:00",
			configTimeZone: `timeZone = ""`,
			localZoneName:  "America/Los_Angeles",
			expectedFull:   "6:00:00 am ",
			expectedLong:   "6:00:00 am ",
			comment:        "environ set but timezone in config set to empty string; no zoneinfo",
		},

		{
			input:          "2021-08-16T06:00:00+00:00",
			configTimeZone: `timeZone = "America/Los_Angeles"`,
			localZoneName:  "",
			expectedFull:   "6:00:00 am UTC",
			expectedLong:   "6:00:00 am UTC",
			comment:        `config set to "America/Los_Angeles", offset unmatched; show "UTC" as default`,
		},
		{
			input:          "2021-08-16T06:00:00-07:00",
			configTimeZone: `timeZone = "America/Los_Angeles"`,
			localZoneName:  "",
			expectedFull:   "6:00:00 am Pacific Daylight Time",
			expectedLong:   "6:00:00 am PDT",
			comment:        `config set to "America/Los_Angeles", offset matched; show zoneinfo`,
		},
		{
			input:          "2021-08-16T06:00:00+09:00",
			configTimeZone: `timeZone = "America/Los_Angeles"`,
			localZoneName:  "",
			expectedFull:   "6:00:00 am ",
			expectedLong:   "6:00:00 am ",
			comment:        `config set to "America/Los_Angeles", offset unmatched; no zoneinfo`,
		},
		{
			input:          "2021-08-16T06:00:00+02:00",
			configTimeZone: `timeZone = "Europe/Oslo"`,
			localZoneName:  "",
			expectedFull:   "6:00:00 am CEST",
			expectedLong:   "6:00:00 am CEST",
			comment:        `config set to "Europe/Oslo", offset matched (summer time); show zoneinfo`,
		},
		{
			input:          "2021-12-16T06:00:00+02:00",
			configTimeZone: `timeZone = "Europe/Oslo"`,
			localZoneName:  "",
			expectedFull:   "6:00:00 am ",
			expectedLong:   "6:00:00 am ",
			comment:        `config set to "Europe/Oslo", offset unmatched (standard time); no zoneinfo`,
		},
		{
			input:          "2021-12-16T06:00:00+01:00",
			configTimeZone: `timeZone = "Europe/Oslo"`,
			localZoneName:  "",
			expectedFull:   "6:00:00 am CET",
			expectedLong:   "6:00:00 am CET",
			comment:        `config set to "Europe/Oslo", offset matched (standard time); show zoneinfo`,
		},
		{
			input:          "2021-08-16T06:00:00-07:00",
			configTimeZone: `timeZone = "Mexico/BajaNorte"`,
			localZoneName:  "",
			expectedFull:   "6:00:00 am Pacific Daylight Time",
			expectedLong:   "6:00:00 am PDT",
			comment:        `config set to "Mexico/BajaNorte", offset matched; show zoneinfo`,
		},
		{
			input:          "2021-08-16T06:00:00+00:00",
			configTimeZone: `timeZone = "Etc/Greenwich"`,
			localZoneName:  "",
			expectedFull:   "6:00:00 am Greenwich Mean Time",
			expectedLong:   "6:00:00 am GMT",
			comment:        `zoneinfo "UTC" could be overwritten if time zone offset is +00:00:00 and matched`,
		},

		// no effect of "TZ"
		{
			input:          "2021-08-16T06:00:00+00:00",
			configTimeZone: `timeZone = "America/Los_Angeles"`,
			localZoneName:  "Asia/Tokyo",
			expectedFull:   "6:00:00 am UTC",
			expectedLong:   "6:00:00 am UTC",
			comment:        `config set to "America/Los_Angeles", offset unmatched; show "UTC" as default`,
		},
		{
			input:          "2021-08-16T06:00:00-07:00",
			configTimeZone: `timeZone = "America/Los_Angeles"`,
			localZoneName:  "Asia/Tokyo",
			expectedFull:   "6:00:00 am Pacific Daylight Time",
			expectedLong:   "6:00:00 am PDT",
			comment:        `config set to "America/Los_Angeles", offset matched; show zoneinfo`,
		},
		{
			input:          "2021-08-16T06:00:00+09:00",
			configTimeZone: `timeZone = "America/Los_Angeles"`,
			localZoneName:  "Asia/Tokyo",
			expectedFull:   "6:00:00 am ",
			expectedLong:   "6:00:00 am ",
			comment:        `config set to "America/Los_Angeles", offset unmatched; no zoneinfo`,
		},

		// toml string of time is same as above
		{
			input:          `"2021-08-16T06:00:00+00:00"`,
			configTimeZone: `timeZone = "America/Los_Angeles"`,
			localZoneName:  "Asia/Tokyo",
			expectedFull:   "6:00:00 am ",
			expectedLong:   "6:00:00 am ",
			comment:        `config set to "America/Los_Angeles", offset unmatched; show "UTC"`,
		},
		{
			input:          `"2021-08-16T06:00:00-07:00"`,
			configTimeZone: `timeZone = "America/Los_Angeles"`,
			localZoneName:  "Asia/Tokyo",
			expectedFull:   "6:00:00 am Pacific Daylight Time",
			expectedLong:   "6:00:00 am PDT",
			comment:        `config set to "America/Los_Angeles", offset matched, show zoneinfo`,
		},
		{
			input:          `"2021-08-16T06:00:00+09:00"`,
			configTimeZone: `timeZone = "America/Los_Angeles"`,
			localZoneName:  "Asia/Tokyo",
			expectedFull:   "6:00:00 am ",
			expectedLong:   "6:00:00 am ",
			comment:        `config set to "America/Los_Angeles", offset unmatched, no zoneinfo`,
		},

		// without offset
		{
			input:          `"2021-08-16T06:00:00"`,
			configTimeZone: `timeZone = "America/Los_Angeles"`,
			localZoneName:  "Asia/Tokyo",
			expectedFull:   "6:00:00 am Pacific Daylight Time",
			expectedLong:   "6:00:00 am PDT",
			comment:        "NOTE: This does not show `UTC` because `cast.ToTimeInDefaultLocationE` uses `timeZone` value",
		},
		{
			input:          `2021-08-16T06:00:00`,
			configTimeZone: `timeZone = "America/Los_Angeles"`,
			localZoneName:  "Asia/Tokyo",
			expectedFull:   "6:00:00 am Pacific Daylight Time",
			expectedLong:   "6:00:00 am PDT",
			comment:        "NOTE: This does not show `UTC`. The original value type is AsTimeProvider provided by go-toml.",
		},
		{
			input:          `"2021-08-16"`,
			configTimeZone: `timeZone = "America/Los_Angeles"`,
			localZoneName:  "Asia/Tokyo",
			expectedFull:   "0:00:00 am Pacific Daylight Time",
			expectedLong:   "0:00:00 am PDT",
			comment:        "Date only string type, `cast` uses timzeZone value; show zoneinfo",
		},
		{
			input:          `2021-08-16`,
			configTimeZone: `timeZone = "America/Los_Angeles"`,
			localZoneName:  "Asia/Tokyo",
			expectedFull:   "0:00:00 am Pacific Daylight Time",
			expectedLong:   "0:00:00 am PDT",
			comment:        "Date only time type, AsTimeProvider; show zoneinfo",
		},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf(`%s;%s;%s`, tc.input, tc.configTimeZone, tc.localZoneName)
		t.Run(name, func(t *testing.T) {
			c := qt.New(t)

			// "TZ" (or the other) environment variables detmerine global `time.Local`.
			// What we really want to test is them (#9996). However,
			// This initialization can be called only once and hard to test.
			// Instead of "TZ", we use and change `time.Local`.
			// It is ok because both `go-toml` and `cast` use the value.
			origLocal := *time.Local
			testLoc, err := time.LoadLocation(tc.localZoneName)
			c.Assert(err, qt.IsNil)
			time.Local = testLoc
			t.Cleanup(func () {
				 time.Local = &origLocal
			})

			files := `
-- config.toml --
%s
-- content/p.md --
+++
title = "title"
date = %s
+++
-- layouts/index.html --
{{ $t := (site.GetPage "p").Date }}
_{{ $t | time.Format ":time_full" }}_
_{{ $t | time.Format ":time_long" }}_
`
			txtar := fmt.Sprintf(
				files,
				tc.configTimeZone,
				tc.input,
			)

			b := NewIntegrationTestBuilder(
				IntegrationTestConfig{
					T:           t,
					TxtarString: txtar,
				},
			).Build()

			content := fmt.Sprintf("_%s_\n_%s_", tc.expectedFull, tc.expectedLong)
			b.AssertFileContent("public/index.html", content)

		})
	}

}

func TestTOMLMultilingualTimeZone(t *testing.T) {
	t.Parallel()

	files := `
-- config.toml --
defaultContentLanguage = 'en'

[languages]
  [languages.en]
    title = 'English site'
    timeZone = "America/Los_Angeles"
  [languages.nn]
    title = 'Norwegian site'
    timeZone = "Europe/Oslo"
  [languages.ja]
    title = 'Japanese site'
    timeZone = "Asia/Tokyo"
-- content/matched.en.md --
+++
title = "English Content"
date = 2021-08-16T06:00:00-07:00
+++
-- content/matched.nn.md --
+++
title = "Norwegian Content"
date = 2021-08-16T06:00:00+02:00
+++
-- content/matched.ja.md --
+++
title = "Japanese Content"
date = 2021-08-16T06:00:00+09:00
+++
-- content/unmatched.en.md --
+++
title = "English Content"
date = 2021-08-16T06:00:00+09:00
+++
-- content/unmatched.nn.md --
+++
title = "Norwegian Content"
date = 2021-08-16T06:00:00+09:00
+++
-- content/unmatched.ja.md --
+++
title = "Japanese Content"
date = 2021-08-16T06:00:00-07:00
+++
-- content/utc.en.md --
+++
title = "English Content"
date = 2021-08-16T06:00:00+00:00
+++
-- content/utc.nn.md --
+++
title = "Norwegian Content"
date = 2021-08-16T06:00:00+00:00
+++
-- content/utc.ja.md --
+++
title = "Japanese Content"
date = 2021-08-16T06:00:00+00:00
+++
-- content/nooffset.en.md --
+++
title = "English Content"
date = "2021-08-16T06:00:00"
+++
-- content/nooffset.nn.md --
+++
title = "Norwegian Content"
date = "2021-08-16T06:00:00"
+++
-- content/nooffset.ja.md --
+++
title = "Japanese Content"
date = "2021-08-16T06:00:00"
+++
-- layouts/_default/single.html --
_{{ .Date | time.Format ":time_full" }}_
`

	b := NewIntegrationTestBuilder(
		IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	).Build()

	b.AssertFileContent("public/matched/index.html", `6:00:00 am Pacific Daylight Time`)
	b.AssertFileContent("public/nn/matched/index.html", `kl. 06:00:00 CEST`)
	b.AssertFileContent("public/ja/matched/index.html", `6時00分00秒 日本標準時`)

	b.AssertFileContent("public/unmatched/index.html", `6:00:00 am `)
	b.AssertFileContent("public/nn/unmatched/index.html", `kl. 06:00:00 `)
	b.AssertFileContent("public/ja/unmatched/index.html", `6時00分00秒`)

	b.AssertFileContent("public/utc/index.html", `6:00:00 am UTC`)
	b.AssertFileContent("public/nn/utc/index.html", `kl. 06:00:00 UTC`)
	b.AssertFileContent("public/ja/utc/index.html", `6時00分00秒`)

	b.AssertFileContent("public/nooffset/index.html", `6:00:00 am Pacific Daylight Time`)
	b.AssertFileContent("public/nn/nooffset/index.html", `kl. 06:00:00 CEST`)
	b.AssertFileContent("public/ja/nooffset/index.html", `6時00分00秒 日本標準時`)
}
