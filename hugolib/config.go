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

	"github.com/spf13/afero"
	"github.com/spf13/hugo/helpers"
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

	v.AutomaticEnv()
	v.SetEnvPrefix("hugo")
	v.SetConfigFile(configFilename)
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

	v.RegisterAlias("indexes", "taxonomies")

	// Remove these in Hugo 0.23.
	if v.IsSet("disable404") {
		helpers.Deprecated("site config", "disable404", "Use disableKinds=[\"404\"]", false)
	}

	if v.IsSet("disableRSS") {
		helpers.Deprecated("site config", "disableRSS", "Use disableKinds=[\"RSS\"]", false)
	}

	if v.IsSet("disableSitemap") {
		// NOTE: Do not remove this until Hugo 0.24, ERROR in 0.23.
		helpers.Deprecated("site config", "disableSitemap", "Use disableKinds= [\"sitemap\"]", false)
	}

	if v.IsSet("disableRobotsTXT") {
		helpers.Deprecated("site config", "disableRobotsTXT", "Use disableKinds= [\"robotsTXT\"]", false)
	}

	loadDefaultSettingsFor(v)

	return v, nil
}

func loadDefaultSettingsFor(v *viper.Viper) {

	c := helpers.NewContentSpec(v)

	v.SetDefault("cleanDestinationDir", false)
	v.SetDefault("watch", false)
	v.SetDefault("metaDataFormat", "toml")
	v.SetDefault("disable404", false)
	v.SetDefault("disableRSS", false)
	v.SetDefault("disableSitemap", false)
	v.SetDefault("disableRobotsTXT", false)
	v.SetDefault("contentDir", "content")
	v.SetDefault("layoutDir", "layouts")
	v.SetDefault("staticDir", "static")
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
	v.SetDefault("taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	v.SetDefault("permalinks", make(PermalinkOverrides, 0))
	v.SetDefault("sitemap", Sitemap{Priority: -1, Filename: "sitemap.xml"})
	v.SetDefault("pygmentsStyle", "monokai")
	v.SetDefault("pygmentsUseClasses", false)
	v.SetDefault("pygmentsCodeFences", false)
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
	v.SetDefault("blackfriday", c.NewBlackfriday())
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
}
