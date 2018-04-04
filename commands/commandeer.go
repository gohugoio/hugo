// Copyright 2017 The Hugo Authors. All rights reserved.
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
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/gohugoio/hugo/utils"

	"github.com/spf13/afero"

	"github.com/gohugoio/hugo/hugolib"

	"github.com/bep/debounce"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	src "github.com/gohugoio/hugo/source"
)

type commandeer struct {
	*deps.DepsCfg

	subCmdVs []*cobra.Command

	pathSpec    *helpers.PathSpec
	visitedURLs *types.EvictingStringQueue

	staticDirsConfig []*src.Dirs

	// We watch these for changes.
	configFiles []string

	doWithCommandeer func(c *commandeer) error

	// We can do this only once.
	fsCreate sync.Once

	// Used in cases where we get flooded with events in server mode.
	debounce func(f func())

	serverPorts []int
	languages   helpers.Languages

	configured bool
}

func (c *commandeer) Set(key string, value interface{}) {
	if c.configured {
		panic("commandeer cannot be changed")
	}
	c.Cfg.Set(key, value)
}

// PathSpec lazily creates a new PathSpec, as all the paths must
// be configured before it is created.
func (c *commandeer) PathSpec() *helpers.PathSpec {
	c.configured = true
	return c.pathSpec
}

func (c *commandeer) initFs(fs *hugofs.Fs) error {
	c.DepsCfg.Fs = fs
	ps, err := helpers.NewPathSpec(fs, c.Cfg)
	if err != nil {
		return err
	}
	c.pathSpec = ps

	dirsConfig, err := c.createStaticDirsConfig()
	if err != nil {
		return err
	}
	c.staticDirsConfig = dirsConfig

	return nil
}

func newCommandeer(running bool, doWithCommandeer func(c *commandeer) error, subCmdVs ...*cobra.Command) (*commandeer, error) {

	var rebuildDebouncer func(f func())
	if running {
		// The time value used is tested with mass content replacements in a fairly big Hugo site.
		// It is better to wait for some seconds in those cases rather than get flooded
		// with rebuilds.
		rebuildDebouncer, _ = debounce.New(4 * time.Second)
	}

	c := &commandeer{
		doWithCommandeer: doWithCommandeer,
		subCmdVs:         append([]*cobra.Command{hugoCmdV}, subCmdVs...),
		visitedURLs:      types.NewEvictingStringQueue(10),
		debounce:         rebuildDebouncer,
	}

	return c, c.loadConfig(running)
}

func (c *commandeer) loadConfig(running bool) error {

	if c.DepsCfg == nil {
		c.DepsCfg = &deps.DepsCfg{}
	}

	cfg := c.DepsCfg
	c.configured = false
	cfg.Running = running

	var dir string
	if source != "" {
		dir, _ = filepath.Abs(source)
	} else {
		dir, _ = os.Getwd()
	}

	var sourceFs afero.Fs = hugofs.Os
	if c.DepsCfg.Fs != nil {
		sourceFs = c.DepsCfg.Fs.Source
	}

	config, configFiles, err := hugolib.LoadConfig(hugolib.ConfigSourceDescriptor{Fs: sourceFs, Path: source, WorkingDir: dir, Filename: cfgFile})
	if err != nil {
		return err
	}

	c.Cfg = config
	c.configFiles = configFiles

	for _, cmdV := range c.subCmdVs {
		c.initializeFlags(cmdV)
	}

	if l, ok := c.Cfg.Get("languagesSorted").(helpers.Languages); ok {
		c.languages = l
	}

	if baseURL != "" {
		config.Set("baseURL", baseURL)
	}

	if c.doWithCommandeer != nil {
		err = c.doWithCommandeer(c)
	}

	if err != nil {
		return err
	}

	if len(disableKinds) > 0 {
		c.Set("disableKinds", disableKinds)
	}

	logger, err := createLogger(cfg.Cfg)
	if err != nil {
		return err
	}

	cfg.Logger = logger

	config.Set("logI18nWarnings", logI18nWarnings)

	if theme != "" {
		config.Set("theme", theme)
	}

	if themesDir != "" {
		config.Set("themesDir", themesDir)
	}

	if destination != "" {
		config.Set("publishDir", destination)
	}

	config.Set("workingDir", dir)

	if contentDir != "" {
		config.Set("contentDir", contentDir)
	}

	if layoutDir != "" {
		config.Set("layoutDir", layoutDir)
	}

	if cacheDir != "" {
		config.Set("cacheDir", cacheDir)
	}

	createMemFs := config.GetBool("renderToMemory")

	if createMemFs {
		// Rendering to memoryFS, publish to Root regardless of publishDir.
		config.Set("publishDir", "/")
	}

	c.fsCreate.Do(func() {
		fs := hugofs.NewFrom(sourceFs, config)

		// Hugo writes the output to memory instead of the disk.
		if createMemFs {
			fs.Destination = new(afero.MemMapFs)
		}

		err = c.initFs(fs)
	})

	if err != nil {
		return err
	}

	cacheDir = config.GetString("cacheDir")
	if cacheDir != "" {
		if helpers.FilePathSeparator != cacheDir[len(cacheDir)-1:] {
			cacheDir = cacheDir + helpers.FilePathSeparator
		}
		isDir, err := helpers.DirExists(cacheDir, sourceFs)
		utils.CheckErr(cfg.Logger, err)
		if !isDir {
			mkdir(cacheDir)
		}
		config.Set("cacheDir", cacheDir)
	} else {
		config.Set("cacheDir", helpers.GetTempDir("hugo_cache", sourceFs))
	}

	cfg.Logger.INFO.Println("Using config file:", config.ConfigFileUsed())

	themeDir := c.PathSpec().GetThemeDir()
	if themeDir != "" {
		if _, err := sourceFs.Stat(themeDir); os.IsNotExist(err) {
			return newSystemError("Unable to find theme Directory:", themeDir)
		}
	}

	themeVersionMismatch, minVersion := c.isThemeVsHugoVersionMismatch(sourceFs)

	if themeVersionMismatch {
		cfg.Logger.ERROR.Printf("Current theme does not support Hugo version %s. Minimum version required is %s\n",
			helpers.CurrentHugoVersion.ReleaseVersion(), minVersion)
	}

	return nil

}
