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

package commands

import (
	"bytes"
	"io"
	"net/url"

	"github.com/gohugoio/hugo/transform"
	"github.com/gohugoio/hugo/transform/livereloadinject"
)

func injectLiveReloadScript(src io.Reader, baseURL url.URL) string {
	var b bytes.Buffer
	chain := transform.Chain{livereloadinject.New(baseURL)}
	chain.Apply(&b, src)

	return b.String()
}
