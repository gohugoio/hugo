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
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/docshelper"
	"github.com/gohugoio/hugo/parser"
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
		KeepDocumentTags:        true,
		KeepConditionalComments: true,
		KeepEndTags:             true,
		KeepDefaultAttrVals:     true,
		KeepWhitespace:          false,
	},
	CSS: css.Minifier{
		Precision: 0,
		KeepCSS2:  true,
	},
	JS: js.Minifier{
		Version: 2022,
	},
	JSON: json.Minifier{},
	SVG: svg.Minifier{
		KeepComments: false,
		Precision:    0,
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

	// Handle upstream renames.
	if td, found := m["tdewolff"]; found {
		tdm := maps.ToStringMap(td)
		for _, key := range []string{"css", "svg"} {
			if v, found := tdm[key]; found {
				vm := maps.ToStringMap(v)
				if vv, found := vm["decimal"]; found {
					vvi := cast.ToInt(vv)
					if vvi > 0 {
						vm["precision"] = vvi
					}
				}
			}
		}
	}

	err = mapstructure.WeakDecode(m, &conf)

	if err != nil {
		return
	}

	return
}

func init() {
	docsProvider := func() docshelper.DocProvider {
		return docshelper.DocProvider{"config": map[string]any{"minify": parser.LowerCaseCamelJSONMarshaller{Value: defaultConfig}}}
	}
	docshelper.AddDocProviderFunc(docsProvider)
}
