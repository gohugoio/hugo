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

package js

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

// This test is added to test/warn against breaking the "stability" of the
// cache key. It's sometimes needed to break this, but should be avoided if possible.
func TestOptionKey(t *testing.T) {
	c := qt.New(t)

	opts := internalOptions{
		TargetPath: "foo",
	}

	key := (&buildTransformation{options: opts}).Key()

	c.Assert(key.Value(), qt.Equals, "jsbuild_9405671309963492201")
}

func TestToInternalOptions(t *testing.T) {
	c := qt.New(t)

	o := Options{
		TargetPath:  "v1",
		Target:      "v2",
		JSXFactory:  "v3",
		JSXFragment: "v4",
		Externals:   []string{"react"},
		Defines:     map[string]interface{}{"process.env.NODE_ENV": "production"},
		Minify:      true,
	}

	c.Assert(toInternalOptions(o), qt.DeepEquals, internalOptions{
		TargetPath:  "v1",
		Minify:      true,
		Target:      "v2",
		JSXFactory:  "v3",
		JSXFragment: "v4",
		Externals:   []string{"react"},
		Defines:     map[string]string{"process.env.NODE_ENV": "production"},
		TSConfig:    "",
	})

	c.Assert(toInternalOptions(Options{}), qt.DeepEquals, internalOptions{
		TargetPath:  "",
		Minify:      false,
		Target:      "esnext",
		JSXFactory:  "",
		JSXFragment: "",
		Externals:   nil,
		Defines:     nil,
		TSConfig:    "",
	})
}
