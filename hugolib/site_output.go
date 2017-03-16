// Copyright 2017-present The Hugo Authors. All rights reserved.
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
	"path"
	"strings"

	"github.com/spf13/hugo/config"
	"github.com/spf13/hugo/output"
)

type siteOutputDefinitions []siteOutputDefinition

type siteOutputDefinition struct {
	// What Kinds of pages are excluded in this definition.
	// A blank strings means NONE.
	// Comma separated list (for now).
	ExcludedKinds string

	Outputs []output.Format
}

func (defs siteOutputDefinitions) ForKind(kind string) []output.Format {
	var result []output.Format

	for _, def := range defs {
		if def.ExcludedKinds == "" || !strings.Contains(def.ExcludedKinds, kind) {
			result = append(result, def.Outputs...)
		}
	}

	return result
}

func createSiteOutputDefinitions(cfg config.Provider) siteOutputDefinitions {

	var defs siteOutputDefinitions

	// All have HTML
	defs = append(defs, siteOutputDefinition{ExcludedKinds: "", Outputs: []output.Format{output.HTMLType}})

	// TODO(bep) output deprecate rssURI
	rssBase := cfg.GetString("rssURI")
	if rssBase == "" {
		rssBase = "index"
	}

	// RSS has now a well defined media type, so strip any suffix provided
	rssBase = strings.TrimSuffix(rssBase, path.Ext(rssBase))
	rssType := output.RSSType
	rssType.BaseName = rssBase

	// Some have RSS
	defs = append(defs, siteOutputDefinition{ExcludedKinds: "page", Outputs: []output.Format{rssType}})

	return defs
}
