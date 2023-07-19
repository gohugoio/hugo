// Copyright 2021 The Hugo Authors. All rights reserved.
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

package htime

import (
	"log"
	"strings"
	"time"

	"github.com/bep/clocks"
	"github.com/spf13/cast"

	"github.com/gohugoio/locales"
)

var (
	longDayNames = []string{
		"Sunday",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
	}

	shortDayNames = []string{
		"Sun",
		"Mon",
		"Tue",
		"Wed",
		"Thu",
		"Fri",
		"Sat",
	}

	shortMonthNames = []string{
		"Jan",
		"Feb",
		"Mar",
		"Apr",
		"May",
		"Jun",
		"Jul",
		"Aug",
		"Sep",
		"Oct",
		"Nov",
		"Dec",
	}

	longMonthNames = []string{
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}

	Clock = clocks.System()
)

func NewTimeFormatter(ltr locales.Translator) TimeFormatter {
	if ltr == nil {
		panic("must provide a locales.Translator")
	}
	return TimeFormatter{
		ltr: ltr,
	}
}

// TimeFormatter is locale aware.
type TimeFormatter struct {
	ltr locales.Translator
}

func (f TimeFormatter) Format(t time.Time, layout string) string {
	if layout == "" {
		return ""
	}

	if layout[0] == ':' {
		// It may be one of Hugo's custom layouts.
		switch strings.ToLower(layout[1:]) {
		case "date_full":
			return f.ltr.FmtDateFull(t)
		case "date_long":
			return f.ltr.FmtDateLong(t)
		case "date_medium":
			return f.ltr.FmtDateMedium(t)
		case "date_short":
			return f.ltr.FmtDateShort(t)
		case "time_full":
			return f.ltr.FmtTimeFull(t)
		case "time_long":
			return f.ltr.FmtTimeLong(t)
		case "time_medium":
			return f.ltr.FmtTimeMedium(t)
		case "time_short":
			return f.ltr.FmtTimeShort(t)
		}
	}

	s := t.Format(layout)

	monthIdx := t.Month() - 1 // Month() starts at 1.
	dayIdx := t.Weekday()

	if strings.Contains(layout, "January") {
		s = strings.ReplaceAll(s, longMonthNames[monthIdx], f.ltr.MonthWide(t.Month()))
	} else if strings.Contains(layout, "Jan") {
		s = strings.ReplaceAll(s, shortMonthNames[monthIdx], f.ltr.MonthAbbreviated(t.Month()))
	}

	if strings.Contains(layout, "Monday") {
		s = strings.ReplaceAll(s, longDayNames[dayIdx], f.ltr.WeekdayWide(t.Weekday()))
	} else if strings.Contains(layout, "Mon") {
		s = strings.ReplaceAll(s, shortDayNames[dayIdx], f.ltr.WeekdayAbbreviated(t.Weekday()))
	}

	return s
}

func ToTimeInDefaultLocationE(i any, location *time.Location) (tim time.Time, err error) {
	switch vv := i.(type) {
	case AsTimeProvider:
		return vv.AsTime(location), nil
	// issue #8895
	// datetimes parsed by `go-toml` have empty zone name
	// convert back them into string and use `cast`
	// TODO(bep) add tests, make sure we really need this.
	case time.Time:
		i = vv.Format(time.RFC3339)
	}
	return cast.ToTimeInDefaultLocationE(i, location)
}

// Now returns time.Now() or time value based on the `clock` flag.
// Use this function to fake time inside hugo.
func Now() time.Time {
	return Clock.Now()
}

func Since(t time.Time) time.Duration {
	return Clock.Since(t)
}

// AsTimeProvider is implemented by go-toml's LocalDate and LocalDateTime.
type AsTimeProvider interface {
	AsTime(zone *time.Location) time.Time
}

// StopWatch is a simple helper to measure time during development.
func StopWatch(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("StopWatch %q took %s", name, time.Since(start))
	}
}
