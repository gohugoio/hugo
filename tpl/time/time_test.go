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

package time

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		layout string
		value  interface{}
		expect interface{}
	}{
		{"Monday, Jan 2, 2006", "2015-01-21", "Wednesday, Jan 21, 2015"},
		{"Monday, Jan 2, 2006", time.Date(2015, time.January, 21, 0, 0, 0, 0, time.UTC), "Wednesday, Jan 21, 2015"},
		{"This isn't a date layout string", "2015-01-21", "This isn't a date layout string"},
		// The following test case gives either "Tuesday, Jan 20, 2015" or "Monday, Jan 19, 2015" depending on the local time zone
		{"Monday, Jan 2, 2006", 1421733600, time.Unix(1421733600, 0).Format("Monday, Jan 2, 2006")},
		{"Monday, Jan 2, 2006", 1421733600.123, false},
		{time.RFC3339, time.Date(2016, time.March, 3, 4, 5, 0, 0, time.UTC), "2016-03-03T04:05:00Z"},
		{time.RFC1123, time.Date(2016, time.March, 3, 4, 5, 0, 0, time.UTC), "Thu, 03 Mar 2016 04:05:00 UTC"},
		{time.RFC3339, "Thu, 03 Mar 2016 04:05:00 UTC", "2016-03-03T04:05:00Z"},
		{time.RFC1123, "2016-03-03T04:05:00Z", "Thu, 03 Mar 2016 04:05:00 UTC"},
	} {
		result, err := ns.Format(test.layout, test.value)
		if b, ok := test.expect.(bool); ok && !b {
			if err == nil {
				t.Errorf("[%d] DateFormat didn't return an expected error, got %v", i, result)
			}
		} else {
			if err != nil {
				t.Errorf("[%d] DateFormat failed: %s", i, err)
				continue
			}
			if result != test.expect {
				t.Errorf("[%d] DateFormat got %v but expected %v", i, result, test.expect)
			}
		}
	}
}

func TestDuration(t *testing.T) {
	t.Parallel()

	ns := New()

	for i, test := range []struct {
		unit   interface{}
		num    interface{}
		expect interface{}
	}{
		{"nanosecond", 10, 10 * time.Nanosecond},
		{"ns", 10, 10 * time.Nanosecond},
		{"microsecond", 20, 20 * time.Microsecond},
		{"us", 20, 20 * time.Microsecond},
		{"µs", 20, 20 * time.Microsecond},
		{"millisecond", 20, 20 * time.Millisecond},
		{"ms", 20, 20 * time.Millisecond},
		{"second", 30, 30 * time.Second},
		{"s", 30, 30 * time.Second},
		{"minute", 20, 20 * time.Minute},
		{"m", 20, 20 * time.Minute},
		{"hour", 20, 20 * time.Hour},
		{"h", 20, 20 * time.Hour},
		{"hours", 20, false},
		{"hour", "30", 30 * time.Hour},
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
