// Copyright 2018 The Hugo Authors. All rights reserved.
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

package htesting

import (
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/langs"
	"github.com/spf13/viper"
)

type testSite struct {
	h hugo.Info
	l *langs.Language
}

func (t testSite) Hugo() hugo.Info {
	return t.h
}

func (t testSite) IsServer() bool {
	return false
}

func (t testSite) Language() *langs.Language {
	return t.l
}

// NewTestHugoSite creates a new minimal test site.
func NewTestHugoSite() hugo.Site {
	return testSite{
		h: hugo.NewInfo(hugo.EnvironmentProduction),
		l: langs.NewLanguage("en", newTestConfig()),
	}
}

func newTestConfig() *viper.Viper {
	v := viper.New()
	v.Set("contentDir", "content")
	return v
}
