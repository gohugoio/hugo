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
	"encoding/json"
	"testing"

	"path/filepath"

	"github.com/gohugoio/hugo/deps"
)

// Issue #1123
// Testing prevention of cyclic refs in JSON encoding
// May be smart to run with: -timeout 4000ms
func TestEncodePage(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()

	// borrowed from menu_test.go
	for _, src := range menuPageSources {
		writeSource(t, fs, filepath.Join("content", src.Name), string(src.Content))

	}

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	_, err := json.Marshal(s)
	check(t, err)

	_, err = json.Marshal(s.RegularPages[0])
	check(t, err)
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Failed %s", err)
	}
}
