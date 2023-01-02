// Copyright 2022 The Hugo Authors. All rights reserved.
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

package sass

import (
	"fmt"
	"sort"
	"strings"
)

const (
	HugoVarsNamespace = "hugo:vars"
)

func CreateVarsStyleSheet(vars map[string]string) string {
	if vars == nil {
		return ""
	}
	var varsStylesheet string

	var varsSlice []string
	for k, v := range vars {
		var prefix string
		if !strings.HasPrefix(k, "$") {
			prefix = "$"
		}
		// These variables can be a combination of Sass identifiers (e.g. sans-serif), which
		// should not be quoted, and URLs et, which should be quoted.
		// unquote() is knowing what to do with each.
		varsSlice = append(varsSlice, fmt.Sprintf("%s%s: unquote(%q);", prefix, k, v))
	}
	sort.Strings(varsSlice)
	varsStylesheet = strings.Join(varsSlice, "\n")

	return varsStylesheet

}
