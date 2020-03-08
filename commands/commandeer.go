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

package commands

import (
	"bytes"
	"errors"
	"sync"

	hconfig "github.com/gohugoio/hugo/config"

	"golang.org/x/sync/semaphore"

	"io/ioutil"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugo"

	jww "github.com/spf13/jwalterweatherman"

	"os"
	"path/filepath"
	"regexp"
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
	hugoSites *hugolib.HugoSites
	fsCreate  sync.Once
	created   chan struct{}
}

type commandeer struct {
	*commandeerHugoState

	logger       *loggers.Logger
	serverConfig *config.Server

	// Currently only set when in "fast render mode". But it seems to
	// be fast enough that we could maybe just add it for all server modes.
	changeDetector *fileChangeDetector

	// We need to reuse this on server rebuilds.
	destinationFs afero.Fs

	h    *hugoBuilderCommon
	ftch flagsToConfigHandler

	visitedURLs *types.EvictingStringQueue

	cfgInit func(c *commandeer) error

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
	wasError            bool

	configured bool
	paused     bool

	fullRebuildSem *semaphore.Weighted

	// Any error from the last build.
	buildErr error
}

func newCommandeerHugoState() *commandeerHugoState {
	return &commandeerHugoState{
		created: make(chan struct{}),
	}
}

func (c *commandeerHugoState) hugo() *hugolib.HugoSites {
	<-c.created
	return c.hugoSites
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
		herrors.FprintStackTraceFromErr(&b, c.buildErr)
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

func newCommandeer(mustHaveConfigFile, running bool, h *hugoBuilderCommon, f flagsToConfigHandler, cfgInit func(c *commandeer) error, subCmdVs ...*cobra.Command) (*commandeer, error) {

	var rebuildDebouncer func(f func())
	if running {
		// The time value used is tested with mass content replacements in a fairly big Hugo site.
		// It is better to wait for some seconds in those cases rather than get flooded
		// with rebuilds.
		rebuildDebouncer = debounce.New(4 * time.Second)
	}

	out := ioutil.Discard
	if !h.quiet {
		out = os.Stdout
	}

	c := &commandeer{
		h:                   h,
		ftch:                f,
		commandeerHugoState: newCommandeerHugoState(),
		cfgInit:             cfgInit,
		visitedURLs:         types.NewEvictingStringQueue(10),
		debounce:            rebuildDebouncer,
		fullRebuildSem:      semaphore.NewWeighted(1),
		// This will be replaced later, but we need something to log to before the configuration is read.
		logger: loggers.NewLogger(jww.LevelError, jww.LevelError, out, ioutil.Discard, running),
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

	cfgSetAndInit := func(cfg config.Provider) error {
		c.Cfg = cfg
		if c.cfgInit == nil {
			return nil
		}
		err := c.cfgInit(c)
		return err
	}

	configPath := c.h.source
	if configPath == "" {
		configPath = dir
	}
	config, configFiles, err := hugolib.LoadConfig(
		hugolib.ConfigSourceDescriptor{
			Fs:           sourceFs,
			Logger:       c.logger,
			Path:         configPath,
			WorkingDir:   dir,
			Filename:     c.h.cfgFile,
			AbsConfigDir: c.h.getConfigDir(dir),
			Environ:      os.Environ(),
			Environment:  environment},
		cfgSetAndInit,
		doWithConfig)

	if err != nil && mustHaveConfigFile {
		return err
	} else if mustHaveConfigFile && len(configFiles) == 0 {
		return hugolib.ErrNoConfigFile
	}

	c.configFiles = configFiles

	if l, ok := c.Cfg.Get("languagesSorted").(langs.Languages); ok {
		c.languagesConfigured = true
		c.languages = l
	}

	// Set some commonly used flags
	c.doLiveReload = running && !c.Cfg.GetBool("disableLiveReload")
	c.fastRenderMode = c.doLiveReload && !c.Cfg.GetBool("disableFastRender")
	c.showErrorInBrowser = c.doLiveReload && !c.Cfg.GetBool("disableBrowserError")

	// This is potentially double work, but we need to do this one more time now
	// that all the languages have been configured.
	if c.cfgInit != nil {
		if err := c.cfgInit(c); err != nil {
			return err
		}
	}

	logger, err := c.createLogger(config, running)
	if err != nil {
		return err
	}

	cfg.Logger = logger
	c.logger = logger
	c.serverConfig = hconfig.DecodeServer(cfg.Cfg)

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

		if c.Cfg.GetBool("logPathWarnings") {
			fs.Destination = hugofs.NewCreateCountingFs(fs.Destination)
		}

		// To debug hard-to-find path issues.
		//fs.Destination = hugofs.NewStacktracerFs(fs.Destination, `fr/fr`)

		err = c.initFs(fs)
		if err != nil {
			close(c.created)
			return
		}

		var h *hugolib.HugoSites

		h, err = hugolib.NewHugoSites(*c.DepsCfg)
		c.hugoSites = h
		close(c.created)

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

	return nil

}
