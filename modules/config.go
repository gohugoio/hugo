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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/config"
	"github.com/mitchellh/mapstructure"
)

const WorkspaceDisabled = "off"

var DefaultModuleConfig = Config{
	// Default to direct, which means "git clone" and similar. We
	// will investigate proxy settings in more depth later.
	// See https://github.com/golang/go/issues/26334
	Proxy: "direct",

	// Comma separated glob list matching paths that should not use the
	// proxy configured above.
	NoProxy: "none",

	// Comma separated glob list matching paths that should be
	// treated as private.
	Private: "*.*",

	// Default is no workspace resolution.
	Workspace: WorkspaceDisabled,

	// A list of replacement directives mapping a module path to a directory
	// or a theme component in the themes folder.
	// Note that this will turn the component into a traditional theme component
	// that does not partake in vendoring etc.
	// The syntax is the similar to the replacement directives used in go.mod, e.g:
	//    github.com/mod1 -> ../mod1,github.com/mod2 -> ../mod2
	Replacements: nil,
}

// ApplyProjectConfigDefaults applies default/missing module configuration for
// the main project.
func ApplyProjectConfigDefaults(mod Module, cfgs ...config.AllProvider) error {
	moda := mod.(*moduleAdapter)

	// To bridge between old and new configuration format we need
	// a way to make sure all of the core components are configured on
	// the basic level.
	componentsConfigured := make(map[string]bool)
	for _, mnt := range moda.mounts {
		if !strings.HasPrefix(mnt.Target, files.JsConfigFolderMountPrefix) {
			componentsConfigured[mnt.Component()] = true
		}
	}

	var mounts []Mount

	for _, component := range []string{
		files.ComponentFolderContent,
		files.ComponentFolderData,
		files.ComponentFolderLayouts,
		files.ComponentFolderI18n,
		files.ComponentFolderArchetypes,
		files.ComponentFolderAssets,
		files.ComponentFolderStatic,
	} {
		if componentsConfigured[component] {
			continue
		}

		first := cfgs[0]
		dirsBase := first.DirsBase()
		isMultihost := first.IsMultihost()

		for i, cfg := range cfgs {
			dirs := cfg.Dirs()
			var dir string
			var dropLang bool
			switch component {
			case files.ComponentFolderContent:
				dir = dirs.ContentDir
				dropLang = dir == dirsBase.ContentDir
			case files.ComponentFolderData:
				//lint:ignore SA1019 Keep as adapter for now.
				dir = dirs.DataDir
			case files.ComponentFolderLayouts:
				//lint:ignore SA1019 Keep as adapter for now.
				dir = dirs.LayoutDir
			case files.ComponentFolderI18n:
				//lint:ignore SA1019 Keep as adapter for now.
				dir = dirs.I18nDir
			case files.ComponentFolderArchetypes:
				//lint:ignore SA1019 Keep as adapter for now.
				dir = dirs.ArcheTypeDir
			case files.ComponentFolderAssets:
				//lint:ignore SA1019 Keep as adapter for now.
				dir = dirs.AssetDir
			case files.ComponentFolderStatic:
				// For static dirs, we only care about the language in multihost setups.
				dropLang = !isMultihost
			}

			var perLang bool
			switch component {
			case files.ComponentFolderContent, files.ComponentFolderStatic:
				perLang = true
			default:
			}
			if i > 0 && !perLang {
				continue
			}

			var lang string
			if perLang && !dropLang {
				lang = cfg.Language().Lang
			}

			// Static mounts are a little special.
			if component == files.ComponentFolderStatic {
				staticDirs := cfg.StaticDirs()
				for _, dir := range staticDirs {
					mounts = append(mounts, Mount{Lang: lang, Source: dir, Target: component})
				}
				continue
			}

			if dir != "" {
				mounts = append(mounts, Mount{Lang: lang, Source: dir, Target: component})
			}
		}
	}

	moda.mounts = append(moda.mounts, mounts...)

	// Temporary: Remove duplicates.
	seen := make(map[string]bool)
	var newMounts []Mount
	for _, m := range moda.mounts {
		key := m.Source + m.Target + m.Lang
		if seen[key] {
			continue
		}
		seen[key] = true
		newMounts = append(newMounts, m)
	}
	moda.mounts = newMounts

	return nil
}

// DecodeConfig creates a modules Config from a given Hugo configuration.
func DecodeConfig(cfg config.Provider) (Config, error) {
	return decodeConfig(cfg, nil)
}

func decodeConfig(cfg config.Provider, pathReplacements map[string]string) (Config, error) {
	c := DefaultModuleConfig
	c.replacementsMap = pathReplacements

	if cfg == nil {
		return c, nil
	}

	themeSet := cfg.IsSet("theme")
	moduleSet := cfg.IsSet("module")

	if moduleSet {
		m := cfg.GetStringMap("module")
		if err := mapstructure.WeakDecode(m, &c); err != nil {
			return c, err
		}

		if c.replacementsMap == nil {

			if len(c.Replacements) == 1 {
				c.Replacements = strings.Split(c.Replacements[0], ",")
			}

			for i, repl := range c.Replacements {
				c.Replacements[i] = strings.TrimSpace(repl)
			}

			c.replacementsMap = make(map[string]string)
			for _, repl := range c.Replacements {
				parts := strings.Split(repl, "->")
				if len(parts) != 2 {
					return c, fmt.Errorf(`invalid module.replacements: %q; configure replacement pairs on the form "oldpath->newpath" `, repl)
				}

				c.replacementsMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		if c.replacementsMap != nil && c.Imports != nil {
			for i, imp := range c.Imports {
				if newImp, found := c.replacementsMap[imp.Path]; found {
					imp.Path = newImp
					imp.pathProjectReplaced = true
					c.Imports[i] = imp
				}
			}
		}

		for i, mnt := range c.Mounts {
			mnt.Source = filepath.Clean(mnt.Source)
			mnt.Target = filepath.Clean(mnt.Target)
			c.Mounts[i] = mnt
		}

		if c.Workspace == "" {
			c.Workspace = WorkspaceDisabled
		}
		if c.Workspace != WorkspaceDisabled {
			c.Workspace = filepath.Clean(c.Workspace)
			if !filepath.IsAbs(c.Workspace) {
				workingDir := cfg.GetString("workingDir")
				c.Workspace = filepath.Join(workingDir, c.Workspace)
			}
			if _, err := os.Stat(c.Workspace); err != nil {
				//lint:ignore ST1005 end user message.
				return c, fmt.Errorf("module workspace %q does not exist. Check your module.workspace setting (or HUGO_MODULE_WORKSPACE env var).", c.Workspace)
			}
		}
	}

	if themeSet {
		imports := config.GetStringSlicePreserveString(cfg, "theme")
		for _, imp := range imports {
			c.Imports = append(c.Imports, Import{
				Path: imp,
			})
		}
	}

	return c, nil
}

// Config holds a module config.
type Config struct {
	// File system mounts.
	Mounts []Mount

	// Module imports.
	Imports []Import

	// Meta info about this module (license information etc.).
	Params map[string]any

	// Will be validated against the running Hugo version.
	HugoVersion HugoVersion

	// Optional Glob pattern matching module paths to skip when vendoring, e.g. “github.com/**”
	NoVendor string

	// When enabled, we will pick the vendored module closest to the module
	// using it.
	// The default behavior is to pick the first.
	// Note that there can still be only one dependency of a given module path,
	// so once it is in use it cannot be redefined.
	VendorClosest bool

	// A comma separated (or a slice) list of module path to directory replacement mapping,
	// e.g. github.com/bep/my-theme -> ../..,github.com/bep/shortcodes -> /some/path.
	// This is mostly useful for temporary locally development of a module, and then it makes sense to set it as an
	// OS environment variable, e.g: env HUGO_MODULE_REPLACEMENTS="github.com/bep/my-theme -> ../..".
	// Any relative path is relate to themesDir, and absolute paths are allowed.
	Replacements    []string
	replacementsMap map[string]string

	// Defines the proxy server to use to download remote modules. Default is direct, which means “git clone” and similar.
	// Configures GOPROXY when running the Go command for module operations.
	Proxy string

	// Comma separated glob list matching paths that should not use the proxy configured above.
	// Configures GONOPROXY when running the Go command for module operations.
	NoProxy string

	// Comma separated glob list matching paths that should be treated as private.
	// Configures GOPRIVATE when running the Go command for module operations.
	Private string

	// Defaults to "off".
	// Set to a work file, e.g. hugo.work, to enable Go "Workspace" mode.
	// Can be relative to the working directory or absolute.
	// Requires Go 1.18+.
	// Note that this can also be set via OS env, e.g. export HUGO_MODULE_WORKSPACE=/my/hugo.work.
	Workspace string
}

// hasModuleImport reports whether the project config have one or more
// modules imports, e.g. github.com/bep/myshortcodes.
func (c Config) hasModuleImport() bool {
	for _, imp := range c.Imports {
		if isProbablyModule(imp.Path) {
			return true
		}
	}
	return false
}

// HugoVersion holds Hugo binary version requirements for a module.
type HugoVersion struct {
	// The minimum Hugo version that this module works with.
	Min hugo.VersionString

	// The maximum Hugo version that this module works with.
	Max hugo.VersionString

	// Set if the extended version is needed.
	Extended bool
}

func (v HugoVersion) String() string {
	extended := ""
	if v.Extended {
		extended = " extended"
	}

	if v.Min != "" && v.Max != "" {
		return fmt.Sprintf("%s/%s%s", v.Min, v.Max, extended)
	}

	if v.Min != "" {
		return fmt.Sprintf("Min %s%s", v.Min, extended)
	}

	if v.Max != "" {
		return fmt.Sprintf("Max %s%s", v.Max, extended)
	}

	return extended
}

// IsValid reports whether this version is valid compared to the running
// Hugo binary.
func (v HugoVersion) IsValid() bool {
	current := hugo.CurrentVersion.Version()
	if v.Extended && !hugo.IsExtended {
		return false
	}

	isValid := true

	if v.Min != "" && current.Compare(v.Min) > 0 {
		isValid = false
	}

	if v.Max != "" && current.Compare(v.Max) < 0 {
		isValid = false
	}

	return isValid
}

type Import struct {
	// Module path
	Path string
	// Set when Path is replaced in project config.
	pathProjectReplaced bool
	// Ignore any config in config.toml (will still follow imports).
	IgnoreConfig bool
	// Do not follow any configured imports.
	IgnoreImports bool
	// Do not mount any folder in this import.
	NoMounts bool
	// Never vendor this import (only allowed in main project).
	NoVendor bool
	// Turn off this module.
	Disable bool
	// File mounts.
	Mounts []Mount
}

type Mount struct {
	// Relative path in source repo, e.g. "scss".
	Source string

	// Relative target path, e.g. "assets/bootstrap/scss".
	Target string

	// Any file in this mount will be associated with this language.
	Lang string

	// Include only files matching the given Glob patterns (string or slice).
	IncludeFiles any

	// Exclude all files matching the given Glob patterns (string or slice).
	ExcludeFiles any

	// Disable watching in watch mode for this mount.
	DisableWatch bool
}

// Used as key to remove duplicates.
func (m Mount) key() string {
	return strings.Join([]string{m.Lang, m.Source, m.Target}, "/")
}

func (m Mount) Component() string {
	return strings.Split(m.Target, fileSeparator)[0]
}

func (m Mount) ComponentAndName() (string, string) {
	c, n, _ := strings.Cut(m.Target, fileSeparator)
	return c, n
}
