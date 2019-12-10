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

package template

import (
	"io"
	"reflect"

	"github.com/gohugoio/hugo/tpl/internal/go_templates/texttemplate/parse"
)

/*

This files contains the Hugo related addons. All the other files in this
package is auto generated.

*/

// TODO1 name
type Preparer interface {
	Prepare() (*Template, error)
}

type TemplateExecutor struct {
}

func (t *TemplateExecutor) Execute(p Preparer, wr io.Writer, data interface{}) error {
	tmpl, err := p.Prepare()
	if err != nil {
		return err
	}

	value, ok := data.(reflect.Value)
	if !ok {
		value = reflect.ValueOf(data)
	}

	state := &state{
		tmpl: tmpl,
		wr:   wr,
		vars: []variable{{"$", value}},
	}

	return tmpl.executeWithState(state, value)

}

// Prepare returns a template ready for execution.
func (t *Template) Prepare() (*Template, error) {

	return t, nil
}

// Below are modifed structs etc.

// state represents the state of an execution. It's not part of the
// template so that multiple executions of the same template
// can execute in parallel.
type state struct {
	tmpl  *Template
	wr    io.Writer
	node  parse.Node // current node, for errors
	vars  []variable // push-down stack of variable values.
	depth int        // the height of the stack of executing templates.
}

func (t *Template) executeWithState(state *state, value reflect.Value) (err error) {
	defer errRecover(&err)
	if t.Tree == nil || t.Root == nil {
		state.errorf("%q is an incomplete or empty template", t.Name())
	}
	state.walk(value, t.Root)
	return
}
