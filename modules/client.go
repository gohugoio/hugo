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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/common/loggers"

	"strings"
	"time"

	"github.com/gohugoio/hugo/config"

	"github.com/rogpeppe/go-internal/module"

	"github.com/gohugoio/hugo/common/hugio"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

var (
	fileSeparator = string(os.PathSeparator)
)

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

	env := os.Environ()
	mcfg := cfg.ModuleConfig

	config.SetEnvVars(&env,
		"PWD", cfg.WorkingDir,
		"GO111MODULE", "on",
		"GOPROXY", mcfg.Proxy,
		"GOPRIVATE", mcfg.Private,
		"GONOPROXY", mcfg.NoProxy)

	if cfg.CacheDir != "" {
		// Module cache stored below $GOPATH/pkg
		config.SetEnvVars(&env, "GOPATH", cfg.CacheDir)

	}

	logger := cfg.Logger
	if logger == nil {
		logger = loggers.NewWarningLogger()
	}

	return &Client{
		fs:                fs,
		ccfg:              cfg,
		logger:            logger,
		moduleConfig:      mcfg,
		environ:           env,
		GoModulesFilename: goModFilename}
}

// Client contains most of the API provided by this package.
type Client struct {
	fs     afero.Fs
	logger *loggers.Logger

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

// Graph writes a module dependenchy graph to the given writer.
func (c *Client) Graph(w io.Writer) error {
	mc, coll := c.collect(true)
	if coll.err != nil {
		return coll.err
	}
	for _, module := range mc.AllModules {
		if module.Owner() == nil {
			continue
		}

		prefix := ""
		if module.Disabled() {
			prefix = "DISABLED "
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
		fmt.Fprintln(w, prefix+dep)
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
//    hugo --ignoreVendor
//
// Given a module tree, Hugo will pick the first module for a given path,
// meaning that if the top-level module is vendored, that will be the full
// set of dependencies.
func (c *Client) Vendor() error {
	vendorDir := filepath.Join(c.ccfg.WorkingDir, vendord)
	if err := c.rmVendorDir(vendorDir); err != nil {
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
		// We respect the --ignoreVendor flag even for the vendor command.
		if !t.IsGoMod() && !t.Vendor() {
			// We currently do not vendor components living in the
			// theme directory, see https://github.com/gohugoio/hugo/issues/5993
			continue
		}

		fmt.Fprintln(&modulesContent, "# "+t.Path()+" "+t.Version())

		dir := t.Dir()

		for _, mount := range t.Mounts() {
			if err := hugio.CopyDir(c.fs, filepath.Join(dir, mount.Source), filepath.Join(vendorDir, t.Path(), mount.Source), nil); err != nil {
				return errors.Wrap(err, "failed to copy module to vendor dir")
			}
		}

		// Include the resource cache if present.
		resourcesDir := filepath.Join(dir, files.FolderResources)
		_, err := c.fs.Stat(resourcesDir)
		if err == nil {
			if err := hugio.CopyDir(c.fs, resourcesDir, filepath.Join(vendorDir, t.Path(), files.FolderResources), nil); err != nil {
				return errors.Wrap(err, "failed to copy resources to vendor dir")
			}
		}

		// Also include any theme.toml or config.* files in the root.
		configFiles, _ := afero.Glob(c.fs, filepath.Join(dir, "config.*"))
		configFiles = append(configFiles, filepath.Join(dir, "theme.toml"))
		for _, configFile := range configFiles {
			if err := hugio.CopyFile(c.fs, configFile, filepath.Join(vendorDir, t.Path(), filepath.Base(configFile))); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			}
		}
	}

	if modulesContent.Len() > 0 {
		if err := afero.WriteFile(c.fs, filepath.Join(vendorDir, vendorModulesFilename), modulesContent.Bytes(), 0666); err != nil {
			return err
		}
	}

	return nil
}

// Get runs "go get" with the supplied arguments.
func (c *Client) Get(args ...string) error {
	if err := c.runGo(context.Background(), c.logger.Out, append([]string{"get"}, args...)...); err != nil {
		errors.Wrapf(err, "failed to get %q", args)
	}
	return nil
}

// Init initializes this as a Go Module with the given path.
// If path is empty, Go will try to guess.
// If this succeeds, this project will be marked as Go Module.
func (c *Client) Init(path string) error {
	err := c.runGo(context.Background(), c.logger.Out, "mod", "init", path)
	if err != nil {
		return errors.Wrap(err, "failed to init modules")
	}

	c.GoModulesFilename = filepath.Join(c.ccfg.WorkingDir, goModFilename)

	return nil
}

func isProbablyModule(path string) bool {
	return module.CheckPath(path) == nil
}

func (c *Client) listGoMods() (goModules, error) {
	if c.GoModulesFilename == "" || !c.moduleConfig.hasModuleImport() {
		return nil, nil
	}

	out := ioutil.Discard
	err := c.runGo(context.Background(), out, "mod", "download")
	if err != nil {
		return nil, errors.Wrap(err, "failed to download modules")
	}

	b := &bytes.Buffer{}
	err = c.runGo(context.Background(), b, "list", "-m", "-json", "all")
	if err != nil {
		return nil, errors.Wrap(err, "failed to list modules")
	}

	var modules goModules

	dec := json.NewDecoder(b)
	for {
		m := &goModule{}
		if err := dec.Decode(m); err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Wrap(err, "failed to decode modules list")
		}

		modules = append(modules, m)
	}

	return modules, err

}

func (c *Client) rewriteGoMod(name string, isGoMod map[string]bool) error {
	data, err := c.rewriteGoModRewrite(name, isGoMod)
	if err != nil {
		return err
	}
	if data != nil {
		if err := afero.WriteFile(c.fs, filepath.Join(c.ccfg.WorkingDir, name), data, 0666); err != nil {
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
		if os.IsNotExist(err) {
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
	args ...string) error {

	if c.goBinaryStatus != 0 {
		return nil
	}

	//defer c.logger.PrintTimer(time.Now(), fmt.Sprint(args))

	stderr := new(bytes.Buffer)
	cmd := exec.CommandContext(ctx, "go", args...)

	cmd.Env = c.environ
	cmd.Dir = c.ccfg.WorkingDir
	cmd.Stdout = stdout
	cmd.Stderr = io.MultiWriter(stderr, os.Stderr)

	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.Error); ok && ee.Err == exec.ErrNotFound {
			c.goBinaryStatus = goBinaryStatusNotFound
			return nil
		}

		_, ok := err.(*exec.ExitError)
		if !ok {
			return errors.Errorf("failed to execute 'go %v': %s %T", args, err, err)
		}

		// Too old Go version
		if strings.Contains(stderr.String(), "flag provided but not defined") {
			c.goBinaryStatus = goBinaryStatusTooOld
			return nil
		}

		return errors.Errorf("go command failed: %s", stderr)

	}

	return nil
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

// ClientConfig configures the module Client.
type ClientConfig struct {
	Fs     afero.Fs
	Logger *loggers.Logger

	// If set, it will be run before we do any duplicate checks for modules
	// etc.
	HookBeforeFinalize func(m *ModulesConfig) error

	// Ignore any _vendor directory.
	IgnoreVendor bool

	// Absolute path to the project dir.
	WorkingDir string

	// Absolute path to the project's themes dir.
	ThemesDir string

	CacheDir     string // Module cache
	ModuleConfig Config
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
