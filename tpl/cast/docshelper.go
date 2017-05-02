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
	"github.com/spf13/hugo/deps"
	"github.com/spf13/hugo/docshelper"
	"github.com/spf13/hugo/tpl/internal"
)

// This file provides documentation support and is randomly put into this package.
func init() {
	docsProvider := func() map[string]interface{} {
		docs := make(map[string]interface{})
		d := &deps.Deps{}

		var namespaces internal.TemplateFuncsNamespaces

		for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
			nf := nsf(d)
			namespaces = append(namespaces, nf)

		}

		docs["funcs"] = namespaces
		return docs
	}

	docshelper.AddDocProvider("tpl", docsProvider)
}
