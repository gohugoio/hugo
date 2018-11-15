// Copyright 2018 The Hugo Authors. All rights reserved.
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

package commands

import (
	"bytes"
	"errors"

	"io/ioutil"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugo"

	jww "github.com/spf13/jwalterweatherman"

	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"

	"github.com/spf13/cobra"

	"github.com/gohugoio/hugo/hugolib"
	"github.com/spf13/afero"

	"github.com/bep/debounce"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/langs"
)

type commandeerHugoState struct {
	*deps.DepsCfg
	hugo     *hugolib.HugoSites
	fsCreate sync.Once
}

type commandeer struct {
	*commandeerHugoState

	logger *loggers.Logger

	// Currently only set when in "fast render mode". But it seems to
	// be fast enough that we could maybe just add it for all server modes.
	changeDetector *fileChangeDetector

	// We need to reuse this on server rebuilds.
	destinationFs afero.Fs

	h    *hugoBuilderCommon
	ftch flagsToConfigHandler

	visitedURLs *types.EvictingStringQueue

	doWithCommandeer func(c *commandeer) error

	// We watch these for changes.
	configFiles []string

	// Used in cases where we get flooded with events in server mode.
	debounce func(f func())

	serverPorts         []int
	languagesConfigured bool
	languages           langs.Languages
	doLiveReload        bool
	fastRenderMode      bool
	showErrorInBrowser  bool

	configured bool
	paused     bool

	// Any error from the last build.
	buildErr error
}

func (c *commandeer) errCount() int {
	return int(c.logger.ErrorCounter.Count())
}

func (c *commandeer) getErrorWithContext() interface{} {
	errCount := c.errCount()

	if errCount == 0 {
		return nil
	}

	m := make(map[string]interface{})

	m["Error"] = errors.New(removeErrorPrefixFromLog(c.logger.Errors()))
	m["Version"] = hugo.BuildVersionString()

	fe := herrors.UnwrapErrorWithFileContext(c.buildErr)
	if fe != nil {
		m["File"] = fe
	}

	if c.h.verbose {
		var b bytes.Buffer
		herrors.FprintStackTrace(&b, c.buildErr)
		m["StackTrace"] = b.String()
	}

	return m
}

func (c *commandeer) Set(key string, value interface{}) {
	if c.configured {
		panic("commandeer cannot be changed")
	}
	c.Cfg.Set(key, value)
}

func (c *commandeer) initFs(fs *hugofs.Fs) error {
	c.destinationFs = fs.Destination
	c.DepsCfg.Fs = fs

	return nil
}

func newCommandeer(mustHaveConfigFile, running bool, h *hugoBuilderCommon, f flagsToConfigHandler, doWithCommandeer func(c *commandeer) error, subCmdVs ...*cobra.Command) (*commandeer, error) {

	var rebuildDebouncer func(f func())
	if running {
		// The time value used is tested with mass content replacements in a fairly big Hugo site.
		// It is better to wait for some seconds in those cases rather than get flooded
		// with rebuilds.
		rebuildDebouncer, _, _ = debounce.New(4 * time.Second)
	}

	c := &commandeer{
		h:                   h,
		ftch:                f,
		commandeerHugoState: &commandeerHugoState{},
		doWithCommandeer:    doWithCommandeer,
		visitedURLs:         types.NewEvictingStringQueue(10),
		debounce:            rebuildDebouncer,
		// This will be replaced later, but we need something to log to before the configuration is read.
		logger: loggers.NewLogger(jww.LevelError, jww.LevelError, os.Stdout, ioutil.Discard, running),
	}

	return c, c.loadConfig(mustHaveConfigFile, running)
}

type fileChangeDetector struct {
	sync.Mutex
	current map[string]string
	prev    map[string]string

	irrelevantRe *regexp.Regexp
}

func (f *fileChangeDetector) OnFileClose(name, md5sum string) {
	f.Lock()
	defer f.Unlock()
	f.current[name] = md5sum
}

func (f *fileChangeDetector) changed() []string {
	if f == nil {
		return nil
	}
	f.Lock()
	defer f.Unlock()
	var c []string
	for k, v := range f.current {
		vv, found := f.prev[k]
		if !found || v != vv {
			c = append(c, k)
		}
	}

	return f.filterIrrelevant(c)
}

func (f *fileChangeDetector) filterIrrelevant(in []string) []string {
	var filtered []string
	for _, v := range in {
		if !f.irrelevantRe.MatchString(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func (f *fileChangeDetector) PrepareNew() {
	if f == nil {
		return
	}

	f.Lock()
	defer f.Unlock()

	if f.current == nil {
		f.current = make(map[string]string)
		f.prev = make(map[string]string)
		return
	}

	f.prev = make(map[string]string)
	for k, v := range f.current {
		f.prev[k] = v
	}
	f.current = make(map[string]string)
}

func (c *commandeer) loadConfig(mustHaveConfigFile, running bool) error {

	if c.DepsCfg == nil {
		c.DepsCfg = &deps.DepsCfg{}
	}

	if c.logger != nil {
		// Truncate the error log if this is a reload.
		c.logger.Reset()
	}

	cfg := c.DepsCfg
	c.configured = false
	cfg.Running = running

	var dir string
	if c.h.source != "" {
		dir, _ = filepath.Abs(c.h.source)
	} else {
		dir, _ = os.Getwd()
	}

	var sourceFs afero.Fs = hugofs.Os
	if c.DepsCfg.Fs != nil {
		sourceFs = c.DepsCfg.Fs.Source
	}

	environment := c.h.getEnvironment(running)

	doWithConfig := func(cfg config.Provider) error {

		if c.ftch != nil {
			c.ftch.flagsToConfig(cfg)
		}

		cfg.Set("workingDir", dir)
		cfg.Set("environment", environment)
		return nil
	}

	doWithCommandeer := func(cfg config.Provider) error {
		c.Cfg = cfg
		if c.doWithCommandeer == nil {
			return nil
		}
		err := c.doWithCommandeer(c)
		return err
	}

	configPath := c.h.source
	if configPath == "" {
		configPath = dir
	}
	config, configFiles, err := hugolib.LoadConfig(
		hugolib.ConfigSourceDescriptor{
			Fs:           sourceFs,
			Path:         configPath,
			WorkingDir:   dir,
			Filename:     c.h.cfgFile,
			AbsConfigDir: c.h.getConfigDir(dir),
			Environment:  environment},
		doWithCommandeer,
		doWithConfig)

	if err != nil {
		if mustHaveConfigFile {
			return err
		}
		if err != hugolib.ErrNoConfigFile {
			return err
		}

	}

	c.configFiles = configFiles

	if l, ok := c.Cfg.Get("languagesSorted").(langs.Languages); ok {
		c.languagesConfigured = true
		c.languages = l
	}

	// Set some commonly used flags
	c.doLiveReload = !c.h.buildWatch && !c.Cfg.GetBool("disableLiveReload")
	c.fastRenderMode = c.doLiveReload && !c.Cfg.GetBool("disableFastRender")
	c.showErrorInBrowser = c.doLiveReload && !c.Cfg.GetBool("disableBrowserError")

	// This is potentially double work, but we need to do this one more time now
	// that all the languages have been configured.
	if c.doWithCommandeer != nil {
		if err := c.doWithCommandeer(c); err != nil {
			return err
		}
	}

	logger, err := c.createLogger(config, running)
	if err != nil {
		return err
	}

	cfg.Logger = logger
	c.logger = logger

	createMemFs := config.GetBool("renderToMemory")

	if createMemFs {
		// Rendering to memoryFS, publish to Root regardless of publishDir.
		config.Set("publishDir", "/")
	}

	c.fsCreate.Do(func() {
		fs := hugofs.NewFrom(sourceFs, config)

		if c.destinationFs != nil {
			// Need to reuse the destination on server rebuilds.
			fs.Destination = c.destinationFs
		} else if createMemFs {
			// Hugo writes the output to memory instead of the disk.
			fs.Destination = new(afero.MemMapFs)
		}

		if c.fastRenderMode {
			// For now, fast render mode only. It should, however, be fast enough
			// for the full variant, too.
			changeDetector := &fileChangeDetector{
				// We use this detector to decide to do a Hot reload of a single path or not.
				// We need to filter out source maps and possibly some other to be able
				// to make that decision.
				irrelevantRe: regexp.MustCompile(`\.map$`),
			}
			changeDetector.PrepareNew()
			fs.Destination = hugofs.NewHashingFs(fs.Destination, changeDetector)
			c.changeDetector = changeDetector
		}

		err = c.initFs(fs)
		if err != nil {
			return
		}

		var h *hugolib.HugoSites

		h, err = hugolib.NewHugoSites(*c.DepsCfg)
		c.hugo = h

	})

	if err != nil {
		return err
	}

	cacheDir, err := helpers.GetCacheDir(sourceFs, config)
	if err != nil {
		return err
	}
	config.Set("cacheDir", cacheDir)

	cfg.Logger.INFO.Println("Using config file:", config.ConfigFileUsed())

	themeDir := c.hugo.PathSpec.GetFirstThemeDir()
	if themeDir != "" {
		if _, err := sourceFs.Stat(themeDir); os.IsNotExist(err) {
			return newSystemError("Unable to find theme Directory:", themeDir)
		}
	}

	dir, themeVersionMismatch, minVersion := c.isThemeVsHugoVersionMismatch(sourceFs)

	if themeVersionMismatch {
		name := filepath.Base(dir)
		cfg.Logger.ERROR.Printf("%s theme does not support Hugo version %s. Minimum version required is %s\n",
			strings.ToUpper(name), hugo.CurrentVersion.ReleaseVersion(), minVersion)
	}

	return nil

}
