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

package hugolib

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
)

func Test073(t *testing.T) {
	assertDisabledTaxonomyAndTerm := func(b *sitesBuilder, taxonomy, term bool) {
		b.Assert(b.CheckExists("public/tags/index.html"), qt.Equals, taxonomy)
		b.Assert(b.CheckExists("public/tags/tag1/index.html"), qt.Equals, term)
	}

	assertOutputTaxonomyAndTerm := func(b *sitesBuilder, taxonomy, term bool) {
		b.Assert(b.CheckExists("public/tags/index.json"), qt.Equals, taxonomy)
		b.Assert(b.CheckExists("public/tags/tag1/index.json"), qt.Equals, term)
	}

	for _, this := range []struct {
		name   string
		config string
		assert func(err error, out string, b *sitesBuilder)
	}{
		{
			"Outputs for both taxonomy and taxonomyTerm",
			`[outputs]
 taxonomy = ["JSON"]
 taxonomyTerm = ["JSON"]

`,
			func(err error, out string, b *sitesBuilder) {
				b.Assert(err, qt.IsNil)
				assertOutputTaxonomyAndTerm(b, true, true)
			},
		},
		{
			"Outputs for taxonomyTerm",
			`[outputs]
taxonomyTerm = ["JSON"]

`,
			func(err error, out string, b *sitesBuilder) {
				b.Assert(err, qt.IsNil)
				assertOutputTaxonomyAndTerm(b, true, false)
			},
		},
		{
			"Outputs for taxonomy only",
			`[outputs]
taxonomy = ["JSON"]

`,
			func(err error, out string, b *sitesBuilder) {
				b.Assert(err, qt.Not(qt.IsNil))
				b.Assert(out, qt.Contains, `ignoreErrors = ["error-output-taxonomy"]`)
			},
		},
		{
			"Outputs for taxonomy only, ignore error",
			`
ignoreErrors = ["error-output-taxonomy"]
[outputs]
taxonomy = ["JSON"]

`,
			func(err error, out string, b *sitesBuilder) {
				b.Assert(err, qt.IsNil)
				assertOutputTaxonomyAndTerm(b, true, false)
			},
		},
		{
			"Disable both taxonomy and taxonomyTerm",
			`disableKinds = ["taxonomy", "taxonomyTerm"]`,
			func(err error, out string, b *sitesBuilder) {
				b.Assert(err, qt.IsNil)
				assertDisabledTaxonomyAndTerm(b, false, false)
			},
		},
		{
			"Disable only taxonomyTerm",
			`disableKinds = ["taxonomyTerm"]`,
			func(err error, out string, b *sitesBuilder) {
				b.Assert(err, qt.IsNil)
				assertDisabledTaxonomyAndTerm(b, false, true)
			},
		},
		{
			"Disable only taxonomy",
			`disableKinds = ["taxonomy"]`,
			func(err error, out string, b *sitesBuilder) {
				b.Assert(err, qt.Not(qt.IsNil))
				b.Assert(out, qt.Contains, `ignoreErrors = ["error-disable-taxonomy"]`)
			},
		},
		{
			"Disable only taxonomy, ignore error",
			`disableKinds = ["taxonomy"]
			ignoreErrors = ["error-disable-taxonomy"]`,
			func(err error, out string, b *sitesBuilder) {
				b.Assert(err, qt.IsNil)
				assertDisabledTaxonomyAndTerm(b, false, true)
			},
		},
	} {
		t.Run(this.name, func(t *testing.T) {
			b := newTestSitesBuilder(t).WithConfigFile("toml", this.config)
			b.WithTemplatesAdded("_default/list.json", "JSON")
			out, err := captureStdout(func() error {
				return b.BuildE(BuildCfg{})
			})
			fmt.Println(out)
			this.assert(err, out, b)
		})
	}
}
