// Copyright 2015 The Hugo Authors. All rights reserved.
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
	"strings"
	"testing"
)

// testdataPermalinks is used by a couple of tests; the expandsTo content is
// subject to the data in SIMPLE_PAGE_JSON.
var testdataPermalinks = []struct {
	spec      string
	valid     bool
	expandsTo string
}{
	//{"/:year/:month/:title/", true, "/2012/04/spf13-vim-3.0-release-and-new-website/"},
	//{"/:title", true, "/spf13-vim-3.0-release-and-new-website"},
	//{":title", true, "spf13-vim-3.0-release-and-new-website"},
	//{"/blog/:year/:yearday/:title", true, "/blog/2012/97/spf13-vim-3.0-release-and-new-website"},
	{"/:year-:month-:title", true, "/2012-04-spf13-vim-3.0-release-and-new-website"},
	{"/blog/:year-:month-:title", true, "/blog/2012-04-spf13-vim-3.0-release-and-new-website"},
	{"/blog-:year-:month-:title", true, "/blog-2012-04-spf13-vim-3.0-release-and-new-website"},
	//{"/blog/:fred", false, ""},
	//{"/:year//:title", false, ""},
	//{
	//"/:section/:year/:month/:day/:weekdayname/:yearday/:title",
	//true,
	//"/blue/2012/04/06/Friday/97/spf13-vim-3.0-release-and-new-website",
	//},
	//{
	//"/:weekday/:weekdayname/:month/:monthname",
	//true,
	//"/5/Friday/04/April",
	//},
	//{
	//"/:slug/:title",
	//true,
	//"/spf13-vim-3-0-release-and-new-website/spf13-vim-3.0-release-and-new-website",
	//},
}

func TestPermalinkValidation(t *testing.T) {
	t.Parallel()
	for _, item := range testdataPermalinks {
		pp := pathPattern(item.spec)
		have := pp.validate()
		if have == item.valid {
			continue
		}
		var howBad string
		if have {
			howBad = "validates but should not have"
		} else {
			howBad = "should have validated but did not"
		}
		t.Errorf("permlink spec %q %s", item.spec, howBad)
	}
}

func TestPermalinkExpansion(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	page, err := s.NewPageFrom(strings.NewReader(simplePageJSON), "blue/test-page.md")

	if err != nil {
		t.Fatalf("failed before we began, could not parse SIMPLE_PAGE_JSON: %s", err)
	}
	for _, item := range testdataPermalinks {
		if !item.valid {
			continue
		}
		pp := pathPattern(item.spec)
		result, err := pp.Expand(page)
		if err != nil {
			t.Errorf("failed to expand page: %s", err)
			continue
		}
		if result != item.expandsTo {
			t.Errorf("expansion mismatch!\n\tExpected: %q\n\tReceived: %q", item.expandsTo, result)
		}
	}
}
