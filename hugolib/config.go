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
	viper.SetDefault("Watch", false)
	viper.SetDefault("MetaDataFormat", "toml")
	viper.SetDefault("Disable404", false)
	viper.SetDefault("DisableRSS", false)
	viper.SetDefault("DisableSitemap", false)
	viper.SetDefault("DisableRobotsTXT", false)
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("LayoutDir", "layouts")
	viper.SetDefault("StaticDir", "static")
	viper.SetDefault("ArchetypeDir", "archetypes")
	viper.SetDefault("PublishDir", "public")
	viper.SetDefault("DataDir", "data")
	viper.SetDefault("I18nDir", "i18n")
	viper.SetDefault("ThemesDir", "themes")
	viper.SetDefault("DefaultLayout", "post")
	viper.SetDefault("BuildDrafts", false)
	viper.SetDefault("BuildFuture", false)
	viper.SetDefault("BuildExpired", false)
	viper.SetDefault("UglyURLs", false)
	viper.SetDefault("Verbose", false)
	viper.SetDefault("IgnoreCache", false)
	viper.SetDefault("CanonifyURLs", false)
	viper.SetDefault("RelativeURLs", false)
	viper.SetDefault("RemovePathAccents", false)
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	viper.SetDefault("Permalinks", make(PermalinkOverrides, 0))
	viper.SetDefault("Sitemap", Sitemap{Priority: -1, Filename: "sitemap.xml"})
	viper.SetDefault("DefaultExtension", "html")
	viper.SetDefault("PygmentsStyle", "monokai")
	viper.SetDefault("PygmentsUseClasses", false)
	viper.SetDefault("PygmentsCodeFences", false)
	viper.SetDefault("PygmentsOptions", "")
	viper.SetDefault("DisableLiveReload", false)
	viper.SetDefault("PluralizeListTitles", true)
	viper.SetDefault("PreserveTaxonomyNames", false)
	viper.SetDefault("ForceSyncStatic", false)
	viper.SetDefault("FootnoteAnchorPrefix", "")
	viper.SetDefault("FootnoteReturnLinkContents", "")
	viper.SetDefault("NewContentEditor", "")
	viper.SetDefault("Paginate", 10)
	viper.SetDefault("PaginatePath", "page")
	viper.SetDefault("Blackfriday", helpers.NewBlackfriday())
	viper.SetDefault("RSSUri", "index.xml")
	viper.SetDefault("SectionPagesMenu", "")
	viper.SetDefault("DisablePathToLower", false)
	viper.SetDefault("HasCJKLanguage", false)
	viper.SetDefault("EnableEmoji", false)
	viper.SetDefault("PygmentsCodeFencesGuessSyntax", false)
	viper.SetDefault("UseModTimeAsFallback", false)
	viper.SetDefault("Multilingual", false)
	viper.SetDefault("DefaultContentLanguage", "en")
}
