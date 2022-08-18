// Copyright 2021 The Hugo Authors. All rights reserved.
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

package security

import (
	"github.com/gohugoio/hugo/docshelper"
)

func init() {
	docsProvider := func() docshelper.DocProvider {

		return docshelper.DocProvider{"config": DefaultConfig.ToSecurityMap()}
	}
	docshelper.AddDocProviderFunc(docsProvider)
}
