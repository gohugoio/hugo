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

package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bep/logg"
	"github.com/bep/simplecobra"
	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/common/terminal"
	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/identity"
	"github.com/gohugoio/hugo/livereload"
	"github.com/gohugoio/hugo/resources/page"
	"github.com/gohugoio/hugo/watcher"
	"github.com/spf13/fsync"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type hugoBuilder struct {
	r *rootCommand

	confmu sync.Mutex
	conf   *commonConfig

	// May be nil.
	s *serverCommand

	// Currently only set when in "fast render mode".
	changeDetector *fileChangeDetector
	visitedURLs    *types.EvictingStringQueue

	fullRebuildSem *semaphore.Weighted
	debounce       func(f func())

	onConfigLoaded func(reloaded bool) error

	fastRenderMode     bool
	showErrorInBrowser bool

	errState hugoBuilderErrState
}

var errConfigNotSet = errors.New("config not set")

func (c *hugoBuilder) withConfE(fn func(conf *commonConfig) error) error {
	c.confmu.Lock()
	defer c.confmu.Unlock()
	if c.conf == nil {
		return errConfigNotSet
	}
	return fn(c.conf)
}

func (c *hugoBuilder) withConf(fn func(conf *commonConfig)) {
	c.confmu.Lock()
	defer c.confmu.Unlock()
	fn(c.conf)
}

type hugoBuilderErrState struct {
	mu       sync.Mutex
	paused   bool
	builderr error
	waserr   bool
}

func (e *hugoBuilderErrState) setPaused(p bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.paused = p
}

func (e *hugoBuilderErrState) isPaused() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.paused
}

func (e *hugoBuilderErrState) setBuildErr(err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.builderr = err
}

func (e *hugoBuilderErrState) buildErr() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.builderr
}

func (e *hugoBuilderErrState) setWasErr(w bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.waserr = w
}

func (e *hugoBuilderErrState) wasErr() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.waserr
}

func (c *hugoBuilder) errCount() int {
	return c.r.logger.LoggCount(logg.LevelError) + loggers.Log().LoggCount(logg.LevelError)
}

// getDirList provides NewWatcher() with a list of directories to watch for changes.
func (c *hugoBuilder) getDirList() ([]string, error) {
	h, err := c.hugo()
	if err != nil {
		return nil, err
	}

	return helpers.UniqueStringsSorted(h.PathSpec.BaseFs.WatchFilenames()), nil
}

func (c *hugoBuilder) initCPUProfile() (func(), error) {
	if c.r.cpuprofile == "" {
		return nil, nil
	}

	f, err := os.Create(c.r.cpuprofile)
	if err != nil {
		return nil, fmt.Errorf("failed to create CPU profile: %w", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		return nil, fmt.Errorf("failed to start CPU profile: %w", err)
	}
	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}, nil
}

func (c *hugoBuilder) initMemProfile() {
	if c.r.memprofile == "" {
		return
	}

	f, err := os.Create(c.r.memprofile)
	if err != nil {
		c.r.logger.Errorf("could not create memory profile: ", err)
	}
	defer f.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		c.r.logger.Errorf("could not write memory profile: ", err)
	}
}

func (c *hugoBuilder) initMemTicker() func() {
	memticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	printMem := func() {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("\n\nAlloc = %v\nTotalAlloc = %v\nSys = %v\nNumGC = %v\n\n", formatByteCount(m.Alloc), formatByteCount(m.TotalAlloc), formatByteCount(m.Sys), m.NumGC)
	}

	go func() {
		for {
			select {
			case <-memticker.C:
				printMem()
			case <-quit:
				memticker.Stop()
				printMem()
				return
			}
		}
	}()

	return func() {
		close(quit)
	}
}

func (c *hugoBuilder) initMutexProfile() (func(), error) {
	if c.r.mutexprofile == "" {
		return nil, nil
	}

	f, err := os.Create(c.r.mutexprofile)
	if err != nil {
		return nil, err
	}

	runtime.SetMutexProfileFraction(1)

	return func() {
		pprof.Lookup("mutex").WriteTo(f, 0)
		f.Close()
	}, nil
}

func (c *hugoBuilder) initProfiling() (func(), error) {
	stopCPUProf, err := c.initCPUProfile()
	if err != nil {
		return nil, err
	}

	stopMutexProf, err := c.initMutexProfile()
	if err != nil {
		return nil, err
	}

	stopTraceProf, err := c.initTraceProfile()
	if err != nil {
		return nil, err
	}

	var stopMemTicker func()
	if c.r.printm {
		stopMemTicker = c.initMemTicker()
	}

	return func() {
		c.initMemProfile()

		if stopCPUProf != nil {
			stopCPUProf()
		}
		if stopMutexProf != nil {
			stopMutexProf()
		}

		if stopTraceProf != nil {
			stopTraceProf()
		}

		if stopMemTicker != nil {
			stopMemTicker()
		}
	}, nil
}

func (c *hugoBuilder) initTraceProfile() (func(), error) {
	if c.r.traceprofile == "" {
		return nil, nil
	}

	f, err := os.Create(c.r.traceprofile)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace file: %w", err)
	}

	if err := trace.Start(f); err != nil {
		return nil, fmt.Errorf("failed to start trace: %w", err)
	}

	return func() {
		trace.Stop()
		f.Close()
	}, nil
}

// newWatcher creates a new watcher to watch filesystem events.
func (c *hugoBuilder) newWatcher(pollIntervalStr string, dirList ...string) (*watcher.Batcher, error) {
	staticSyncer := &staticSyncer{c: c}

	var pollInterval time.Duration
	poll := pollIntervalStr != ""
	if poll {
		pollInterval, err := types.ToDurationE(pollIntervalStr)
		if err != nil {
			return nil, fmt.Errorf("invalid value for flag poll: %s", err)
		}
		c.r.logger.Printf("Use watcher with poll interval %v", pollInterval)
	}

	if pollInterval == 0 {
		pollInterval = 500 * time.Millisecond
	}

	watcher, err := watcher.New(500*time.Millisecond, pollInterval, poll)
	if err != nil {
		return nil, err
	}

	h, err := c.hugo()
	if err != nil {
		return nil, err
	}
	spec := h.Deps.SourceSpec

	for _, d := range dirList {
		if d != "" {
			if spec.IgnoreFile(d) {
				continue
			}
			_ = watcher.Add(d)
		}
	}

	// Identifies changes to config (config.toml) files.
	configSet := make(map[string]bool)
	var configFiles []string
	c.withConf(func(conf *commonConfig) {
		configFiles = conf.configs.LoadingInfo.ConfigFiles
	})

	c.r.Println("Watching for config changes in", strings.Join(configFiles, ", "))
	for _, configFile := range configFiles {
		watcher.Add(configFile)
		configSet[configFile] = true
	}

	go func() {
		for {
			select {
			case changes := <-c.r.changesFromBuild:
				c.errState.setBuildErr(nil)
				unlock, err := h.LockBuild()
				if err != nil {
					c.r.logger.Errorln("Failed to acquire a build lock: %s", err)
					return
				}
				c.changeDetector.PrepareNew()
				err = c.rebuildSitesForChanges(changes)
				if err != nil {
					c.r.logger.Errorln("Error while watching:", err)
				}
				if c.s != nil && c.s.doLiveReload {
					doReload := c.changeDetector == nil || len(c.changeDetector.changed()) > 0
					doReload = doReload || c.showErrorInBrowser && c.errCount() > 0
					if doReload {
						livereload.ForceRefresh()
					}
				}
				unlock()

			case evs := <-watcher.Events:
				unlock, err := h.LockBuild()
				if err != nil {
					c.r.logger.Errorln("Failed to acquire a build lock: %s", err)
					return
				}
				c.handleEvents(watcher, staticSyncer, evs, configSet)
				if c.showErrorInBrowser && c.errCount() > 0 {
					// Need to reload browser to show the error
					livereload.ForceRefresh()
				}
				unlock()
			case err := <-watcher.Errors():
				if err != nil && !herrors.IsNotExist(err) {
					c.r.logger.Errorln("Error while watching:", err)
				}
			}
		}
	}()

	return watcher, nil
}

func (c *hugoBuilder) build() error {
	stopProfiling, err := c.initProfiling()
	if err != nil {
		return err
	}

	defer func() {
		if stopProfiling != nil {
			stopProfiling()
		}
	}()

	if err := c.fullBuild(false); err != nil {
		return err
	}

	if !c.r.quiet {
		c.r.Println()
		h, err := c.hugo()
		if err != nil {
			return err
		}

		h.PrintProcessingStats(os.Stdout)
		c.r.Println()
	}

	return nil
}

func (c *hugoBuilder) buildSites(noBuildLock bool) (err error) {
	h, err := c.hugo()
	if err != nil {
		return err
	}
	return h.Build(hugolib.BuildCfg{NoBuildLock: noBuildLock})
}

func (c *hugoBuilder) copyStatic() (map[string]uint64, error) {
	m, err := c.doWithPublishDirs(c.copyStaticTo)
	if err == nil || herrors.IsNotExist(err) {
		return m, nil
	}
	return m, err
}

func (c *hugoBuilder) copyStaticTo(sourceFs *filesystems.SourceFilesystem) (uint64, error) {
	infol := c.r.logger.InfoCommand("static")
	publishDir := helpers.FilePathSeparator

	if sourceFs.PublishFolder != "" {
		publishDir = filepath.Join(publishDir, sourceFs.PublishFolder)
	}

	fs := &countingStatFs{Fs: sourceFs.Fs}

	syncer := fsync.NewSyncer()
	c.withConf(func(conf *commonConfig) {
		syncer.NoTimes = conf.configs.Base.NoTimes
		syncer.NoChmod = conf.configs.Base.NoChmod
		syncer.ChmodFilter = chmodFilter

		syncer.DestFs = conf.fs.PublishDirStatic
		// Now that we are using a unionFs for the static directories
		// We can effectively clean the publishDir on initial sync
		syncer.Delete = conf.configs.Base.CleanDestinationDir
	})

	syncer.SrcFs = fs

	if syncer.Delete {
		infol.Logf("removing all files from destination that don't exist in static dirs")

		syncer.DeleteFilter = func(f fsync.FileInfo) bool {
			return f.IsDir() && strings.HasPrefix(f.Name(), ".")
		}
	}
	start := time.Now()

	// because we are using a baseFs (to get the union right).
	// set sync src to root
	err := syncer.Sync(publishDir, helpers.FilePathSeparator)
	if err != nil {
		return 0, err
	}
	loggers.TimeTrackf(infol, start, nil, "syncing static files to %s", publishDir)

	// Sync runs Stat 2 times for every source file.
	numFiles := fs.statCounter / 2

	return numFiles, err
}

func (c *hugoBuilder) doWithPublishDirs(f func(sourceFs *filesystems.SourceFilesystem) (uint64, error)) (map[string]uint64, error) {
	langCount := make(map[string]uint64)

	h, err := c.hugo()
	if err != nil {
		return nil, err
	}
	staticFilesystems := h.BaseFs.SourceFilesystems.Static

	if len(staticFilesystems) == 0 {
		c.r.logger.Infoln("No static directories found to sync")
		return langCount, nil
	}

	for lang, fs := range staticFilesystems {
		cnt, err := f(fs)
		if err != nil {
			return langCount, err
		}
		if lang == "" {
			// Not multihost
			c.withConf(func(conf *commonConfig) {
				for _, l := range conf.configs.Languages {
					langCount[l.Lang] = cnt
				}
			})
		} else {
			langCount[lang] = cnt
		}
	}

	return langCount, nil
}

func (c *hugoBuilder) fullBuild(noBuildLock bool) error {
	var (
		g         errgroup.Group
		langCount map[string]uint64
	)

	c.r.logger.Println("Start building sites … ")
	c.r.logger.Println(hugo.BuildVersionString())
	c.r.logger.Println()
	if terminal.IsTerminal(os.Stdout) {
		defer func() {
			fmt.Print(showCursor + clearLine)
		}()
	}

	copyStaticFunc := func() error {
		cnt, err := c.copyStatic()
		if err != nil {
			return fmt.Errorf("error copying static files: %w", err)
		}
		langCount = cnt
		return nil
	}
	buildSitesFunc := func() error {
		if err := c.buildSites(noBuildLock); err != nil {
			return fmt.Errorf("error building site: %w", err)
		}
		return nil
	}
	// Do not copy static files and build sites in parallel if cleanDestinationDir is enabled.
	// This flag deletes all static resources in /public folder that are missing in /static,
	// and it does so at the end of copyStatic() call.
	var cleanDestinationDir bool
	c.withConf(func(conf *commonConfig) {
		cleanDestinationDir = conf.configs.Base.CleanDestinationDir
	})
	if cleanDestinationDir {
		if err := copyStaticFunc(); err != nil {
			return err
		}
		if err := buildSitesFunc(); err != nil {
			return err
		}
	} else {
		g.Go(copyStaticFunc)
		g.Go(buildSitesFunc)
		if err := g.Wait(); err != nil {
			return err
		}
	}

	h, err := c.hugo()
	if err != nil {
		return err
	}
	for _, s := range h.Sites {
		s.ProcessingStats.Static = langCount[s.Language().Lang]
	}

	if c.r.gc {
		count, err := h.GC()
		if err != nil {
			return err
		}
		for _, s := range h.Sites {
			// We have no way of knowing what site the garbage belonged to.
			s.ProcessingStats.Cleaned = uint64(count)
		}
	}

	return nil
}

func (c *hugoBuilder) fullRebuild(changeType string) {
	if changeType == configChangeGoMod {
		// go.mod may be changed during the build itself, and
		// we really want to prevent superfluous builds.
		if !c.fullRebuildSem.TryAcquire(1) {
			return
		}
		c.fullRebuildSem.Release(1)
	}

	c.fullRebuildSem.Acquire(context.Background(), 1)

	go func() {
		defer c.fullRebuildSem.Release(1)

		c.printChangeDetected(changeType)

		defer func() {
			// Allow any file system events to arrive basimplecobra.
			// This will block any rebuild on config changes for the
			// duration of the sleep.
			time.Sleep(2 * time.Second)
		}()

		defer c.postBuild("Rebuilt", time.Now())

		err := c.reloadConfig()
		if err != nil {
			// Set the processing on pause until the state is recovered.
			c.errState.setPaused(true)
			c.handleBuildErr(err, "Failed to reload config")
		} else {
			c.errState.setPaused(false)
		}

		if !c.errState.isPaused() {
			_, err := c.copyStatic()
			if err != nil {
				c.r.logger.Errorln(err)
				return
			}
			err = c.buildSites(false)
			if err != nil {
				c.r.logger.Errorln(err)
			} else if c.s != nil && c.s.doLiveReload {
				livereload.ForceRefresh()
			}
		}
	}()
}

func (c *hugoBuilder) handleBuildErr(err error, msg string) {
	c.errState.setBuildErr(err)
	c.r.logger.Errorln(msg + ": " + cleanErrorLog(err.Error()))
}

func (c *hugoBuilder) handleEvents(watcher *watcher.Batcher,
	staticSyncer *staticSyncer,
	evs []fsnotify.Event,
	configSet map[string]bool,
) {
	defer func() {
		c.errState.setWasErr(false)
	}()

	var isHandled bool

	// Filter out ghost events (from deleted, renamed directories).
	// This seems to be a bug in fsnotify, or possibly MacOS.
	var n int
	for _, ev := range evs {
		keep := true
		if ev.Has(fsnotify.Create) || ev.Has(fsnotify.Write) {
			if _, err := os.Stat(ev.Name); err != nil {
				keep = false
			}
		}
		if keep {
			evs[n] = ev
			n++
		}
	}
	evs = evs[:n]

	for _, ev := range evs {
		isConfig := configSet[ev.Name]
		configChangeType := configChangeConfig
		if isConfig {
			if strings.Contains(ev.Name, "go.mod") {
				configChangeType = configChangeGoMod
			}
			if strings.Contains(ev.Name, ".work") {
				configChangeType = configChangeGoWork
			}
		}
		if !isConfig {
			// It may be one of the /config folders
			dirname := filepath.Dir(ev.Name)
			if dirname != "." && configSet[dirname] {
				isConfig = true
			}
		}

		if isConfig {
			isHandled = true

			if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
				continue
			}

			if ev.Op&fsnotify.Remove == fsnotify.Remove || ev.Op&fsnotify.Rename == fsnotify.Rename {
				c.withConf(func(conf *commonConfig) {
					for _, configFile := range conf.configs.LoadingInfo.ConfigFiles {
						counter := 0
						for watcher.Add(configFile) != nil {
							counter++
							if counter >= 100 {
								break
							}
							time.Sleep(100 * time.Millisecond)
						}
					}
				})
			}

			// Config file(s) changed. Need full rebuild.
			c.fullRebuild(configChangeType)

			return
		}
	}

	if isHandled {
		return
	}

	if c.errState.isPaused() {
		// Wait for the server to get into a consistent state before
		// we continue with processing.
		return
	}

	if len(evs) > 50 {
		// This is probably a mass edit of the content dir.
		// Schedule a full rebuild for when it slows down.
		c.debounce(func() {
			c.fullRebuild("")
		})
		return
	}

	c.r.logger.Debugln("Received System Events:", evs)

	staticEvents := []fsnotify.Event{}
	dynamicEvents := []fsnotify.Event{}

	filterDuplicateEvents := func(evs []fsnotify.Event) []fsnotify.Event {
		seen := make(map[string]bool)
		var n int
		for _, ev := range evs {
			if seen[ev.Name] {
				continue
			}
			seen[ev.Name] = true
			evs[n] = ev
			n++
		}
		return evs[:n]
	}

	h, err := c.hugo()
	if err != nil {
		c.r.logger.Errorln("Error getting the Hugo object:", err)
		return
	}
	n = 0
	for _, ev := range evs {
		if h.ShouldSkipFileChangeEvent(ev) {
			continue
		}
		evs[n] = ev
		n++
	}
	evs = evs[:n]

	for _, ev := range evs {
		ext := filepath.Ext(ev.Name)
		baseName := filepath.Base(ev.Name)
		istemp := strings.HasSuffix(ext, "~") ||
			(ext == ".swp") || // vim
			(ext == ".swx") || // vim
			(ext == ".bck") || // helix
			(ext == ".tmp") || // generic temp file
			(ext == ".DS_Store") || // OSX Thumbnail
			baseName == "4913" || // vim
			strings.HasPrefix(ext, ".goutputstream") || // gnome
			strings.HasSuffix(ext, "jb_old___") || // intelliJ
			strings.HasSuffix(ext, "jb_tmp___") || // intelliJ
			strings.HasSuffix(ext, "jb_bak___") || // intelliJ
			strings.HasPrefix(ext, ".sb-") || // byword
			strings.HasPrefix(baseName, ".#") || // emacs
			strings.HasPrefix(baseName, "#") // emacs
		if istemp {
			continue
		}

		if h.Deps.SourceSpec.IgnoreFile(ev.Name) {
			continue
		}
		// Sometimes during rm -rf operations a '"": REMOVE' is triggered. Just ignore these
		if ev.Name == "" {
			continue
		}

		// Write and rename operations are often followed by CHMOD.
		// There may be valid use cases for rebuilding the site on CHMOD,
		// but that will require more complex logic than this simple conditional.
		// On OS X this seems to be related to Spotlight, see:
		// https://github.com/go-fsnotify/fsnotify/issues/15
		// A workaround is to put your site(s) on the Spotlight exception list,
		// but that may be a little mysterious for most end users.
		// So, for now, we skip reload on CHMOD.
		// We do have to check for WRITE though. On slower laptops a Chmod
		// could be aggregated with other important events, and we still want
		// to rebuild on those
		if ev.Op&(fsnotify.Chmod|fsnotify.Write|fsnotify.Create) == fsnotify.Chmod {
			continue
		}

		walkAdder := func(path string, f hugofs.FileMetaInfo) error {
			if f.IsDir() {
				c.r.logger.Println("adding created directory to watchlist", path)
				if err := watcher.Add(path); err != nil {
					return err
				}
			} else if !staticSyncer.isStatic(h, path) {
				// Hugo's rebuilding logic is entirely file based. When you drop a new folder into
				// /content on OSX, the above logic will handle future watching of those files,
				// but the initial CREATE is lost.
				dynamicEvents = append(dynamicEvents, fsnotify.Event{Name: path, Op: fsnotify.Create})
			}
			return nil
		}

		// recursively add new directories to watch list
		if ev.Has(fsnotify.Create) || ev.Has(fsnotify.Rename) {
			c.withConf(func(conf *commonConfig) {
				if s, err := conf.fs.Source.Stat(ev.Name); err == nil && s.Mode().IsDir() {
					_ = helpers.Walk(conf.fs.Source, ev.Name, walkAdder)
				}
			})
		}

		if staticSyncer.isStatic(h, ev.Name) {
			staticEvents = append(staticEvents, ev)
		} else {
			dynamicEvents = append(dynamicEvents, ev)
		}
	}

	lrl := c.r.logger.InfoCommand("livereload")

	staticEvents = filterDuplicateEvents(staticEvents)
	dynamicEvents = filterDuplicateEvents(dynamicEvents)

	if len(staticEvents) > 0 {
		c.printChangeDetected("Static files")

		if c.r.forceSyncStatic {
			c.r.logger.Printf("Syncing all static files\n")
			_, err := c.copyStatic()
			if err != nil {
				c.r.logger.Errorln("Error copying static files to publish dir:", err)
				return
			}
		} else {
			if err := staticSyncer.syncsStaticEvents(staticEvents); err != nil {
				c.r.logger.Errorln("Error syncing static files to publish dir:", err)
				return
			}
		}

		if c.s != nil && c.s.doLiveReload {
			// Will block forever trying to write to a channel that nobody is reading if livereload isn't initialized

			if !c.errState.wasErr() && len(staticEvents) == 1 {
				h, err := c.hugo()
				if err != nil {
					c.r.logger.Errorln("Error getting the Hugo object:", err)
					return
				}

				path := h.BaseFs.SourceFilesystems.MakeStaticPathRelative(staticEvents[0].Name)
				path = h.RelURL(paths.ToSlashTrimLeading(path), false)

				lrl.Logf("refreshing static file %q", path)
				livereload.RefreshPath(path)
			} else {
				lrl.Logf("got %d static file change events, force refresh", len(staticEvents))
				livereload.ForceRefresh()
			}
		}
	}

	if len(dynamicEvents) > 0 {
		partitionedEvents := partitionDynamicEvents(
			h.BaseFs.SourceFilesystems,
			dynamicEvents)

		onePageName := pickOneWriteOrCreatePath(h.Conf.ContentTypes(), partitionedEvents.ContentEvents)

		c.printChangeDetected("")
		c.changeDetector.PrepareNew()

		func() {
			defer c.postBuild("Total", time.Now())
			if err := c.rebuildSites(dynamicEvents); err != nil {
				c.handleBuildErr(err, "Rebuild failed")
			}
		}()

		if c.s != nil && c.s.doLiveReload {
			if c.errState.wasErr() {
				livereload.ForceRefresh()
				return
			}

			changed := c.changeDetector.changed()
			if c.changeDetector != nil {
				lrl.Logf("build changed %d files", len(changed))
				if len(changed) == 0 {
					// Nothing has changed.
					return
				}
			}

			// If this change set also contains one or more CSS files, we need to
			// refresh these as well.
			var cssChanges []string
			var otherChanges []string

			for _, ev := range changed {
				if strings.HasSuffix(ev, ".css") {
					cssChanges = append(cssChanges, ev)
				} else {
					otherChanges = append(otherChanges, ev)
				}
			}

			if len(partitionedEvents.ContentEvents) > 0 {
				navigate := c.s != nil && c.s.navigateToChanged
				// We have fetched the same page above, but it may have
				// changed.
				var p page.Page

				if navigate {
					if onePageName != "" {
						p = h.GetContentPage(onePageName)
					}
				}

				if p != nil && p.RelPermalink() != "" {
					link, port := p.RelPermalink(), p.Site().ServerPort()
					lrl.Logf("navigating to %q using port %d", link, port)
					livereload.NavigateToPathForPort(link, port)
				} else {
					lrl.Logf("no page to navigate to, force refresh")
					livereload.ForceRefresh()
				}
			} else if len(otherChanges) > 0 {
				if len(otherChanges) == 1 {
					// Allow single changes to be refreshed without a full page reload.
					pathToRefresh := h.PathSpec.RelURL(paths.ToSlashTrimLeading(otherChanges[0]), false)
					lrl.Logf("refreshing %q", pathToRefresh)
					livereload.RefreshPath(pathToRefresh)
				} else if len(cssChanges) == 0 {
					lrl.Logf("force refresh")
					livereload.ForceRefresh()
				}
			}

			if len(cssChanges) > 0 {
				// Allow some time for the live reload script to get reconnected.
				if len(otherChanges) > 0 {
					time.Sleep(200 * time.Millisecond)
				}
				for _, ev := range cssChanges {
					pathToRefresh := h.PathSpec.RelURL(paths.ToSlashTrimLeading(ev), false)
					lrl.Logf("refreshing CSS %q", pathToRefresh)
					livereload.RefreshPath(pathToRefresh)
				}
			}
		}
	}
}

func (c *hugoBuilder) postBuild(what string, start time.Time) {
	if h, err := c.hugo(); err == nil && h.Conf.Running() {
		h.LogServerAddresses()
	}
	c.r.timeTrack(start, what)
}

func (c *hugoBuilder) hugo() (*hugolib.HugoSites, error) {
	var h *hugolib.HugoSites
	if err := c.withConfE(func(conf *commonConfig) error {
		var err error
		h, err = c.r.HugFromConfig(conf)
		return err
	}); err != nil {
		return nil, err
	}

	if c.s != nil {
		// A running server, register the media types.
		for _, s := range h.Sites {
			s.RegisterMediaTypes()
		}
	}
	return h, nil
}

func (c *hugoBuilder) hugoTry() *hugolib.HugoSites {
	var h *hugolib.HugoSites
	c.withConf(func(conf *commonConfig) {
		h, _ = c.r.HugFromConfig(conf)
	})
	return h
}

func (c *hugoBuilder) loadConfig(cd *simplecobra.Commandeer, running bool) error {
	cfg := config.New()
	cfg.Set("renderToMemory", c.r.renderToMemory)
	watch := c.r.buildWatch || (c.s != nil && c.s.serverWatch)
	if c.r.environment == "" {
		// We need to set the environment as early as possible because we need it to load the correct config.
		// Check if the user has set it in env.
		if env := os.Getenv("HUGO_ENVIRONMENT"); env != "" {
			c.r.environment = env
		} else if env := os.Getenv("HUGO_ENV"); env != "" {
			c.r.environment = env
		} else {
			if c.s != nil {
				// The server defaults to development.
				c.r.environment = hugo.EnvironmentDevelopment
			} else {
				c.r.environment = hugo.EnvironmentProduction
			}
		}
	}
	cfg.Set("environment", c.r.environment)

	cfg.Set("internal", maps.Params{
		"running":        running,
		"watch":          watch,
		"verbose":        c.r.isVerbose(),
		"fastRenderMode": c.fastRenderMode,
	})

	conf, err := c.r.ConfigFromProvider(configKey{counter: c.r.configVersionID.Load()}, flagsToCfg(cd, cfg))
	if err != nil {
		return err
	}

	if len(conf.configs.LoadingInfo.ConfigFiles) == 0 {
		//lint:ignore ST1005 end user message.
		return errors.New("Unable to locate config file or config directory. Perhaps you need to create a new site.\nRun `hugo help new` for details.")
	}

	c.conf = conf
	if c.onConfigLoaded != nil {
		if err := c.onConfigLoaded(false); err != nil {
			return err
		}
	}

	return nil
}

var rebuildCounter atomic.Uint64

func (c *hugoBuilder) printChangeDetected(typ string) {
	msg := "\nChange"
	if typ != "" {
		msg += " of " + typ
	}
	msg += fmt.Sprintf(" detected, rebuilding site (#%d).", rebuildCounter.Add(1))

	c.r.logger.Println(msg)
	const layout = "2006-01-02 15:04:05.000 -0700"
	c.r.logger.Println(htime.Now().Format(layout))
}

func (c *hugoBuilder) rebuildSites(events []fsnotify.Event) error {
	if err := c.errState.buildErr(); err != nil {
		ferrs := herrors.UnwrapFileErrorsWithErrorContext(err)
		for _, err := range ferrs {
			events = append(events, fsnotify.Event{Name: err.Position().Filename, Op: fsnotify.Write})
		}
	}
	c.errState.setBuildErr(nil)
	h, err := c.hugo()
	if err != nil {
		return err
	}

	return h.Build(hugolib.BuildCfg{NoBuildLock: true, RecentlyVisited: c.visitedURLs, ErrRecovery: c.errState.wasErr()}, events...)
}

func (c *hugoBuilder) rebuildSitesForChanges(ids []identity.Identity) error {
	c.errState.setBuildErr(nil)
	h, err := c.hugo()
	if err != nil {
		return err
	}
	whatChanged := &hugolib.WhatChanged{}
	whatChanged.Add(ids...)
	err = h.Build(hugolib.BuildCfg{NoBuildLock: true, WhatChanged: whatChanged, RecentlyVisited: c.visitedURLs, ErrRecovery: c.errState.wasErr()})
	c.errState.setBuildErr(err)
	return err
}

func (c *hugoBuilder) reloadConfig() error {
	c.r.Reset()
	c.r.configVersionID.Add(1)

	if err := c.withConfE(func(conf *commonConfig) error {
		oldConf := conf
		newConf, err := c.r.ConfigFromConfig(configKey{counter: c.r.configVersionID.Load()}, conf)
		if err != nil {
			return err
		}
		sameLen := len(oldConf.configs.Languages) == len(newConf.configs.Languages)
		if !sameLen {
			if oldConf.configs.IsMultihost || newConf.configs.IsMultihost {
				return errors.New("multihost change detected, please restart server")
			}
		}
		c.conf = newConf
		return nil
	}); err != nil {
		return err
	}

	if c.onConfigLoaded != nil {
		if err := c.onConfigLoaded(true); err != nil {
			return err
		}
	}

	return nil
}
