// Copyright 2017 The Hugo Authors. All rights reserved.
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

package collections

import (
	"fmt"
	"html/template"
	"testing"

	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/tpl/strings"
	"github.com/stretchr/testify/assert"
)

func TestApply(t *testing.T) {
	t.Parallel()

	hstrings := strings.New(&deps.Deps{})

	ns := New(&deps.Deps{})
	ns.Funcs(template.FuncMap{
		"apply":   ns.Apply,
		"chomp":   hstrings.Chomp,
		"strings": hstrings,
		"print":   fmt.Sprint,
	})

	strings := []interface{}{"a\n", "b\n"}
	noStringers := []interface{}{tstNoStringer{}, tstNoStringer{}}

	result, _ := ns.Apply(strings, "chomp", ".")
	assert.Equal(t, []interface{}{template.HTML("a"), template.HTML("b")}, result)

	result, _ = ns.Apply(strings, "chomp", "c\n")
	assert.Equal(t, []interface{}{template.HTML("c"), template.HTML("c")}, result)

	result, _ = ns.Apply(strings, "strings.Chomp", "c\n")
	assert.Equal(t, []interface{}{template.HTML("c"), template.HTML("c")}, result)

	result, _ = ns.Apply(strings, "print", "a", "b", "c")
	assert.Equal(t, []interface{}{"abc", "abc"}, result, "testing variadic")

	result, _ = ns.Apply(nil, "chomp", ".")
	assert.Equal(t, []interface{}{}, result)

	_, err := ns.Apply(strings, "apply", ".")
	if err == nil {
		t.Errorf("apply with apply should fail")
	}

	var nilErr *error
	_, err = ns.Apply(nilErr, "chomp", ".")
	if err == nil {
		t.Errorf("apply with nil in seq should fail")
	}

	_, err = ns.Apply(strings, "dobedobedo", ".")
	if err == nil {
		t.Errorf("apply with unknown func should fail")
	}

	_, err = ns.Apply(noStringers, "chomp", ".")
	if err == nil {
		t.Errorf("apply when func fails should fail")
	}

	_, err = ns.Apply(tstNoStringer{}, "chomp", ".")
	if err == nil {
		t.Errorf("apply with non-sequence should fail")
	}

	_, err = ns.Apply(strings, "foo.Chomp", "c\n")
	if err == nil {
		t.Errorf("apply with unknown namespace should fail")
	}

	_, err = ns.Apply(strings, "strings.Foo", "c\n")
	if err == nil {
		t.Errorf("apply with unknown namespace method should fail")
	}
}
