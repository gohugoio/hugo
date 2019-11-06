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

package mmark

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/gohugoio/hugo/common/loggers"

	qt "github.com/frankban/quicktest"
	"github.com/gohugoio/hugo/markup/blackfriday/blackfriday_config"
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/miekg/mmark"
)

func TestGetMmarkExtensions(t *testing.T) {
	b := blackfriday_config.Default

	//TODO: This is doing the same just with different marks...
	type data struct {
		testFlag int
	}

	b.Extensions = []string{"tables"}
	b.ExtensionsMask = []string{""}
	allExtensions := []data{
		{mmark.EXTENSION_TABLES},
		{mmark.EXTENSION_FENCED_CODE},
		{mmark.EXTENSION_AUTOLINK},
		{mmark.EXTENSION_SPACE_HEADERS},
		{mmark.EXTENSION_CITATION},
		{mmark.EXTENSION_TITLEBLOCK_TOML},
		{mmark.EXTENSION_HEADER_IDS},
		{mmark.EXTENSION_AUTO_HEADER_IDS},
		{mmark.EXTENSION_UNIQUE_HEADER_IDS},
		{mmark.EXTENSION_FOOTNOTES},
		{mmark.EXTENSION_SHORT_REF},
		{mmark.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK},
		{mmark.EXTENSION_INCLUDE},
	}

	actualFlags := getMmarkExtensions(b)
	for _, e := range allExtensions {
		if actualFlags&e.testFlag != e.testFlag {
			t.Errorf("Flag %v was not found in the list of extensions.", e)
		}
	}
}

func TestConvert(t *testing.T) {
	c := qt.New(t)
	p, err := Provider.New(converter.ProviderConfig{Cfg: viper.New(), Logger: loggers.NewErrorLogger()})
	c.Assert(err, qt.IsNil)
	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)
	b, err := conv.Convert(converter.RenderContext{Src: []byte("testContent")})
	c.Assert(err, qt.IsNil)
	c.Assert(string(b.Bytes()), qt.Equals, "<p>testContent</p>\n")
}
