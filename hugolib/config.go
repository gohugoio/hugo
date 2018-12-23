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

	"os"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib/paths"
	"github.com/pkg/errors"
	_errors "github.com/pkg/errors"

	"github.com/gohugoio/hugo/langs"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/privacy"
	"github.com/gohugoio/hugo/config/services"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// SiteConfig represents the config in .Site.Config.
type SiteConfig struct {
	// This contains all privacy related settings that can be used to
	// make the YouTube template etc. GDPR compliant.
	Privacy privacy.Config

	// Services contains config for services such as Google Analytics etc.
	Services services.Config
}

func loadSiteConfig(cfg config.Provider) (scfg SiteConfig, err error) {
	privacyConfig, err := privacy.DecodeConfig(cfg)
	if err != nil {
		return
	}

	servicesConfig, err := services.DecodeConfig(cfg)
	if err != nil {
		return
	}

	scfg.Privacy = privacyConfig
	scfg.Services = servicesConfig

	return
}

// ConfigSourceDescriptor describes where to find the config (e.g. config.toml etc.).
type ConfigSourceDescriptor struct {
	Fs afero.Fs

	// Path to the config file to use, e.g. /my/project/config.toml
	Filename string

	// The path to the directory to look for configuration. Is used if Filename is not
	// set or if it is set to a relative filename.
	Path string

	// The project's working dir. Is used to look for additional theme config.
	WorkingDir string

	// The (optional) directory for additional configuration files.
	AbsConfigDir string

	// production, development
	Environment string
}

func (d ConfigSourceDescriptor) configFilenames() []string {
	if d.Filename == "" {
		return []string{"config"}
	}
	return strings.Split(d.Filename, ",")
}

func (d ConfigSourceDescriptor) configFileDir() string {
	if d.Path != "" {
		return d.Path
	}
	return d.WorkingDir
}

// LoadConfigDefault is a convenience method to load the default "config.toml" config.
func LoadConfigDefault(fs afero.Fs) (*viper.Viper, error) {
	v, _, err := LoadConfig(ConfigSourceDescriptor{Fs: fs, Filename: "config.toml"})
	return v, err
}

var ErrNoConfigFile = errors.New("Unable to locate config file or config directory. Perhaps you need to create a new site.\n       Run `hugo help new` for details.\n")

// LoadConfig loads Hugo configuration into a new Viper and then adds
// a set of defaults.
func LoadConfig(d ConfigSourceDescriptor, doWithConfig ...func(cfg config.Provider) error) (*viper.Viper, []string, error) {
	if d.Environment == "" {
		d.Environment = hugo.EnvironmentProduction
	}

	var configFiles []string

	v := viper.New()
	l := configLoader{ConfigSourceDescriptor: d}

	v.AutomaticEnv()
	v.SetEnvPrefix("hugo")

	var cerr error

	for _, name := range d.configFilenames() {
		var filename string
		if filename, cerr = l.loadConfig(name, v); cerr != nil && cerr != ErrNoConfigFile {
			return nil, nil, cerr
		}
		configFiles = append(configFiles, filename)
	}

	if d.AbsConfigDir != "" {
		dirnames, err := l.loadConfigFromConfigDir(v)
		if err == nil {
			configFiles = append(configFiles, dirnames...)
		}
		cerr = err
	}

	if err := loadDefaultSettingsFor(v); err != nil {
		return v, configFiles, err
	}

	if cerr == nil {
		themeConfigFiles, err := l.loadThemeConfig(v)
		if err != nil {
			return v, configFiles, err
		}

		if len(themeConfigFiles) > 0 {
			configFiles = append(configFiles, themeConfigFiles...)
		}
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

	return v, configFiles, cerr

}

type configLoader struct {
	ConfigSourceDescriptor
}

func (l configLoader) wrapFileInfoError(err error, fi os.FileInfo) error {
	rfi, ok := fi.(hugofs.RealFilenameInfo)
	if !ok {
		return err
	}
	return l.wrapFileError(err, rfi.RealFilename())
}

func (l configLoader) loadConfig(configName string, v *viper.Viper) (string, error) {
	baseDir := l.configFileDir()
	var baseFilename string
	if filepath.IsAbs(configName) {
		baseFilename = configName
	} else {
		baseFilename = filepath.Join(baseDir, configName)
	}

	var filename string
	fileExt := helpers.ExtNoDelimiter(configName)
	if fileExt != "" {
		exists, _ := helpers.Exists(baseFilename, l.Fs)
		if exists {
			filename = baseFilename
		}
	} else {
		for _, ext := range []string{"toml", "yaml", "yml", "json"} {
			filenameToCheck := baseFilename + "." + ext
			exists, _ := helpers.Exists(filenameToCheck, l.Fs)
			if exists {
				filename = filenameToCheck
				fileExt = ext
				break
			}
		}
	}

	if filename == "" {
		return "", ErrNoConfigFile
	}

	m, err := config.FromFileToMap(l.Fs, filename)
	if err != nil {
		return "", l.wrapFileError(err, filename)
	}

	if err = v.MergeConfigMap(m); err != nil {
		return "", l.wrapFileError(err, filename)
	}

	return filename, nil

}

func (l configLoader) wrapFileError(err error, filename string) error {
	err, _ = herrors.WithFileContextForFile(
		err,
		filename,
		filename,
		l.Fs,
		herrors.SimpleLineMatcher)
	return err
}

func (l configLoader) newRealBaseFs(path string) afero.Fs {
	return hugofs.NewBasePathRealFilenameFs(afero.NewBasePathFs(l.Fs, path).(*afero.BasePathFs))

}

func (l configLoader) loadConfigFromConfigDir(v *viper.Viper) ([]string, error) {
	sourceFs := l.Fs
	configDir := l.AbsConfigDir

	if _, err := sourceFs.Stat(configDir); err != nil {
		// Config dir does not exist.
		return nil, nil
	}

	defaultConfigDir := filepath.Join(configDir, "_default")
	environmentConfigDir := filepath.Join(configDir, l.Environment)

	var configDirs []string
	// Merge from least to most specific.
	for _, dir := range []string{defaultConfigDir, environmentConfigDir} {
		if _, err := sourceFs.Stat(dir); err == nil {
			configDirs = append(configDirs, dir)
		}
	}

	if len(configDirs) == 0 {
		return nil, nil
	}

	// Keep track of these so we can watch them for changes.
	var dirnames []string

	for _, configDir := range configDirs {
		err := afero.Walk(sourceFs, configDir, func(path string, fi os.FileInfo, err error) error {
			if fi == nil {
				return nil
			}

			if fi.IsDir() {
				dirnames = append(dirnames, path)
				return nil
			}

			name := helpers.Filename(filepath.Base(path))

			item, err := metadecoders.Default.UnmarshalFileToMap(sourceFs, path)
			if err != nil {
				return l.wrapFileError(err, path)
			}

			var keyPath []string

			if name != "config" {
				// Can be params.jp, menus.en etc.
				name, lang := helpers.FileAndExtNoDelimiter(name)

				keyPath = []string{name}

				if lang != "" {
					keyPath = []string{"languages", lang}
					switch name {
					case "menu", "menus":
						keyPath = append(keyPath, "menus")
					case "params":
						keyPath = append(keyPath, "params")
					}
				}
			}

			root := item
			if len(keyPath) > 0 {
				root = make(map[string]interface{})
				m := root
				for i, key := range keyPath {
					if i >= len(keyPath)-1 {
						m[key] = item
					} else {
						nm := make(map[string]interface{})
						m[key] = nm
						m = nm
					}
				}
			}

			// Migrate menu => menus etc.
			config.RenameKeys(root)

			if err := v.MergeConfigMap(root); err != nil {
				return l.wrapFileError(err, path)
			}

			return nil

		})

		if err != nil {
			return nil, err
		}

	}

	return dirnames, nil
}

func loadLanguageSettings(cfg config.Provider, oldLangs langs.Languages) error {

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
		languages2 langs.Languages
		err        error
	)

	if len(languages) == 0 {
		languages2 = append(languages2, langs.NewDefaultLanguage(cfg))
	} else {
		languages2, err = toSortedLanguages(cfg, languages)
		if err != nil {
			return _errors.Wrap(err, "Failed to parse multilingual config")
		}
	}

	if oldLangs != nil {
		// When in multihost mode, the languages are mapped to a server, so
		// some structural language changes will need a restart of the dev server.
		// The validation below isn't complete, but should cover the most
		// important cases.
		var invalid bool
		if languages2.IsMultihost() != oldLangs.IsMultihost() {
			invalid = true
		} else {
			if languages2.IsMultihost() && len(languages2) != len(oldLangs) {
				invalid = true
			}
		}

		if invalid {
			return errors.New("language change needing a server restart detected")
		}

		if languages2.IsMultihost() {
			// We need to transfer any server baseURL to the new language
			for i, ol := range oldLangs {
				nl := languages2[i]
				nl.Set("baseURL", ol.GetString("baseURL"))
			}
		}
	}

	// The defaultContentLanguage is something the user has to decide, but it needs
	// to match a language in the language definition list.
	langExists := false
	for _, lang := range languages2 {
		if lang.Lang == defaultLang {
			langExists = true
			break
		}
	}

	if !langExists {
		return fmt.Errorf("site config value %q for defaultContentLanguage does not match any language definition", defaultLang)
	}

	cfg.Set("languagesSorted", languages2)
	cfg.Set("multilingual", len(languages2) > 1)

	multihost := languages2.IsMultihost()

	if multihost {
		cfg.Set("defaultContentLanguageInSubdir", true)
		cfg.Set("multihost", true)
	}

	if multihost {
		// The baseURL may be provided at the language level. If that is true,
		// then every language must have a baseURL. In this case we always render
		// to a language sub folder, which is then stripped from all the Permalink URLs etc.
		for _, l := range languages2 {
			burl := l.GetLocal("baseURL")
			if burl == nil {
				return errors.New("baseURL must be set on all or none of the languages")
			}
		}

	}

	return nil
}

func (l configLoader) loadThemeConfig(v1 *viper.Viper) ([]string, error) {
	themesDir := paths.AbsPathify(l.WorkingDir, v1.GetString("themesDir"))
	themes := config.GetStringSlicePreserveString(v1, "theme")

	themeConfigs, err := paths.CollectThemes(l.Fs, themesDir, themes)
	if err != nil {
		return nil, err
	}

	if len(themeConfigs) == 0 {
		return nil, nil
	}

	v1.Set("allThemes", themeConfigs)

	var configFilenames []string
	for _, tc := range themeConfigs {
		if tc.ConfigFilename != "" {
			configFilenames = append(configFilenames, tc.ConfigFilename)
			if err := l.applyThemeConfig(v1, tc); err != nil {
				return nil, err
			}
		}
	}

	return configFilenames, nil

}

func (l configLoader) applyThemeConfig(v1 *viper.Viper, theme paths.ThemeConfig) error {

	const (
		paramsKey    = "params"
		languagesKey = "languages"
		menuKey      = "menus"
	)

	v2 := theme.Cfg

	for _, key := range []string{paramsKey, "outputformats", "mediatypes"} {
		l.mergeStringMapKeepLeft("", key, v1, v2)
	}

	themeLower := strings.ToLower(theme.Name)
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
		for k := range v1Langs {
			langParamsKey := languagesKey + "." + k + "." + paramsKey
			l.mergeStringMapKeepLeft(paramsKey, langParamsKey, v1, v2)
		}
		v2Langs := v2.GetStringMap(languagesKey)
		for k := range v2Langs {
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
	if v2.IsSet(menuKey) {
		v2menus := v2.GetStringMap(menuKey)
		for k, v := range v2menus {
			menuEntry := menuKey + "." + k
			if !v1.IsSet(menuEntry) {
				v1.SetDefault(menuEntry, v)
			}
		}
	}

	return nil

}

func (configLoader) mergeStringMapKeepLeft(rootKey, key string, v1, v2 config.Provider) {
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
	v.SetDefault("assetDir", "assets")
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
	v.SetDefault("environment", hugo.EnvironmentProduction)
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
	v.SetDefault("defaultContentLanguage", "en")
	v.SetDefault("defaultContentLanguageInSubdir", false)
	v.SetDefault("enableMissingTranslationPlaceholders", false)
	v.SetDefault("enableGitInfo", false)
	v.SetDefault("ignoreFiles", make([]string, 0))
	v.SetDefault("disableAliases", false)
	v.SetDefault("debug", false)
	v.SetDefault("disableFastRender", false)
	v.SetDefault("timeout", 10000) // 10 seconds
	v.SetDefault("enableInlineShortcodes", false)
	return nil
}
