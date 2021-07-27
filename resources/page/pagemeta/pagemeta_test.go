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

package pagemeta

import (
	"fmt"
	"testing"

	"github.com/gohugoio/hugo/htesting/hqt"

	"github.com/gohugoio/hugo/config"

	qt "github.com/frankban/quicktest"
)

func TestDecodeBuildConfig(t *testing.T) {
	t.Parallel()

	c := qt.New(t)

	configTempl := `
[_build]
render = %s
list = %s
publishResources = true`

	for _, test := range []struct {
		args   []interface{}
		expect BuildConfig
	}{
		{
			[]interface{}{"true", "true"},
			BuildConfig{
				Render:           Always,
				List:             Always,
				PublishResources: true,
				set:              true,
			},
		},
		{[]interface{}{"true", "false"}, BuildConfig{
			Render:           Always,
			List:             Never,
			PublishResources: true,
			set:              true,
		}},
		{[]interface{}{`"always"`, `"always"`}, BuildConfig{
			Render:           Always,
			List:             Always,
			PublishResources: true,
			set:              true,
		}},
		{[]interface{}{`"never"`, `"never"`}, BuildConfig{
			Render:           Never,
			List:             Never,
			PublishResources: true,
			set:              true,
		}},
		{[]interface{}{`"link"`, `"local"`}, BuildConfig{
			Render:           Link,
			List:             ListLocally,
			PublishResources: true,
			set:              true,
		}},
		{[]interface{}{`"always"`, `"asdfadf"`}, BuildConfig{
			Render:           Always,
			List:             Always,
			PublishResources: true,
			set:              true,
		}},
	} {
		cfg, err := config.FromConfigString(fmt.Sprintf(configTempl, test.args...), "toml")
		c.Assert(err, qt.IsNil)
		bcfg, err := DecodeBuildConfig(cfg.Get("_build"))
		c.Assert(err, qt.IsNil)

		eq := qt.CmpEquals(hqt.DeepAllowUnexported(BuildConfig{}))

		c.Assert(bcfg, eq, test.expect)

	}
}
