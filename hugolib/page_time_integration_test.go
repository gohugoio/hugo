package hugolib

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

var PAGE_WITH_INVALID_DATE = `---
date: 2010-05-02_15:29:31+08:00
---
Page With Invalid Date (replace T with _ for RFC 3339)`

var PAGE_WITH_DATE_RFC3339 = `---
date: 2010-05-02T15:29:31+08:00
---
Page With Date RFC3339`

var PAGE_WITH_DATE_RFC3339_NO_T = `---
date: 2010-05-02 15:29:31+08:00
---
Page With Date RFC3339_NO_T`

var PAGE_WITH_DATE_RFC1123 = `---
date: Sun, 02 May 2010 15:29:31 PST
---
Page With Date RFC1123`

var PAGE_WITH_DATE_RFC1123Z = `---
date: Sun, 02 May 2010 15:29:31 +0800
---
Page With Date RFC1123Z`

var PAGE_WITH_DATE_RFC822 = `---
date: 02 May 10 15:29 PST
---
Page With Date RFC822`

var PAGE_WITH_DATE_RFC822Z = `---
date: 02 May 10 15:29 +0800
---
Page With Date RFC822Z`

var PAGE_WITH_DATE_ANSIC = `---
date: Sun May 2 15:29:31 2010
---
Page With Date ANSIC`

var PAGE_WITH_DATE_UnixDate = `---
date: Sun May 2 15:29:31 PST 2010
---
Page With Date UnixDate`

var PAGE_WITH_DATE_RubyDate = `---
date: Sun May 02 15:29:31 +0800 2010
---
Page With Date RubyDate`

var PAGE_WITH_DATE_HugoYearNumeric = `---
date: 2010-05-02
---
Page With Date HugoYearNumeric`

var PAGE_WITH_DATE_HugoYear = `---
date: 02 May 2010
---
Page With Date HugoYear`

var PAGE_WITH_DATE_HugoLong = `---
date: 02 May 2010 15:29 PST
---
Page With Date HugoLong`

func TestDegenerateDateFrontMatter(t *testing.T) {
	p, _ := NewPageFrom(strings.NewReader(PAGE_WITH_INVALID_DATE), "page/with/invalid/date")
	if p.Date != *new(time.Time) {
		t.Fatalf("Date should be set to time.Time zero value.  Got: %s", p.Date)
	}
}

func TestParsingDateInFrontMatter(t *testing.T) {
	tests := []struct {
		buf string
		dt  string
	}{
		{PAGE_WITH_DATE_RFC3339, "2010-05-02T15:29:31+08:00"},
		{PAGE_WITH_DATE_RFC3339_NO_T, "2010-05-02T15:29:31+08:00"},
		{PAGE_WITH_DATE_RFC1123Z, "2010-05-02T15:29:31+08:00"},
		{PAGE_WITH_DATE_RFC822Z, "2010-05-02T15:29:00+08:00"},
		{PAGE_WITH_DATE_ANSIC, "2010-05-02T15:29:31Z"},
		{PAGE_WITH_DATE_RubyDate, "2010-05-02T15:29:31+08:00"},
		{PAGE_WITH_DATE_HugoYearNumeric, "2010-05-02T00:00:00Z"},
		{PAGE_WITH_DATE_HugoYear, "2010-05-02T00:00:00Z"},
	}

	tzShortCodeTests := []struct {
		buf string
		dt  string
	}{
		{PAGE_WITH_DATE_RFC1123, "2010-05-02T15:29:31-08:00"},
		{PAGE_WITH_DATE_RFC822, "2010-05-02T15:29:00-08:00Z"},
		{PAGE_WITH_DATE_UnixDate, "2010-05-02T15:29:31-08:00"},
		{PAGE_WITH_DATE_HugoLong, "2010-05-02T15:21:00+08:00"},
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
		p, err := NewPageFrom(strings.NewReader(test.buf), "page/with/date")
		if err != nil {
			t.Fatalf("Expected to be able to parse page.")
		}
		if !dt.Equal(p.Date) {
			t.Errorf("Date does not equal frontmatter:\n%s\nExpecting: %s\n      Got: %s. Diff: %s\n internal: %#v\n           %#v", test.buf, dt, p.Date, dt.Sub(p.Date), dt, p.Date)
		}
	}
}
