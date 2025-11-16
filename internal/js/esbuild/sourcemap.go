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

package esbuild

import (
	"encoding/json"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/gohugoio/hugo/common/paths"
)

type sourceMap struct {
	Version        int      `json:"version"`
	Sources        []string `json:"sources"`
	SourcesContent []string `json:"sourcesContent"`
	Mappings       string   `json:"mappings"`
	Names          []string `json:"names"`
}

func fixOutputFile(o *api.OutputFile, resolve func(string) string) error {
	if strings.HasSuffix(o.Path, ".map") {
		b, err := fixSourceMap(o.Contents, resolve)
		if err != nil {
			return err
		}
		o.Contents = b
	}
	return nil
}

func fixSourceMap(s []byte, resolve func(string) string) ([]byte, error) {
	var sm sourceMap
	if err := json.Unmarshal([]byte(s), &sm); err != nil {
		return nil, err
	}

	sm.Sources = fixSourceMapSources(sm.Sources, resolve)

	b, err := json.Marshal(sm)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func fixSourceMapSources(s []string, resolve func(string) string) []string {
	var result []string
	for _, src := range s {
		if s := resolve(src); s != "" {
			// Absolute filenames works fine on U*ix (tested in Chrome on MacOs), but works very poorly on Windows (again Chrome).
			// So, convert it to a URL.
			if u, err := paths.UrlFromFilename(s); err == nil {
				result = append(result, u.String())
			}
		}
	}
	return result
}

// Used in tests.
func SourcesFromSourceMap(s string) []string {
	var sm sourceMap
	if err := json.Unmarshal([]byte(s), &sm); err != nil {
		return nil
	}
	return sm.Sources
}
