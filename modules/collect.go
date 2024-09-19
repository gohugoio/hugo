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

package modules

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bep/debounce"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/paths"

	"github.com/spf13/cast"

	"github.com/gohugoio/hugo/common/maps"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/parser/metadecoders"

	"github.com/gohugoio/hugo/hugofs/files"

	"golang.org/x/mod/module"

	"github.com/gohugoio/hugo/config"
	"github.com/spf13/afero"
)

var ErrNotExist = errors.New("module does not exist")

const vendorModulesFilename = "modules.txt"

func (h *Client) Collect() (ModulesConfig, error) {
	mc, coll := h.collect(true)
	if coll.err != nil {
		return mc, coll.err
	}

	if err := (&mc).setActiveMods(h.logger); err != nil {
		return mc, err
	}

	if h.ccfg.HookBeforeFinalize != nil {
		if err := h.ccfg.HookBeforeFinalize(&mc); err != nil {
			return mc, err
		}
	}

	if err := (&mc).finalize(h.logger); err != nil {
		return mc, err
	}

	return mc, nil
}

func (h *Client) collect(tidy bool) (ModulesConfig, *collector) {
	if h == nil {
		panic("nil client")
	}
	c := &collector{
		Client: h,
	}

	c.collect()
	if c.err != nil {
		return ModulesConfig{}, c
	}

	// https://github.com/gohugoio/hugo/issues/6115
	/*if !c.skipTidy && tidy {
		if err := h.tidy(c.modules, true); err != nil {
			c.err = err
			return ModulesConfig{}, c
		}
	}*/

	var workspaceFilename string
	if h.ccfg.ModuleConfig.Workspace != WorkspaceDisabled {
		workspaceFilename = h.ccfg.ModuleConfig.Workspace
	}

	return ModulesConfig{
		AllModules:          c.modules,
		GoModulesFilename:   c.GoModulesFilename,
		GoWorkspaceFilename: workspaceFilename,
	}, c
}

type ModulesConfig struct {
	// All active modules.
	AllModules Modules

	// Set if this is a Go modules enabled project.
	GoModulesFilename string

	// Set if a Go workspace file is configured.
	GoWorkspaceFilename string
}

func (m ModulesConfig) HasConfigFile() bool {
	for _, mod := range m.AllModules {
		if len(mod.ConfigFilenames()) > 0 {
			return true
		}
	}
	return false
}

func (m *ModulesConfig) setActiveMods(logger loggers.Logger) error {
	for _, mod := range m.AllModules {
		if !mod.Config().HugoVersion.IsValid() {
			logger.Warnf(`Module %q is not compatible with this Hugo version: %s; run "hugo mod graph" for more information.`, mod.Path(), mod.Config().HugoVersion)
		}
	}

	return nil
}

func (m *ModulesConfig) finalize(logger loggers.Logger) error {
	for _, mod := range m.AllModules {
		m := mod.(*moduleAdapter)
		m.mounts = filterUnwantedMounts(m.mounts)
	}
	return nil
}

func filterUnwantedMounts(mounts []Mount) []Mount {
	// Remove duplicates
	seen := make(map[string]bool)
	tmp := mounts[:0]
	for _, m := range mounts {
		if !seen[m.key()] {
			tmp = append(tmp, m)
		}
		seen[m.key()] = true
	}
	return tmp
}

type collected struct {
	// Pick the first and prevent circular loops.
	seen map[string]bool

	// Maps module path to a _vendor dir. These values are fetched from
	// _vendor/modules.txt, and the first (top-most) will win.
	vendored map[string]vendoredModule

	// Set if a Go modules enabled project.
	gomods goModules

	// Ordered list of collected modules, including Go Modules and theme
	// components stored below /themes.
	modules Modules
}

// Collects and creates a module tree.
type collector struct {
	*Client

	// Store away any non-fatal error and return at the end.
	err error

	// Set to disable any Tidy operation in the end.
	skipTidy bool

	*collected
}

func (c *collector) initModules() error {
	c.collected = &collected{
		seen:     make(map[string]bool),
		vendored: make(map[string]vendoredModule),
		gomods:   goModules{},
	}

	// If both these are true, we don't even need Go installed to build.
	if c.ccfg.IgnoreVendor == nil && c.isVendored(c.ccfg.WorkingDir) {
		return nil
	}

	// We may fail later if we don't find the mods.
	return c.loadModules()
}

func (c *collector) isSeen(path string) bool {
	key := pathKey(path)
	if c.seen[key] {
		return true
	}
	c.seen[key] = true
	return false
}

func (c *collector) getVendoredDir(path string) (vendoredModule, bool) {
	v, found := c.vendored[path]
	return v, found
}

func (c *collector) add(owner *moduleAdapter, moduleImport Import) (*moduleAdapter, error) {
	var (
		mod       *goModule
		moduleDir string
		version   string
		vendored  bool
	)

	modulePath := moduleImport.Path
	var realOwner Module = owner

	if !c.ccfg.shouldIgnoreVendor(modulePath) {
		if err := c.collectModulesTXT(owner); err != nil {
			return nil, err
		}

		// Try _vendor first.
		var vm vendoredModule
		vm, vendored = c.getVendoredDir(modulePath)
		if vendored {
			moduleDir = vm.Dir
			realOwner = vm.Owner
			version = vm.Version

			if owner.projectMod {
				// We want to keep the go.mod intact with the versions and all.
				c.skipTidy = true
			}

		}
	}

	if moduleDir == "" {
		var versionQuery string
		mod = c.gomods.GetByPath(modulePath)
		if mod != nil {
			moduleDir = mod.Dir
			versionQuery = mod.Version
		}

		if moduleDir == "" {
			if c.GoModulesFilename != "" && isProbablyModule(modulePath) {
				// Try to "go get" it and reload the module configuration.
				if versionQuery == "" {
					// See https://golang.org/ref/mod#version-queries
					// This will select the latest release-version (not beta etc.).
					versionQuery = "upgrade"
				}

				// Note that we cannot use c.Get for this, as that may
				// trigger a new module collection and potentially create a infinite loop.
				if err := c.get(fmt.Sprintf("%s@%s", modulePath, versionQuery)); err != nil {
					return nil, err
				}
				if err := c.loadModules(); err != nil {
					return nil, err
				}

				mod = c.gomods.GetByPath(modulePath)
				if mod != nil {
					moduleDir = mod.Dir
				}
			}

			// Fall back to project/themes/<mymodule>
			if moduleDir == "" {
				var err error
				moduleDir, err = c.createThemeDirname(modulePath, owner.projectMod || moduleImport.pathProjectReplaced)
				if err != nil {
					c.err = err
					return nil, nil
				}
				if found, _ := afero.Exists(c.fs, moduleDir); !found {
					//lint:ignore ST1005 end user message.
					c.err = c.wrapModuleNotFound(fmt.Errorf(`module %q not found in %q; either add it as a Hugo Module or store it in %q.`, modulePath, moduleDir, c.ccfg.ThemesDir))
					return nil, nil
				}
			}
		}
	}

	if found, _ := afero.Exists(c.fs, moduleDir); !found {
		c.err = c.wrapModuleNotFound(fmt.Errorf("%q not found", moduleDir))
		return nil, nil
	}

	if !strings.HasSuffix(moduleDir, fileSeparator) {
		moduleDir += fileSeparator
	}

	ma := &moduleAdapter{
		dir:     moduleDir,
		vendor:  vendored,
		gomod:   mod,
		version: version,
		// This may be the owner of the _vendor dir
		owner: realOwner,
	}

	if mod == nil {
		ma.path = modulePath
	}

	if !moduleImport.IgnoreConfig {
		if err := c.applyThemeConfig(ma); err != nil {
			return nil, err
		}
	}

	if err := c.applyMounts(moduleImport, ma); err != nil {
		return nil, err
	}

	c.modules = append(c.modules, ma)
	return ma, nil
}

func (c *collector) addAndRecurse(owner *moduleAdapter) error {
	moduleConfig := owner.Config()
	if owner.projectMod {
		if err := c.applyMounts(Import{}, owner); err != nil {
			return fmt.Errorf("failed to apply mounts for project: %w", err)
		}
	}

	for _, moduleImport := range moduleConfig.Imports {
		if moduleImport.Disable {
			continue
		}
		if !c.isSeen(moduleImport.Path) {
			tc, err := c.add(owner, moduleImport)
			if err != nil {
				return err
			}
			if tc == nil || moduleImport.IgnoreImports {
				continue
			}
			if err := c.addAndRecurse(tc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *collector) applyMounts(moduleImport Import, mod *moduleAdapter) error {
	if moduleImport.NoMounts {
		mod.mounts = nil
		return nil
	}

	mounts := moduleImport.Mounts

	modConfig := mod.Config()

	if len(mounts) == 0 {
		// Mounts not defined by the import.
		mounts = modConfig.Mounts
	}

	if !mod.projectMod && len(mounts) == 0 {
		// Create default mount points for every component folder that
		// exists in the module.
		for _, componentFolder := range files.ComponentFolders {
			sourceDir := filepath.Join(mod.Dir(), componentFolder)
			_, err := c.fs.Stat(sourceDir)
			if err == nil {
				mounts = append(mounts, Mount{
					Source: componentFolder,
					Target: componentFolder,
				})
			}
		}
	}

	var err error
	mounts, err = c.normalizeMounts(mod, mounts)
	if err != nil {
		return err
	}

	mounts, err = c.mountCommonJSConfig(mod, mounts)
	if err != nil {
		return err
	}

	mod.mounts = mounts
	return nil
}

func (c *collector) applyThemeConfig(tc *moduleAdapter) error {
	var (
		configFilename string
		themeCfg       map[string]any
		hasConfigFile  bool
		err            error
	)

LOOP:
	for _, configBaseName := range config.DefaultConfigNames {
		for _, configFormats := range config.ValidConfigFileExtensions {
			configFilename = filepath.Join(tc.Dir(), configBaseName+"."+configFormats)
			hasConfigFile, _ = afero.Exists(c.fs, configFilename)
			if hasConfigFile {
				break LOOP
			}
		}
	}

	// The old theme information file.
	themeTOML := filepath.Join(tc.Dir(), "theme.toml")

	hasThemeTOML, _ := afero.Exists(c.fs, themeTOML)
	if hasThemeTOML {
		data, err := afero.ReadFile(c.fs, themeTOML)
		if err != nil {
			return err
		}
		themeCfg, err = metadecoders.Default.UnmarshalToMap(data, metadecoders.TOML)
		if err != nil {
			c.logger.Warnf("Failed to read module config for %q in %q: %s", tc.Path(), themeTOML, err)
		} else {
			maps.PrepareParams(themeCfg)
		}
	}

	if hasConfigFile {
		if configFilename != "" {
			var err error
			tc.cfg, err = config.FromFile(c.fs, configFilename)
			if err != nil {
				return err
			}
		}

		tc.configFilenames = append(tc.configFilenames, configFilename)

	}

	// Also check for a config dir, which we overlay on top of the file configuration.
	configDir := filepath.Join(tc.Dir(), "config")
	dcfg, dirnames, err := config.LoadConfigFromDir(c.fs, configDir, c.ccfg.Environment)
	if err != nil {
		return err
	}

	if len(dirnames) > 0 {
		tc.configFilenames = append(tc.configFilenames, dirnames...)

		if hasConfigFile {
			// Set will overwrite existing keys.
			tc.cfg.Set("", dcfg.Get(""))
		} else {
			tc.cfg = dcfg
		}
	}

	config, err := decodeConfig(tc.cfg, c.moduleConfig.replacementsMap)
	if err != nil {
		return err
	}

	const oldVersionKey = "min_version"

	if hasThemeTOML {

		// Merge old with new
		if minVersion, found := themeCfg[oldVersionKey]; found {
			if config.HugoVersion.Min == "" {
				config.HugoVersion.Min = hugo.VersionString(cast.ToString(minVersion))
			}
		}

		if config.Params == nil {
			config.Params = make(map[string]any)
		}

		for k, v := range themeCfg {
			if k == oldVersionKey {
				continue
			}
			config.Params[k] = v
		}

	}

	tc.config = config

	return nil
}

func (c *collector) collect() {
	defer c.logger.PrintTimerIfDelayed(time.Now(), "hugo: collected modules")
	d := debounce.New(2 * time.Second)
	d(func() {
		c.logger.Println("hugo: downloading modules …")
	})
	defer d(func() {})

	if err := c.initModules(); err != nil {
		c.err = err
		return
	}

	projectMod := createProjectModule(c.gomods.GetMain(), c.ccfg.WorkingDir, c.moduleConfig)

	if err := c.addAndRecurse(projectMod); err != nil {
		c.err = err
		return
	}

	// Add the project mod on top.
	c.modules = append(Modules{projectMod}, c.modules...)
}

func (c *collector) isVendored(dir string) bool {
	_, err := c.fs.Stat(filepath.Join(dir, vendord, vendorModulesFilename))
	return err == nil
}

func (c *collector) collectModulesTXT(owner Module) error {
	vendorDir := filepath.Join(owner.Dir(), vendord)
	filename := filepath.Join(vendorDir, vendorModulesFilename)

	f, err := c.fs.Open(filename)
	if err != nil {
		if herrors.IsNotExist(err) {
			return nil
		}

		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// # github.com/alecthomas/chroma v0.6.3
		line := scanner.Text()
		line = strings.Trim(line, "# ")
		line = strings.TrimSpace(line)
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return fmt.Errorf("invalid modules list: %q", filename)
		}
		path := parts[0]

		shouldAdd := c.Client.moduleConfig.VendorClosest

		if !shouldAdd {
			if _, found := c.vendored[path]; !found {
				shouldAdd = true
			}
		}

		if shouldAdd {
			c.vendored[path] = vendoredModule{
				Owner:   owner,
				Dir:     filepath.Join(vendorDir, path),
				Version: parts[1],
			}
		}

	}
	return nil
}

func (c *collector) loadModules() error {
	modules, err := c.listGoMods()
	if err != nil {
		return err
	}
	c.gomods = modules
	return nil
}

// Matches postcss.config.js etc.
var commonJSConfigs = regexp.MustCompile(`(babel|postcss|tailwind)\.config\.js`)

func (c *collector) mountCommonJSConfig(owner *moduleAdapter, mounts []Mount) ([]Mount, error) {
	for _, m := range mounts {
		if strings.HasPrefix(m.Target, files.JsConfigFolderMountPrefix) {
			// This follows the convention of the other component types (assets, content, etc.),
			// if one or more is specified by the user, we skip the defaults.
			// These mounts were added to Hugo in 0.75.
			return mounts, nil
		}
	}

	// Mount the common JS config files.
	d, err := c.fs.Open(owner.Dir())
	if err != nil {
		return mounts, fmt.Errorf("failed to open dir %q: %q", owner.Dir(), err)
	}
	defer d.Close()
	fis, err := d.(fs.ReadDirFile).ReadDir(-1)
	if err != nil {
		return mounts, fmt.Errorf("failed to read dir %q: %q", owner.Dir(), err)
	}

	for _, fi := range fis {
		n := fi.Name()

		should := n == files.FilenamePackageHugoJSON || n == files.FilenamePackageJSON
		should = should || commonJSConfigs.MatchString(n)

		if should {
			mounts = append(mounts, Mount{
				Source: n,
				Target: filepath.Join(files.ComponentFolderAssets, files.FolderJSConfig, n),
			})
		}

	}

	return mounts, nil
}

func (c *collector) normalizeMounts(owner *moduleAdapter, mounts []Mount) ([]Mount, error) {
	var out []Mount
	dir := owner.Dir()

	for _, mnt := range mounts {
		errMsg := fmt.Sprintf("invalid module config for %q", owner.Path())

		if mnt.Source == "" || mnt.Target == "" {
			return nil, errors.New(errMsg + ": both source and target must be set")
		}

		mnt.Source = filepath.Clean(mnt.Source)
		mnt.Target = filepath.Clean(mnt.Target)
		var sourceDir string

		if owner.projectMod && filepath.IsAbs(mnt.Source) {
			// Abs paths in the main project is allowed.
			sourceDir = mnt.Source
		} else {
			sourceDir = filepath.Join(dir, mnt.Source)
		}

		// Verify that Source exists
		_, err := c.fs.Stat(sourceDir)
		if err != nil {
			if paths.IsSameFilePath(sourceDir, c.ccfg.PublishDir) {
				// This is a little exotic, but there are use cases for mounting the public folder.
				// This will typically also be in .gitingore, so create it.
				if err := c.fs.MkdirAll(sourceDir, 0o755); err != nil {
					return nil, fmt.Errorf("%s: %q", errMsg, err)
				}
			} else if strings.HasSuffix(sourceDir, files.FilenameHugoStatsJSON) {
				// A common pattern for Tailwind 3 is to mount that file to get it on the server watch list.

				// A common pattern is also to add hugo_stats.json to .gitignore.

				// Create an empty file.
				f, err := c.fs.Create(sourceDir)
				if err != nil {
					return nil, fmt.Errorf("%s: %q", errMsg, err)
				}
				f.Close()
			} else {
				// TODO(bep) commenting out for now, as this will create to much noise.
				// c.logger.Warnf("module %q: mount source %q does not exist", owner.Path(), sourceDir)
				continue
			}
		}

		// Verify that target points to one of the predefined component dirs
		targetBase := mnt.Target
		idxPathSep := strings.Index(mnt.Target, string(os.PathSeparator))
		if idxPathSep != -1 {
			targetBase = mnt.Target[0:idxPathSep]
		}
		if !files.IsComponentFolder(targetBase) {
			return nil, fmt.Errorf("%s: mount target must be one of: %v", errMsg, files.ComponentFolders)
		}

		out = append(out, mnt)
	}

	return out, nil
}

func (c *collector) wrapModuleNotFound(err error) error {
	if c.Client.ccfg.IgnoreModuleDoesNotExist {
		return nil
	}
	err = fmt.Errorf(err.Error()+": %w", ErrNotExist)
	if c.GoModulesFilename == "" {
		return err
	}

	baseMsg := "we found a go.mod file in your project, but"

	switch c.goBinaryStatus {
	case goBinaryStatusNotFound:
		return fmt.Errorf(baseMsg+" you need to install Go to use it. See https://golang.org/dl/ : %q", err)
	case goBinaryStatusTooOld:
		return fmt.Errorf(baseMsg+" you need to a newer version of Go to use it. See https://golang.org/dl/ : %w", err)
	}

	return err
}

type vendoredModule struct {
	Owner   Module
	Dir     string
	Version string
}

func createProjectModule(gomod *goModule, workingDir string, conf Config) *moduleAdapter {
	// Create a pseudo module for the main project.
	var path string
	if gomod == nil {
		path = "project"
	}

	return &moduleAdapter{
		path:       path,
		dir:        workingDir,
		gomod:      gomod,
		projectMod: true,
		config:     conf,
	}
}

// In the first iteration of Hugo Modules, we do not support multiple
// major versions running at the same time, so we pick the first (upper most).
// We will investigate namespaces in future versions.
// TODO(bep) add a warning when the above happens.
func pathKey(p string) string {
	prefix, _, _ := module.SplitPathVersion(p)

	return strings.ToLower(prefix)
}
