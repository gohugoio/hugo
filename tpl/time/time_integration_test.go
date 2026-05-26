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

package time_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// See issue 14948.
func TestTimeFormatMonthAbbreviationsEnGB(t *testing.T) {
	t.Parallel()

	files := `
-- hugo.toml --
locale = 'en-GB'
disableKinds = ['page','rss','section','sitemap','taxonomy','term']
-- layouts/home.html --
{{ $dates := slice
	"2025-01-17T10:42:00-08:00"
	"2025-02-17T10:42:00-08:00"
	"2025-03-17T10:42:00-07:00"
	"2025-04-17T10:42:00-07:00"
	"2025-05-17T10:42:00-07:00"
	"2025-06-17T10:42:00-07:00"
	"2025-07-17T10:42:00-07:00"
	"2025-08-17T10:42:00-07:00"
	"2025-09-17T10:42:00-07:00"
	"2025-10-17T10:42:00-07:00"
	"2025-11-17T10:42:00-08:00"
	"2025-12-17T10:42:00-08:00"
}}
{{- range $dates }}
	{{- time.Format "Jan" . }}|
{{- end }}
`

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html", "Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sept|Oct|Nov|Dec|")
}
