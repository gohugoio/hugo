// Copyright 2025 The Hugo Authors. All rights reserved.
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

package metadecoders_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestYAMLIntegerSortIssue14078(t *testing.T) {
	files := `
-- assets/mydata.yaml --
a:
   weight: 1
x:
  weight: 2
c:
  weight: 3
t:
  weight: 4

-- layouts/all.html --
{{ $mydata := resources.Get "mydata.yaml" | transform.Unmarshal }}
Sorted: {{ sort $mydata "weight" }}|

`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "Sorted: [map[weight:1] map[weight:2] map[weight:3] map[weight:4]]|")
}
