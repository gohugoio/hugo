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
	"strings"
	"testing"
	gtime "time"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/tpl/time"

	translators "github.com/gohugoio/localescompressed"
)

func TestTimeLocation(t *testing.T) {
	t.Parallel()

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: t},
	).Build()

	loc, _ := gtime.LoadLocation("America/Antigua")
	ns := time.New(htime.NewTimeFormatter(translators.GetTranslator("en")), loc, b.H.Deps)

	for i, test := range []struct {
		name     string
		value    string
		location any
		expect   any
	}{
		{"Empty location", "2020-10-20", "", "2020-10-20 00:00:00 +0000 UTC"},
		{"New location", "2020-10-20", nil, "2020-10-20 00:00:00 -0400 AST"},
		{"New York EDT", "2020-10-20", "America/New_York", "2020-10-20 00:00:00 -0400 EDT"},
		{"New York EST", "2020-01-20", "America/New_York", "2020-01-20 00:00:00 -0500 EST"},
		{"Empty location, time", "2020-10-20 20:33:59", "", "2020-10-20 20:33:59 +0000 UTC"},
		{"New York, time", "2020-10-20 20:33:59", "America/New_York", "2020-10-20 20:33:59 -0400 EDT"},
		// The following have an explicit offset specified. In this case, it overrides timezone
		{"Offset minus 0700, empty location", "2020-09-23T20:33:44-0700", "", "2020-09-23 20:33:44 -0700 -0700"},
		{"Offset plus 0200, empty location", "2020-09-23T20:33:44+0200", "", "2020-09-23 20:33:44 +0200 +0200"},

		{"Offset, New York", "2020-09-23T20:33:44-0700", "America/New_York", "2020-09-23 20:33:44 -0700 -0700"},
		{"Offset, Oslo", "2020-09-23T20:33:44+0200", "Europe/Oslo", "2020-09-23 20:33:44 +0200 +0200"},

		// Failures.
		{"Invalid time zone", "2020-01-20", "invalid-timezone", false},
		{"Invalid time value", "invalid-value", "", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			var args []any
			if test.location != nil {
				args = append(args, test.location)
			}
			result, err := ns.AsTime(test.value, args...)
			if b, ok := test.expect.(bool); ok && !b {
				if err == nil {
					t.Errorf("[%d] AsTime didn't return an expected error, got %v", i, result)
				}
			} else {
				if err != nil {
					t.Errorf("[%d] AsTime failed: %s", i, err)
					return
				}

				// See https://github.com/gohugoio/hugo/issues/8843#issuecomment-891551447
				// Drop the location string (last element) when comparing,
				// as that may change depending on the local locale.
				timeStr := result.(gtime.Time).String()
				timeStr = timeStr[:strings.LastIndex(timeStr, " ")]
				if !strings.HasPrefix(test.expect.(string), timeStr) {
					t.Errorf("[%d] AsTime got %v but expected %v", i, timeStr, test.expect)
				}
			}
		})
	}
}

func TestFormat(t *testing.T) {
	c := qt.New(t)

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: t},
	).Build()

	c.Run("UTC", func(c *qt.C) {
		c.Parallel()

		ns := time.New(htime.NewTimeFormatter(translators.GetTranslator("en")), gtime.UTC, b.H.Deps)

		for i, test := range []struct {
			layout string
			value  any
			expect any
		}{
			{"Monday, Jan 2, 2006", "2015-01-21", "Wednesday, Jan 21, 2015"},
			{"Monday, Jan 2, 2006", gtime.Date(2015, gtime.January, 21, 0, 0, 0, 0, gtime.UTC), "Wednesday, Jan 21, 2015"},
			{"This isn't a date layout string", "2015-01-21", "This isn't a date layout string"},
			// The following test case gives either "Tuesday, Jan 20, 2015" or "Monday, Jan 19, 2015" depending on the local time zone
			{"Monday, Jan 2, 2006", 1421733600, gtime.Unix(1421733600, 0).Format("Monday, Jan 2, 2006")},
			{"Monday, Jan 2, 2006", 1421733600.123, false},
			{gtime.RFC3339, gtime.Date(2016, gtime.March, 3, 4, 5, 0, 0, gtime.UTC), "2016-03-03T04:05:00Z"},
			{gtime.RFC1123, gtime.Date(2016, gtime.March, 3, 4, 5, 0, 0, gtime.UTC), "Thu, 03 Mar 2016 04:05:00 UTC"},
			{gtime.RFC3339, "Thu, 03 Mar 2016 04:05:00 UTC", "2016-03-03T04:05:00Z"},
			{gtime.RFC1123, "2016-03-03T04:05:00Z", "Thu, 03 Mar 2016 04:05:00 UTC"},
			// Custom layouts, as introduced in Hugo 0.87.
			{":date_medium", "2015-01-21", "Jan 21, 2015"},
		} {
			result, err := ns.Format(test.layout, test.value)
			if b, ok := test.expect.(bool); ok && !b {
				if err == nil {
					c.Errorf("[%d] DateFormat didn't return an expected error, got %v", i, result)
				}
			} else {
				if err != nil {
					c.Errorf("[%d] DateFormat failed: %s", i, err)
					continue
				}
				if result != test.expect {
					c.Errorf("[%d] DateFormat got %v but expected %v", i, result, test.expect)
				}
			}
		}
	})

	// Issue #9084
	c.Run("TZ America/Los_Angeles", func(c *qt.C) {
		c.Parallel()

		loc, err := gtime.LoadLocation("America/Los_Angeles")
		c.Assert(err, qt.IsNil)
		ns := time.New(htime.NewTimeFormatter(translators.GetTranslator("en")), loc, b.H.Deps)

		d, err := ns.Format(":time_full", "2020-03-09T11:00:00")

		c.Assert(err, qt.IsNil)
		c.Assert(d, qt.Equals, "11:00:00 am Pacific Daylight Time")
	})
}

func TestDuration(t *testing.T) {
	t.Parallel()

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: t},
	).Build()

	ns := time.New(htime.NewTimeFormatter(translators.GetTranslator("en")), gtime.UTC, b.H.Deps)

	for i, test := range []struct {
		unit   any
		num    any
		expect any
	}{
		{"nanosecond", 10, 10 * gtime.Nanosecond},
		{"ns", 10, 10 * gtime.Nanosecond},
		{"microsecond", 20, 20 * gtime.Microsecond},
		{"us", 20, 20 * gtime.Microsecond},
		{"µs", 20, 20 * gtime.Microsecond},
		{"millisecond", 20, 20 * gtime.Millisecond},
		{"ms", 20, 20 * gtime.Millisecond},
		{"second", 30, 30 * gtime.Second},
		{"s", 30, 30 * gtime.Second},
		{"minute", 20, 20 * gtime.Minute},
		{"m", 20, 20 * gtime.Minute},
		{"hour", 20, 20 * gtime.Hour},
		{"h", 20, 20 * gtime.Hour},
		{"hours", 20, false},
		{"hour", "30", 30 * gtime.Hour},
	} {
		result, err := ns.Duration(test.unit, test.num)
		if b, ok := test.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] Duration didn't return an expected error, got %v", i, result)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] Duration failed: %s", i, err)
				continue
			}
			if result != test.expect {
				t.Errorf("[%d] Duration got %v but expected %v", i, result, test.expect)
			}
		}
	}
}

func TestIn(t *testing.T) {
	t.Parallel()

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: t},
	).Build()

	ns := time.New(htime.NewTimeFormatter(translators.GetTranslator("en")), gtime.UTC, b.H.Deps)

	in := gtime.Date(2025, gtime.March, 31, 15, 0, 0, 0, gtime.UTC)

	tests := []struct {
		name    string
		tzn     string // time zone name
		want    string
		wantErr bool
	}{
		{name: "A", tzn: "America/Denver", want: "2025-03-31T09:00:00-06:00", wantErr: false},
		{name: "B", tzn: "Australia/Adelaide", want: "2025-04-01T01:30:00+10:30", wantErr: false},
		{name: "C", tzn: "Europe/Oslo", want: "2025-03-31T17:00:00+02:00", wantErr: false},
		{name: "D", tzn: "UTC", want: "2025-03-31T15:00:00+00:00", wantErr: false},
		{name: "E", tzn: "", want: "2025-03-31T15:00:00+00:00", wantErr: false},
		{name: "F", tzn: "InvalidTimeZoneName", want: "0001-01-01T00:00:00+00:00", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ns.In(tt.tzn, in)
			if (err != nil) != tt.wantErr {
				t.Errorf("time.In() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := result.Format("2006-01-02T15:04:05-07:00")
			if got != tt.want {
				t.Errorf("time.In() = %v, want %v", got, tt.want)
			}
		})
	}
}

// For benchmark tests below.
var timeZoneNames []string = []string{"America/New_York", "Europe/Oslo", "Australia/Sydney", "UTC", "Local"}

func BenchmarkInWithCaching(b *testing.B) {
	bb := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{T: b},
	).Build()

	ns := time.New(htime.NewTimeFormatter(translators.GetTranslator("en")), gtime.UTC, bb.H.Deps)

	for i := 0; i < b.N; i++ {
		timeZoneName := timeZoneNames[i%len(timeZoneNames)]
		_, err := ns.In(timeZoneName, gtime.Now())
		if err != nil {
			b.Fatalf("Error during benchmark: %v", err)
		}
	}
}

func BenchmarkInWithoutCaching(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timeZoneName := timeZoneNames[i%len(timeZoneNames)]
		location, err := gtime.LoadLocation(timeZoneName)
		if err != nil {
			b.Fatalf("Error during benchmark: %v", err)
		}

		_ = gtime.Now().In(location)
	}
}
