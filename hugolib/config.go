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
	"fmt"

	"github.com/spf13/hugo/helpers"
	"github.com/spf13/viper"
)

// LoadGlobalConfig loads Hugo configuration into the global Viper.
func LoadGlobalConfig(relativeSourcePath, configFilename string) error {
	if relativeSourcePath == "" {
		relativeSourcePath = "."
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("hugo")
	viper.SetConfigFile(configFilename)
	// See https://github.com/spf13/viper/issues/73#issuecomment-126970794
	if relativeSourcePath == "" {
		viper.AddConfigPath(".")
	} else {
		viper.AddConfigPath(relativeSourcePath)
	}
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return err
		}
		return fmt.Errorf("Unable to locate Config file. Perhaps you need to create a new site.\n       Run `hugo help new` for details. (%s)\n", err)
	}

	viper.RegisterAlias("indexes", "taxonomies")

	loadDefaultSettings()

	return nil
}

func loadDefaultSettings() {
	viper.SetDefault("cleanDestinationDir", false)
	viper.SetDefault("watch", false)
	viper.SetDefault("metaDataFormat", "toml")
	viper.SetDefault("disable404", false)
	viper.SetDefault("disableRSS", false)
	viper.SetDefault("disableSitemap", false)
	viper.SetDefault("disableRobotsTXT", false)
	viper.SetDefault("contentDir", "content")
	viper.SetDefault("layoutDir", "layouts")
	viper.SetDefault("staticDir", "static")
	viper.SetDefault("archetypeDir", "archetypes")
	viper.SetDefault("publishDir", "public")
	viper.SetDefault("dataDir", "data")
	viper.SetDefault("i18nDir", "i18n")
	viper.SetDefault("themesDir", "themes")
	viper.SetDefault("defaultLayout", "post")
	viper.SetDefault("buildDrafts", false)
	viper.SetDefault("buildFuture", false)
	viper.SetDefault("buildExpired", false)
	viper.SetDefault("uglyURLs", false)
	viper.SetDefault("verbose", false)
	viper.SetDefault("ignoreCache", false)
	viper.SetDefault("canonifyURLs", false)
	viper.SetDefault("relativeURLs", false)
	viper.SetDefault("removePathAccents", false)
	viper.SetDefault("taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	viper.SetDefault("permalinks", make(PermalinkOverrides, 0))
	viper.SetDefault("sitemap", Sitemap{Priority: -1, Filename: "sitemap.xml"})
	viper.SetDefault("defaultExtension", "html")
	viper.SetDefault("pygmentsStyle", "monokai")
	viper.SetDefault("pygmentsUseClasses", false)
	viper.SetDefault("pygmentsCodeFences", false)
	viper.SetDefault("pygmentsOptions", "")
	viper.SetDefault("disableLiveReload", false)
	viper.SetDefault("pluralizeListTitles", true)
	viper.SetDefault("preserveTaxonomyNames", false)
	viper.SetDefault("forceSyncStatic", false)
	viper.SetDefault("footnoteAnchorPrefix", "")
	viper.SetDefault("footnoteReturnLinkContents", "")
	viper.SetDefault("newContentEditor", "")
	viper.SetDefault("paginate", 10)
	viper.SetDefault("paginatePath", "page")
	viper.SetDefault("blackfriday", helpers.NewBlackfriday(viper.GetViper()))
	viper.SetDefault("rSSUri", "index.xml")
	viper.SetDefault("sectionPagesMenu", "")
	viper.SetDefault("disablePathToLower", false)
	viper.SetDefault("hasCJKLanguage", false)
	viper.SetDefault("enableEmoji", false)
	viper.SetDefault("pygmentsCodeFencesGuessSyntax", false)
	viper.SetDefault("useModTimeAsFallback", false)
	viper.SetDefault("currentContentLanguage", helpers.NewDefaultLanguage())
	viper.SetDefault("defaultContentLanguage", "en")
	viper.SetDefault("defaultContentLanguageInSubdir", false)
	viper.SetDefault("enableMissingTranslationPlaceholders", false)
	viper.SetDefault("enableGitInfo", false)
}
