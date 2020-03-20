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

package cast

import (
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/docshelper"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/tpl/internal"
	"github.com/spf13/viper"
)

// This file provides documentation support and is randomly put into this package.
func init() {
	docsProvider := func() docshelper.DocProvider {
		d := &deps.Deps{
			Cfg:                 viper.New(),
			Log:                 loggers.NewErrorLogger(),
			BuildStartListeners: &deps.Listeners{},
			Site:                page.NewDummyHugoSite(newTestConfig()),
		}

		var namespaces internal.TemplateFuncsNamespaces

		for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
			nf := nsf(d)
			namespaces = append(namespaces, nf)

		}

		return docshelper.DocProvider{"tpl": map[string]interface{}{"funcs": namespaces}}

	}

	docshelper.AddDocProviderFunc(docsProvider)
}

func newTestConfig() *viper.Viper {
	v := viper.New()
	v.Set("contentDir", "content")
	return v
}
