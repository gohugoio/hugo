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

	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/types"

	"github.com/gohugoio/hugo/common/maps"
	cpaths "github.com/gohugoio/hugo/common/paths"

	"github.com/gobwas/glob"
	hglob "github.com/gohugoio/hugo/hugofs/glob"

	"github.com/gohugoio/hugo/common/loggers"

	"github.com/gohugoio/hugo/cache/filecache"

	"github.com/gohugoio/hugo/parser/metadecoders"

	"errors"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/langs"
	"github.com/gohugoio/hugo/modules"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/config/privacy"
	"github.com/gohugoio/hugo/config/security"
	"github.com/gohugoio/hugo/config/services"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/afero"
)

var ErrNoConfigFile = errors.New("Unable to locate config file or config directory. Perhaps you need to create a new site.\n       Run `hugo help new` for details.\n")

// LoadConfig loads Hugo configuration into a new Viper and then adds
// a set of defaults.
func LoadConfig(d ConfigSourceDescriptor, doWithConfig ...func(cfg config.Provider) error) (config.Provider, []string, error) {
	if d.Environment == "" {
		d.Environment = hugo.EnvironmentProduction
	}

	if len(d.Environ) == 0 && !hugo.IsRunningAsTest() {
		d.Environ = os.Environ()
	}

	var configFiles []string

	l := configLoader{ConfigSourceDescriptor: d, cfg: config.New()}
	// Make sure we always do this, even in error situations,
	// as we have commands (e.g. "hugo mod init") that will
	// use a partial configuration to do its job.
	defer l.deleteMergeStrategies()

	names := d.configFilenames()

	if names != nil {
		for _, name := range names {
			var filename string
			filename, err := l.loadConfig(name)
			if err == nil {
				configFiles = append(configFiles, filename)
			} else if err != ErrNoConfigFile {
				return nil, nil, l.wrapFileError(err, filename)
			}
		}
	} else {
		for _, name := range config.DefaultConfigNames {
			var filename string
			filename, err := l.loadConfig(name)
			if err == nil {
				configFiles = append(configFiles, filename)
				break
			} else if err != ErrNoConfigFile {
				return nil, nil, l.wrapFileError(err, filename)
			}
		}
	}

	if d.AbsConfigDir != "" {

		dcfg, dirnames, err := config.LoadConfigFromDir(l.Fs, d.AbsConfigDir, l.Environment)

		if err == nil {
			if len(dirnames) > 0 {
				l.cfg.Set("", dcfg.Get(""))
				configFiles = append(configFiles, dirnames...)
			}
		} else if err != ErrNoConfigFile {
			if len(dirnames) > 0 {
				return nil, nil, l.wrapFileError(err, dirnames[0])
			}
			return nil, nil, err
		}
	}

	if err := l.applyConfigDefaults(); err != nil {
		return l.cfg, configFiles, err
	}

	l.cfg.SetDefaultMergeStrategy()

	// We create languages based on the settings, so we need to make sure that
	// all configuration is loaded/set before doing that.
	for _, d := range doWithConfig {
		if err := d(l.cfg); err != nil {
			return l.cfg, configFiles, err
		}
	}

	// Some settings are used before we're done collecting all settings,
	// so apply OS environment both before and after.
	if err := l.applyOsEnvOverrides(d.Environ); err != nil {
		return l.cfg, configFiles, err
	}

	modulesConfig, err := l.loadModulesConfig()
	if err != nil {
		return l.cfg, configFiles, err
	}

	// Need to run these after the modules are loaded, but before
	// they are finalized.
	collectHook := func(m *modules.ModulesConfig) error {
		// We don't need the merge strategy configuration anymore,
		// remove it so it doesn't accidentally show up in other settings.
		l.deleteMergeStrategies()

		if err := l.loadLanguageSettings(nil); err != nil {
			return err
		}

		mods := m.ActiveModules

		// Apply default project mounts.
		if err := modules.ApplyProjectConfigDefaults(l.cfg, mods[0]); err != nil {
			return err
		}

		return nil
	}

	_, modulesConfigFiles, modulesCollectErr := l.collectModules(modulesConfig, l.cfg, collectHook)
	if err != nil {
		return l.cfg, configFiles, err
	}

	configFiles = append(configFiles, modulesConfigFiles...)

	if err := l.applyOsEnvOverrides(d.Environ); err != nil {
		return l.cfg, configFiles, err
	}

	if err = l.applyConfigAliases(); err != nil {
		return l.cfg, configFiles, err
	}

	if err == nil {
		err = modulesCollectErr
	}

	return l.cfg, configFiles, err
}

// LoadConfigDefault is a convenience method to load the default "hugo.toml" config.
func LoadConfigDefault(fs afero.Fs) (config.Provider, error) {
	v, _, err := LoadConfig(ConfigSourceDescriptor{Fs: fs})
	return v, err
}

// ConfigSourceDescriptor describes where to find the config (e.g. config.toml etc.).
type ConfigSourceDescriptor struct {
	Fs     afero.Fs
	Logger loggers.Logger

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

func (d ConfigSourceDescriptor) configFileDir() string {
	if d.Path != "" {
		return d.Path
	}
	return d.WorkingDir
}

func (d ConfigSourceDescriptor) configFilenames() []string {
	if d.Filename == "" {
		return nil
	}
	return strings.Split(d.Filename, ",")
}

// SiteConfig represents the config in .Site.Config.
type SiteConfig struct {
	// This contains all privacy related settings that can be used to
	// make the YouTube template etc. GDPR compliant.
	Privacy privacy.Config

	// Services contains config for services such as Google Analytics etc.
	Services services.Config
}

type configLoader struct {
	cfg config.Provider
	ConfigSourceDescriptor
}

// Handle some legacy values.
func (l configLoader) applyConfigAliases() error {
	aliases := []types.KeyValueStr{{Key: "taxonomies", Value: "indexes"}}

	for _, alias := range aliases {
		if l.cfg.IsSet(alias.Key) {
			vv := l.cfg.Get(alias.Key)
			l.cfg.Set(alias.Value, vv)
		}
	}

	return nil
}

func (l configLoader) applyConfigDefaults() error {
	defaultSettings := maps.Params{
		"cleanDestinationDir":                  false,
		"watch":                                false,
		"resourceDir":                          "resources",
		"publishDir":                           "public",
		"publishDirOrig":                       "public",
		"themesDir":                            "themes",
		"buildDrafts":                          false,
		"buildFuture":                          false,
		"buildExpired":                         false,
		"environment":                          hugo.EnvironmentProduction,
		"uglyURLs":                             false,
		"verbose":                              false,
		"ignoreCache":                          false,
		"canonifyURLs":                         false,
		"relativeURLs":                         false,
		"removePathAccents":                    false,
		"titleCaseStyle":                       "AP",
		"taxonomies":                           maps.Params{"tag": "tags", "category": "categories"},
		"permalinks":                           maps.Params{},
		"sitemap":                              maps.Params{"priority": -1, "filename": "sitemap.xml"},
		"disableLiveReload":                    false,
		"pluralizeListTitles":                  true,
		"forceSyncStatic":                      false,
		"footnoteAnchorPrefix":                 "",
		"footnoteReturnLinkContents":           "",
		"newContentEditor":                     "",
		"paginate":                             10,
		"paginatePath":                         "page",
		"summaryLength":                        70,
		"rssLimit":                             -1,
		"sectionPagesMenu":                     "",
		"disablePathToLower":                   false,
		"hasCJKLanguage":                       false,
		"enableEmoji":                          false,
		"defaultContentLanguage":               "en",
		"defaultContentLanguageInSubdir":       false,
		"enableMissingTranslationPlaceholders": false,
		"enableGitInfo":                        false,
		"ignoreFiles":                          make([]string, 0),
		"disableAliases":                       false,
		"debug":                                false,
		"disableFastRender":                    false,
		"timeout":                              "30s",
		"enableInlineShortcodes":               false,
	}

	l.cfg.SetDefaults(defaultSettings)

	return nil
}

func (l configLoader) applyOsEnvOverrides(environ []string) error {
	if len(environ) == 0 {
		return nil
	}

	const delim = "__env__delim"

	// Extract all that start with the HUGO prefix.
	// The delimiter is the following rune, usually "_".
	const hugoEnvPrefix = "HUGO"
	var hugoEnv []types.KeyValueStr
	for _, v := range environ {
		key, val := config.SplitEnvVar(v)
		if strings.HasPrefix(key, hugoEnvPrefix) {
			delimiterAndKey := strings.TrimPrefix(key, hugoEnvPrefix)
			if len(delimiterAndKey) < 2 {
				continue
			}
			// Allow delimiters to be case sensitive.
			// It turns out there isn't that many allowed special
			// chars in environment variables when used in Bash and similar,
			// so variables on the form HUGOxPARAMSxFOO=bar is one option.
			key := strings.ReplaceAll(delimiterAndKey[1:], delimiterAndKey[:1], delim)
			key = strings.ToLower(key)
			hugoEnv = append(hugoEnv, types.KeyValueStr{
				Key:   key,
				Value: val,
			})

		}
	}

	for _, env := range hugoEnv {
		existing, nestedKey, owner, err := maps.GetNestedParamFn(env.Key, delim, l.cfg.Get)
		if err != nil {
			return err
		}

		if existing != nil {
			val, err := metadecoders.Default.UnmarshalStringTo(env.Value, existing)
			if err != nil {
				continue
			}

			if owner != nil {
				owner[nestedKey] = val
			} else {
				l.cfg.Set(env.Key, val)
			}
		} else if nestedKey != "" {
			owner[nestedKey] = env.Value
		} else {
			// The container does not exist yet.
			l.cfg.Set(strings.ReplaceAll(env.Key, delim, "."), env.Value)
		}
	}

	return nil
}

func (l configLoader) collectModules(modConfig modules.Config, v1 config.Provider, hookBeforeFinalize func(m *modules.ModulesConfig) error) (modules.Modules, []string, error) {
	workingDir := l.WorkingDir
	if workingDir == "" {
		workingDir = v1.GetString("workingDir")
	}

	themesDir := cpaths.AbsPathify(l.WorkingDir, v1.GetString("themesDir"))

	var ignoreVendor glob.Glob
	if s := v1.GetString("ignoreVendorPaths"); s != "" {
		ignoreVendor, _ = hglob.GetGlob(hglob.NormalizePath(s))
	}

	filecacheConfigs, err := filecache.DecodeConfig(l.Fs, v1)
	if err != nil {
		return nil, nil, err
	}

	secConfig, err := security.DecodeConfig(v1)
	if err != nil {
		return nil, nil, err
	}
	ex := hexec.New(secConfig)

	v1.Set("filecacheConfigs", filecacheConfigs)

	var configFilenames []string

	hook := func(m *modules.ModulesConfig) error {
		for _, tc := range m.ActiveModules {
			if len(tc.ConfigFilenames()) > 0 {
				if tc.Watch() {
					configFilenames = append(configFilenames, tc.ConfigFilenames()...)
				}

				// Merge from theme config into v1 based on configured
				// merge strategy.
				v1.Merge("", tc.Cfg().Get(""))

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
		Exec:               ex,
		HookBeforeFinalize: hook,
		WorkingDir:         workingDir,
		ThemesDir:          themesDir,
		Environment:        l.Environment,
		CacheDir:           filecacheConfigs.CacheDirModules(),
		ModuleConfig:       modConfig,
		IgnoreVendor:       ignoreVendor,
	})

	v1.Set("modulesClient", modulesClient)

	moduleConfig, err := modulesClient.Collect()

	// Avoid recreating these later.
	v1.Set("allModules", moduleConfig.ActiveModules)

	// We want to watch these for changes and trigger rebuild on version
	// changes etc.
	if moduleConfig.GoModulesFilename != "" {

		configFilenames = append(configFilenames, moduleConfig.GoModulesFilename)
	}

	if moduleConfig.GoWorkspaceFilename != "" {
		configFilenames = append(configFilenames, moduleConfig.GoWorkspaceFilename)

	}

	return moduleConfig.ActiveModules, configFilenames, err
}

func (l configLoader) loadConfig(configName string) (string, error) {
	baseDir := l.configFileDir()
	var baseFilename string
	if filepath.IsAbs(configName) {
		baseFilename = configName
	} else {
		baseFilename = filepath.Join(baseDir, configName)
	}

	var filename string
	if cpaths.ExtNoDelimiter(configName) != "" {
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
		return filename, err
	}

	// Set overwrites keys of the same name, recursively.
	l.cfg.Set("", m)

	return filename, nil
}

func (l configLoader) deleteMergeStrategies() {
	l.cfg.WalkParams(func(params ...config.KeyParams) bool {
		params[len(params)-1].Params.DeleteMergeStrategy()
		return false
	})
}

func (l configLoader) loadLanguageSettings(oldLangs langs.Languages) error {
	_, err := langs.LoadLanguageSettings(l.cfg, oldLangs)
	return err
}

func (l configLoader) loadModulesConfig() (modules.Config, error) {
	modConfig, err := modules.DecodeConfig(l.cfg)
	if err != nil {
		return modules.Config{}, err
	}

	return modConfig, nil
}

func (configLoader) loadSiteConfig(cfg config.Provider) (scfg SiteConfig, err error) {
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

func (l configLoader) wrapFileError(err error, filename string) error {
	fe := herrors.UnwrapFileError(err)
	if fe != nil {
		pos := fe.Position()
		pos.Filename = filename
		fe.UpdatePosition(pos)
		return err
	}
	return herrors.NewFileErrorFromFile(err, filename, l.Fs, nil)
}
