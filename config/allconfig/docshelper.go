// Copyright 2024 The Hugo Authors. All rights reserved.
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

package allconfig

import (
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/docshelper"
)

// This is is just some helpers used to create some JSON used in the Hugo docs.
func init() {
	docsProvider := func() docshelper.DocProvider {
		cfg := config.New()
		for configRoot, v := range allDecoderSetups {
			if v.internalOrDeprecated {
				continue
			}
			cfg.Set(configRoot, make(maps.Params))
		}
		lang := maps.Params{
			"en": maps.Params{
				"menus":  maps.Params{},
				"params": maps.Params{},
			},
		}
		cfg.Set("languages", lang)
		cfg.SetDefaultMergeStrategy()

		configHelpers := map[string]any{
			"mergeStrategy": cfg.Get(""),
		}
		return docshelper.DocProvider{"config_helpers": configHelpers}
	}

	docshelper.AddDocProviderFunc(docsProvider)
}
