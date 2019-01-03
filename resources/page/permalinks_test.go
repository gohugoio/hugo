// Copyright 2019 The Hugo Authors. All rights reserved.
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

package page

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// testdataPermalinks is used by a couple of tests; the expandsTo content is
// subject to the data in simplePageJSON.
var testdataPermalinks = []struct {
	spec      string
	valid     bool
	expandsTo string
}{
	{":title", true, "spf13-vim-3.0-release-and-new-website"},
	{"/:year-:month-:title", true, "/2012-04-spf13-vim-3.0-release-and-new-website"},
	{"/:year/:yearday/:month/:monthname/:day/:weekday/:weekdayname/", true, "/2012/97/04/April/06/5/Friday/"}, // Dates
	{"/:section/", true, "/blue/"},                                // Section
	{"/:title/", true, "/spf13-vim-3.0-release-and-new-website/"}, // Title
	{"/:slug/", true, "/the-slug/"},                               // Slug
	// TODO(bep) page {"/:filename/", true, "/test-page/"},                          // Filename
	// TODO(moorereason): need test scaffolding for this.
	//{"/:sections/", false, "/blue/"},                              // Sections

	// Failures
	{"/blog/:fred", false, ""},
	{"/:year//:title", false, ""},
}

func TestPermalinkExpansion(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	page := newTestPage()
	page.title = "Spf13 Vim 3.0 Release and new website"
	d, _ := time.Parse("2006-01-02", "2012-04-06")
	page.date = d
	page.section = "blue"
	page.slug = "The Slug"

	for i, item := range testdataPermalinks {

		msg := fmt.Sprintf("Test %d", i)

		if !item.valid {
			continue
		}

		permalinksConfig := map[string]string{
			"posts": item.spec,
		}

		ps := newTestPathSpec()
		ps.Cfg.Set("permalinks", permalinksConfig)

		expander, err := NewPermalinkExpander(ps)
		assert.NoError(err)

		expanded, err := expander.Expand("posts", page)
		assert.NoError(err)
		assert.Equal(item.expandsTo, expanded, msg)

	}
}

func TestPermalinkExpansionMultiSection(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	page := newTestPage()
	page.title = "Page Title"
	d, _ := time.Parse("2006-01-02", "2012-04-06")
	page.date = d
	page.section = "blue"
	page.slug = "The Slug"

	permalinksConfig := map[string]string{
		"posts": "/:slug",
		"blog":  "/:section/:year",
	}

	ps := newTestPathSpec()
	ps.Cfg.Set("permalinks", permalinksConfig)

	expander, err := NewPermalinkExpander(ps)
	assert.NoError(err)

	expanded, err := expander.Expand("posts", page)
	assert.NoError(err)
	assert.Equal("/the-slug", expanded)

	expanded, err = expander.Expand("blog", page)
	assert.NoError(err)
	assert.Equal("/blue/2012", expanded)

}

func TestPermalinkExpansionConcurrent(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	permalinksConfig := map[string]string{
		"posts": "/:slug/",
	}

	ps := newTestPathSpec()
	ps.Cfg.Set("permalinks", permalinksConfig)

	expander, err := NewPermalinkExpander(ps)
	assert.NoError(err)

	var wg sync.WaitGroup

	for i := 1; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			page := newTestPage()
			for j := 1; j < 20; j++ {
				page.slug = fmt.Sprintf("slug%d", i+j)
				expanded, err := expander.Expand("posts", page)
				assert.NoError(err)
				assert.Equal(fmt.Sprintf("/%s/", page.slug), expanded)
			}
		}(i)
	}

	wg.Wait()
}

func BenchmarkPermalinkExpand(b *testing.B) {
	page := newTestPage()
	page.title = "Hugo Rocks"
	d, _ := time.Parse("2006-01-02", "2019-02-28")
	page.date = d

	permalinksConfig := map[string]string{
		"posts": "/:year-:month-:title",
	}

	ps := newTestPathSpec()
	ps.Cfg.Set("permalinks", permalinksConfig)

	expander, err := NewPermalinkExpander(ps)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, err := expander.Expand("posts", page)
		if err != nil {
			b.Fatal(err)
		}
		if s != "/2019-02-hugo-rocks" {
			b.Fatal(s)
		}

	}
}
