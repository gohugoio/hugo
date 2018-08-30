// Copyright 2015 The Hugo Authors. All rights reserved.
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

package hugolib

import (
	"path/filepath"
	"strings"
	"testing"
)

var simplePageYAML = `---
contenttype: ""
---
Sample Text
`

func TestDegenerateMissingFolderInPageFilename(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	p, err := s.NewPageFrom(strings.NewReader(simplePageYAML), filepath.Join("foobar"))
	if err != nil {
		t.Fatalf("Error in NewPageFrom")
	}
	if p.Section() != "" {
		t.Fatalf("No section should be set for a file path: foobar")
	}
}
