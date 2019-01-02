// Copyright 2019 The Hugo Authors. All rights reserved.
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

// Package page contains the core interfaces and types for the Page resource,
// a core component in Hugo.
package page

import (
	"fmt"
)

// Data represents the .Data element in a Page in Hugo. We make this
// a type so we can do lazy loading of .Data.Pages
type Data map[string]interface{}

// Pages returns the pages stored with key "pages". If this is a func,
// it will be invoked.
func (d Data) Pages() Pages {
	v, found := d["pages"]
	if !found {
		return nil
	}

	switch vv := v.(type) {
	case Pages:
		return vv
	case func() Pages:
		return vv()
	default:
		panic(fmt.Sprintf("%T is not Pages", v))
	}
}
