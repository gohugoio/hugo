// Copyright 2018 The Hugo Authors. All rights reserved.
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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gohugoio/hugo/deps"
	"github.com/stretchr/testify/require"
)

// pages with various combinations of front matter items
const (
	NO______ = "Page With empty front matter"

	fm______ = `---

---
Page With empty front matter`

	fm_D____ = `---
title: Simple
date: '2013-10-15T06:16:13'
---
Page With Date only`

	fm__P___ = `---
title: Simple
publishdate: '1969-01-10T09:17:42'
---
Page With PublishDate only`

	fm___L__ = `---
title: Simple
lastmod: '2017-09-03T22:22:22'
---
Page With Date and PublishDate`

	fm____M_ = `---
title: Simple
modified: '2018-01-24T12:21:39'
---
Page With Date and PublishDate`

	fm_____E = `---
title: Simple
expirydate: '2025-12-31T23:59:59'
---
Page With Date and PublishDate`

	fm_DP___ = `---
title: Simple
date: '2013-10-15T06:16:13'
publishdate: '1969-01-10T09:17:42'
---
Page With Date and PublishDate`

	fm__PL__ = `---
title: Simple
publishdate: '1969-01-10T09:17:42'
lastmod: '2017-09-03T22:22:22'
---
Page With Date and PublishDate`

	fm_D_L__ = `---
title: Simple
date: '2013-10-15T06:16:13'
lastmod: '2017-09-03T22:22:22'
---
Page With Date and PublishDate`

	fm__PLME = `---
title: Simple
publishdate: '1969-01-10T09:17:42'
lastmod: '2017-09-03T22:22:22'
modified: '2018-01-24T12:21:39'
expirydate: '2025-12-31T23:59:59'
---
Page With Date, PublishDate and LastMod`

	fm_D_LME = `---
title: Simple
date: '2013-10-15T06:16:13'
lastmod: '2017-09-03T22:22:22'
modified: '2018-01-24T12:21:39'
expirydate: '2025-12-31T23:59:59'
---
Page With Date, PublishDate and LastMod`

	fm_DP_ME = `---
title: Simple
date: '2013-10-15T06:16:13'
publishdate: '1969-01-10T09:17:42'
modified: '2018-01-24T12:21:39'
expirydate: '2025-12-31T23:59:59'
---
Page With Date, PublishDate and LastMod`

	fm_DPL_E = `---
title: Simple
date: '2013-10-15T06:16:13'
publishdate: '1969-01-10T09:17:42'
lastmod: '2017-09-03T22:22:22'
expirydate: '2025-12-31T23:59:59'
---
Page With Date, PublishDate and LastMod`

	fm_DPLM_ = `---
title: Simple
date: '2013-10-15T06:16:13'
publishdate: '1969-01-10T09:17:42'
lastmod: '2017-09-03T22:22:22'
modified: '2018-01-24T12:21:39'
---
Page With Date, PublishDate and LastMod`

	fm_DPL__ = `---
title: Simple
date: '2013-10-15T06:16:13'
publishdate: '1969-01-10T09:17:42'
lastmod: '2017-09-03T22:22:22'
---
Page With Date, PublishDate and LastMod`

	fm_DPLME = `---
title: Simple
date: '2013-10-15T06:16:13'
publishdate: '1969-01-10T09:17:42'
lastmod: '2017-09-03T22:22:22'
modified: '2018-01-24T12:21:39'
expirydate: '2025-12-31T23:59:59'
---
Page With Date, PublishDate and LastMod`
)

// the value for each kind of date
// field if set in the above pages, keyed
// by a single char name for reference
// in the test case table below.
// Since a given type (e.g. PublishDate)
// is always set to the same value, it
// makes it easy to visualize and test
// which input date field gets copied to
// which output page variable.
const (
	D = "2013-10-15T06:16:13Z" // front matter date
	P = "1969-01-10T09:17:42Z" // front matter publishdate
	L = "2017-09-03T22:22:22Z" // front matter lastmod
	M = "2018-01-24T12:21:39Z" // front matter modified
	E = "2025-12-31T23:59:59Z" // front matter expirydate
	//F = "2016-06-09T00:00:00Z"  // filename date
	o = "0001-01-01T00:00:00Z" // zero value of type Time, default for some date fields
	x = "nil"                  // nil date value, default for some date fields
	//c = "fs create timestamp" // file system timestamp, actual value determined at test runtime
	m = "fs mod timestamp" // file system timestamp, actual value determined at test runtime
)

// enumeration of defaultDate config options
const (
	NO = "none" // explicitly disabled
	MT = "file modification timestamp"
	//FN = "filename"             // NOT (YET) SUPPORTED (See issue #285, pr #3310, #3762, #3859)

	// Represents the defaultDate option not being set,
	// used to test the "default defaultDate" is as expected.
	// DO NOT manually create tests for this configuration
	// case, as they are auto-generated based on the
	// defaultDefaultDate value below
	NS = "not set"
)

const (
	defaultDefaultDate = NO
)

type DateDefaultTestCombination struct {
	pageText         string
	filename         string
	defaultDate      string
	expectedDate     string
	expectedPubDate  string
	expectedLastMod  string
	expectedModified string
	expectedExpiry   string
}

func TestDateDefaultsUnderVariousInputCombinations(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	var tests = map[string]DateDefaultTestCombination{
		// The INPUT columns on the left are the test case inputs:
		//
		//   ID:           Used to name individual tests (for failure messages,
		//                 and for human discussions)
		//
		//   front matter: The name indicates whether front matter exists (NO/fm)
		//                 and if it does, which dates are set:
		//                   - D: date
		//                   - P: publishdate
		//                   - L: lastmod
		//                   - M: modified
		//                   - E: expirydate
		//
		//                 For example, "fm__PL__" means content with publishdate
		//                 and lastMod specified in the front matter.
		//
		//   filename:     For when we start deriving metadata from it
		//
		//   defaultDate:  Whether or not useModTimeAsFallback is enabled.
		//
		// The OUTPUT columns on the right are the expected outputs to each of
		// the date Page Variables and Page Params.
		//
		//   Column heading: The Page Variable and Page Param to test
		//                   - D: Date variable, date page param
		//                   - P: PublishDate variable, publishdate page param
		//                   - L: LastMod variable, lastmod page param
		//                   - M: modified page param
		//                   - E: ExpiryDate variable, expirydate page param
		//
		//   Column value: the expected value for that variable/param, denoted
		//                 using the key for the source of the value. The
		//                 keys are the same as used in the front matter column,
		//				   with the addition of:
		//                   - x: nil
		//                   - o: zero value for type Time, "0001-01-01T00:00:00Z"
		//                   - m: the value of the file's modification timestamp
		//                   - c: the value of the file's creation timestamp
		//                   - F = filename
		//                 This makes it easy to visualize and test which
		//                 input date field gets copied to which output page
		//                 variable.
		//
		//   |-------------- inputs --------------|-- outputs --|
		//ID  frontmatter      filename       def  D  P  L  M  E
		"N1": {NO______, "___________test.md", NO, o, o, o, x, o}, // 4 year-one dates, 1 nil date (inconsistent)
		"N2": {fm______, "___________test.md", NO, o, o, o, x, o}, //   do we really want to any defaults to be 0001-01-01?

		"1D": {fm_D____, "___________test.md", NO, D, D, D, x, o}, // date fills 3
		"1P": {fm__P___, "___________test.md", NO, P, P, P, x, o}, // publishdate fills 3
		"1L": {fm___L__, "___________test.md", NO, o, o, L, x, o}, // lastmod doesn't fill modified (compare 1M)
		"1M": {fm____M_, "___________test.md", NO, o, o, M, M, o}, // modified fills lastmod
		"1E": {fm_____E, "___________test.md", NO, o, o, o, x, E},

		"21": {fm_DP___, "___________test.md", NO, D, P, D, x, o}, // date overrides publishdate to fill lastmod (compare 1D, 1P)
		"22": {fm__PL__, "___________test.md", NO, P, P, L, x, o},
		"23": {fm_D_L__, "___________test.md", NO, D, D, L, x, o},

		"4D": {fm__PLME, "___________test.md", NO, P, P, L, M, E},
		"4P": {fm_D_LME, "___________test.md", NO, D, D, L, M, E},
		"4L": {fm_DP_ME, "___________test.md", NO, D, P, M, M, E}, // modified fills lastmod
		"4M": {fm_DPL_E, "___________test.md", NO, D, P, L, x, E}, // lastmod doesn't fill modified (compare 4E)
		"4E": {fm_DPLM_, "___________test.md", NO, D, P, L, M, o},

		"x1": {fm_DPL__, "___________test.md", NO, D, P, L, x, o},
		"x2": {fm_DPLME, "___________test.md", NO, D, P, L, M, E},

		"M0": {NO______, "___________test.md", MT, m, o, m, x, o}, // MT fills date, lastmod, but not modified. Filled date doesn't flow to publishdate per 1D
		"ML": {fm___L__, "___________test.md", MT, m, o, L, x, o}, // MT files date, beating lastmod in front matter! Front matter should always take precedence.
		"MM": {fm____M_, "___________test.md", MT, m, o, M, M, o}, // MT files date, while modified fills lastmod and modified
		"MD": {fm_D____, "___________test.md", MT, D, D, D, x, o}, // date overrides MT, but modified not set either by date or MT
		"MP": {fm__P___, "___________test.md", MT, P, P, P, x, o}, // same issues as MD
		"Ma": {fm_DP___, "___________test.md", MT, D, P, D, x, o}, // MT does not fill modified
		"Mb": {fm__PL__, "___________test.md", MT, P, P, L, x, o}, //    "
		"Mc": {fm_DPL__, "___________test.md", MT, D, P, L, x, o}, //    "

		// NOT YET SUPPORTED
		//{"F0", NO______, "2016-06-09_test.md", NO, o, o, o, x, x}, // date in filename not used
		//{"F1", NO______, "2016-06-09_test.md", FN, F, F, F, x, x}, // date in filename used
		//{"F2", fm__P___, "2016-06-09_test.md", FN, F, P, F, x, x}, // date in both frontmatter and filename
	}
	fmt.Println("Base default date test suite: ", len(tests))

	// Automatically add test cases to guard against Hugo's default behavior
	// changing unintentionally.
	for id, test := range tests {
		if test.defaultDate == defaultDefaultDate {
			dtest := DateDefaultTestCombination{
				pageText:         test.pageText,
				filename:         test.filename,
				defaultDate:      NS,
				expectedDate:     test.expectedDate,
				expectedPubDate:  test.expectedPubDate,
				expectedLastMod:  test.expectedLastMod,
				expectedModified: test.expectedModified,
				expectedExpiry:   test.expectedExpiry,
			}
			tests[id+"(default cfg variant)"] = dtest
		}
	}
	fmt.Println("Base default date test suite + automatic default cfg tests: ", len(tests))

	for id, test := range tests {
		var (
			cfg, fs = newTestCfg()
		)

		writeToFs(t, fs.Source, filepath.Join("content", test.filename), test.pageText)

		switch test.defaultDate {
		case NS:
			// we're testing Hugo's default setting, whatever that is,
			// to make sure it isn't changed unintentionally
		case NO:
			cfg.Set("useModTimeAsFallback", false)
		case MT:
			cfg.Set("useModTimeAsFallback", true)
		default:
			t.Errorf("test %s using unsupported option", id)
		}

		s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{SkipRender: true})

		assert.Equal(1, len(s.RegularPages))
		p := s.RegularPages[0]
		fi := p.Source.FileInfo()

		// check Page Variables
		checkDate(t, id, "Date", p.Date, test.expectedDate, fi)
		checkDate(t, id, "PublishDate", p.PublishDate, test.expectedPubDate, fi)
		checkDate(t, id, "LastMod", p.Lastmod, test.expectedLastMod, fi)
		checkDate(t, id, "ExpiryDate", p.ExpiryDate, test.expectedExpiry, fi)

		// check Page Params
		checkDate(t, id, "param date", p.params["date"], test.expectedDate, fi)
		checkDate(t, id, "param publishdate", p.params["publishdate"], test.expectedPubDate, fi)
		checkDate(t, id, "param modified", p.params["modified"], test.expectedModified, fi)
		checkDate(t, id, "param expirydate", p.params["expirydate"], test.expectedExpiry, fi)
	}
}

func checkDate(t *testing.T, testId string, dateType string, given interface{}, expected string, fi os.FileInfo) {
	switch given.(type) {
	case nil:
		if expected != x {
			t.Errorf("test %s: %s is nil but expected \"%s\"", testId, dateType, expected)
		}
	case string:
		if given != expected {
			t.Errorf("test %s: %s is \"%s\" but expected \"%s\"", testId, dateType, given, expected)
		}
	case time.Time:
		// convert expected to time.Time for comparison
		var expectedTime time.Time
		if expected == m {
			expectedTime = fi.ModTime()
		} else {
			et, err := time.Parse(time.RFC3339, expected)
			if err != nil {
				t.Errorf("test %s: %s is %s but expected \"%s\"", testId, dateType, given, expected)
				return
			}
			expectedTime = et
		}
		if given != expectedTime {
			t.Errorf("test %s: %s is %s but expected %s", testId, dateType, given, expectedTime)
		}
	default:
		t.Errorf("test %s: %s is unexpected type: %T", testId, dateType, given)
	}
}
