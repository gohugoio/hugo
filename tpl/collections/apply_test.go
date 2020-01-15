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

package collections

import (
	"io"
	"reflect"
	"testing"

	"fmt"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/output"
	"github.com/gohugoio/hugo/tpl"
)

type templateFinder int

func (templateFinder) Lookup(name string) (tpl.Template, bool) {
	return nil, false
}

func (templateFinder) HasTemplate(name string) bool {
	return false
}

func (templateFinder) LookupVariant(name string, variants tpl.TemplateVariants) (tpl.Template, bool, bool) {
	return nil, false, false
}

func (templateFinder) LookupLayout(d output.LayoutDescriptor, f output.Format) (tpl.Template, bool, error) {
	return nil, false, nil
}

func (templateFinder) Execute(t tpl.Template, wr io.Writer, data interface{}) error {
	return nil
}

func (templateFinder) GetFunc(name string) (reflect.Value, bool) {
	if name == "dobedobedo" {
		return reflect.Value{}, false
	}

	return reflect.ValueOf(fmt.Sprint), true

}

func TestApply(t *testing.T) {
	t.Parallel()
	c := qt.New(t)
	d := &deps.Deps{}
	d.SetTmpl(new(templateFinder))
	ns := New(d)

	strings := []interface{}{"a\n", "b\n"}

	result, err := ns.Apply(strings, "print", "a", "b", "c")
	c.Assert(err, qt.IsNil)
	c.Assert(result, qt.DeepEquals, []interface{}{"abc", "abc"})

	_, err = ns.Apply(strings, "apply", ".")
	c.Assert(err, qt.Not(qt.IsNil))

	var nilErr *error
	_, err = ns.Apply(nilErr, "chomp", ".")
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = ns.Apply(strings, "dobedobedo", ".")
	c.Assert(err, qt.Not(qt.IsNil))

	_, err = ns.Apply(strings, "foo.Chomp", "c\n")
	if err == nil {
		t.Errorf("apply with unknown func should fail")
	}

}
