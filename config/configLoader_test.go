// Copyright 2019 The Hugo Authors. All rights reserved.
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

package config

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestIsValidConfigFileName(t *testing.T) {
	c := qt.New(t)

	for _, ext := range ValidConfigFileExtensions {
		filename := "config." + ext
		c.Assert(IsValidConfigFilename(filename), qt.Equals, true)
		c.Assert(IsValidConfigFilename(strings.ToUpper(filename)), qt.Equals, true)
	}

	c.Assert(IsValidConfigFilename(""), qt.Equals, false)
	c.Assert(IsValidConfigFilename("config.toml.swp"), qt.Equals, false)
}
