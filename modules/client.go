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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gohugoio/hugo/common/collections"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hexec"
	"github.com/gohugoio/hugo/common/loggers"

	hglob "github.com/gohugoio/hugo/hugofs/glob"

	"github.com/gobwas/glob"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/config"

	"golang.org/x/mod/module"

	"github.com/gohugoio/hugo/common/hugio"

	"github.com/spf13/afero"
)

var fileSeparator = string(os.PathSeparator)

const (
	goBinaryStatusOK goBinaryStatus = iota
	goBinaryStatusNotFound
	goBinaryStatusTooOld
)

// The "vendor" dir is reserved for Go Modules.
const vendord = "_vendor"

const (
	goModFilename = "go.mod"
	goSumFilename = "go.sum"
)

// NewClient creates a new Client that can be used to manage the Hugo Components
// in a given workingDir.
// The Client will resolve the dependencies recursively, but needs the top
// level imports to start out.
func NewClient(cfg ClientConfig) *Client {
	fs := cfg.Fs
	n := filepath.Join(cfg.WorkingDir, goModFilename)
	goModEnabled, _ := afero.Exists(fs, n)
	var goModFilename string
	if goModEnabled {
		goModFilename = n
	}

	var env []string
	mcfg := cfg.ModuleConfig

	config.SetEnvVars(&env,
		"PWD", cfg.WorkingDir,
		"GO111MODULE", "on",
		"GOPROXY", mcfg.Proxy,
		"GOPRIVATE", mcfg.Private,
		"GONOPROXY", mcfg.NoProxy,
		"GOPATH", cfg.CacheDir,
		"GOWORK", mcfg.Workspace, // Requires Go 1.18, see https://tip.golang.org/doc/go1.18
		// GOCACHE was introduced in Go 1.15. This matches the location derived from GOPATH above.
		"GOCACHE", filepath.Join(cfg.CacheDir, "pkg", "mod"),
	)

	logger := cfg.Logger
	if logger == nil {
		logger = loggers.NewDefault()
	}

	var noVendor glob.Glob
	if cfg.ModuleConfig.NoVendor != "" {
		noVendor, _ = hglob.GetGlob(hglob.NormalizePath(cfg.ModuleConfig.NoVendor))
	}

	return &Client{
		fs:                fs,
		ccfg:              cfg,
		logger:            logger,
		noVendor:          noVendor,
		moduleConfig:      mcfg,
		environ:           env,
		GoModulesFilename: goModFilename,
	}
}

// Client contains most of the API provided by this package.
type Client struct {
	fs     afero.Fs
	logger loggers.Logger

	noVendor glob.Glob

	ccfg ClientConfig

	// The top level module config
	moduleConfig Config

	// Environment variables used in "go get" etc.
	environ []string

	// Set when Go modules are initialized in the current repo, that is:
	// a go.mod file exists.
	GoModulesFilename string

	// Set if we get a exec.ErrNotFound when running Go, which is most likely
	// due to being run on a system without Go installed. We record it here
	// so we can give an instructional error at the end if module/theme
	// resolution fails.
	goBinaryStatus goBinaryStatus
}

// Graph writes a module dependency graph to the given writer.
func (c *Client) Graph(w io.Writer) error {
	mc, coll := c.collect(true)
	if coll.err != nil {
		return coll.err
	}
	for _, module := range mc.AllModules {
		if module.Owner() == nil {
			continue
		}

		dep := pathVersion(module.Owner()) + " " + pathVersion(module)
		if replace := module.Replace(); replace != nil {
			if replace.Version() != "" {
				dep += " => " + pathVersion(replace)
			} else {
				// Local dir.
				dep += " => " + replace.Dir()
			}
		}
		fmt.Fprintln(w, dep)
	}

	return nil
}

// Tidy can be used to remove unused dependencies from go.mod and go.sum.
func (c *Client) Tidy() error {
	tc, coll := c.collect(false)
	if coll.err != nil {
		return coll.err
	}

	if coll.skipTidy {
		return nil
	}

	return c.tidy(tc.AllModules, false)
}

// Vendor writes all the module dependencies to a _vendor folder.
//
// Unlike Go, we support it for any level.
//
// We, by default, use the /_vendor folder first, if found. To disable,
// run with
//
//	hugo --ignoreVendorPaths=".*"
//
// Given a module tree, Hugo will pick the first module for a given path,
// meaning that if the top-level module is vendored, that will be the full
// set of dependencies.
func (c *Client) Vendor() error {
	vendorDir := filepath.Join(c.ccfg.WorkingDir, vendord)
	if err := c.rmVendorDir(vendorDir); err != nil {
		return err
	}
	if err := c.fs.MkdirAll(vendorDir, 0o755); err != nil {
		return err
	}

	// Write the modules list to modules.txt.
	//
	// On the form:
	//
	// # github.com/alecthomas/chroma v0.6.3
	//
	// This is how "go mod vendor" does it. Go also lists
	// the packages below it, but that is currently not applicable to us.
	//
	var modulesContent bytes.Buffer

	tc, coll := c.collect(true)
	if coll.err != nil {
		return coll.err
	}

	for _, t := range tc.AllModules {
		if t.Owner() == nil {
			// This is the project.
			continue
		}

		if !c.shouldVendor(t.Path()) {
			continue
		}

		if !t.IsGoMod() && !t.Vendor() {
			// We currently do not vendor components living in the
			// theme directory, see https://github.com/gohugoio/hugo/issues/5993
			continue
		}

		// See https://github.com/gohugoio/hugo/issues/8239
		// This is an error situation. We need something to vendor.
		if t.Mounts() == nil {
			return fmt.Errorf("cannot vendor module %q, need at least one mount", t.Path())
		}

		fmt.Fprintln(&modulesContent, "# "+t.Path()+" "+t.Version())

		dir := t.Dir()

		for _, mount := range t.Mounts() {
			sourceFilename := filepath.Join(dir, mount.Source)
			targetFilename := filepath.Join(vendorDir, t.Path(), mount.Source)
			fi, err := c.fs.Stat(sourceFilename)
			if err != nil {
				return fmt.Errorf("failed to vendor module: %w", err)
			}

			if fi.IsDir() {
				if err := hugio.CopyDir(c.fs, sourceFilename, targetFilename, nil); err != nil {
					return fmt.Errorf("failed to copy module to vendor dir: %w", err)
				}
			} else {
				targetDir := filepath.Dir(targetFilename)

				if err := c.fs.MkdirAll(targetDir, 0o755); err != nil {
					return fmt.Errorf("failed to make target dir: %w", err)
				}

				if err := hugio.CopyFile(c.fs, sourceFilename, targetFilename); err != nil {
					return fmt.Errorf("failed to copy module file to vendor: %w", err)
				}
			}
		}

		// Include the resource cache if present.
		resourcesDir := filepath.Join(dir, files.FolderResources)
		_, err := c.fs.Stat(resourcesDir)
		if err == nil {
			if err := hugio.CopyDir(c.fs, resourcesDir, filepath.Join(vendorDir, t.Path(), files.FolderResources), nil); err != nil {
				return fmt.Errorf("failed to copy resources to vendor dir: %w", err)
			}
		}

		// Include the config directory if present.
		configDir := filepath.Join(dir, "config")
		_, err = c.fs.Stat(configDir)
		if err == nil {
			if err := hugio.CopyDir(c.fs, configDir, filepath.Join(vendorDir, t.Path(), "config"), nil); err != nil {
				return fmt.Errorf("failed to copy config dir to vendor dir: %w", err)
			}
		}

		// Also include any theme.toml or config.* or hugo.* files in the root.
		configFiles, _ := afero.Glob(c.fs, filepath.Join(dir, "config.*"))
		configFiles2, _ := afero.Glob(c.fs, filepath.Join(dir, "hugo.*"))
		configFiles = append(configFiles, configFiles2...)
		configFiles = append(configFiles, filepath.Join(dir, "theme.toml"))
		for _, configFile := range configFiles {
			if err := hugio.CopyFile(c.fs, configFile, filepath.Join(vendorDir, t.Path(), filepath.Base(configFile))); err != nil {
				if !herrors.IsNotExist(err) {
					return err
				}
			}
		}
	}

	if modulesContent.Len() > 0 {
		if err := afero.WriteFile(c.fs, filepath.Join(vendorDir, vendorModulesFilename), modulesContent.Bytes(), 0o666); err != nil {
			return err
		}
	}

	return nil
}

// Get runs "go get" with the supplied arguments.
func (c *Client) Get(args ...string) error {
	if len(args) == 0 || (len(args) == 1 && strings.Contains(args[0], "-u")) {
		update := len(args) != 0
		patch := update && (args[0] == "-u=patch") //

		// We need to be explicit about the modules to get.
		var modules []string
		// Update all active modules if the -u flag presents.
		if update {
			mc, coll := c.collect(true)
			if coll.err != nil {
				return coll.err
			}
			for _, m := range mc.AllModules {
				if m.Owner() == nil || !isProbablyModule(m.Path()) {
					continue
				}
				modules = append(modules, m.Path())
			}
		} else {
			for _, m := range c.moduleConfig.Imports {
				if !isProbablyModule(m.Path) {
					// Skip themes/components stored below /themes etc.
					// There may be false positives in the above, but those
					// should be rare, and they will fail below with an
					// "cannot find module providing ..." message.
					continue
				}
				modules = append(modules, m.Path)
			}
		}

		for _, m := range modules {
			var args []string

			if update && !patch {
				args = append(args, "-u")
			} else if update && patch {
				args = append(args, "-u=patch")
			}
			args = append(args, m)

			if err := c.get(args...); err != nil {
				return err
			}
		}

		return nil
	}

	return c.get(args...)
}

func (c *Client) get(args ...string) error {
	if err := c.runGo(context.Background(), c.logger.Out(), append([]string{"get"}, args...)...); err != nil {
		return fmt.Errorf("failed to get %q: %w", args, err)
	}
	return nil
}

// Init initializes this as a Go Module with the given path.
// If path is empty, Go will try to guess.
// If this succeeds, this project will be marked as Go Module.
func (c *Client) Init(path string) error {
	err := c.runGo(context.Background(), c.logger.Out(), "mod", "init", path)
	if err != nil {
		return fmt.Errorf("failed to init modules: %w", err)
	}

	c.GoModulesFilename = filepath.Join(c.ccfg.WorkingDir, goModFilename)

	return nil
}

var verifyErrorDirRe = regexp.MustCompile(`dir has been modified \((.*?)\)`)

// Verify checks that the dependencies of the current module,
// which are stored in a local downloaded source cache, have not been
// modified since being downloaded.
func (c *Client) Verify(clean bool) error {
	// TODO(bep) add path to mod clean
	err := c.runVerify()
	if err != nil {
		if clean {
			m := verifyErrorDirRe.FindAllStringSubmatch(err.Error(), -1)
			if m != nil {
				for i := 0; i < len(m); i++ {
					c, err := hugofs.MakeReadableAndRemoveAllModulePkgDir(c.fs, m[i][1])
					if err != nil {
						return err
					}
					fmt.Println("Cleaned", c)
				}
			}
			// Try to verify it again.
			err = c.runVerify()
		}
	}
	return err
}

func (c *Client) Clean(pattern string) error {
	mods, err := c.listGoMods()
	if err != nil {
		return err
	}

	var g glob.Glob

	if pattern != "" {
		var err error
		g, err = hglob.GetGlob(pattern)
		if err != nil {
			return err
		}
	}

	for _, m := range mods {
		if m.Replace != nil || m.Main {
			continue
		}

		if g != nil && !g.Match(m.Path) {
			continue
		}
		dirCount, err := hugofs.MakeReadableAndRemoveAllModulePkgDir(c.fs, m.Dir)
		if err == nil {
			c.logger.Printf("hugo: removed %d dirs in module cache for %q", dirCount, m.Path)
		}
	}
	return err
}

func (c *Client) runVerify() error {
	return c.runGo(context.Background(), io.Discard, "mod", "verify")
}

func isProbablyModule(path string) bool {
	return module.CheckPath(path) == nil
}

func (c *Client) listGoMods() (goModules, error) {
	if c.GoModulesFilename == "" || !c.moduleConfig.hasModuleImport() {
		return nil, nil
	}

	downloadModules := func(modules ...string) error {
		args := []string{"mod", "download", "-modcacherw"}
		args = append(args, modules...)
		out := io.Discard
		err := c.runGo(context.Background(), out, args...)
		if err != nil {
			return fmt.Errorf("failed to download modules: %w", err)
		}
		return nil
	}

	if err := downloadModules(); err != nil {
		return nil, err
	}

	listAndDecodeModules := func(handle func(m *goModule) error, modules ...string) error {
		b := &bytes.Buffer{}
		args := []string{"list", "-m", "-json"}
		if len(modules) > 0 {
			args = append(args, modules...)
		} else {
			args = append(args, "all")
		}
		err := c.runGo(context.Background(), b, args...)
		if err != nil {
			return fmt.Errorf("failed to list modules: %w", err)
		}

		dec := json.NewDecoder(b)
		for {
			m := &goModule{}
			if err := dec.Decode(m); err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("failed to decode modules list: %w", err)
			}

			if err := handle(m); err != nil {
				return err
			}
		}
		return nil
	}

	var modules goModules
	err := listAndDecodeModules(func(m *goModule) error {
		modules = append(modules, m)
		return nil
	})
	if err != nil {
		return nil, err
	}

	// From Go 1.17, go lazy loads transitive dependencies.
	// That does not work for us.
	// So, download these modules and update the Dir in the modules list.
	var modulesToDownload []string
	for _, m := range modules {
		if m.Dir == "" {
			modulesToDownload = append(modulesToDownload, fmt.Sprintf("%s@%s", m.Path, m.Version))
		}
	}

	if len(modulesToDownload) > 0 {
		if err := downloadModules(modulesToDownload...); err != nil {
			return nil, err
		}
		err := listAndDecodeModules(func(m *goModule) error {
			if mm := modules.GetByPath(m.Path); mm != nil {
				mm.Dir = m.Dir
			}
			return nil
		}, modulesToDownload...)
		if err != nil {
			return nil, err
		}
	}

	return modules, err
}

func (c *Client) rewriteGoMod(name string, isGoMod map[string]bool) error {
	data, err := c.rewriteGoModRewrite(name, isGoMod)
	if err != nil {
		return err
	}
	if data != nil {
		if err := afero.WriteFile(c.fs, filepath.Join(c.ccfg.WorkingDir, name), data, 0o666); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) rewriteGoModRewrite(name string, isGoMod map[string]bool) ([]byte, error) {
	if name == goModFilename && c.GoModulesFilename == "" {
		// Already checked.
		return nil, nil
	}

	modlineSplitter := getModlineSplitter(name == goModFilename)

	b := &bytes.Buffer{}
	f, err := c.fs.Open(filepath.Join(c.ccfg.WorkingDir, name))
	if err != nil {
		if herrors.IsNotExist(err) {
			// It's been deleted.
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var dirty bool

	for scanner.Scan() {
		line := scanner.Text()
		var doWrite bool

		if parts := modlineSplitter(line); parts != nil {
			modname, modver := parts[0], parts[1]
			modver = strings.TrimSuffix(modver, "/"+goModFilename)
			modnameVer := modname + " " + modver
			doWrite = isGoMod[modnameVer]
		} else {
			doWrite = true
		}

		if doWrite {
			fmt.Fprintln(b, line)
		} else {
			dirty = true
		}
	}

	if !dirty {
		// Nothing changed
		return nil, nil
	}

	return b.Bytes(), nil
}

func (c *Client) rmVendorDir(vendorDir string) error {
	modulestxt := filepath.Join(vendorDir, vendorModulesFilename)

	if _, err := c.fs.Stat(vendorDir); err != nil {
		return nil
	}

	_, err := c.fs.Stat(modulestxt)
	if err != nil {
		// If we have a _vendor dir without modules.txt it sounds like
		// a _vendor dir created by others.
		return errors.New("found _vendor dir without modules.txt, skip delete")
	}

	return c.fs.RemoveAll(vendorDir)
}

func (c *Client) runGo(
	ctx context.Context,
	stdout io.Writer,
	args ...string,
) error {
	if c.goBinaryStatus != 0 {
		return nil
	}

	stderr := new(bytes.Buffer)

	argsv := collections.StringSliceToInterfaceSlice(args)
	argsv = append(argsv, hexec.WithEnviron(c.environ))
	argsv = append(argsv, hexec.WithStderr(goOutputReplacerWriter{w: io.MultiWriter(stderr, os.Stderr)}))
	argsv = append(argsv, hexec.WithStdout(stdout))
	argsv = append(argsv, hexec.WithDir(c.ccfg.WorkingDir))
	argsv = append(argsv, hexec.WithContext(ctx))

	cmd, err := c.ccfg.Exec.New("go", argsv...)
	if err != nil {
		return err
	}

	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.Error); ok && ee.Err == exec.ErrNotFound {
			c.goBinaryStatus = goBinaryStatusNotFound
			return nil
		}

		if strings.Contains(stderr.String(), "invalid version: unknown revision") {
			// See https://github.com/gohugoio/hugo/issues/6825
			c.logger.Println(`An unknown revision most likely means that someone has deleted the remote ref (e.g. with a force push to GitHub).
To resolve this, you need to manually edit your go.mod file and replace the version for the module in question with a valid ref.

The easiest is to just enter a valid branch name there, e.g. master, which would be what you put in place of 'v0.5.1' in the example below.

require github.com/gohugoio/hugo-mod-jslibs/instantpage v0.5.1

If you then run 'hugo mod graph' it should resolve itself to the most recent version (or commit if no semver versions are available).`)
		}

		_, ok := err.(*exec.ExitError)
		if !ok {
			return fmt.Errorf("failed to execute 'go %v': %s %T", args, err, err)
		}

		// Too old Go version
		if strings.Contains(stderr.String(), "flag provided but not defined") {
			c.goBinaryStatus = goBinaryStatusTooOld
			return nil
		}

		return fmt.Errorf("go command failed: %s", stderr)

	}

	return nil
}

var goOutputReplacer = strings.NewReplacer(
	"go: to add module requirements and sums:", "hugo: to add module requirements and sums:",
	"go mod tidy", "hugo mod tidy",
)

type goOutputReplacerWriter struct {
	w io.Writer
}

func (w goOutputReplacerWriter) Write(p []byte) (n int, err error) {
	s := goOutputReplacer.Replace(string(p))
	_, err = w.w.Write([]byte(s))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *Client) tidy(mods Modules, goModOnly bool) error {
	isGoMod := make(map[string]bool)
	for _, m := range mods {
		if m.Owner() == nil {
			continue
		}
		if m.IsGoMod() {
			// Matching the format in go.mod
			pathVer := m.Path() + " " + m.Version()
			isGoMod[pathVer] = true
		}
	}

	if err := c.rewriteGoMod(goModFilename, isGoMod); err != nil {
		return err
	}

	if goModOnly {
		return nil
	}

	if err := c.rewriteGoMod(goSumFilename, isGoMod); err != nil {
		return err
	}

	return nil
}

func (c *Client) shouldVendor(path string) bool {
	return c.noVendor == nil || !c.noVendor.Match(path)
}

func (c *Client) createThemeDirname(modulePath string, isProjectMod bool) (string, error) {
	invalid := fmt.Errorf("invalid module path %q; must be relative to themesDir when defined outside of the project", modulePath)

	modulePath = filepath.Clean(modulePath)
	if filepath.IsAbs(modulePath) {
		if isProjectMod {
			return modulePath, nil
		}
		return "", invalid
	}

	moduleDir := filepath.Join(c.ccfg.ThemesDir, modulePath)
	if !isProjectMod && !strings.HasPrefix(moduleDir, c.ccfg.ThemesDir) {
		return "", invalid
	}
	return moduleDir, nil
}

// ClientConfig configures the module Client.
type ClientConfig struct {
	Fs     afero.Fs
	Logger loggers.Logger

	// If set, it will be run before we do any duplicate checks for modules
	// etc.
	HookBeforeFinalize func(m *ModulesConfig) error

	// Ignore any _vendor directory for module paths matching the given pattern.
	// This can be nil.
	IgnoreVendor glob.Glob

	// Ignore any module not found errors.
	IgnoreModuleDoesNotExist bool

	// Absolute path to the project dir.
	WorkingDir string

	// Absolute path to the project's themes dir.
	ThemesDir string

	// The publish dir.
	PublishDir string

	// Eg. "production"
	Environment string

	Exec *hexec.Exec

	CacheDir     string // Module cache
	ModuleConfig Config
}

func (c ClientConfig) shouldIgnoreVendor(path string) bool {
	return c.IgnoreVendor != nil && c.IgnoreVendor.Match(path)
}

type goBinaryStatus int

type goModule struct {
	Path     string         // module path
	Version  string         // module version
	Versions []string       // available module versions (with -versions)
	Replace  *goModule      // replaced by this module
	Time     *time.Time     // time version was created
	Update   *goModule      // available update, if any (with -u)
	Main     bool           // is this the main module?
	Indirect bool           // is this module only an indirect dependency of main module?
	Dir      string         // directory holding files for this module, if any
	GoMod    string         // path to go.mod file for this module, if any
	Error    *goModuleError // error loading module
}

type goModuleError struct {
	Err string // the error itself
}

type goModules []*goModule

func (modules goModules) GetByPath(p string) *goModule {
	if modules == nil {
		return nil
	}

	for _, m := range modules {
		if strings.EqualFold(p, m.Path) {
			return m
		}
	}

	return nil
}

func (modules goModules) GetMain() *goModule {
	for _, m := range modules {
		if m.Main {
			return m
		}
	}

	return nil
}

func getModlineSplitter(isGoMod bool) func(line string) []string {
	if isGoMod {
		return func(line string) []string {
			if strings.HasPrefix(line, "require (") {
				return nil
			}
			if !strings.HasPrefix(line, "require") && !strings.HasPrefix(line, "\t") {
				return nil
			}
			line = strings.TrimPrefix(line, "require")
			line = strings.TrimSpace(line)
			line = strings.TrimSuffix(line, "// indirect")

			return strings.Fields(line)
		}
	}

	return func(line string) []string {
		return strings.Fields(line)
	}
}

func pathVersion(m Module) string {
	versionStr := m.Version()
	if m.Vendor() {
		versionStr += "+vendor"
	}
	if versionStr == "" {
		return m.Path()
	}
	return fmt.Sprintf("%s@%s", m.Path(), versionStr)
}
