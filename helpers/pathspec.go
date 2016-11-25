// Copyright 2016-present The Hugo Authors. All rights reserved.
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

package helpers

import "github.com/spf13/viper"

// PathSpec holds methods that decides how paths in URLs and files in Hugo should look like.
type PathSpec struct {
	disablePathToLower bool
	removePathAccents  bool
	uglyURLs           bool
	canonifyURLs       bool

	currentContentLanguage *Language

	// pagination path handling
	paginatePath string

	// The PathSpec looks up its config settings in both the current language
	// and then in the global Viper config.
	// Some settings, the settings listed below, does not make sense to be set
	// on per-language-basis. We have no good way of protecting against this
	// other than a "white-list". See language.go.
	defaultContentLanguageInSubdir bool
	defaultContentLanguage         string
	multilingual                   bool
}

// NewPathSpecFromCurrentLanguage creates a new PathSpec from the
// current language (currentContentLanguage in viper).
func NewPathSpecFromCurrentLanguage() *PathSpec {
	return NewPathSpecFromLanguage(viper.Get("currentContentLanguage").(*Language))
}

// NewPathSpecFromViper creates a new PathSpec from the global viper instance.
func NewPathSpecFromViper() *PathSpec {
	return &PathSpec{
		disablePathToLower:             viper.GetBool("disablePathToLower"),
		removePathAccents:              viper.GetBool("removePathAccents"),
		uglyURLs:                       viper.GetBool("uglyURLs"),
		canonifyURLs:                   viper.GetBool("canonifyURLs"),
		multilingual:                   viper.GetBool("multilingual"),
		defaultContentLanguageInSubdir: viper.GetBool("defaultContentLanguageInSubdir"),
		defaultContentLanguage:         viper.GetString("defaultContentLanguage"),
		currentContentLanguage:         viper.Get("currentContentLanguage").(*Language),
		paginatePath:                   viper.GetString("paginatePath"),
	}
}

// NewPathSpecFromLanguage creates a new PathSpec from the given Language.
func NewPathSpecFromLanguage(l *Language) *PathSpec {
	return &PathSpec{
		disablePathToLower:             l.GetBool("disablePathToLower"),
		removePathAccents:              l.GetBool("removePathAccents"),
		uglyURLs:                       l.GetBool("uglyURLs"),
		canonifyURLs:                   l.GetBool("canonifyURLs"),
		multilingual:                   l.GetBool("multilingual"),
		defaultContentLanguageInSubdir: l.GetBool("defaultContentLanguageInSubdir"),
		defaultContentLanguage:         l.GetString("defaultContentLanguage"),
		currentContentLanguage:         l.Get("currentContentLanguage").(*Language),
		paginatePath:                   l.GetString("paginatePath"),
	}
}

// PaginatePath returns the configured root path used for paginator pages.
func (p *PathSpec) PaginatePath() string {
	return p.paginatePath
}

var currentPathSpec *PathSpec

func InitCurrentPathSpec() {
	currentPathSpec = NewPathSpecFromCurrentLanguage()
}

// CurrentPathSpec returns the current PathSpec.
// If it is not set, a new will be created based in the currently active language.
func CurrentPathSpec() *PathSpec {
	if currentPathSpec != nil {
		return currentPathSpec
	}
	return NewPathSpecFromCurrentLanguage()
}

// ResetCurrentPathSpec is used in tests.
func ResetCurrentPathSpec() {
	currentPathSpec = nil
}
