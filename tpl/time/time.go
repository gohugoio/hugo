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

// Package time provides template functions for measuring and displaying time.
package time

import (
	"errors"
	"fmt"
	"time"
	_time "time"

	"github.com/gohugoio/hugo/common/htime"

	"github.com/spf13/cast"

	translators "github.com/gohugoio/localescompressed"
)

// New returns a new instance of the time-namespaced template functions.
func New(langCode string, timeFormatter htime.TimeFormatter, location *time.Location) *Namespace {
	timeFormatters := make(map[string]htime.TimeFormatter)
	timeFormatters[langCode] = timeFormatter
	timeFormatters["default"] = timeFormatter
	return &Namespace{
		timeFormatters: timeFormatters,
		location:       location,
	}
}

// Namespace provides template functions for the "time" namespace.
type Namespace struct {
	timeFormatters map[string]htime.TimeFormatter
	location       *time.Location
}

// AsTime converts the textual representation of the datetime string into
// a time.Time interface.
func (ns *Namespace) AsTime(v any, args ...any) (any, error) {
	loc := ns.location
	if len(args) > 0 {
		locStr, err := cast.ToStringE(args[0])
		if err != nil {
			return nil, err
		}
		loc, err = _time.LoadLocation(locStr)
		if err != nil {
			return nil, err
		}
	}

	return htime.ToTimeInDefaultLocationE(v, loc)

}

// Format converts the textual representation of the datetime string in v into
// time.Time if needed and formats it with the given layout.
func (ns *Namespace) Format(layout string, args ...any) (string, error) {
	var v any
	var locale any

	if len(args) == 0 {
		return "", errors.New("missing date/time argument")
	}
	v = args[0]
	if len(args) == 2 {
		locale = args[1]
	}
	if len(args) > 2 {
		return "", errors.New("missing date/time argument")
	}

	t, err := htime.ToTimeInDefaultLocationE(v, ns.location)
	if err != nil {
		return "", err
	}

	localeStr := ""
	switch val := locale.(type) {
	case string:
		localeStr = val
	case *string:
		localeStr = *val
	case nil:
		localeStr = "default"
	default:
		return "", errors.New("locale must be a string or nil")
	}

	formatter, ok := ns.timeFormatters[localeStr]
	if ok {
		return formatter.Format(t, layout), nil
	}

	translator := translators.GetTranslator(localeStr)
	if translator != nil {
		formatter = htime.NewTimeFormatter(translator)
		ns.timeFormatters[localeStr] = formatter
		return formatter.Format(t, layout), nil
	}

	return "", errors.New("no time formatter for language '" + localeStr + "'")
}

// Now returns the current local time or `clock` time
func (ns *Namespace) Now() _time.Time {
	return htime.Now()
}

// ParseDuration parses the duration string s.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// See https://golang.org/pkg/time/#ParseDuration
func (ns *Namespace) ParseDuration(s any) (_time.Duration, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return 0, err
	}

	return _time.ParseDuration(ss)
}

var durationUnits = map[string]_time.Duration{
	"nanosecond":  _time.Nanosecond,
	"ns":          _time.Nanosecond,
	"microsecond": _time.Microsecond,
	"us":          _time.Microsecond,
	"µs":          _time.Microsecond,
	"millisecond": _time.Millisecond,
	"ms":          _time.Millisecond,
	"second":      _time.Second,
	"s":           _time.Second,
	"minute":      _time.Minute,
	"m":           _time.Minute,
	"hour":        _time.Hour,
	"h":           _time.Hour,
}

// Duration converts the given number to a time.Duration.
// Unit is one of nanosecond/ns, microsecond/us/µs, millisecond/ms, second/s, minute/m or hour/h.
func (ns *Namespace) Duration(unit any, number any) (_time.Duration, error) {
	unitStr, err := cast.ToStringE(unit)
	if err != nil {
		return 0, err
	}
	unitDuration, found := durationUnits[unitStr]
	if !found {
		return 0, fmt.Errorf("%q is not a valid duration unit", unit)
	}
	n, err := cast.ToInt64E(number)
	if err != nil {
		return 0, err
	}
	return _time.Duration(n) * unitDuration, nil
}
