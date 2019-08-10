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
	"testing"
	"time"

	"github.com/gohugoio/hugo/common/types"

	qt "github.com/frankban/quicktest"
)

func TestRelated(t *testing.T) {
	c := qt.New(t)

	t.Parallel()

	pages := Pages{
		&testPage{
			title:   "Page 1",
			pubDate: mustParseDate("2017-01-03"),
			params: map[string]interface{}{
				"keywords": []string{"hugo", "says"},
			},
		},
		&testPage{
			title:   "Page 2",
			pubDate: mustParseDate("2017-01-02"),
			params: map[string]interface{}{
				"keywords": []string{"hugo", "rocks"},
			},
		},
		&testPage{
			title:   "Page 3",
			pubDate: mustParseDate("2017-01-01"),
			params: map[string]interface{}{
				"keywords": []string{"bep", "says"},
			},
		},
	}

	result, err := pages.RelatedTo(types.NewKeyValuesStrings("keywords", "hugo", "rocks"))

	c.Assert(err, qt.IsNil)
	c.Assert(len(result), qt.Equals, 2)
	c.Assert(result[0].Title(), qt.Equals, "Page 2")
	c.Assert(result[1].Title(), qt.Equals, "Page 1")

	result, err = pages.Related(pages[0])
	c.Assert(err, qt.IsNil)
	c.Assert(len(result), qt.Equals, 2)
	c.Assert(result[0].Title(), qt.Equals, "Page 2")
	c.Assert(result[1].Title(), qt.Equals, "Page 3")

	result, err = pages.RelatedIndices(pages[0], "keywords")
	c.Assert(err, qt.IsNil)
	c.Assert(len(result), qt.Equals, 2)
	c.Assert(result[0].Title(), qt.Equals, "Page 2")
	c.Assert(result[1].Title(), qt.Equals, "Page 3")

	result, err = pages.RelatedTo(types.NewKeyValuesStrings("keywords", "bep", "rocks"))
	c.Assert(err, qt.IsNil)
	c.Assert(len(result), qt.Equals, 2)
	c.Assert(result[0].Title(), qt.Equals, "Page 2")
	c.Assert(result[1].Title(), qt.Equals, "Page 3")
}

func mustParseDate(s string) time.Time {
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return d
}
