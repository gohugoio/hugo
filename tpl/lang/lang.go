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

package lang

import (
	"github.com/spf13/cast"
	"github.com/spf13/hugo/deps"
)

// New returns a new instance of the lang-namespaced template functions.
func New(deps *deps.Deps) *Namespace {
	return &Namespace{
		deps: deps,
	}
}

// Namespace provides template functions for the "lang" namespace.
type Namespace struct {
	deps *deps.Deps
}

// Namespace returns a pointer to the current namespace instance.
// TODO(bep) namespace remove this and other unused when done.
func (ns *Namespace) Namespace() *Namespace { return ns }

// Translate ...
func (ns *Namespace) Translate(id interface{}, args ...interface{}) (string, error) {
	sid, err := cast.ToStringE(id)
	if err != nil {
		return "", nil
	}

	return ns.deps.Translate(sid, args...), nil
}
