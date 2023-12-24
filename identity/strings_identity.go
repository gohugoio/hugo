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

package identity

import (
	"strings"
)

// StringPrefixIdentity that is probably dependent on another string if
// if the other string is below this.
type StringPrefixIdentity string

func (s StringPrefixIdentity) IdentifierBase() string {
	return string(s)
}

func (s StringPrefixIdentity) IsProbablyDependent(id Identity) bool {
	return strings.HasPrefix(id.IdentifierBase(), string(s))
}
