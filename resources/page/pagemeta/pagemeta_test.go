// Copyright 2020 The Hugo Authors. All rights reserved.
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

package pagemeta

import (
	"fmt"
	"testing"
	"time"

	"github.com/gohugoio/hugo/htesting/hqt"

	"github.com/gohugoio/hugo/config"

	qt "github.com/frankban/quicktest"
)

func TestDecodeBuildConfig(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	configTempl := `
[build]
render = %s
list = %s
publishResources = true`

	for _, test := range []struct {
		args   []any
		expect BuildConfig
	}{
		{
			[]any{"true", "true"},
			BuildConfig{
				Render:           Always,
				List:             Always,
				PublishResources: true,
				set:              true,
			},
		},
		{[]any{"true", "false"}, BuildConfig{
			Render:           Always,
			List:             Never,
			PublishResources: true,
			set:              true,
		}},
		{[]any{`"always"`, `"always"`}, BuildConfig{
			Render:           Always,
			List:             Always,
			PublishResources: true,
			set:              true,
		}},
		{[]any{`"never"`, `"never"`}, BuildConfig{
			Render:           Never,
			List:             Never,
			PublishResources: true,
			set:              true,
		}},
		{[]any{`"link"`, `"local"`}, BuildConfig{
			Render:           Link,
			List:             ListLocally,
			PublishResources: true,
			set:              true,
		}},
		{[]any{`"always"`, `"asdfadf"`}, BuildConfig{
			Render:           Always,
			List:             Always,
			PublishResources: true,
			set:              true,
		}},
	} {
		cfg, err := config.FromConfigString(fmt.Sprintf(configTempl, test.args...), "toml")
		c.Assert(err, qt.IsNil)
		bcfg, err := DecodeBuildConfig(cfg.Get("build"))
		c.Assert(err, qt.IsNil)

		eq := qt.CmpEquals(hqt.DeepAllowUnexported(BuildConfig{}))

		c.Assert(bcfg, eq, test.expect)

	}
}

func TestDateAndSlugFromBaseFilename(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	tests := []struct {
		name string
		date string
		slug string
	}{
		// date
		{"2025-07-04 page.md", "2025-07-04T00:00:00+02:00", "page"},
		{"2025-07-04-page.md", "2025-07-04T00:00:00+02:00", "page"},
		{"2025-07-04_page.md", "2025-07-04T00:00:00+02:00", "page"},
		{"2025-07-04page.md", "2025-07-04T00:00:00+02:00", "page"},
		{"2025-07-04", "2025-07-04T00:00:00+02:00", ""},
		{"2025-07-04-.md", "2025-07-04T00:00:00+02:00", ""},
		{"2025-07-04.md", "2025-07-04T00:00:00+02:00", ""},
		// date and time
		{"2025-07-04-22-17-13 page.md", "2025-07-04T22:17:13+02:00", "page"},
		{"2025-07-04-22-17-13-page.md", "2025-07-04T22:17:13+02:00", "page"},
		{"2025-07-04-22-17-13_page.md", "2025-07-04T22:17:13+02:00", "page"},
		{"2025-07-04-22-17-13page.md", "2025-07-04T22:17:13+02:00", "page"},
		{"2025-07-04-22-17-13", "2025-07-04T22:17:13+02:00", ""},
		{"2025-07-04-22-17-13-.md", "2025-07-04T22:17:13+02:00", ""},
		{"2025-07-04-22-17-13.md", "2025-07-04T22:17:13+02:00", ""},
		// date and time with other separators between the two
		{"2025-07-04T22-17-13.md", "2025-07-04T22:17:13+02:00", ""},
		{"2025-07-04 22-17-13.md", "2025-07-04T22:17:13+02:00", ""},
		// no date or time
		{"something.md", "0001-01-01T00:00:00+00:00", ""},                // 9 chars
		{"some-thing-.md", "0001-01-01T00:00:00+00:00", ""},              // 10 chars
		{"somethingsomething.md", "0001-01-01T00:00:00+00:00", ""},       // 18 chars
		{"something-something.md", "0001-01-01T00:00:00+00:00", ""},      // 19 chars
		{"something-something-else.md", "0001-01-01T00:00:00+00:00", ""}, // 27 chars
		// invalid
		{"2025-07-4-page.md", "0001-01-01T00:00:00+00:00", ""},
		{"2025-07-4-22-17-13-page.md", "0001-01-01T00:00:00+00:00", ""},
		{"asdfasdf.md", "0001-01-01T00:00:00+00:00", ""},
	}

	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		t.Error("Unable to determine location from given time zone")
	}
	for _, test := range tests {

		gotDate, gotSlug := dateAndSlugFromBaseFilename(location, test.name)

		c.Assert(gotDate.Format("2006-01-02T15:04:05-07:00"), qt.Equals, test.date)
		c.Assert(gotSlug, qt.Equals, test.slug)

	}
}

func TestExpandDefaultValues(t *testing.T) {
	c := qt.New(t)
	c.Assert(expandDefaultValues([]string{"a", ":default", "d"}, []string{"b", "c"}), qt.DeepEquals, []string{"a", "b", "c", "d"})
	c.Assert(expandDefaultValues([]string{"a", "b", "c"}, []string{"a", "b", "c"}), qt.DeepEquals, []string{"a", "b", "c"})
	c.Assert(expandDefaultValues([]string{":default", "a", ":default", "d"}, []string{"b", "c"}), qt.DeepEquals, []string{"b", "c", "a", "b", "c", "d"})
}
