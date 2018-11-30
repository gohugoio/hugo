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
	_time "time"

	"github.com/spf13/cast"
)

// New returns a new instance of the time-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "time" namespace.
type Namespace struct{}

// AsTime converts the textual representation of the datetime string into
// a time.Time interface.
func (ns *Namespace) AsTime(v interface{}) (interface{}, error) {
	t, err := cast.ToTimeE(v)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Format converts the textual representation of the datetime string into
// the other form or returns it of the time.Time value. These are formatted
// with the layout string
func (ns *Namespace) Format(layout string, v interface{}) (string, error) {
	t, err := cast.ToTimeE(v)
	if err != nil {
		return "", err
	}

	return t.Format(layout), nil
}

// Now returns the current local time.
func (ns *Namespace) Now() _time.Time {
	return _time.Now()
}

// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// See https://golang.org/pkg/time/#ParseDuration
func (ns *Namespace) ParseDuration(in interface{}) (_time.Duration, error) {
	s, err := cast.ToStringE(in)
	if err != nil {
		return 0, err
	}

	return _time.ParseDuration(s)
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
func (ns *Namespace) Duration(unit interface{}, number interface{}) (_time.Duration, error) {
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
