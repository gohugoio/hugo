// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"os"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/cache/filecache"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/hugolib/paths"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/modules"
	"github.com/pkg/errors"

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
	Fs     afero.Fs
	Logger *loggers.Logger

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

	// Defaults to os.Environ if not set.
	Environ []string
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

	if len(d.Environ) == 0 {
		d.Environ = os.Environ()
	}

	var configFiles []string

	v := viper.New()
	l := configLoader{ConfigSourceDescriptor: d}

	for _, name := range d.configFilenames() {
		var filename string
		filename, err := l.loadConfig(name, v)
		if err == nil {
			configFiles = append(configFiles, filename)
		} else if err != ErrNoConfigFile {
			return nil, nil, err
		}
	}

	if d.AbsConfigDir != "" {
		dirnames, err := l.loadConfigFromConfigDir(v)
		if err == nil {
			configFiles = append(configFiles, dirnames...)
		} else if err != ErrNoConfigFile {
			return nil, nil, err
		}
	}

	if err := loadDefaultSettingsFor(v); err != nil {
		return v, configFiles, err
	}

	// We create languages based on the settings, so we need to make sure that
	// all configuration is loaded/set before doing that.
	for _, d := range doWithConfig {
		if err := d(v); err != nil {
			return v, configFiles, err
		}
	}

	// Apply environment overrides
	if len(d.Environ) > 0 {
		// Extract all that start with the HUGO_ prefix
		const hugoEnvPrefix = "HUGO_"
		var hugoEnv []string
		for _, v := range d.Environ {
			key, val := config.SplitEnvVar(v)
			if strings.HasPrefix(key, hugoEnvPrefix) {
				hugoEnv = append(hugoEnv, strings.ToLower(strings.TrimPrefix(key, hugoEnvPrefix)), val)
			}
		}

		if len(hugoEnv) > 0 {
			for i := 0; i < len(hugoEnv); i += 2 {
				key, valStr := strings.ToLower(hugoEnv[i]), hugoEnv[i+1]

				existing, nestedKey, owner, err := maps.GetNestedParamFn(key, "_", v.Get)
				if err != nil {
					return v, configFiles, err
				}

				if existing != nil {
					val, err := metadecoders.Default.UnmarshalStringTo(valStr, existing)
					if err != nil {
						continue
					}

					if owner != nil {
						owner[nestedKey] = val
					} else {
						v.Set(key, val)
					}
				} else {
					v.Set(key, valStr)
				}
			}
		}
	}

	modulesConfig, err := l.loadModulesConfig(v)
	if err != nil {
		return v, configFiles, err
	}

	// Need to run these after the modules are loaded, but before
	// they are finalized.
	collectHook := func(m *modules.ModulesConfig) error {
		if err := loadLanguageSettings(v, nil); err != nil {
			return err
		}

		mods := m.ActiveModules

		// Apply default project mounts.
		if err := modules.ApplyProjectConfigDefaults(v, mods[0]); err != nil {
			return err
		}

		return nil
	}

	_, modulesConfigFiles, err := l.collectModules(modulesConfig, v, collectHook)

	if err == nil && len(modulesConfigFiles) > 0 {
		configFiles = append(configFiles, modulesConfigFiles...)
	}

	return v, configFiles, err

}

func loadLanguageSettings(cfg config.Provider, oldLangs langs.Languages) error {
	_, err := langs.LoadLanguageSettings(cfg, oldLangs)
	return err
}

type configLoader struct {
	ConfigSourceDescriptor
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
	if helpers.ExtNoDelimiter(configName) != "" {
		exists, _ := helpers.Exists(baseFilename, l.Fs)
		if exists {
			filename = baseFilename
		}
	} else {
		for _, ext := range config.ValidConfigFileExtensions {
			filenameToCheck := baseFilename + "." + ext
			exists, _ := helpers.Exists(filenameToCheck, l.Fs)
			if exists {
				filename = filenameToCheck
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
			if fi == nil || err != nil {
				return nil
			}

			if fi.IsDir() {
				dirnames = append(dirnames, path)
				return nil
			}

			if !config.IsValidConfigFilename(path) {
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

func (l configLoader) loadModulesConfig(v1 *viper.Viper) (modules.Config, error) {

	modConfig, err := modules.DecodeConfig(v1)
	if err != nil {
		return modules.Config{}, err
	}

	return modConfig, nil
}

func (l configLoader) collectModules(modConfig modules.Config, v1 *viper.Viper, hookBeforeFinalize func(m *modules.ModulesConfig) error) (modules.Modules, []string, error) {
	workingDir := l.WorkingDir
	if workingDir == "" {
		workingDir = v1.GetString("workingDir")
	}

	themesDir := paths.AbsPathify(l.WorkingDir, v1.GetString("themesDir"))

	ignoreVendor := v1.GetBool("ignoreVendor")

	filecacheConfigs, err := filecache.DecodeConfig(l.Fs, v1)
	if err != nil {
		return nil, nil, err
	}

	v1.Set("filecacheConfigs", filecacheConfigs)

	var configFilenames []string

	hook := func(m *modules.ModulesConfig) error {
		for _, tc := range m.ActiveModules {
			if tc.ConfigFilename() != "" {
				if tc.Watch() {
					configFilenames = append(configFilenames, tc.ConfigFilename())
				}
				if err := l.applyThemeConfig(v1, tc); err != nil {
					return err
				}
			}
		}

		if hookBeforeFinalize != nil {
			return hookBeforeFinalize(m)
		}

		return nil

	}

	modulesClient := modules.NewClient(modules.ClientConfig{
		Fs:                 l.Fs,
		Logger:             l.Logger,
		HookBeforeFinalize: hook,
		WorkingDir:         workingDir,
		ThemesDir:          themesDir,
		CacheDir:           filecacheConfigs.CacheDirModules(),
		ModuleConfig:       modConfig,
		IgnoreVendor:       ignoreVendor,
	})

	v1.Set("modulesClient", modulesClient)

	moduleConfig, err := modulesClient.Collect()

	// Avoid recreating these later.
	v1.Set("allModules", moduleConfig.ActiveModules)

	if moduleConfig.GoModulesFilename != "" {
		// We want to watch this for changes and trigger rebuild on version
		// changes etc.
		configFilenames = append(configFilenames, moduleConfig.GoModulesFilename)
	}

	return moduleConfig.ActiveModules, configFilenames, err

}

func (l configLoader) applyThemeConfig(v1 *viper.Viper, theme modules.Module) error {

	const (
		paramsKey    = "params"
		languagesKey = "languages"
		menuKey      = "menus"
	)

	v2 := theme.Cfg()

	for _, key := range []string{paramsKey, "outputformats", "mediatypes"} {
		l.mergeStringMapKeepLeft("", key, v1, v2)
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

	v.RegisterAlias("indexes", "taxonomies")

	/*

		TODO(bep) from 0.56 these are configured as module mounts.
			v.SetDefault("contentDir", "content")
			v.SetDefault("layoutDir", "layouts")
			v.SetDefault("assetDir", "assets")
			v.SetDefault("staticDir", "static")
			v.SetDefault("dataDir", "data")
			v.SetDefault("i18nDir", "i18n")
			v.SetDefault("archetypeDir", "archetypes")
	*/

	v.SetDefault("cleanDestinationDir", false)
	v.SetDefault("watch", false)
	v.SetDefault("resourceDir", "resources")
	v.SetDefault("publishDir", "public")
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
	v.SetDefault("permalinks", make(map[string]string))
	v.SetDefault("sitemap", config.Sitemap{Priority: -1, Filename: "sitemap.xml"})
	v.SetDefault("disableLiveReload", false)
	v.SetDefault("pluralizeListTitles", true)
	v.SetDefault("forceSyncStatic", false)
	v.SetDefault("footnoteAnchorPrefix", "")
	v.SetDefault("footnoteReturnLinkContents", "")
	v.SetDefault("newContentEditor", "")
	v.SetDefault("paginate", 10)
	v.SetDefault("paginatePath", "page")
	v.SetDefault("summaryLength", 70)
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
	v.SetDefault("timeout", "30s")
	v.SetDefault("enableInlineShortcodes", false)

	return nil
}
