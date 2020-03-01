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
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/docshelper"
	"github.com/gohugoio/hugo/parser"

	"github.com/mitchellh/mapstructure"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

var defaultTdewolffConfig = tdewolffConfig{
	HTML: html.Minifier{
		KeepDocumentTags:        true,
		KeepConditionalComments: true,
		KeepEndTags:             true,
		KeepDefaultAttrVals:     true,
		KeepWhitespace:          false,
		// KeepQuotes:              false, >= v2.6.2
	},
	CSS: css.Minifier{
		Decimals: -1, // will be deprecated
		// Precision: 0,  // use Precision with >= v2.7.0
		KeepCSS2: true,
	},
	JS:   js.Minifier{},
	JSON: json.Minifier{},
	SVG: svg.Minifier{
		Decimals: -1, // will be deprecated
		// Precision: 0,  // use Precision with >= v2.7.0
	},
	XML: xml.Minifier{
		KeepWhitespace: false,
	},
}

type tdewolffConfig struct {
	HTML html.Minifier
	CSS  css.Minifier
	JS   js.Minifier
	JSON json.Minifier
	SVG  svg.Minifier
	XML  xml.Minifier
}

type minifiersConfig struct {
	EnableHTML bool
	EnableCSS  bool
	EnableJS   bool
	EnableJSON bool
	EnableSVG  bool
	EnableXML  bool

	Tdewolff tdewolffConfig
}

var defaultConfig = minifiersConfig{
	EnableHTML: true,
	EnableCSS:  true,
	EnableJS:   true,
	EnableJSON: true,
	EnableSVG:  true,
	EnableXML:  true,

	Tdewolff: defaultTdewolffConfig,
}

func decodeConfig(cfg config.Provider) (conf minifiersConfig, err error) {
	conf = defaultConfig

	m := cfg.GetStringMap("minifiers")
	if m == nil {
		return
	}

	err = mapstructure.WeakDecode(m, &conf)

	if err != nil {
		return
	}

	return
}

func init() {
	docsProvider := func() map[string]interface{} {
		docs := make(map[string]interface{})
		docs["minifiers"] = parser.LowerCaseCamelJSONMarshaller{Value: defaultConfig}
		return docs

	}
	docshelper.AddDocProvider("config", docsProvider)
}
