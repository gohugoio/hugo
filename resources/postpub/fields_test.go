// Copyright 2020 The Hugo Authors. All rights reserved.
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

package postpub

import (
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/gohugoio/hugo/media"
)

func TestCreatePlaceholders(t *testing.T) {
	c := qt.New(t)

	m := structToMap(media.CSSType)

	insertFieldPlaceholders("foo", m, func(s string) string {
		return "pre_" + s + "_post"
	})

	c.Assert(m, qt.DeepEquals, map[string]interface{}{
		"IsZero":      "pre_foo.IsZero_post",
		"MarshalJSON": "pre_foo.MarshalJSON_post",
		"Suffixes":    "pre_foo.Suffixes_post",
		"Delimiter":   "pre_foo.Delimiter_post",
		"FirstSuffix": "pre_foo.FirstSuffix_post",
		"IsText":      "pre_foo.IsText_post",
		"String":      "pre_foo.String_post",
		"Type":        "pre_foo.Type_post",
		"MainType":    "pre_foo.MainType_post",
		"SubType":     "pre_foo.SubType_post",
	})
}
