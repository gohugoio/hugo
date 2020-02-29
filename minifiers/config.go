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

var DefaultConfig = minifiersConfig{
	EnableHtml: true,
	EnableCss:  true,
	EnableJs:   true,
	EnableJson: true,
	EnableSvg:  true,
	EnableXml:  true,

	Html: html.Minifier{
		KeepDocumentTags:        true,
		KeepConditionalComments: true,
		KeepEndTags:             true,
		KeepDefaultAttrVals:     true,
		KeepWhitespace:          false,
		KeepQuotes:              false,
	},
	Css: css.Minifier{
		Precision: 0,
		KeepCSS2:  true,
	},
	Js:   js.Minifier{},
	Json: json.Minifier{},
	Svg: svg.Minifier{
		Precision: 0,
	},
	Xml: xml.Minifier{
		KeepWhitespace: false,
	},
}

type minifiersConfig struct {
	EnableHtml bool
	EnableCss  bool
	EnableJs   bool
	EnableJson bool
	EnableSvg  bool
	EnableXml  bool

	Html html.Minifier
	Css  css.Minifier
	Js   js.Minifier
	Json json.Minifier
	Svg  svg.Minifier
	Xml  xml.Minifier
}

func decodeConfig(cfg config.Provider) (conf minifiersConfig, err error) {
	conf = DefaultConfig

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
		docs["minifiers"] = parser.LowerCaseCamelJSONMarshaller{Value: DefaultConfig}
		return docs

	}
	docshelper.AddDocProvider("config", docsProvider)
}
