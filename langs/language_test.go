// Copyright 2018 The Hugo Authors. All rights reserved.
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

package langs

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/viper"
)

func TestGetGlobalOnlySetting(t *testing.T) {
	c := qt.New(t)
	v := viper.New()
	v.Set("defaultContentLanguageInSubdir", true)
	v.Set("contentDir", "content")
	v.Set("paginatePath", "page")
	lang := NewDefaultLanguage(v)
	lang.Set("defaultContentLanguageInSubdir", false)
	lang.Set("paginatePath", "side")

	c.Assert(lang.GetBool("defaultContentLanguageInSubdir"), qt.Equals, true)
	c.Assert(lang.GetString("paginatePath"), qt.Equals, "side")
}

func TestLanguageParams(t *testing.T) {
	c := qt.New(t)

	v := viper.New()
	v.Set("p1", "p1cfg")
	v.Set("contentDir", "content")

	lang := NewDefaultLanguage(v)
	lang.SetParam("p1", "p1p")

	c.Assert(lang.Params()["p1"], qt.Equals, "p1p")
	c.Assert(lang.Get("p1"), qt.Equals, "p1cfg")
}
