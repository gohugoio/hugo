// Copyright 2023 The Hugo Authors. All rights reserved.
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

package create_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

func TestGetRemoteHead(t *testing.T) {

	files := `
-- config.toml --
[security]
  [security.http]
    methods = ['(?i)GET|POST|HEAD']
    urls = ['.*gohugo\.io.*']

-- layouts/index.html --
{{ $url := "https://gohugo.io/img/hugo.png" }}
{{ $opts := dict "method" "head" }}
{{ with resources.GetRemote $url $opts }}
  {{ with .Err }}
    {{ errorf "Unable to get remote resource: %s" . }}
  {{ else }}
    Head Content: {{ .Content }}. Head Data: {{ .Data }}
  {{ end }}
{{ else }}
  {{ errorf "Unable to get remote resource: %s" $url }}
{{ end }}
`

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	)

	b.Build()

	b.AssertFileContent("public/index.html",
		"Head Content: .",
		"Head Data: map[ContentLength:18210 ContentType:image/png Status:200 OK StatusCode:200 TransferEncoding:[]]",
	)

}

func TestGetRemoteRetry(t *testing.T) {
	t.Parallel()

	temporaryHTTPCodes := []int{408, 429, 500, 502, 503, 504}
	numPages := 30

	handler := func(w http.ResponseWriter, r *http.Request) {
		if rand.Intn(4) == 0 {
			w.WriteHeader(temporaryHTTPCodes[rand.Intn(len(temporaryHTTPCodes))])
			return
		}
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Response for " + r.URL.Path + "."))
	}

	srv := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(func() { srv.Close() })

	files := `
-- hugo.toml --
disableKinds = ["home", "taxonomy", "term"]
[security]
[security.http]
urls = ['.*']
mediaTypes = ['text/plain']
-- layouts/_default/single.html --
{{ $url := printf "%s%s" "URL" .RelPermalink}}
{{ $opts := dict }}
{{ with resources.GetRemote $url $opts }}
  {{ with .Err }}
    {{ errorf "Unable to get remote resource: %s" . }}
  {{ else }}
    Content: {{ .Content }}
  {{ end }}
{{ else }}
  {{ errorf "Unable to get remote resource: %s" $url }}
{{ end }}
`

	for i := 0; i < numPages; i++ {
		files += fmt.Sprintf("-- content/post/p%d.md --\n", i)
	}

	files = strings.ReplaceAll(files, "URL", srv.URL)

	b := hugolib.NewIntegrationTestBuilder(
		hugolib.IntegrationTestConfig{
			T:           t,
			TxtarString: files,
		},
	)

	b.Build()

	for i := 0; i < numPages; i++ {
		b.AssertFileContent(fmt.Sprintf("public/post/p%d/index.html", i), fmt.Sprintf("Content: Response for /post/p%d/.", i))
	}

}
