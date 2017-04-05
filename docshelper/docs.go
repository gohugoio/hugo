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

import (
	"encoding/json"
)

var DocProviders = make(map[string]DocProvider)

func AddDocProvider(name string, provider DocProvider) {
	DocProviders[name] = provider
}

type DocProvider func() map[string]interface{}

func (d DocProvider) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(d(), "", "  ")
}
