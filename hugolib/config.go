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
	"path/filepath"

	"io"
	"strings"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// ConfigSourceDescriptor describes where to find the config (e.g. config.toml etc.).
type ConfigSourceDescriptor struct {
	Fs afero.Fs

	// Full path to the config file to use, i.e. /my/project/config.toml
	Filename string

	// The path to the directory to look for configuration. Is used if Filename is not
	// set.
	Path string

	// The project's working dir. Is used to look for additional theme config.
	WorkingDir string
}

func (d ConfigSourceDescriptor) configFilenames() []string {
	return strings.Split(d.Filename, ",")
}

// LoadConfigDefault is a convenience method to load the default "config.toml" config.
func LoadConfigDefault(fs afero.Fs) (*viper.Viper, error) {
	v, _, err := LoadConfig(ConfigSourceDescriptor{Fs: fs, Filename: "config.toml"})
	return v, err
}

// LoadConfig loads Hugo configuration into a new Viper and then adds
// a set of defaults.
func LoadConfig(d ConfigSourceDescriptor, doWithConfig ...func(cfg config.Provider) error) (*viper.Viper, []string, error) {
	var configFiles []string

	fs := d.Fs
	v := viper.New()
	v.SetFs(fs)

	if d.Path == "" {
		d.Path = "."
	}

	configFilenames := d.configFilenames()
	v.AutomaticEnv()
	v.SetEnvPrefix("hugo")
	v.SetConfigFile(configFilenames[0])
	v.AddConfigPath(d.Path)

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return nil, configFiles, err
		}
		return nil, configFiles, fmt.Errorf("Unable to locate Config file. Perhaps you need to create a new site.\n       Run `hugo help new` for details. (%s)\n", err)
	}

	if cf := v.ConfigFileUsed(); cf != "" {
		configFiles = append(configFiles, cf)
	}

	for _, configFile := range configFilenames[1:] {
		var r io.Reader
		var err error
		if r, err = fs.Open(configFile); err != nil {
			return nil, configFiles, fmt.Errorf("Unable to open Config file.\n (%s)\n", err)
		}
		if err = v.MergeConfig(r); err != nil {
			return nil, configFiles, fmt.Errorf("Unable to parse/merge Config file (%s).\n (%s)\n", configFile, err)
		}
		configFiles = append(configFiles, configFile)
	}

	if err := loadDefaultSettingsFor(v); err != nil {
		return v, configFiles, err
	}

	themeConfigFile, err := loadThemeConfig(d, v)
	if err != nil {
		return v, configFiles, err
	}

	if themeConfigFile != "" {
		configFiles = append(configFiles, themeConfigFile)
	}

	// We create languages based on the settings, so we need to make sure that
	// all configuration is loaded/set before doing that.
	for _, d := range doWithConfig {
		if err := d(v); err != nil {
			return v, configFiles, err
		}
	}

	if err := loadLanguageSettings(v, nil); err != nil {
		return v, configFiles, err
	}

	return v, configFiles, nil

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
			for _, disabled := range disableLanguages {
				if disabled == defaultLang {
					return fmt.Errorf("cannot disable default language %q", defaultLang)
				}

				if strings.EqualFold(k, disabled) {
					v.(map[string]interface{})["disabled"] = true
					break
				}
			}
			languages[k] = v
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

func loadThemeConfig(d ConfigSourceDescriptor, v1 *viper.Viper) (string, error) {

	theme := v1.GetString("theme")
	if theme == "" {
		return "", nil
	}

	themesDir := helpers.AbsPathify(d.WorkingDir, v1.GetString("themesDir"))
	configDir := filepath.Join(themesDir, theme)

	var (
		configPath string
		exists     bool
		err        error
	)

	// Viper supports more, but this is the sub-set supported by Hugo.
	for _, configFormats := range []string{"toml", "yaml", "yml", "json"} {
		configPath = filepath.Join(configDir, "config."+configFormats)
		exists, err = helpers.Exists(configPath, d.Fs)
		if err != nil {
			return "", err
		}
		if exists {
			break
		}
	}

	if !exists {
		// No theme config set.
		return "", nil
	}

	v2 := viper.New()
	v2.SetFs(d.Fs)
	v2.AutomaticEnv()
	v2.SetEnvPrefix("hugo")
	v2.SetConfigFile(configPath)

	err = v2.ReadInConfig()
	if err != nil {
		return "", err
	}

	const (
		paramsKey    = "params"
		languagesKey = "languages"
		menuKey      = "menu"
	)

	for _, key := range []string{paramsKey, "outputformats", "mediatypes"} {
		mergeStringMapKeepLeft("", key, v1, v2)
	}

	themeLower := strings.ToLower(theme)
	themeParamsNamespace := paramsKey + "." + themeLower

	// Set namespaced params
	if v2.IsSet(paramsKey) && !v1.IsSet(themeParamsNamespace) {
		// Set it in the default store to make sure it gets in the same or
		// behind the others.
		v1.SetDefault(themeParamsNamespace, v2.Get(paramsKey))
	}

	// Only add params and new menu entries, we do not add language definitions.
	if v1.IsSet(languagesKey) && v2.IsSet(languagesKey) {
		v1Langs := v1.GetStringMap(languagesKey)
		for k, _ := range v1Langs {
			langParamsKey := languagesKey + "." + k + "." + paramsKey
			mergeStringMapKeepLeft(paramsKey, langParamsKey, v1, v2)
		}
		v2Langs := v2.GetStringMap(languagesKey)
		for k, _ := range v2Langs {
			if k == "" {
				continue
			}
			langParamsKey := languagesKey + "." + k + "." + paramsKey
			langParamsThemeNamespace := langParamsKey + "." + themeLower
			// Set namespaced params
			if v2.IsSet(langParamsKey) && !v1.IsSet(langParamsThemeNamespace) {
				v1.SetDefault(langParamsThemeNamespace, v2.Get(langParamsKey))
			}

			langMenuKey := languagesKey + "." + k + "." + menuKey
			if v2.IsSet(langMenuKey) {
				// Only add if not in the main config.
				v2menus := v2.GetStringMap(langMenuKey)
				for k, v := range v2menus {
					menuEntry := menuKey + "." + k
					menuLangEntry := langMenuKey + "." + k
					if !v1.IsSet(menuEntry) && !v1.IsSet(menuLangEntry) {
						v1.Set(menuLangEntry, v)
					}
				}
			}
		}
	}

	// Add menu definitions from theme not found in project
	if v2.IsSet("menu") {
		v2menus := v2.GetStringMap(menuKey)
		for k, v := range v2menus {
			menuEntry := menuKey + "." + k
			if !v1.IsSet(menuEntry) {
				v1.SetDefault(menuEntry, v)
			}
		}
	}

	return v2.ConfigFileUsed(), nil

}

func mergeStringMapKeepLeft(rootKey, key string, v1, v2 *viper.Viper) {
	if !v2.IsSet(key) {
		return
	}

	if !v1.IsSet(key) && !(rootKey != "" && rootKey != key && v1.IsSet(rootKey)) {
		v1.Set(key, v2.Get(key))
		return
	}

	m1 := v1.GetStringMap(key)
	m2 := v2.GetStringMap(key)

	for k, v := range m2 {
		if _, found := m1[k]; !found {
			if rootKey != "" && v1.IsSet(rootKey+"."+k) {
				continue
			}
			m1[k] = v
		}
	}
}

func loadDefaultSettingsFor(v *viper.Viper) error {

	c, err := helpers.NewContentSpec(v)
	if err != nil {
		return err
	}

	v.RegisterAlias("indexes", "taxonomies")

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
	v.SetDefault("timeout", 10000) // 10 seconds

	// Remove in Hugo 0.39

	if v.GetBool("useModTimeAsFallback") {

		helpers.Deprecated("Site config", "useModTimeAsFallback", `Replace with this in your config.toml:
    
[frontmatter]
date = [ "date",":fileModTime", ":default"]
lastmod = ["lastmod" ,":fileModTime", ":default"]
`, false)

	}

	return nil
}
