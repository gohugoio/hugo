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
	"strings"
	"time"

	"github.com/spf13/cast"

	toml "github.com/pelletier/go-toml/v2"

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

	s = strings.ReplaceAll(s, longMonthNames[monthIdx], f.ltr.MonthWide(t.Month()))
	if !strings.Contains(s, f.ltr.MonthWide(t.Month())) {
		s = strings.ReplaceAll(s, shortMonthNames[monthIdx], f.ltr.MonthAbbreviated(t.Month()))
	}
	s = strings.ReplaceAll(s, longDayNames[dayIdx], f.ltr.WeekdayWide(t.Weekday()))
	if !strings.Contains(s, f.ltr.WeekdayWide(t.Weekday())) {
		s = strings.ReplaceAll(s, shortDayNames[dayIdx], f.ltr.WeekdayAbbreviated(t.Weekday()))
	}

	return s
}

func ToTimeInDefaultLocationE(i interface{}, location *time.Location) (tim time.Time, err error) {
	switch vv := i.(type) {
	case toml.LocalDate:
		return vv.AsTime(location), nil
	case toml.LocalDateTime:
		return vv.AsTime(location), nil
	// issue #8895
	// datetimes parsed by `go-toml` have empty zone name
	// convert back them into string and use `cast`
	case time.Time:
		i = vv.Format(time.RFC3339)
	}
	return cast.ToTimeInDefaultLocationE(i, location)
}
