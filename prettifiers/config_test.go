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

package prettifiers

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/yosssi/gohtml"

	qt "github.com/frankban/quicktest"
)

func TestConfig(t *testing.T) {
	c := qt.New(t)
	v := viper.New()

	v.Set("prettifyOutput", true)
	v.Set("prettify", map[string]interface{}{
		"disableHTML": true,
	})

	conf, err := decodeConfig(v)

	c.Assert(err, qt.IsNil)

	c.Assert(conf.PrettifyOutput, qt.Equals, true)

	// `enable` flags
	c.Assert(conf.DisableHTML, qt.Equals, true)
}

func TestConfigCondensedHTML(t *testing.T)   { testHTMLCondense(t, true) }
func TestConfigUncondensedHTML(t *testing.T) { testHTMLCondense(t, false) }

func testHTMLCondense(t *testing.T, condense bool) {
	c := qt.New(t)
	v := viper.New()

	v.Set("prettify", map[string]interface{}{
		"html": map[string]interface{}{
			"condense": condense,
		},
	})

	conf, err := decodeConfig(v)

	c.Assert(err, qt.IsNil)

	c.Assert(conf.HTML.Condense, qt.Equals, condense)
	c.Assert(gohtml.Condense, qt.Equals, condense)

}
