// Copyright 2019 The Hugo Authors. All rights reserved.
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

package prettifiers

import (
	"sort"
	"strings"

	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/docshelper"
	"github.com/gohugoio/hugo/parser"

	"github.com/mitchellh/mapstructure"
	"github.com/yosssi/gohtml"
)

type prettifyConfig struct {
	// Whether to prettify the published output (the HTML written to /public).
	PrettifyOutput bool

	DisableHTML bool

	HTML htmlConfig
}

type htmlConfig struct {
	Condense           bool
	InlineTags         []string
	InlineTagMaxLength int
}

var defaultConfig = prettifyConfig{
	HTML: htmlConfig{
		// Copy the defaults of gohtml
		Condense:           gohtml.Condense,
		InlineTags:         boolSetToSlice(gohtml.InlineTags),
		InlineTagMaxLength: gohtml.InlineTagMaxLength,
	},
}

func decodeConfig(cfg config.Provider) (conf prettifyConfig, err error) {
	conf = defaultConfig

	// May be set by CLI.
	conf.PrettifyOutput = cfg.GetBool("prettifyOutput")

	v := cfg.Get("prettify")
	if v == nil {
		return
	}

	m := maps.ToStringMap(v)

	err = mapstructure.WeakDecode(m, &conf)

	if err != nil {
		return
	}

	// Set some global properties for the HTML formatter
	gohtml.Condense = conf.HTML.Condense
	gohtml.InlineTags = sliceToBoolSet(conf.HTML.InlineTags)
	gohtml.InlineTagMaxLength = conf.HTML.InlineTagMaxLength

	return
}

// boolSetToSlice converts a map[string]bool to a sorted list of keys.
func boolSetToSlice(set map[string]bool) []string {
	slice := make([]string, 0, len(set))
	for tag, isShort := range set {
		if isShort {
			slice = append(slice, tag)
		}
	}
	sort.Strings(slice) // Ensure consistent ordering
	return slice
}

// sliceToBoolSet converts a list of strings to a map[string]bool mapping the items in the list to true.
func sliceToBoolSet(items []string) map[string]bool {
	set := make(map[string]bool)
	for _, tag := range items {
		set[strings.ToLower(tag)] = true
	}
	return set
}

func init() {
	docsProvider := func() docshelper.DocProvider {
		return docshelper.DocProvider{"config": map[string]interface{}{"prettify": parser.LowerCaseCamelJSONMarshaller{Value: defaultConfig}}}
	}
	docshelper.AddDocProviderFunc(docsProvider)
}
