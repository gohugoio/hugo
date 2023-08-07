// Copyright 2017-present The Hugo Authors. All rights reserved.
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

// Package docshelper provides some helpers for the Hugo documentation, and
// is of limited interest for the general Hugo user.
package docshelper

import "fmt"

type (
	DocProviderFunc = func() DocProvider
	DocProvider     map[string]any
)

var docProviderFuncs []DocProviderFunc

func AddDocProviderFunc(fn DocProviderFunc) {
	docProviderFuncs = append(docProviderFuncs, fn)
}

func GetDocProvider() DocProvider {
	provider := make(DocProvider)

	for _, fn := range docProviderFuncs {
		p := fn()
		for k, v := range p {
			if _, found := provider[k]; found {
				// We use to merge config, but not anymore.
				// These constructs will eventually go away, so just make it simple.
				panic(fmt.Sprintf("Duplicate doc provider key: %q", k))
			}
			provider[k] = v
		}
	}

	return provider
}
