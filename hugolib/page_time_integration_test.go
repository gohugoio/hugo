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
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/spf13/cast"
)

const (
	pageWithInvalidDate = `---
date: 2010-05-02_15:29:31+08:00
---
Page With Invalid Date (replace T with _ for RFC 3339)`

	pageWithDateRFC3339 = `---
date: 2010-05-02T15:29:31+08:00
---
Page With Date RFC3339`

	pageWithDateRFC3339NoT = `---
date: 2010-05-02 15:29:31+08:00
---
Page With Date RFC3339_NO_T`

	pageWithRFC1123 = `---
date: Sun, 02 May 2010 15:29:31 PST
---
Page With Date RFC1123`

	pageWithDateRFC1123Z = `---
date: Sun, 02 May 2010 15:29:31 +0800
---
Page With Date RFC1123Z`

	pageWithDateRFC822 = `---
date: 02 May 10 15:29 PST
---
Page With Date RFC822`

	pageWithDateRFC822Z = `---
date: 02 May 10 15:29 +0800
---
Page With Date RFC822Z`

	pageWithDateANSIC = `---
date: Sun May 2 15:29:31 2010
---
Page With Date ANSIC`

	pageWithDateUnixDate = `---
date: Sun May 2 15:29:31 PST 2010
---
Page With Date UnixDate`

	pageWithDateRubyDate = `---
date: Sun May 02 15:29:31 +0800 2010
---
Page With Date RubyDate`

	pageWithDateHugoYearNumeric = `---
date: 2010-05-02
---
Page With Date HugoYearNumeric`

	pageWithDateHugoYear = `---
date: 02 May 2010
---
Page With Date HugoYear`

	pageWithDateHugoLong = `---
date: 02 May 2010 15:29 PST
---
Page With Date HugoLong`
)

func TestDegenerateDateFrontMatter(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	p, _ := s.NewPageFrom(strings.NewReader(pageWithInvalidDate), "page/with/invalid/date")
	if p.Date != *new(time.Time) {
		t.Fatalf("Date should be set to time.Time zero value.  Got: %s", p.Date)
	}
}

func TestParsingDateInFrontMatter(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	tests := []struct {
		buf string
		dt  string
	}{
		{pageWithDateRFC3339, "2010-05-02T15:29:31+08:00"},
		{pageWithDateRFC3339NoT, "2010-05-02T15:29:31+08:00"},
		{pageWithDateRFC1123Z, "2010-05-02T15:29:31+08:00"},
		{pageWithDateRFC822Z, "2010-05-02T15:29:00+08:00"},
		{pageWithDateANSIC, "2010-05-02T15:29:31Z"},
		{pageWithDateRubyDate, "2010-05-02T15:29:31+08:00"},
		{pageWithDateHugoYearNumeric, "2010-05-02T00:00:00Z"},
		{pageWithDateHugoYear, "2010-05-02T00:00:00Z"},
	}

	tzShortCodeTests := []struct {
		buf string
		dt  string
	}{
		{pageWithRFC1123, "2010-05-02T15:29:31-08:00"},
		{pageWithDateRFC822, "2010-05-02T15:29:00-08:00Z"},
		{pageWithDateUnixDate, "2010-05-02T15:29:31-08:00"},
		{pageWithDateHugoLong, "2010-05-02T15:21:00+08:00"},
	}

	if _, err := time.LoadLocation("PST"); err == nil {
		tests = append(tests, tzShortCodeTests...)
	} else {
		fmt.Fprintf(os.Stderr, "Skipping shortname timezone tests.\n")
	}

	for _, test := range tests {
		dt, e := time.Parse(time.RFC3339, test.dt)
		if e != nil {
			t.Fatalf("Unable to parse date time (RFC3339) for running the test: %s", e)
		}
		p, err := s.NewPageFrom(strings.NewReader(test.buf), "page/with/date")
		if err != nil {
			t.Fatalf("Expected to be able to parse page.")
		}
		if !dt.Equal(p.Date) {
			t.Errorf("Date does not equal frontmatter:\n%s\nExpecting: %s\n      Got: %s. Diff: %s\n internal: %#v\n           %#v", test.buf, dt, p.Date, dt.Sub(p.Date), dt, p.Date)
		}
	}
}

// Temp test https://github.com/gohugoio/hugo/issues/3059
func TestParsingDateParallel(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup

	for j := 0; j < 100; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				dateStr := "2010-05-02 15:29:31 +08:00"

				dt, err := time.Parse("2006-01-02 15:04:05 -07:00", dateStr)
				if err != nil {
					t.Fatal(err)
				}

				if dt.Year() != 2010 {
					t.Fatal("time.Parse: Invalid date:", dt)
				}

				dt2 := cast.ToTime(dateStr)

				if dt2.Year() != 2010 {
					t.Fatal("cast.ToTime: Invalid date:", dt2.Year())
				}
			}
		}()
	}
	wg.Wait()

}
