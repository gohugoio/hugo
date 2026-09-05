// Copyright 2026 The Hugo Authors. All rights reserved.
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

package metadecoders

import (
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

// See issue 14596.
func TestMarshalYAMLRepeatedEmptySlices(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	in := map[string]any{
		"alpha": map[string]any{"tags": []any{}},
		"beta":  map[string]any{"tags": []any{}},
		"gamma": map[string]any{"tags": []any{}},
	}

	out, err := MarshalYAML(in)
	c.Assert(err, qt.IsNil)
	got := string(out)
	c.Assert(strings.Contains(got, "&"), qt.IsFalse)
	c.Assert(strings.Contains(got, "*"), qt.IsFalse)

	var decoded any
	c.Assert(UnmarshalYaml(out, &decoded), qt.IsNil)
}

// See issue 14596.
func TestMarshalYAMLRepeatedEmptyStringSlices(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	in := map[string]any{
		"a": []string{},
		"b": []string{},
		"c": []string{"item"},
	}

	out, err := MarshalYAML(in)
	c.Assert(err, qt.IsNil)

	var decoded any
	c.Assert(UnmarshalYaml(out, &decoded), qt.IsNil)
}

func TestMarshalYAMLKeepsAliasesCompact(t *testing.T) {
	t.Parallel()
	c := qt.New(t)

	in := []byte(`
a: &a [_, _, _, _, _, _, _, _, _, _, _, _, _, _, _]
b: &b [*a, *a, *a, *a, *a, *a, *a, *a, *a, *a]
c: &c [*b, *b, *b, *b, *b, *b, *b, *b, *b, *b]
`)
	m, err := Default.Unmarshal(in, YAML)
	c.Assert(err, qt.IsNil)

	out, err := MarshalYAML(m)
	c.Assert(err, qt.IsNil)
	c.Assert(len(out) < 512, qt.IsTrue)
	c.Assert(strings.Contains(string(out), "*"), qt.IsTrue)

	var decoded any
	c.Assert(UnmarshalYaml(out, &decoded), qt.IsNil)
}
