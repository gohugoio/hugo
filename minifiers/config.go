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

package minifiers

import (
	"fmt"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/spf13/cast"

	"github.com/mitchellh/mapstructure"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

var defaultTdewolffConfig = TdewolffConfig{
	HTML: html.Minifier{
		KeepDocumentTags:    true,
		KeepSpecialComments: true,
		KeepEndTags:         true,
		KeepDefaultAttrVals: true,
		KeepWhitespace:      false,
	},
	CSS: css.Minifier{
		Precision: 0, // 0 means no trimming
		Version:   0, // 0 means the latest CSS version
	},
	JS: js.Minifier{
		Version: 2022,
	},
	JSON: json.Minifier{},
	SVG: svg.Minifier{
		KeepComments: false,
		Precision:    0, // 0 means no trimming
	},
	XML: xml.Minifier{
		KeepWhitespace: false,
	},
}

type TdewolffConfig struct {
	HTML html.Minifier
	CSS  css.Minifier
	JS   js.Minifier
	JSON json.Minifier
	SVG  svg.Minifier
	XML  xml.Minifier
}

type MinifyConfig struct {
	// Whether to minify the published output (the HTML written to /public).
	MinifyOutput bool

	DisableHTML bool
	DisableCSS  bool
	DisableJS   bool
	DisableJSON bool
	DisableSVG  bool
	DisableXML  bool

	Tdewolff TdewolffConfig
}

var defaultConfig = MinifyConfig{
	Tdewolff: defaultTdewolffConfig,
}

func DecodeConfig(v any) (conf MinifyConfig, err error) {
	conf = defaultConfig

	if v == nil {
		return
	}

	m := maps.ToStringMap(v)

	// Handle deprecations.
	if td, found := m["tdewolff"]; found {
		tdm := maps.ToStringMap(td)

		// Decimals was renamed to Precision in tdewolff/minify v2.7.0.
		// https://github.com/tdewolff/minify/commit/2fed4401348ce36bd6c20e77335a463e69d94386
		for _, key := range []string{"css", "svg"} {
			if v, found := tdm[key]; found {
				vm := maps.ToStringMap(v)
				ko := "decimals"
				kn := "precision"
				if vv, found := vm[ko]; found {
					hugo.Deprecate(
						fmt.Sprintf("site config key minify.tdewolff.%s.%s", key, ko),
						fmt.Sprintf("Use config key minify.tdewolff.%s.%s instead.", key, kn),
						"v0.150.0",
					)
					if _, found = vm[kn]; !found {
						vvi := cast.ToInt(vv)
						if vvi > 0 {
							vm[kn] = vvi
						}
					}
					delete(vm, ko)
				}
			}
		}

		// KeepConditionalComments was renamed to KeepSpecialComments in tdewolff/minify v2.20.13.
		// https://github.com/tdewolff/minify/commit/342cbc1974162db0ad3327f7a42a623b2cd3ebbc
		if v, found := tdm["html"]; found {
			vm := maps.ToStringMap(v)
			ko := "keepconditionalcomments"
			kn := "keepspecialcomments"
			if vv, found := vm[ko]; found {
				hugo.Deprecate(
					fmt.Sprintf("site config key minify.tdewolff.html.%s", ko),
					fmt.Sprintf("Use config key minify.tdewolff.html.%s instead.", kn),
					"v0.150.0",
				)
				if _, found := vm[kn]; !found {
					vm[kn] = cast.ToBool(vv)
				}
				delete(vm, ko)
			}
		}

		// KeepCSS2 was deprecated in favor of Version in tdewolff/minify v2.24.1.
		// https://github.com/tdewolff/minify/commit/57e3ebe0e6914b82c9ab0849a14f86bc29cd2ebf
		if v, found := tdm["css"]; found {
			vm := maps.ToStringMap(v)
			ko := "keepcss2"
			kn := "version"
			if vv, found := vm[ko]; found {
				hugo.Deprecate(
					fmt.Sprintf("site config key minify.tdewolff.css.%s", ko),
					fmt.Sprintf("Use config key minify.tdewolff.css.%s instead.", kn),
					"v0.150.0",
				)
				if _, found := vm[kn]; !found {
					if cast.ToBool(vv) {
						vm[kn] = 2
					}
				}
				delete(vm, ko)
			}
		}
	}

	err = mapstructure.WeakDecode(m, &conf)
	if err != nil {
		return
	}

	return
}
