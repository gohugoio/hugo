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

package org_test

import (
	"testing"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config/testconfig"
	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/markup/converter"
	"github.com/gohugoio/hugo/markup/org"

	qt "github.com/frankban/quicktest"
)

func TestConvert(t *testing.T) {
	c := qt.New(t)
	p, err := org.Provider.New(converter.ProviderConfig{
		Logger: loggers.NewDefault(),
		Conf:   testconfig.GetTestConfig(afero.NewMemMapFs(), nil),
	})
	c.Assert(err, qt.IsNil)
	conv, err := p.New(converter.DocumentContext{})
	c.Assert(err, qt.IsNil)
	b, err := conv.Convert(converter.RenderContext{Src: []byte("testContent")})
	c.Assert(err, qt.IsNil)
	c.Assert(string(b.Bytes()), qt.Equals, "<p>testContent</p>\n")
}
