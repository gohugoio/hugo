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
	"fmt"
	"time"

	"github.com/gohugoio/hugo/cache/dynacache"
	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/deps"

	"github.com/spf13/cast"
)

// New returns a new instance of the time-namespaced template functions.
func New(timeFormatter htime.TimeFormatter, location *time.Location, deps *deps.Deps) *Namespace {
	if deps.MemCache == nil {
		panic("must provide MemCache")
	}

	return &Namespace{
		timeFormatter: timeFormatter,
		location:      location,
		deps:          deps,
		cacheIn: dynacache.GetOrCreatePartition[string, *time.Location](
			deps.MemCache,
			"/tmpl/time/in",
			dynacache.OptionsPartition{Weight: 30, ClearWhen: dynacache.ClearNever},
		),
	}
}

// Namespace provides template functions for the "time" namespace.
type Namespace struct {
	timeFormatter htime.TimeFormatter
	location      *time.Location
	deps          *deps.Deps
	cacheIn       *dynacache.Partition[string, *time.Location]
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
		loc, err = time.LoadLocation(locStr)
		if err != nil {
			return nil, err
		}
	}

	return htime.ToTimeInDefaultLocationE(v, loc)
}

// Format converts the textual representation of the datetime string in v into
// time.Time if needed and formats it with the given layout.
func (ns *Namespace) Format(layout string, v any) (string, error) {
	t, err := htime.ToTimeInDefaultLocationE(v, ns.location)
	if err != nil {
		return "", err
	}

	return ns.timeFormatter.Format(t, layout), nil
}

// Now returns the current local time or `clock` time
func (ns *Namespace) Now() time.Time {
	return htime.Now()
}

// In returns the time t in the IANA time zone specified by timeZoneName.
// If timeZoneName is "" or "UTC", the time is returned in UTC.
// If timeZoneName is "Local", the time is returned in the system's local time zone.
// Otherwise, timeZoneName must be a valid IANA location name (e.g., "Europe/Oslo").
func (ns *Namespace) In(timeZoneName string, t time.Time) (time.Time, error) {
	location, err := ns.cacheIn.GetOrCreate(dynacache.CleanKey(timeZoneName), func(string) (*time.Location, error) {
		return time.LoadLocation(timeZoneName)
	})
	if err != nil {
		return time.Time{}, err
	}

	return t.In(location), nil
}

// ParseDuration parses the duration string s.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// See https://golang.org/pkg/time/#ParseDuration
func (ns *Namespace) ParseDuration(s any) (time.Duration, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return 0, err
	}

	return time.ParseDuration(ss)
}

var durationUnits = map[string]time.Duration{
	"nanosecond":  time.Nanosecond,
	"ns":          time.Nanosecond,
	"microsecond": time.Microsecond,
	"us":          time.Microsecond,
	"µs":          time.Microsecond,
	"millisecond": time.Millisecond,
	"ms":          time.Millisecond,
	"second":      time.Second,
	"s":           time.Second,
	"minute":      time.Minute,
	"m":           time.Minute,
	"hour":        time.Hour,
	"h":           time.Hour,
}

// Duration converts the given number to a time.Duration.
// Unit is one of nanosecond/ns, microsecond/us/µs, millisecond/ms, second/s, minute/m or hour/h.
func (ns *Namespace) Duration(unit any, number any) (time.Duration, error) {
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
	return time.Duration(n) * unitDuration, nil
}
