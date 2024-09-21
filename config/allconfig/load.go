// Copyright 2024 The Hugo Authors. All rights reserved.
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

// Package allconfig contains the full configuration for Hugo.
package allconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	hglob "github.com/gohugoio/hugo/hugofs/glob"
	"github.com/gohugoio/hugo/modules"
	"github.com/gohugoio/hugo/parser/metadecoders"
	"github.com/spf13/afero"
)

//lint:ignore ST1005 end user message.
var ErrNoConfigFile = errors.New("Unable to locate config file or config directory. Perhaps you need to create a new site.\n       Run `hugo help new` for details.\n")

func LoadConfig(d ConfigSourceDescriptor) (*Configs, error) {
	if len(d.Environ) == 0 && !hugo.IsRunningAsTest() {
		d.Environ = os.Environ()
	}

	if d.Logger == nil {
		d.Logger = loggers.NewDefault()
	}

	l := &configLoader{ConfigSourceDescriptor: d, cfg: config.New()}
	// Make sure we always do this, even in error situations,
	// as we have commands (e.g. "hugo mod init") that will
	// use a partial configuration to do its job.
	defer l.deleteMergeStrategies()
	res, _, err := l.loadConfigMain(d)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	configs, err := fromLoadConfigResult(d.Fs, d.Logger, res)
	if err != nil {
		return nil, fmt.Errorf("failed to create config from result: %w", err)
	}

	moduleConfig, modulesClient, err := l.loadModules(configs, d.IgnoreModuleDoesNotExist)
	if err != nil {
		return nil, fmt.Errorf("failed to load modules: %w", err)
	}

	if len(l.ModulesConfigFiles) > 0 {
		// Config merged in from modules.
		// Re-read the config.
		configs, err = fromLoadConfigResult(d.Fs, d.Logger, res)
		if err != nil {
			return nil, fmt.Errorf("failed to create config from modules config: %w", err)
		}
		if err := configs.transientErr(); err != nil {
			return nil, fmt.Errorf("failed to create config from modules config: %w", err)
		}
		configs.LoadingInfo.ConfigFiles = append(configs.LoadingInfo.ConfigFiles, l.ModulesConfigFiles...)
	} else if err := configs.transientErr(); err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	configs.Modules = moduleConfig.AllModules
	configs.ModulesClient = modulesClient

	if err := configs.Init(); err != nil {
		return nil, fmt.Errorf("failed to init config: %w", err)
	}

	loggers.InitGlobalLogger(d.Logger.Level(), configs.Base.PanicOnWarning)

	return configs, nil
}

// ConfigSourceDescriptor describes where to find the config (e.g. config.toml etc.).
type ConfigSourceDescriptor struct {
	Fs     afero.Fs
	Logger loggers.Logger

	// Config received from the command line.
	// These will override any config file settings.
	Flags config.Provider

	// Path to the config file to use, e.g. /my/project/config.toml
	Filename string

	// The (optional) directory for additional configuration files.
	ConfigDir string

	// production, development
	Environment string

	// Defaults to os.Environ if not set.
	Environ []string

	// If set, this will be used to ignore the module does not exist error.
	IgnoreModuleDoesNotExist bool
}

func (d ConfigSourceDescriptor) configFilenames() []string {
	if d.Filename == "" {
		return nil
	}
	return strings.Split(d.Filename, ",")
}

type configLoader struct {
	cfg        config.Provider
	BaseConfig config.BaseConfig
	ConfigSourceDescriptor

	// collected
	ModulesConfig      modules.ModulesConfig
	ModulesConfigFiles []string
}

// Handle some legacy values.
func (l configLoader) applyConfigAliases() error {
	aliases := []types.KeyValueStr{
		{Key: "indexes", Value: "taxonomies"},
		{Key: "logI18nWarnings", Value: "printI18nWarnings"},
		{Key: "logPathWarnings", Value: "printPathWarnings"},
		{Key: "ignoreErrors", Value: "ignoreLogs"},
	}

	for _, alias := range aliases {
		if l.cfg.IsSet(alias.Key) {
			vv := l.cfg.Get(alias.Key)
			l.cfg.Set(alias.Value, vv)
		}
	}

	return nil
}

func (l configLoader) applyDefaultConfig() error {
	defaultSettings := maps.Params{
		"baseURL":                              "",
		"cleanDestinationDir":                  false,
		"watch":                                false,
		"contentDir":                           "content",
		"resourceDir":                          "resources",
		"publishDir":                           "public",
		"publishDirOrig":                       "public",
		"themesDir":                            "themes",
		"assetDir":                             "assets",
		"layoutDir":                            "layouts",
		"i18nDir":                              "i18n",
		"dataDir":                              "data",
		"archetypeDir":                         "archetypes",
		"configDir":                            "config",
		"staticDir":                            "static",
		"buildDrafts":                          false,
		"buildFuture":                          false,
		"buildExpired":                         false,
		"params":                               maps.Params{},
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
		"menus":                                maps.Params{},
		"disableLiveReload":                    false,
		"pluralizeListTitles":                  true,
		"capitalizeListTitles":                 true,
		"forceSyncStatic":                      false,
		"footnoteAnchorPrefix":                 "",
		"footnoteReturnLinkContents":           "",
		"newContentEditor":                     "",
		"paginate":                             0,  // Moved into the paginator struct in Hugo v0.128.0.
		"paginatePath":                         "", // Moved into the paginator struct in Hugo v0.128.0.
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
		"timeZone":                             "",
		"enableInlineShortcodes":               false,
	}

	l.cfg.SetDefaults(defaultSettings)

	return nil
}

func (l configLoader) normalizeCfg(cfg config.Provider) error {
	if b, ok := cfg.Get("minifyOutput").(bool); ok && b {
		cfg.Set("minify.minifyOutput", true)
	} else if b, ok := cfg.Get("minify").(bool); ok && b {
		cfg.Set("minify", maps.Params{"minifyOutput": true})
	}

	return nil
}

func (l configLoader) cleanExternalConfig(cfg config.Provider) error {
	if cfg.IsSet("internal") {
		cfg.Set("internal", nil)
	}
	return nil
}

func (l configLoader) applyFlagsOverrides(cfg config.Provider) error {
	for _, k := range cfg.Keys() {
		l.cfg.Set(k, cfg.Get(k))
	}
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
		} else {
			if nestedKey != "" {
				owner[nestedKey] = env.Value
			} else {
				var val any
				key := strings.ReplaceAll(env.Key, delim, ".")
				_, ok := allDecoderSetups[key]
				if ok {
					// A map.
					if v, err := metadecoders.Default.UnmarshalStringTo(env.Value, map[string]interface{}{}); err == nil {
						val = v
					}
				}
				if val == nil {
					// A string.
					val = l.envStringToVal(key, env.Value)
				}
				l.cfg.Set(key, val)
			}
		}
	}

	return nil
}

func (l *configLoader) envStringToVal(k, v string) any {
	switch k {
	case "disablekinds", "disablelanguages":
		if strings.Contains(v, ",") {
			return strings.Split(v, ",")
		} else {
			return strings.Fields(v)
		}
	default:
		return v
	}
}

func (l *configLoader) loadConfigMain(d ConfigSourceDescriptor) (config.LoadConfigResult, modules.ModulesConfig, error) {
	var res config.LoadConfigResult

	if d.Flags != nil {
		if err := l.normalizeCfg(d.Flags); err != nil {
			return res, l.ModulesConfig, err
		}
	}

	if d.Fs == nil {
		return res, l.ModulesConfig, errors.New("no filesystem provided")
	}

	if d.Flags != nil {
		if err := l.applyFlagsOverrides(d.Flags); err != nil {
			return res, l.ModulesConfig, err
		}
		workingDir := filepath.Clean(l.cfg.GetString("workingDir"))

		l.BaseConfig = config.BaseConfig{
			WorkingDir: workingDir,
			ThemesDir:  paths.AbsPathify(workingDir, l.cfg.GetString("themesDir")),
		}

	}

	names := d.configFilenames()

	if names != nil {
		for _, name := range names {
			var filename string
			filename, err := l.loadConfig(name)
			if err == nil {
				res.ConfigFiles = append(res.ConfigFiles, filename)
			} else if err != ErrNoConfigFile {
				return res, l.ModulesConfig, l.wrapFileError(err, filename)
			}
		}
	} else {
		for _, name := range config.DefaultConfigNames {
			var filename string
			filename, err := l.loadConfig(name)
			if err == nil {
				res.ConfigFiles = append(res.ConfigFiles, filename)
				break
			} else if err != ErrNoConfigFile {
				return res, l.ModulesConfig, l.wrapFileError(err, filename)
			}
		}
	}

	if d.ConfigDir != "" {
		absConfigDir := paths.AbsPathify(l.BaseConfig.WorkingDir, d.ConfigDir)
		dcfg, dirnames, err := config.LoadConfigFromDir(l.Fs, absConfigDir, l.Environment)
		if err == nil {
			if len(dirnames) > 0 {
				if err := l.normalizeCfg(dcfg); err != nil {
					return res, l.ModulesConfig, err
				}
				if err := l.cleanExternalConfig(dcfg); err != nil {
					return res, l.ModulesConfig, err
				}
				l.cfg.Set("", dcfg.Get(""))
				res.ConfigFiles = append(res.ConfigFiles, dirnames...)
			}
		} else if err != ErrNoConfigFile {
			if len(dirnames) > 0 {
				return res, l.ModulesConfig, l.wrapFileError(err, dirnames[0])
			}
			return res, l.ModulesConfig, err
		}
	}

	res.Cfg = l.cfg

	if err := l.applyDefaultConfig(); err != nil {
		return res, l.ModulesConfig, err
	}

	// Some settings are used before we're done collecting all settings,
	// so apply OS environment both before and after.
	if err := l.applyOsEnvOverrides(d.Environ); err != nil {
		return res, l.ModulesConfig, err
	}

	workingDir := filepath.Clean(l.cfg.GetString("workingDir"))

	l.BaseConfig = config.BaseConfig{
		WorkingDir: workingDir,
		CacheDir:   l.cfg.GetString("cacheDir"),
		ThemesDir:  paths.AbsPathify(workingDir, l.cfg.GetString("themesDir")),
	}

	var err error
	l.BaseConfig.CacheDir, err = helpers.GetCacheDir(l.Fs, l.BaseConfig.CacheDir)
	if err != nil {
		return res, l.ModulesConfig, err
	}

	res.BaseConfig = l.BaseConfig

	l.cfg.SetDefaultMergeStrategy()

	res.ConfigFiles = append(res.ConfigFiles, l.ModulesConfigFiles...)

	if d.Flags != nil {
		if err := l.applyFlagsOverrides(d.Flags); err != nil {
			return res, l.ModulesConfig, err
		}
	}

	if err := l.applyOsEnvOverrides(d.Environ); err != nil {
		return res, l.ModulesConfig, err
	}

	if err = l.applyConfigAliases(); err != nil {
		return res, l.ModulesConfig, err
	}

	return res, l.ModulesConfig, err
}

func (l *configLoader) loadModules(configs *Configs, ignoreModuleDoesNotExist bool) (modules.ModulesConfig, *modules.Client, error) {
	bcfg := configs.LoadingInfo.BaseConfig
	conf := configs.Base
	workingDir := bcfg.WorkingDir
	themesDir := bcfg.ThemesDir
	publishDir := bcfg.PublishDir

	cfg := configs.LoadingInfo.Cfg

	var ignoreVendor glob.Glob
	if s := conf.IgnoreVendorPaths; s != "" {
		ignoreVendor, _ = hglob.GetGlob(hglob.NormalizePath(s))
	}

	ex := hexec.New(conf.Security, workingDir)

	hook := func(m *modules.ModulesConfig) error {
		for _, tc := range m.AllModules {
			if len(tc.ConfigFilenames()) > 0 {
				if tc.Watch() {
					l.ModulesConfigFiles = append(l.ModulesConfigFiles, tc.ConfigFilenames()...)
				}

				// Merge in the theme config using the configured
				// merge strategy.
				cfg.Merge("", tc.Cfg().Get(""))

			}
		}

		return nil
	}

	modulesClient := modules.NewClient(modules.ClientConfig{
		Fs:                       l.Fs,
		Logger:                   l.Logger,
		Exec:                     ex,
		HookBeforeFinalize:       hook,
		WorkingDir:               workingDir,
		ThemesDir:                themesDir,
		PublishDir:               publishDir,
		Environment:              l.Environment,
		CacheDir:                 conf.Caches.CacheDirModules(),
		ModuleConfig:             conf.Module,
		IgnoreVendor:             ignoreVendor,
		IgnoreModuleDoesNotExist: ignoreModuleDoesNotExist,
	})

	moduleConfig, err := modulesClient.Collect()

	// We want to watch these for changes and trigger rebuild on version
	// changes etc.
	if moduleConfig.GoModulesFilename != "" {
		l.ModulesConfigFiles = append(l.ModulesConfigFiles, moduleConfig.GoModulesFilename)
	}

	if moduleConfig.GoWorkspaceFilename != "" {
		l.ModulesConfigFiles = append(l.ModulesConfigFiles, moduleConfig.GoWorkspaceFilename)
	}

	return moduleConfig, modulesClient, err
}

func (l configLoader) loadConfig(configName string) (string, error) {
	baseDir := l.BaseConfig.WorkingDir
	var baseFilename string
	if filepath.IsAbs(configName) {
		baseFilename = configName
	} else {
		baseFilename = filepath.Join(baseDir, configName)
	}

	var filename string
	if paths.ExtNoDelimiter(configName) != "" {
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

	if err := l.normalizeCfg(l.cfg); err != nil {
		return filename, err
	}

	if err := l.cleanExternalConfig(l.cfg); err != nil {
		return filename, err
	}

	return filename, nil
}

func (l configLoader) deleteMergeStrategies() {
	l.cfg.WalkParams(func(params ...maps.KeyParams) bool {
		params[len(params)-1].Params.DeleteMergeStrategy()
		return false
	})
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
