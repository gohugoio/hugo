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
[_build]
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
		bcfg, err := DecodeBuildConfig(cfg.Get("_build"))
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
		{"page.md", "0001-01-01", ""},
		{"2012-09-12-page.md", "2012-09-12", "page"},
		{"2018-02-28-page.md", "2018-02-28", "page"},
		{"2018-02-28_page.md", "2018-02-28", "page"},
		{"2018-02-28 page.md", "2018-02-28", "page"},
		{"2018-02-28page.md", "2018-02-28", "page"},
		{"2018-02-28-.md", "2018-02-28", ""},
		{"2018-02-28-.md", "2018-02-28", ""},
		{"2018-02-28.md", "2018-02-28", ""},
		{"2018-02-28-page", "2018-02-28", "page"},
		{"2012-9-12-page.md", "0001-01-01", ""},
		{"asdfasdf.md", "0001-01-01", ""},
	}

	for _, test := range tests {
		expecteFDate, err := time.Parse("2006-01-02", test.date)
		c.Assert(err, qt.IsNil)

		gotDate, gotSlug := dateAndSlugFromBaseFilename(time.UTC, test.name)

		c.Assert(gotDate, qt.Equals, expecteFDate)
		c.Assert(gotSlug, qt.Equals, test.slug)

	}
}

func TestExpandDefaultValues(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	c.Assert(expandDefaultValues([]string{"a", ":default", "d"}, []string{"b", "c"}), qt.DeepEquals, []string{"a", "b", "c", "d"})
	c.Assert(expandDefaultValues([]string{"a", "b", "c"}, []string{"a", "b", "c"}), qt.DeepEquals, []string{"a", "b", "c"})
	c.Assert(expandDefaultValues([]string{":default", "a", ":default", "d"}, []string{"b", "c"}), qt.DeepEquals, []string{"b", "c", "a", "b", "c", "d"})
}
