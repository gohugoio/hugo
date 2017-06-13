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
	"testing"

	"fmt"

	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl"
	"github.com/stretchr/testify/require"
)

type templateFinder int

func (templateFinder) Lookup(name string) *tpl.TemplateAdapter {
	return nil
}

func (templateFinder) GetFuncs() map[string]interface{} {
	return map[string]interface{}{
		"print": fmt.Sprint,
	}
}

func TestApply(t *testing.T) {
	t.Parallel()

	ns := New(&deps.Deps{Tmpl: new(templateFinder)})

	strings := []interface{}{"a\n", "b\n"}

	result, err := ns.Apply(strings, "print", "a", "b", "c")
	require.NoError(t, err)
	require.Equal(t, []interface{}{"abc", "abc"}, result, "testing variadic")

	_, err = ns.Apply(strings, "apply", ".")
	require.Error(t, err)

	var nilErr *error
	_, err = ns.Apply(nilErr, "chomp", ".")
	require.Error(t, err)

	_, err = ns.Apply(strings, "dobedobedo", ".")
	require.Error(t, err)

	_, err = ns.Apply(strings, "foo.Chomp", "c\n")
	if err == nil {
		t.Errorf("apply with unknown func should fail")
	}

}
