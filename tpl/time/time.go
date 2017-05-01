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
