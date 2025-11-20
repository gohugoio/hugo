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

package openapi3_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gohugoio/hugo/hugolib"
)

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	files := `
-- assets/api/myapi.yaml --
openapi: 3.0.0
info:
  title: Sample API
  description: Optional multiline or single-line description in [CommonMark](http://commonmark.org/help/) or HTML.
  version: 0.1.9
servers:
  - url: http://api.example.com/v1
    description: Optional server description, e.g. Main (production) server
  - url: http://staging-api.example.com
    description: Optional server description, e.g. Internal staging server for testing
paths:
  /users:
    get:
      summary: Returns a list of users.
      description: Optional extended description in CommonMark or HTML.
      responses:
        '200':    # status code
          description: A JSON array of user names
          content:
            application/json:
              schema: 
                type: array
                items: 
                  type: string
-- hugo.toml --
baseURL = 'http://example.com/'
disableLiveReload = true
-- layouts/home.html --
{{ $api := resources.Get "api/myapi.yaml" | openapi3.Unmarshal }}
API: {{ $api.Info.Title | safeHTML }}
  `

	b := hugolib.TestRunning(t, files)

	b.AssertFileContent("public/index.html", `API: Sample API`)

	b.
		EditFileReplaceFunc("assets/api/myapi.yaml", func(s string) string { return strings.ReplaceAll(s, "Sample API", "Hugo API") }).
		Build()

	b.AssertFileContent("public/index.html", `API: Hugo API`)
}

// Test data borrowed/adapted from https://github.com/getkin/kin-openapi/tree/master/openapi3/testdata/refInLocalRef
// The templates below would be simpler if kin-openapi's T would serialize to JSON better, see https://github.com/getkin/kin-openapi/issues/561
const refLocalTemplate = `
-- hugo.toml --
baseURL = 'http://example.com/'
disableKinds = ["page", "section", "taxonomy", "term", "sitemap", "robotsTXT", "404"]
disableLiveReload = true
[HTTPCache]
[HTTPCache.cache]
[HTTPCache.cache.for]
includes = ['**']
[[HTTPCache.polls]]
high = '100ms'
low = '50ms'
[HTTPCache.polls.for]
includes = ['**']
-- assets/api/myapi.json --
{
  "openapi": "3.0.3",
  "info": {
    "title": "Reference in reference example",
    "version": "1.0.0"
  },
  "paths": {
    "/api/test/ref/in/ref": {
      "post": {
        "requestBody": {
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties" : {
                  "data": {
                    "$ref": "#/components/schemas/Request"
                  }
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "messages/response.json"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Request": {
        "$ref": "messages/request.json"
      }
    }
  }
}
-- assets/api/messages/request.json --
{
  "type": "object",
  "required": [
    "definition_reference"
  ],
  "properties": {
    "definition_reference": {
      "$ref": "./data.json"
    }
  }
}
-- assets/api/messages/response.json --
{
  "type": "object",
  "properties": {
    "id": {
      "type": "integer",
      "format": "uint64"
    }
  }
}
-- assets/api/messages/data.json --
{
  "type": "object",
  "properties": {
    "id": {
      "type": "integer",
      "format": "int32"
    },
    "ref_prop_part": {
      "$ref": "DATAPART_REF"
    }
  }
}
-- assets/api/messages/dataPart.json --
{
  "type": "object",
  "properties": {
    "idPart": {
      "type": "integer",
      "format": "int64"
    }
  }
}
-- layouts/home.html --
{{ $getremote := dict
  "method" "post"
  "key" time.Now.UnixNano
}}
{{ $opts := dict
  "getremote" $getremote
}}
{{ with resources.Get "api/myapi.json" }}
  {{ with openapi3.Unmarshal . $opts }}
    Title: {{ .Info.Title | safeHTML }}
    {{ range $k, $v := .Paths.Map }}
      {{ $post := index $v.Post }}
      {{ with index $post.RequestBody.Value.Content  }}
        {{ $mt := (index .   "application/json")}}
        {{ $data := index $mt.Schema.Value.Properties "data" }}
        RequestBody: {{ template "printschema" $mt.Schema.Value  }}$
      {{ end }}
      {{ $response := index (index $post.Responses.Map "200").Value.Content "application/json" }}
      Response: {{ template "printschema" $response.Schema.Value }}
    {{ end }}
  {{ end }}
{{ end }}
{{ define "printschema" }}{{ with .Format}}Format: {{ . }}|{{ end }}{{ with .Properties }} Properties: {{ range $k, $v := . }}{{ $k}}: {{ template "printitems" . -}}{{ end }}{{ end -}}${{ end }}
{{ define "printitems" }}{{ if eq (printf "%T" .) "*openapi3.SchemaRef" }}{{ template "printschema" .Value  }}{{ else }}{{ template "printitem" . -}}{{ end -}}{{ end }}
{{ define "printitem" }}{{ printf "%T: %s" . (.|debug.Dump) | safeHTML }}{{ end }}

`

func TestUnmarshalRefLocal(t *testing.T) {
	files := strings.ReplaceAll(refLocalTemplate, "DATAPART_REF", "./dataPart.json")

	b := hugolib.Test(t, files)

	b.AssertFileContent("public/index.html",
		`Reference in reference example`,
		"RequestBody:  Properties: data:  Properties: definition_reference:  Properties: id: Format: int32|$ref_prop_part:  Properties: idPart: Format: int64|$",
		"Response:  Properties: id: Format: uint64|$",
	)
}

func TestUnmarshalRefLocalEdit(t *testing.T) {
	files := strings.ReplaceAll(refLocalTemplate, "DATAPART_REF", "./dataPart.json")

	b := hugolib.TestRunning(t, files)

	b.AssertFileContent("public/index.html",
		"RequestBody:  Properties: data:  Properties: definition_reference:  Properties: id: Format: int32|$ref_prop_part:  Properties: idPart: Format: int64|$",
	)

	b.EditFileReplaceAll("assets/api/messages/dataPart.json", "int64", "int8").Build()

	b.AssertFileContent("public/index.html",
		"RequestBody:  Properties: data:  Properties: definition_reference:  Properties: id: Format: int32|$ref_prop_part:  Properties: idPart: Format: int8|$",
	)
}

func TestUnmarshalRefRemote(t *testing.T) {
	createFiles := func(t *testing.T) string {
		counter := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" && r.URL.Path == "/api/messages/dataPart.json" {
				n := 8

				// This endpoint will also be called by the poller, so we distinguish
				// between the first and subsequent calls.
				if counter > 0 {
					n = 16
				}
				counter++

				w.Header().Set("Content-Type", "application/json")
				response := fmt.Sprintf(`{
  "type": "object",
  "properties": {
    "idPart": {
      "type": "integer",
      "format": "int%d"
    }
  }
}`, n)
				w.Write([]byte(response))
				return
			}

			http.Error(w, "Not found", http.StatusNotFound)
		}))

		t.Cleanup(func() {
			ts.Close()
		})

		dataPartRef := ts.URL + "/api/messages/dataPart.json"
		return strings.ReplaceAll(refLocalTemplate, "DATAPART_REF", dataPartRef)
	}

	t.Run("Build", func(t *testing.T) {
		b := hugolib.Test(t, createFiles(t))

		b.AssertFileContent("public/index.html",
			`Reference in reference example`,
			"RequestBody:  Properties: data:  Properties: definition_reference:  Properties: id: Format: int32|$ref_prop_part:  Properties: idPart: Format: int8|$",
			"Response:  Properties: id: Format: uint64|$",
		)
	})

	t.Run("Rebuild", func(t *testing.T) {
		b := hugolib.TestRunning(t, createFiles(t))

		b.AssertFileContent("public/index.html",
			"idPart: Format: int8|$",
		)

		// Rebuild triggered by remote polling.
		time.Sleep(800 * time.Millisecond)

		b.AssertFileContent("public/index.html",
			"idPart: Format: int16|$",
		)
	})
}
