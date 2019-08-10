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

package crypto

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/htesting/hqt"
	"github.com/gohugoio/hugo/tpl/internal"
)

func TestInit(t *testing.T) {
	c := qt.New(t)
	var found bool
	var ns *internal.TemplateFuncsNamespace

	for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
		ns = nsf(&deps.Deps{})
		if ns.Name == name {
			found = true
			break
		}
	}

	c.Assert(found, qt.Equals, true)
	c.Assert(ns.Context(), hqt.IsSameType, &Namespace{})
}
