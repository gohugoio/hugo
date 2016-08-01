// Copyright 2015 The Hugo Authors. All rights reserved.
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

package tpl

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/i18n/bundle"
	jww "github.com/spf13/jwalterweatherman"
)

var i18nTfunc bundle.TranslateFunc

func SetI18nTfunc(lang string, bndl *bundle.Bundle) {
	tFunc, err := bndl.Tfunc(lang)
	if err == nil {
		i18nTfunc = tFunc
		return
	}

	jww.WARN.Printf("could not load translations for language %q (%s), will not translate!\n", lang, err.Error())
	i18nTfunc = bundle.TranslateFunc(func(id string, args ...interface{}) string {
		// TODO: depending on the site mode, we might want to fall back on the default
		// language's translation.
		// TODO: eventually, we could add --i18n-warnings and print something when
		// such things happen.
		return fmt.Sprintf("[i18n: %s]", id)
	})
}

func I18nTranslate(id string, args ...interface{}) (string, error) {
	if i18nTfunc == nil {
		return "", fmt.Errorf("i18n not initialized, have you configured everything properly?")
	}
	return i18nTfunc(id, args...), nil
}
