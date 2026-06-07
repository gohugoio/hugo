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

package jsconfig

import (
	"encoding/json"
	"path/filepath"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestJsConfigBuilder(t *testing.T) {
	c := qt.New(t)

	b := NewBuilder()
	b.AddSourceRoot("/c/assets")
	b.AddSourceRoot("/d/assets")

	conf := b.Build("/a/b")
	c.Assert(conf.CompilerOptions.Paths["*"], qt.DeepEquals, []string{filepath.FromSlash("../../c/assets/*"), filepath.FromSlash("../../d/assets/*")})

	// baseUrl is deprecated in TypeScript and not needed when paths is set; see issue 14991.
	data, err := json.Marshal(conf)
	c.Assert(err, qt.IsNil)
	c.Assert(string(data), qt.Not(qt.Contains), "baseUrl")

	c.Assert(NewBuilder().Build("/a/b"), qt.IsNil)
}
