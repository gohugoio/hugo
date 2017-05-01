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

package cast

import (
	_cast "github.com/spf13/cast"
)

// New returns a new instance of the cast-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "cast" namespace.
type Namespace struct {
}

func (ns *Namespace) ToInt(v interface{}) (int, error) {
	return _cast.ToIntE(v)
}

func (ns *Namespace) ToString(v interface{}) (string, error) {
	return _cast.ToStringE(v)
}
