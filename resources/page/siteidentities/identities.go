// Copyright 2024 The Hugo Authors. All rights reserved.
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

package siteidentities

import (
	"github.com/gohugoio/hugo/identity"
)

const (
	// Identifies site.Data.
	// The change detection in /data is currently very coarse grained.
	Data = identity.StringIdentity("site.Data")
)

// FromString returns the identity from the given string,
// or identity.Anonymous if not found.
func FromString(name string) (identity.Identity, bool) {
	switch name {
	case "Data":
		return Data, true
	}
	return identity.Anonymous, false
}
