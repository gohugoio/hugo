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

package hugolib

import (
	"testing"
)

func Test404(t *testing.T) {
	t.Parallel()

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithTemplatesAdded("404.html", "<html><body>Not Found!</body></html>")
	b.Build(BuildCfg{})

	// Note: We currently have only 1 404 page. One might think that we should have
	// multiple, to follow the Custom Output scheme, but I don't see how that would work
	// right now.
	b.AssertFileContent("public/404.html", "Not Found")

}
