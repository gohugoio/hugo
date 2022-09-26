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

// Package commands defines and implements command-line commands and flags
// used by Hugo. Commands and flags are implemented using Cobra.
package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gohugoio/hugo/hugofs/files"
	"github.com/gohugoio/hugo/tpl"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/types"

	"github.com/gohugoio/hugo/hugofs"

	"github.com/gohugoio/hugo/resources/page"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/terminal"

	"github.com/gohugoio/hugo/hugolib/filesystems"

	"golang.org/x/sync/errgroup"

	"github.com/gohugoio/hugo/config"

	flag "github.com/spf13/pflag"

	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/livereload"
	"github.com/gohugoio/hugo/watcher"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/fsync"
	jww "github.com/spf13/jwalterweatherman"
)

// The Response value from Execute.
type Response struct {
	// The build Result will only be set in the hugo build command.
	Result *hugolib.HugoSites

	// Err is set when the command failed to execute.
	Err error

	// The command that was executed.
	Cmd *cobra.Command
}

// IsUserError returns true is the Response error is a user error rather than a
// system error.
func (r Response) IsUserError() bool {
	return r.Err != nil && isUserError(r.Err)
}

// Execute adds all child commands to the root command HugoCmd and sets flags appropriately.
// The args are usually filled with os.Args[1:].
func Execute(args []string) Response {
	hugoCmd := newCommandsBuilder().addAll().build()
	cmd := hugoCmd.getCommand()
	cmd.SetArgs(args)

	c, err := cmd.ExecuteC()

	var resp Response

	if c == cmd && hugoCmd.c != nil {
		// Root command executed
		resp.Result = hugoCmd.c.hugo()
	}

	if err == nil {
		errCount := int(loggers.GlobalErrorCounter.Count())
		if errCount > 0 {
			err = fmt.Errorf("logged %d errors", errCount)
		} else if resp.Result != nil {
			errCount = resp.Result.NumLogErrors()
			if errCount > 0 {
				err = fmt.Errorf("logged %d errors", errCount)
			}
		}

	}

	resp.Err = err
	resp.Cmd = c

	return resp
}

// InitializeConfig initializes a config file with sensible default configuration flags.
func initializeConfig(mustHaveConfigFile, failOnInitErr, running bool,
	h *hugoBuilderCommon,
	f flagsToConfigHandler,
	cfgInit func(c *commandeer) error) (*commandeer, error) {
	c, err := newCommandeer(mustHaveConfigFile, failOnInitErr, running, h, f, cfgInit)
	if err != nil {
		return nil, err
	}

	if h := c.hugoTry(); h != nil {
		for _, s := range h.Sites {
			s.RegisterMediaTypes()
		}
	}

	return c, nil
}

func (c *commandeer) createLogger(cfg config.Provider) (loggers.Logger, error) {
	var (
		logHandle       = ioutil.Discard
		logThreshold    = jww.LevelWarn
		logFile         = cfg.GetString("logFile")
		outHandle       = ioutil.Discard
		stdoutThreshold = jww.LevelWarn
	)

	if !c.h.quiet {
		outHandle = os.Stdout
	}

	if c.h.verboseLog || c.h.logging || (c.h.logFile != "") {
		var err error
		if logFile != "" {
			logHandle, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				return nil, newSystemError("Failed to open log file:", logFile, err)
			}
		} else {
			logHandle, err = ioutil.TempFile("", "hugo")
			if err != nil {
				return nil, newSystemError(err)
			}
		}
	} else if !c.h.quiet && cfg.GetBool("verbose") {
		stdoutThreshold = jww.LevelInfo
	}

	if cfg.GetBool("debug") {
		stdoutThreshold = jww.LevelDebug
	}

	if c.h.verboseLog {
		logThreshold = jww.LevelInfo
		if cfg.GetBool("debug") {
			logThreshold = jww.LevelDebug
		}
	}

	loggers.InitGlobalLogger(stdoutThreshold, logThreshold, outHandle, logHandle)
	helpers.InitLoggers()

	return loggers.NewLogger(stdoutThreshold, logThreshold, outHandle, logHandle, c.running), nil
}

func initializeFlags(cmd *cobra.Command, cfg config.Provider) {
	persFlagKeys := []string{
		"debug",
		"verbose",
		"logFile",
		// Moved from vars
	}
	flagKeys := []string{
		"cleanDestinationDir",
		"buildDrafts",
		"buildFuture",
		"buildExpired",
		"clock",
		"uglyURLs",
		"canonifyURLs",
		"enableRobotsTXT",
		"enableGitInfo",
		"pluralizeListTitles",
		"preserveTaxonomyNames",
		"ignoreCache",
		"forceSyncStatic",
		"noTimes",
		"noChmod",
		"noBuildLock",
		"ignoreVendorPaths",
		"templateMetrics",
		"templateMetricsHints",

		// Moved from vars.
		"baseURL",
		"buildWatch",
		"cacheDir",
		"cfgFile",
		"confirm",
		"contentDir",
		"debug",
		"destination",
		"disableKinds",
		"dryRun",
		"force",
		"gc",
		"printI18nWarnings",
		"printUnusedTemplates",
		"invalidateCDN",
		"layoutDir",
		"logFile",
		"maxDeletes",
		"quiet",
		"renderToMemory",
		"source",
		"target",
		"theme",
		"themesDir",
		"verbose",
		"verboseLog",
		"duplicateTargetPaths",
	}

	for _, key := range persFlagKeys {
		setValueFromFlag(cmd.PersistentFlags(), key, cfg, "", false)
	}
	for _, key := range flagKeys {
		setValueFromFlag(cmd.Flags(), key, cfg, "", false)
	}

	setValueFromFlag(cmd.Flags(), "minify", cfg, "minifyOutput", true)

	// Set some "config aliases"
	setValueFromFlag(cmd.Flags(), "destination", cfg, "publishDir", false)
	setValueFromFlag(cmd.Flags(), "printI18nWarnings", cfg, "logI18nWarnings", false)
	setValueFromFlag(cmd.Flags(), "printPathWarnings", cfg, "logPathWarnings", false)
}

func setValueFromFlag(flags *flag.FlagSet, key string, cfg config.Provider, targetKey string, force bool) {
	key = strings.TrimSpace(key)
	if (force && flags.Lookup(key) != nil) || flags.Changed(key) {
		f := flags.Lookup(key)
		configKey := key
		if targetKey != "" {
			configKey = targetKey
		}
		// Gotta love this API.
		switch f.Value.Type() {
		case "bool":
			bv, _ := flags.GetBool(key)
			cfg.Set(configKey, bv)
		case "string":
			cfg.Set(configKey, f.Value.String())
		case "stringSlice":
			bv, _ := flags.GetStringSlice(key)
			cfg.Set(configKey, bv)
		case "int":
			iv, _ := flags.GetInt(key)
			cfg.Set(configKey, iv)
		default:
			panic(fmt.Sprintf("update switch with %s", f.Value.Type()))
		}

	}
}

func (c *commandeer) fullBuild(noBuildLock bool) error {
	var (
		g         errgroup.Group
		langCount map[string]uint64
	)

	if !c.h.quiet {
		fmt.Println("Start building sites … ")
		fmt.Println(hugo.BuildVersionString())
		if terminal.IsTerminal(os.Stdout) {
			defer func() {
				fmt.Print(showCursor + clearLine)
			}()
		}
	}

	copyStaticFunc := func() error {
		cnt, err := c.copyStatic()
		if err != nil {
			return fmt.Errorf("Error copying static files: %w", err)
		}
		langCount = cnt
		return nil
	}
	buildSitesFunc := func() error {
		if err := c.buildSites(noBuildLock); err != nil {
			return fmt.Errorf("Error building site: %w", err)
		}
		return nil
	}
	// Do not copy static files and build sites in parallel if cleanDestinationDir is enabled.
	// This flag deletes all static resources in /public folder that are missing in /static,
	// and it does so at the end of copyStatic() call.
	if c.Cfg.GetBool("cleanDestinationDir") {
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

	for _, s := range c.hugo().Sites {
		s.ProcessingStats.Static = langCount[s.Language().Lang]
	}

	if c.h.gc {
		count, err := c.hugo().GC()
		if err != nil {
			return err
		}
		for _, s := range c.hugo().Sites {
			// We have no way of knowing what site the garbage belonged to.
			s.ProcessingStats.Cleaned = uint64(count)
		}
	}

	return nil
}

func (c *commandeer) initCPUProfile() (func(), error) {
	if c.h.cpuprofile == "" {
		return nil, nil
	}

	f, err := os.Create(c.h.cpuprofile)
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

func (c *commandeer) initMemProfile() {
	if c.h.memprofile == "" {
		return
	}

	f, err := os.Create(c.h.memprofile)
	if err != nil {
		c.logger.Errorf("could not create memory profile: ", err)
	}
	defer f.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		c.logger.Errorf("could not write memory profile: ", err)
	}
}

func (c *commandeer) initTraceProfile() (func(), error) {
	if c.h.traceprofile == "" {
		return nil, nil
	}

	f, err := os.Create(c.h.traceprofile)
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

func (c *commandeer) initMutexProfile() (func(), error) {
	if c.h.mutexprofile == "" {
		return nil, nil
	}

	f, err := os.Create(c.h.mutexprofile)
	if err != nil {
		return nil, err
	}

	runtime.SetMutexProfileFraction(1)

	return func() {
		pprof.Lookup("mutex").WriteTo(f, 0)
		f.Close()
	}, nil
}

func (c *commandeer) initMemTicker() func() {
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

func (c *commandeer) initProfiling() (func(), error) {
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
	if c.h.printm {
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

func (c *commandeer) build() error {
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

	if !c.h.quiet {
		fmt.Println()
		c.hugo().PrintProcessingStats(os.Stdout)
		fmt.Println()

		hugofs.WalkFilesystems(c.publishDirFs, func(fs afero.Fs) bool {
			if dfs, ok := fs.(hugofs.DuplicatesReporter); ok {
				dupes := dfs.ReportDuplicates()
				if dupes != "" {
					c.logger.Warnln("Duplicate target paths:", dupes)
				}
			}
			return false
		})

		unusedTemplates := c.hugo().Tmpl().(tpl.UnusedTemplatesProvider).UnusedTemplates()
		for _, unusedTemplate := range unusedTemplates {
			c.logger.Warnf("Template %s is unused, source file %s", unusedTemplate.Name(), unusedTemplate.Filename())
		}
	}

	if c.h.buildWatch {
		watchDirs, err := c.getDirList()
		if err != nil {
			return err
		}

		baseWatchDir := c.Cfg.GetString("workingDir")
		rootWatchDirs := getRootWatchDirsStr(baseWatchDir, watchDirs)

		c.logger.Printf("Watching for changes in %s%s{%s}\n", baseWatchDir, helpers.FilePathSeparator, rootWatchDirs)
		c.logger.Println("Press Ctrl+C to stop")
		watcher, err := c.newWatcher(c.h.poll, watchDirs...)
		checkErr(c.Logger, err)
		defer watcher.Close()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		<-sigs
	}

	return nil
}

func (c *commandeer) serverBuild() error {
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

	// TODO(bep) Feedback?
	if !c.h.quiet {
		fmt.Println()
		c.hugo().PrintProcessingStats(os.Stdout)
		fmt.Println()
	}

	return nil
}

func (c *commandeer) copyStatic() (map[string]uint64, error) {
	m, err := c.doWithPublishDirs(c.copyStaticTo)
	if err == nil || os.IsNotExist(err) {
		return m, nil
	}
	return m, err
}

func (c *commandeer) doWithPublishDirs(f func(sourceFs *filesystems.SourceFilesystem) (uint64, error)) (map[string]uint64, error) {
	langCount := make(map[string]uint64)

	staticFilesystems := c.hugo().BaseFs.SourceFilesystems.Static

	if len(staticFilesystems) == 0 {
		c.logger.Infoln("No static directories found to sync")
		return langCount, nil
	}

	for lang, fs := range staticFilesystems {
		cnt, err := f(fs)
		if err != nil {
			return langCount, err
		}

		if lang == "" {
			// Not multihost
			for _, l := range c.languages {
				langCount[l.Lang] = cnt
			}
		} else {
			langCount[lang] = cnt
		}
	}

	return langCount, nil
}

type countingStatFs struct {
	afero.Fs
	statCounter uint64
}

func (fs *countingStatFs) Stat(name string) (os.FileInfo, error) {
	f, err := fs.Fs.Stat(name)
	if err == nil {
		if !f.IsDir() {
			atomic.AddUint64(&fs.statCounter, 1)
		}
	}
	return f, err
}

func chmodFilter(dst, src os.FileInfo) bool {
	// Hugo publishes data from multiple sources, potentially
	// with overlapping directory structures. We cannot sync permissions
	// for directories as that would mean that we might end up with write-protected
	// directories inside /public.
	// One example of this would be syncing from the Go Module cache,
	// which have 0555 directories.
	return src.IsDir()
}

func (c *commandeer) copyStaticTo(sourceFs *filesystems.SourceFilesystem) (uint64, error) {
	publishDir := helpers.FilePathSeparator

	if sourceFs.PublishFolder != "" {
		publishDir = filepath.Join(publishDir, sourceFs.PublishFolder)
	}

	fs := &countingStatFs{Fs: sourceFs.Fs}

	syncer := fsync.NewSyncer()
	syncer.NoTimes = c.Cfg.GetBool("noTimes")
	syncer.NoChmod = c.Cfg.GetBool("noChmod")
	syncer.ChmodFilter = chmodFilter
	syncer.SrcFs = fs
	syncer.DestFs = c.Fs.PublishDirStatic
	// Now that we are using a unionFs for the static directories
	// We can effectively clean the publishDir on initial sync
	syncer.Delete = c.Cfg.GetBool("cleanDestinationDir")

	if syncer.Delete {
		c.logger.Infoln("removing all files from destination that don't exist in static dirs")

		syncer.DeleteFilter = func(f os.FileInfo) bool {
			return f.IsDir() && strings.HasPrefix(f.Name(), ".")
		}
	}
	c.logger.Infoln("syncing static files to", publishDir)

	// because we are using a baseFs (to get the union right).
	// set sync src to root
	err := syncer.Sync(publishDir, helpers.FilePathSeparator)
	if err != nil {
		return 0, err
	}

	// Sync runs Stat 3 times for every source file (which sounds much)
	numFiles := fs.statCounter / 3

	return numFiles, err
}

func (c *commandeer) firstPathSpec() *helpers.PathSpec {
	return c.hugo().Sites[0].PathSpec
}

func (c *commandeer) timeTrack(start time.Time, name string) {
	// Note the use of time.Since here and time.Now in the callers.
	// We have a htime.Sinnce, but that may be adjusted to the future,
	// and that does not make sense here, esp. when used before the
	// global Clock is initialized.
	elapsed := time.Since(start)
	c.logger.Printf("%s in %v ms", name, int(1000*elapsed.Seconds()))
}

// getDirList provides NewWatcher() with a list of directories to watch for changes.
func (c *commandeer) getDirList() ([]string, error) {
	var filenames []string

	walkFn := func(path string, fi hugofs.FileMetaInfo, err error) error {
		if err != nil {
			c.logger.Errorln("walker: ", err)
			return nil
		}

		if fi.IsDir() {
			if fi.Name() == ".git" ||
				fi.Name() == "node_modules" || fi.Name() == "bower_components" {
				return filepath.SkipDir
			}

			filenames = append(filenames, fi.Meta().Filename)
		}

		return nil
	}

	watchFiles := c.hugo().PathSpec.BaseFs.WatchDirs()
	for _, fi := range watchFiles {
		if !fi.IsDir() {
			filenames = append(filenames, fi.Meta().Filename)
			continue
		}

		w := hugofs.NewWalkway(hugofs.WalkwayConfig{Logger: c.logger, Info: fi, WalkFn: walkFn})
		if err := w.Walk(); err != nil {
			c.logger.Errorln("walker: ", err)
		}
	}

	filenames = helpers.UniqueStringsSorted(filenames)

	return filenames, nil
}

func (c *commandeer) buildSites(noBuildLock bool) (err error) {
	return c.hugo().Build(hugolib.BuildCfg{NoBuildLock: noBuildLock})
}

func (c *commandeer) handleBuildErr(err error, msg string) {
	c.buildErr = err
	c.logger.Errorln(msg + ": " + cleanErrorLog(err.Error()))
}

func (c *commandeer) rebuildSites(events []fsnotify.Event) error {
	if c.buildErr != nil {
		ferrs := herrors.UnwrapFileErrorsWithErrorContext(c.buildErr)
		for _, err := range ferrs {
			events = append(events, fsnotify.Event{Name: err.Position().Filename, Op: fsnotify.Write})
		}
	}
	c.buildErr = nil
	visited := c.visitedURLs.PeekAllSet()
	if c.fastRenderMode {
		// Make sure we always render the home pages
		for _, l := range c.languages {
			langPath := c.hugo().PathSpec.GetLangSubDir(l.Lang)
			if langPath != "" {
				langPath = langPath + "/"
			}
			home := c.hugo().PathSpec.PrependBasePath("/"+langPath, false)
			visited[home] = true
		}
	}
	return c.hugo().Build(hugolib.BuildCfg{NoBuildLock: true, RecentlyVisited: visited, ErrRecovery: c.wasError}, events...)
}

func (c *commandeer) partialReRender(urls ...string) error {
	defer func() {
		c.wasError = false
	}()
	c.buildErr = nil
	visited := make(map[string]bool)
	for _, url := range urls {
		visited[url] = true
	}

	// Note: We do not set NoBuildLock as the file lock is not acquired at this stage.
	return c.hugo().Build(hugolib.BuildCfg{NoBuildLock: false, RecentlyVisited: visited, PartialReRender: true, ErrRecovery: c.wasError})
}

func (c *commandeer) fullRebuild(changeType string) {
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
			// Allow any file system events to arrive back.
			// This will block any rebuild on config changes for the
			// duration of the sleep.
			time.Sleep(2 * time.Second)
		}()

		defer c.timeTrack(time.Now(), "Rebuilt")

		c.commandeerHugoState = newCommandeerHugoState()
		err := c.loadConfig()
		if err != nil {
			// Set the processing on pause until the state is recovered.
			c.paused = true
			c.handleBuildErr(err, "Failed to reload config")

		} else {
			c.paused = false
		}

		if !c.paused {
			_, err := c.copyStatic()
			if err != nil {
				c.logger.Errorln(err)
				return
			}

			err = c.buildSites(true)
			if err != nil {
				c.logger.Errorln(err)
			} else if !c.h.buildWatch && !c.Cfg.GetBool("disableLiveReload") {
				livereload.ForceRefresh()
			}
		}
	}()
}

// newWatcher creates a new watcher to watch filesystem events.
func (c *commandeer) newWatcher(pollIntervalStr string, dirList ...string) (*watcher.Batcher, error) {
	if runtime.GOOS == "darwin" {
		tweakLimit()
	}

	staticSyncer, err := newStaticSyncer(c)
	if err != nil {
		return nil, err
	}

	var pollInterval time.Duration
	poll := pollIntervalStr != ""
	if poll {
		pollInterval, err = types.ToDurationE(pollIntervalStr)
		if err != nil {
			return nil, fmt.Errorf("invalid value for flag poll: %s", err)
		}
		c.logger.Printf("Use watcher with poll interval %v", pollInterval)
	}

	if pollInterval == 0 {
		pollInterval = 500 * time.Millisecond
	}

	watcher, err := watcher.New(500*time.Millisecond, pollInterval, poll)
	if err != nil {
		return nil, err
	}

	spec := c.hugo().Deps.SourceSpec

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

	c.logger.Println("Watching for config changes in", strings.Join(c.configFiles, ", "))
	for _, configFile := range c.configFiles {
		watcher.Add(configFile)
		configSet[configFile] = true
	}

	go func() {
		for {
			select {
			case evs := <-watcher.Events:
				unlock, err := c.buildLock()
				if err != nil {
					c.logger.Errorln("Failed to acquire a build lock: %s", err)
					return
				}
				c.handleEvents(watcher, staticSyncer, evs, configSet)
				if c.showErrorInBrowser && c.errCount() > 0 {
					// Need to reload browser to show the error
					livereload.ForceRefresh()
				}
				unlock()
			case err := <-watcher.Errors():
				if err != nil && !os.IsNotExist(err) {
					c.logger.Errorln("Error while watching:", err)
				}
			}
		}
	}()

	return watcher, nil
}

func (c *commandeer) printChangeDetected(typ string) {
	msg := "\nChange"
	if typ != "" {
		msg += " of " + typ
	}
	msg += " detected, rebuilding site."

	c.logger.Println(msg)
	const layout = "2006-01-02 15:04:05.000 -0700"
	c.logger.Println(htime.Now().Format(layout))
}

const (
	configChangeConfig = "config file"
	configChangeGoMod  = "go.mod file"
)

func (c *commandeer) handleEvents(watcher *watcher.Batcher,
	staticSyncer *staticSyncer,
	evs []fsnotify.Event,
	configSet map[string]bool) {
	defer func() {
		c.wasError = false
	}()

	var isHandled bool

	for _, ev := range evs {
		isConfig := configSet[ev.Name]
		configChangeType := configChangeConfig
		if isConfig {
			if strings.Contains(ev.Name, "go.mod") {
				configChangeType = configChangeGoMod
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
				for _, configFile := range c.configFiles {
					counter := 0
					for watcher.Add(configFile) != nil {
						counter++
						if counter >= 100 {
							break
						}
						time.Sleep(100 * time.Millisecond)
					}
				}
			}

			// Config file(s) changed. Need full rebuild.
			c.fullRebuild(configChangeType)

			return
		}
	}

	if isHandled {
		return
	}

	if c.paused {
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

	c.logger.Infoln("Received System Events:", evs)

	staticEvents := []fsnotify.Event{}
	dynamicEvents := []fsnotify.Event{}

	filtered := []fsnotify.Event{}
	for _, ev := range evs {
		if c.hugo().ShouldSkipFileChangeEvent(ev) {
			continue
		}
		// Check the most specific first, i.e. files.
		contentMapped := c.hugo().ContentChanges.GetSymbolicLinkMappings(ev.Name)
		if len(contentMapped) > 0 {
			for _, mapped := range contentMapped {
				filtered = append(filtered, fsnotify.Event{Name: mapped, Op: ev.Op})
			}
			continue
		}

		// Check for any symbolic directory mapping.

		dir, name := filepath.Split(ev.Name)

		contentMapped = c.hugo().ContentChanges.GetSymbolicLinkMappings(dir)

		if len(contentMapped) == 0 {
			filtered = append(filtered, ev)
			continue
		}

		for _, mapped := range contentMapped {
			mappedFilename := filepath.Join(mapped, name)
			filtered = append(filtered, fsnotify.Event{Name: mappedFilename, Op: ev.Op})
		}
	}

	evs = filtered

	for _, ev := range evs {
		ext := filepath.Ext(ev.Name)
		baseName := filepath.Base(ev.Name)
		istemp := strings.HasSuffix(ext, "~") ||
			(ext == ".swp") || // vim
			(ext == ".swx") || // vim
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
		if c.hugo().Deps.SourceSpec.IgnoreFile(ev.Name) {
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

		walkAdder := func(path string, f hugofs.FileMetaInfo, err error) error {
			if f.IsDir() {
				c.logger.Println("adding created directory to watchlist", path)
				if err := watcher.Add(path); err != nil {
					return err
				}
			} else if !staticSyncer.isStatic(path) {
				// Hugo's rebuilding logic is entirely file based. When you drop a new folder into
				// /content on OSX, the above logic will handle future watching of those files,
				// but the initial CREATE is lost.
				dynamicEvents = append(dynamicEvents, fsnotify.Event{Name: path, Op: fsnotify.Create})
			}
			return nil
		}

		// recursively add new directories to watch list
		// When mkdir -p is used, only the top directory triggers an event (at least on OSX)
		if ev.Op&fsnotify.Create == fsnotify.Create {
			if s, err := c.Fs.Source.Stat(ev.Name); err == nil && s.Mode().IsDir() {
				_ = helpers.SymbolicWalk(c.Fs.Source, ev.Name, walkAdder)
			}
		}

		if staticSyncer.isStatic(ev.Name) {
			staticEvents = append(staticEvents, ev)
		} else {
			dynamicEvents = append(dynamicEvents, ev)
		}
	}

	if len(staticEvents) > 0 {
		c.printChangeDetected("Static files")

		if c.Cfg.GetBool("forceSyncStatic") {
			c.logger.Printf("Syncing all static files\n")
			_, err := c.copyStatic()
			if err != nil {
				c.logger.Errorln("Error copying static files to publish dir:", err)
				return
			}
		} else {
			if err := staticSyncer.syncsStaticEvents(staticEvents); err != nil {
				c.logger.Errorln("Error syncing static files to publish dir:", err)
				return
			}
		}

		if !c.h.buildWatch && !c.Cfg.GetBool("disableLiveReload") {
			// Will block forever trying to write to a channel that nobody is reading if livereload isn't initialized

			// force refresh when more than one file
			if !c.wasError && len(staticEvents) == 1 {
				ev := staticEvents[0]
				path := c.hugo().BaseFs.SourceFilesystems.MakeStaticPathRelative(ev.Name)
				path = c.firstPathSpec().RelURL(helpers.ToSlashTrimLeading(path), false)

				livereload.RefreshPath(path)
			} else {
				livereload.ForceRefresh()
			}
		}
	}

	if len(dynamicEvents) > 0 {
		partitionedEvents := partitionDynamicEvents(
			c.firstPathSpec().BaseFs.SourceFilesystems,
			dynamicEvents)

		doLiveReload := !c.h.buildWatch && !c.Cfg.GetBool("disableLiveReload")
		onePageName := pickOneWriteOrCreatePath(partitionedEvents.ContentEvents)

		c.printChangeDetected("")
		c.changeDetector.PrepareNew()

		func() {
			defer c.timeTrack(time.Now(), "Total")
			if err := c.rebuildSites(dynamicEvents); err != nil {
				c.handleBuildErr(err, "Rebuild failed")
			}
		}()

		if doLiveReload {
			if len(partitionedEvents.ContentEvents) == 0 && len(partitionedEvents.AssetEvents) > 0 {
				if c.wasError {
					livereload.ForceRefresh()
					return
				}
				changed := c.changeDetector.changed()
				if c.changeDetector != nil && len(changed) == 0 {
					// Nothing has changed.
					return
				} else if len(changed) == 1 {
					pathToRefresh := c.firstPathSpec().RelURL(helpers.ToSlashTrimLeading(changed[0]), false)
					livereload.RefreshPath(pathToRefresh)
				} else {
					livereload.ForceRefresh()
				}
			}

			if len(partitionedEvents.ContentEvents) > 0 {

				navigate := c.Cfg.GetBool("navigateToChanged")
				// We have fetched the same page above, but it may have
				// changed.
				var p page.Page

				if navigate {
					if onePageName != "" {
						p = c.hugo().GetContentPage(onePageName)
					}
				}

				if p != nil {
					livereload.NavigateToPathForPort(p.RelPermalink(), p.Site().ServerPort())
				} else {
					livereload.ForceRefresh()
				}
			}
		}
	}
}

// dynamicEvents contains events that is considered dynamic, as in "not static".
// Both of these categories will trigger a new build, but the asset events
// does not fit into the "navigate to changed" logic.
type dynamicEvents struct {
	ContentEvents []fsnotify.Event
	AssetEvents   []fsnotify.Event
}

func partitionDynamicEvents(sourceFs *filesystems.SourceFilesystems, events []fsnotify.Event) (de dynamicEvents) {
	for _, e := range events {
		if sourceFs.IsAsset(e.Name) {
			de.AssetEvents = append(de.AssetEvents, e)
		} else {
			de.ContentEvents = append(de.ContentEvents, e)
		}
	}
	return
}

func pickOneWriteOrCreatePath(events []fsnotify.Event) string {
	name := ""

	for _, ev := range events {
		if ev.Op&fsnotify.Write == fsnotify.Write || ev.Op&fsnotify.Create == fsnotify.Create {
			if files.IsIndexContentFile(ev.Name) {
				return ev.Name
			}

			if files.IsContentFile(ev.Name) {
				name = ev.Name
			}

		}
	}

	return name
}

func formatByteCount(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
