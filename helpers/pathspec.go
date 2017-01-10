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

import (
	"fmt"

	"github.com/spf13/hugo/hugofs"
)

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

	// The file systems to use
	fs *hugofs.Fs
}

func (p PathSpec) String() string {
	return fmt.Sprintf("PathSpec, language %q, prefix %q, multilingual: %T", p.currentContentLanguage.Lang, p.getLanguagePrefix(), p.multilingual)
}

// NewPathSpec creats a new PathSpec from the given filesystems and ConfigProvider.
func NewPathSpec(fs *hugofs.Fs, config ConfigProvider) *PathSpec {

	currCl, ok := config.Get("currentContentLanguage").(*Language)

	if !ok {
		// TODO(bep) globals
		currCl = NewLanguage("en")
	}

	return &PathSpec{
		fs:                             fs,
		disablePathToLower:             config.GetBool("disablePathToLower"),
		removePathAccents:              config.GetBool("removePathAccents"),
		uglyURLs:                       config.GetBool("uglyURLs"),
		canonifyURLs:                   config.GetBool("canonifyURLs"),
		multilingual:                   config.GetBool("multilingual"),
		defaultContentLanguageInSubdir: config.GetBool("defaultContentLanguageInSubdir"),
		defaultContentLanguage:         config.GetString("defaultContentLanguage"),
		currentContentLanguage:         currCl,
		paginatePath:                   config.GetString("paginatePath"),
	}
}

// PaginatePath returns the configured root path used for paginator pages.
func (p *PathSpec) PaginatePath() string {
	return p.paginatePath
}
