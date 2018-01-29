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

package hugolib

import (
	"errors"
	"fmt"

	"io"
	"strings"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// LoadConfig loads Hugo configuration into a new Viper and then adds
// a set of defaults.
func LoadConfig(fs afero.Fs, relativeSourcePath, configFilename string) (*viper.Viper, error) {
	v := viper.New()
	v.SetFs(fs)
	if relativeSourcePath == "" {
		relativeSourcePath = "."
	}
	configFilenames := strings.Split(configFilename, ",")
	v.AutomaticEnv()
	v.SetEnvPrefix("hugo")
	v.SetConfigFile(configFilenames[0])
	// See https://github.com/spf13/viper/issues/73#issuecomment-126970794
	if relativeSourcePath == "" {
		v.AddConfigPath(".")
	} else {
		v.AddConfigPath(relativeSourcePath)
	}
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return nil, err
		}
		return nil, fmt.Errorf("Unable to locate Config file. Perhaps you need to create a new site.\n       Run `hugo help new` for details. (%s)\n", err)
	}
	for _, configFile := range configFilenames[1:] {
		var r io.Reader
		var err error
		if r, err = fs.Open(configFile); err != nil {
			return nil, fmt.Errorf("Unable to open Config file.\n (%s)\n", err)
		}
		if err = v.MergeConfig(r); err != nil {
			return nil, fmt.Errorf("Unable to parse/merge Config file (%s).\n (%s)\n", configFile, err)
		}
	}

	v.RegisterAlias("indexes", "taxonomies")

	if err := loadDefaultSettingsFor(v); err != nil {
		return v, err
	}

	return v, nil
}

func loadLanguageSettings(cfg config.Provider, oldLangs helpers.Languages) error {

	defaultLang := cfg.GetString("defaultContentLanguage")

	var languages map[string]interface{}

	languagesFromConfig := cfg.GetStringMap("languages")
	disableLanguages := cfg.GetStringSlice("disableLanguages")

	if len(disableLanguages) == 0 {
		languages = languagesFromConfig
	} else {
		languages = make(map[string]interface{})
		for k, v := range languagesFromConfig {
			isDisabled := false
			for _, disabled := range disableLanguages {
				if disabled == defaultLang {
					return fmt.Errorf("cannot disable default language %q", defaultLang)
				}

				if strings.EqualFold(k, disabled) {
					isDisabled = true
					break
				}
			}
			if !isDisabled {
				languages[k] = v
			}

		}
	}

	var (
		langs helpers.Languages
		err   error
	)

	if len(languages) == 0 {
		langs = append(langs, helpers.NewDefaultLanguage(cfg))
	} else {
		langs, err = toSortedLanguages(cfg, languages)
		if err != nil {
			return fmt.Errorf("Failed to parse multilingual config: %s", err)
		}
	}

	if oldLangs != nil {
		// When in multihost mode, the languages are mapped to a server, so
		// some structural language changes will need a restart of the dev server.
		// The validation below isn't complete, but should cover the most
		// important cases.
		var invalid bool
		if langs.IsMultihost() != oldLangs.IsMultihost() {
			invalid = true
		} else {
			if langs.IsMultihost() && len(langs) != len(oldLangs) {
				invalid = true
			}
		}

		if invalid {
			return errors.New("language change needing a server restart detected")
		}

		if langs.IsMultihost() {
			// We need to transfer any server baseURL to the new language
			for i, ol := range oldLangs {
				nl := langs[i]
				nl.Set("baseURL", ol.GetString("baseURL"))
			}
		}
	}

	// The defaultContentLanguage is something the user has to decide, but it needs
	// to match a language in the language definition list.
	langExists := false
	for _, lang := range langs {
		if lang.Lang == defaultLang {
			langExists = true
			break
		}
	}

	if !langExists {
		return fmt.Errorf("site config value %q for defaultContentLanguage does not match any language definition", defaultLang)
	}

	cfg.Set("languagesSorted", langs)
	cfg.Set("multilingual", len(langs) > 1)

	multihost := langs.IsMultihost()

	if multihost {
		cfg.Set("defaultContentLanguageInSubdir", true)
		cfg.Set("multihost", true)
	}

	if multihost {
		// The baseURL may be provided at the language level. If that is true,
		// then every language must have a baseURL. In this case we always render
		// to a language sub folder, which is then stripped from all the Permalink URLs etc.
		for _, l := range langs {
			burl := l.GetLocal("baseURL")
			if burl == nil {
				return errors.New("baseURL must be set on all or none of the languages")
			}
		}

	}

	return nil
}

func loadDefaultSettingsFor(v *viper.Viper) error {

	c, err := helpers.NewContentSpec(v)
	if err != nil {
		return err
	}

	v.SetDefault("cleanDestinationDir", false)
	v.SetDefault("watch", false)
	v.SetDefault("metaDataFormat", "toml")
	v.SetDefault("contentDir", "content")
	v.SetDefault("layoutDir", "layouts")
	v.SetDefault("staticDir", "static")
	v.SetDefault("resourceDir", "resources")
	v.SetDefault("archetypeDir", "archetypes")
	v.SetDefault("publishDir", "public")
	v.SetDefault("dataDir", "data")
	v.SetDefault("i18nDir", "i18n")
	v.SetDefault("themesDir", "themes")
	v.SetDefault("defaultLayout", "post")
	v.SetDefault("buildDrafts", false)
	v.SetDefault("buildFuture", false)
	v.SetDefault("buildExpired", false)
	v.SetDefault("uglyURLs", false)
	v.SetDefault("verbose", false)
	v.SetDefault("ignoreCache", false)
	v.SetDefault("canonifyURLs", false)
	v.SetDefault("relativeURLs", false)
	v.SetDefault("removePathAccents", false)
	v.SetDefault("titleCaseStyle", "AP")
	v.SetDefault("taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	v.SetDefault("permalinks", make(PermalinkOverrides, 0))
	v.SetDefault("sitemap", Sitemap{Priority: -1, Filename: "sitemap.xml"})
	v.SetDefault("pygmentsStyle", "monokai")
	v.SetDefault("pygmentsUseClasses", false)
	v.SetDefault("pygmentsCodeFences", false)
	v.SetDefault("pygmentsUseClassic", false)
	v.SetDefault("pygmentsOptions", "")
	v.SetDefault("disableLiveReload", false)
	v.SetDefault("pluralizeListTitles", true)
	v.SetDefault("preserveTaxonomyNames", false)
	v.SetDefault("forceSyncStatic", false)
	v.SetDefault("footnoteAnchorPrefix", "")
	v.SetDefault("footnoteReturnLinkContents", "")
	v.SetDefault("newContentEditor", "")
	v.SetDefault("paginate", 10)
	v.SetDefault("paginatePath", "page")
	v.SetDefault("summaryLength", 70)
	v.SetDefault("blackfriday", c.BlackFriday)
	v.SetDefault("rSSUri", "index.xml")
	v.SetDefault("rssLimit", -1)
	v.SetDefault("sectionPagesMenu", "")
	v.SetDefault("disablePathToLower", false)
	v.SetDefault("hasCJKLanguage", false)
	v.SetDefault("enableEmoji", false)
	v.SetDefault("pygmentsCodeFencesGuessSyntax", false)
	v.SetDefault("useModTimeAsFallback", false)
	v.SetDefault("defaultContentLanguage", "en")
	v.SetDefault("defaultContentLanguageInSubdir", false)
	v.SetDefault("enableMissingTranslationPlaceholders", false)
	v.SetDefault("enableGitInfo", false)
	v.SetDefault("ignoreFiles", make([]string, 0))
	v.SetDefault("disableAliases", false)
	v.SetDefault("debug", false)
	v.SetDefault("disableFastRender", false)

	// Remove in Hugo 0.37
	if v.GetBool("useModTimeAsFallback") {
		helpers.Deprecated("Site config", "useModTimeAsFallback", "Try --enableGitInfo or set lastMod in front matter", false)

	}

	return loadLanguageSettings(v, nil)
}
