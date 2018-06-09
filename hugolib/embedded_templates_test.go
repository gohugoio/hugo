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

package hugolib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Just some simple test of the embedded templates to avoid
// https://github.com/gohugoio/hugo/issues/4757 and similar.
func TestEmbeddedTemplates(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	assert.True(true)

	home := []string{"index.html", `
GA:
{{ template "_internal/google_analytics.html" . }}

GA async:

{{ template "_internal/google_analytics_async.html" . }}

Disqus:

{{ template "_internal/disqus.html" . }}

`}

	b := newTestSitesBuilder(t)
	b.WithSimpleConfigFile().WithTemplatesAdded(home...)

	b.Build(BuildCfg{})

	// Gheck GA regular and async
	b.AssertFileContent("public/index.html",
		"'anonymizeIp', true",
		"'script','https://www.google-analytics.com/analytics.js','ga');\n\tga('create', 'ga_id', 'auto')",
		"<script async src='https://www.google-analytics.com/analytics.js'>")

	// Disqus
	b.AssertFileContent("public/index.html", "\"disqus_shortname\" + '.disqus.com/embed.js';")
}
