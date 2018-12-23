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

// Package commands defines and implements command-line commands and flags
// used by Hugo. Commands and flags are implemented using Cobra.
package commands

import (
	"fmt"
	"io/ioutil"

	"os/signal"
	"sort"
	"sync/atomic"

	"github.com/gohugoio/hugo/common/hugo"
	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/common/terminal"

	"syscall"

	"github.com/gohugoio/hugo/hugolib/filesystems"

	"golang.org/x/sync/errgroup"

	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/parser/metadecoders"
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
		resp.Result = hugoCmd.c.hugo
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
func initializeConfig(mustHaveConfigFile, running bool,
	h *hugoBuilderCommon,
	f flagsToConfigHandler,
	doWithCommandeer func(c *commandeer) error) (*commandeer, error) {

	c, err := newCommandeer(mustHaveConfigFile, running, h, f, doWithCommandeer)
	if err != nil {
		return nil, err
	}

	return c, nil

}

func (c *commandeer) createLogger(cfg config.Provider, running bool) (*loggers.Logger, error) {
	var (
		logHandle       = ioutil.Discard
		logThreshold    = jww.LevelWarn
		logFile         = cfg.GetString("logFile")
		outHandle       = os.Stdout
		stdoutThreshold = jww.LevelWarn
	)

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

	return loggers.NewLogger(stdoutThreshold, logThreshold, outHandle, logHandle, running), nil
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
		"templateMetrics",
		"templateMetricsHints",

		// Moved from vars.
		"baseURL",
		"buildWatch",
		"cacheDir",
		"cfgFile",
		"contentDir",
		"debug",
		"destination",
		"disableKinds",
		"gc",
		"layoutDir",
		"logFile",
		"i18n-warnings",
		"quiet",
		"renderToMemory",
		"source",
		"theme",
		"themesDir",
		"verbose",
		"verboseLog",
	}

	// Will set a value even if it is the default.
	flagKeysForced := []string{
		"minify",
	}

	for _, key := range persFlagKeys {
		setValueFromFlag(cmd.PersistentFlags(), key, cfg, "", false)
	}
	for _, key := range flagKeys {
		setValueFromFlag(cmd.Flags(), key, cfg, "", false)
	}

	for _, key := range flagKeysForced {
		setValueFromFlag(cmd.Flags(), key, cfg, "", true)
	}

	// Set some "config aliases"
	setValueFromFlag(cmd.Flags(), "destination", cfg, "publishDir", false)
	setValueFromFlag(cmd.Flags(), "i18n-warnings", cfg, "logI18nWarnings", false)

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
		default:
			panic(fmt.Sprintf("update switch with %s", f.Value.Type()))
		}

	}
}

func isTerminal() bool {
	return terminal.IsTerminal(os.Stdout)

}
func ifTerminal(s string) string {
	if !isTerminal() {
		return ""
	}
	return s
}

func (c *commandeer) fullBuild() error {
	var (
		g         errgroup.Group
		langCount map[string]uint64
	)

	if !c.h.quiet {
		fmt.Print(ifTerminal(hideCursor) + "Building sites … ")
		if isTerminal() {
			defer func() {
				fmt.Print(showCursor + clearLine)
			}()
		}
	}

	copyStaticFunc := func() error {
		cnt, err := c.copyStatic()
		if err != nil {
			if !os.IsNotExist(err) {
				return errors.Wrap(err, "Error copying static files")
			}
			c.logger.INFO.Println("No Static directory found")
		}
		langCount = cnt
		langCount = cnt
		return nil
	}
	buildSitesFunc := func() error {
		if err := c.buildSites(); err != nil {
			return errors.Wrap(err, "Error building site")
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

	for _, s := range c.hugo.Sites {
		s.ProcessingStats.Static = langCount[s.Language.Lang]
	}

	if c.h.gc {
		count, err := c.hugo.GC()
		if err != nil {
			return err
		}
		for _, s := range c.hugo.Sites {
			// We have no way of knowing what site the garbage belonged to.
			s.ProcessingStats.Cleaned = uint64(count)
		}
	}

	return nil

}

func (c *commandeer) build() error {
	defer c.timeTrack(time.Now(), "Total")

	if err := c.fullBuild(); err != nil {
		return err
	}

	// TODO(bep) Feedback?
	if !c.h.quiet {
		fmt.Println()
		c.hugo.PrintProcessingStats(os.Stdout)
		fmt.Println()
	}

	if c.h.buildWatch {
		watchDirs, err := c.getDirList()
		if err != nil {
			return err
		}
		c.logger.FEEDBACK.Println("Watching for changes in", c.hugo.PathSpec.AbsPathify(c.Cfg.GetString("contentDir")))
		c.logger.FEEDBACK.Println("Press Ctrl+C to stop")
		watcher, err := c.newWatcher(watchDirs...)
		checkErr(c.Logger, err)
		defer watcher.Close()

		var sigs = make(chan os.Signal)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		<-sigs
	}

	return nil
}

func (c *commandeer) serverBuild() error {
	defer c.timeTrack(time.Now(), "Total")

	if err := c.fullBuild(); err != nil {
		return err
	}

	// TODO(bep) Feedback?
	if !c.h.quiet {
		fmt.Println()
		c.hugo.PrintProcessingStats(os.Stdout)
		fmt.Println()
	}

	return nil
}

func (c *commandeer) copyStatic() (map[string]uint64, error) {
	return c.doWithPublishDirs(c.copyStaticTo)
}

func (c *commandeer) doWithPublishDirs(f func(sourceFs *filesystems.SourceFilesystem) (uint64, error)) (map[string]uint64, error) {

	langCount := make(map[string]uint64)

	staticFilesystems := c.hugo.BaseFs.SourceFilesystems.Static

	if len(staticFilesystems) == 0 {
		c.logger.INFO.Println("No static directories found to sync")
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

func (c *commandeer) copyStaticTo(sourceFs *filesystems.SourceFilesystem) (uint64, error) {
	publishDir := c.hugo.PathSpec.PublishDir
	// If root, remove the second '/'
	if publishDir == "//" {
		publishDir = helpers.FilePathSeparator
	}

	if sourceFs.PublishFolder != "" {
		publishDir = filepath.Join(publishDir, sourceFs.PublishFolder)
	}

	fs := &countingStatFs{Fs: sourceFs.Fs}

	syncer := fsync.NewSyncer()
	syncer.NoTimes = c.Cfg.GetBool("noTimes")
	syncer.NoChmod = c.Cfg.GetBool("noChmod")
	syncer.SrcFs = fs
	syncer.DestFs = c.Fs.Destination
	// Now that we are using a unionFs for the static directories
	// We can effectively clean the publishDir on initial sync
	syncer.Delete = c.Cfg.GetBool("cleanDestinationDir")

	if syncer.Delete {
		c.logger.INFO.Println("removing all files from destination that don't exist in static dirs")

		syncer.DeleteFilter = func(f os.FileInfo) bool {
			return f.IsDir() && strings.HasPrefix(f.Name(), ".")
		}
	}
	c.logger.INFO.Println("syncing static files to", publishDir)

	var err error

	// because we are using a baseFs (to get the union right).
	// set sync src to root
	err = syncer.Sync(publishDir, helpers.FilePathSeparator)
	if err != nil {
		return 0, err
	}

	// Sync runs Stat 3 times for every source file (which sounds much)
	numFiles := fs.statCounter / 3

	return numFiles, err
}

func (c *commandeer) firstPathSpec() *helpers.PathSpec {
	return c.hugo.Sites[0].PathSpec
}

func (c *commandeer) timeTrack(start time.Time, name string) {
	if c.h.quiet {
		return
	}
	elapsed := time.Since(start)
	c.logger.FEEDBACK.Printf("%s in %v ms", name, int(1000*elapsed.Seconds()))
}

// getDirList provides NewWatcher() with a list of directories to watch for changes.
func (c *commandeer) getDirList() ([]string, error) {
	var a []string

	// To handle nested symlinked content dirs
	var seen = make(map[string]bool)
	var nested []string

	newWalker := func(allowSymbolicDirs bool) func(path string, fi os.FileInfo, err error) error {
		return func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}

				c.logger.ERROR.Println("Walker: ", err)
				return nil
			}

			// Skip .git directories.
			// Related to https://github.com/gohugoio/hugo/issues/3468.
			if fi.Name() == ".git" {
				return nil
			}

			if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
				link, err := filepath.EvalSymlinks(path)
				if err != nil {
					c.logger.ERROR.Printf("Cannot read symbolic link '%s', error was: %s", path, err)
					return nil
				}
				linkfi, err := helpers.LstatIfPossible(c.Fs.Source, link)
				if err != nil {
					c.logger.ERROR.Printf("Cannot stat %q: %s", link, err)
					return nil
				}
				if !allowSymbolicDirs && !linkfi.Mode().IsRegular() {
					c.logger.ERROR.Printf("Symbolic links for directories not supported, skipping %q", path)
					return nil
				}

				if allowSymbolicDirs && linkfi.IsDir() {
					// afero.Walk will not walk symbolic links, so wee need to do it.
					if !seen[path] {
						seen[path] = true
						nested = append(nested, path)
					}
					return nil
				}

				fi = linkfi
			}

			if fi.IsDir() {
				if fi.Name() == ".git" ||
					fi.Name() == "node_modules" || fi.Name() == "bower_components" {
					return filepath.SkipDir
				}
				a = append(a, path)
			}
			return nil
		}
	}

	symLinkWalker := newWalker(true)
	regularWalker := newWalker(false)

	// SymbolicWalk will log anny ERRORs
	// Also note that the Dirnames fetched below will contain any relevant theme
	// directories.
	for _, contentDir := range c.hugo.PathSpec.BaseFs.Content.Dirnames {
		_ = helpers.SymbolicWalk(c.Fs.Source, contentDir, symLinkWalker)
	}

	for _, staticDir := range c.hugo.PathSpec.BaseFs.Data.Dirnames {
		_ = helpers.SymbolicWalk(c.Fs.Source, staticDir, regularWalker)
	}

	for _, staticDir := range c.hugo.PathSpec.BaseFs.I18n.Dirnames {
		_ = helpers.SymbolicWalk(c.Fs.Source, staticDir, regularWalker)
	}

	for _, staticDir := range c.hugo.PathSpec.BaseFs.Layouts.Dirnames {
		_ = helpers.SymbolicWalk(c.Fs.Source, staticDir, regularWalker)
	}

	for _, staticFilesystem := range c.hugo.PathSpec.BaseFs.Static {
		for _, staticDir := range staticFilesystem.Dirnames {
			_ = helpers.SymbolicWalk(c.Fs.Source, staticDir, regularWalker)
		}
	}

	for _, assetDir := range c.hugo.PathSpec.BaseFs.Assets.Dirnames {
		_ = helpers.SymbolicWalk(c.Fs.Source, assetDir, regularWalker)
	}

	if len(nested) > 0 {
		for {

			toWalk := nested
			nested = nested[:0]

			for _, d := range toWalk {
				_ = helpers.SymbolicWalk(c.Fs.Source, d, symLinkWalker)
			}

			if len(nested) == 0 {
				break
			}
		}
	}

	a = helpers.UniqueStrings(a)
	sort.Strings(a)

	return a, nil
}

func (c *commandeer) resetAndBuildSites() (err error) {
	if !c.h.quiet {
		c.logger.FEEDBACK.Println("Started building sites ...")
	}
	return c.hugo.Build(hugolib.BuildCfg{ResetState: true})
}

func (c *commandeer) buildSites() (err error) {
	return c.hugo.Build(hugolib.BuildCfg{})
}

func (c *commandeer) handleBuildErr(err error, msg string) {
	c.buildErr = err

	c.logger.ERROR.Print(msg + ":\n\n")
	c.logger.ERROR.Println(helpers.FirstUpper(err.Error()))
	if !c.h.quiet && c.h.verbose {
		herrors.PrintStackTrace(err)
	}
}

func (c *commandeer) rebuildSites(events []fsnotify.Event) error {
	defer c.timeTrack(time.Now(), "Total")

	c.buildErr = nil
	visited := c.visitedURLs.PeekAllSet()
	if c.fastRenderMode {

		// Make sure we always render the home pages
		for _, l := range c.languages {
			langPath := c.hugo.PathSpec.GetLangSubDir(l.Lang)
			if langPath != "" {
				langPath = langPath + "/"
			}
			home := c.hugo.PathSpec.PrependBasePath("/"+langPath, false)
			visited[home] = true
		}

	}
	return c.hugo.Build(hugolib.BuildCfg{RecentlyVisited: visited}, events...)
}

func (c *commandeer) partialReRender(urls ...string) error {
	c.buildErr = nil
	visited := make(map[string]bool)
	for _, url := range urls {
		visited[url] = true
	}
	return c.hugo.Build(hugolib.BuildCfg{RecentlyVisited: visited, PartialReRender: true})
}

func (c *commandeer) fullRebuild() {
	c.commandeerHugoState = &commandeerHugoState{}
	err := c.loadConfig(true, true)
	if err != nil {
		// Set the processing on pause until the state is recovered.
		c.paused = true
		c.handleBuildErr(err, "Failed to reload config")

	} else {
		c.paused = false
	}

	if !c.paused {
		err := c.buildSites()
		if err != nil {
			c.logger.ERROR.Println(err)
		} else if !c.h.buildWatch && !c.Cfg.GetBool("disableLiveReload") {
			livereload.ForceRefresh()
		}
	}
}

// newWatcher creates a new watcher to watch filesystem events.
func (c *commandeer) newWatcher(dirList ...string) (*watcher.Batcher, error) {
	if runtime.GOOS == "darwin" {
		tweakLimit()
	}

	staticSyncer, err := newStaticSyncer(c)
	if err != nil {
		return nil, err
	}

	watcher, err := watcher.New(1 * time.Second)

	if err != nil {
		return nil, err
	}

	for _, d := range dirList {
		if d != "" {
			_ = watcher.Add(d)
		}
	}

	// Identifies changes to config (config.toml) files.
	configSet := make(map[string]bool)

	c.logger.FEEDBACK.Println("Watching for config changes in", strings.Join(c.configFiles, ", "))
	for _, configFile := range c.configFiles {
		watcher.Add(configFile)
		configSet[configFile] = true
	}

	go func() {
		for {
			select {
			case evs := <-watcher.Events:
				c.handleEvents(watcher, staticSyncer, evs, configSet)
				if c.showErrorInBrowser && c.errCount() > 0 {
					// Need to reload browser to show the error
					livereload.ForceRefresh()
				}
			case err := <-watcher.Errors:
				if err != nil {
					c.logger.ERROR.Println("Error while watching:", err)
				}
			}
		}
	}()

	return watcher, nil
}

func (c *commandeer) handleEvents(watcher *watcher.Batcher,
	staticSyncer *staticSyncer,
	evs []fsnotify.Event,
	configSet map[string]bool) {

	for _, ev := range evs {
		isConfig := configSet[ev.Name]
		if !isConfig {
			// It may be one of the /config folders
			dirname := filepath.Dir(ev.Name)
			if dirname != "." && configSet[dirname] {
				isConfig = true
			}

		}

		if isConfig {
			if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
				continue
			}
			if ev.Op&fsnotify.Remove == fsnotify.Remove {
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
			c.fullRebuild()
			break
		}
	}

	if c.paused {
		// Wait for the server to get into a consistent state before
		// we continue with processing.
		return
	}

	if len(evs) > 50 {
		// This is probably a mass edit of the content dir.
		// Schedule a full rebuild for when it slows down.
		c.debounce(c.fullRebuild)
		return
	}

	c.logger.INFO.Println("Received System Events:", evs)

	staticEvents := []fsnotify.Event{}
	dynamicEvents := []fsnotify.Event{}

	// Special handling for symbolic links inside /content.
	filtered := []fsnotify.Event{}
	for _, ev := range evs {
		// Check the most specific first, i.e. files.
		contentMapped := c.hugo.ContentChanges.GetSymbolicLinkMappings(ev.Name)
		if len(contentMapped) > 0 {
			for _, mapped := range contentMapped {
				filtered = append(filtered, fsnotify.Event{Name: mapped, Op: ev.Op})
			}
			continue
		}

		// Check for any symbolic directory mapping.

		dir, name := filepath.Split(ev.Name)

		contentMapped = c.hugo.ContentChanges.GetSymbolicLinkMappings(dir)

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
		if c.hugo.Deps.SourceSpec.IgnoreFile(ev.Name) {
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

		walkAdder := func(path string, f os.FileInfo, err error) error {
			if f.IsDir() {
				c.logger.FEEDBACK.Println("adding created directory to watchlist", path)
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
		c.logger.FEEDBACK.Println("\nStatic file changes detected")
		const layout = "2006-01-02 15:04:05.000 -0700"
		c.logger.FEEDBACK.Println(time.Now().Format(layout))

		if c.Cfg.GetBool("forceSyncStatic") {
			c.logger.FEEDBACK.Printf("Syncing all static files\n")
			_, err := c.copyStatic()
			if err != nil {
				c.logger.ERROR.Println("Error copying static files to publish dir:", err)
				return
			}
		} else {
			if err := staticSyncer.syncsStaticEvents(staticEvents); err != nil {
				c.logger.ERROR.Println("Error syncing static files to publish dir:", err)
				return
			}
		}

		if !c.h.buildWatch && !c.Cfg.GetBool("disableLiveReload") {
			// Will block forever trying to write to a channel that nobody is reading if livereload isn't initialized

			// force refresh when more than one file
			if len(staticEvents) == 1 {
				ev := staticEvents[0]
				path := c.hugo.BaseFs.SourceFilesystems.MakeStaticPathRelative(ev.Name)
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

		c.logger.FEEDBACK.Println("\nChange detected, rebuilding site")
		const layout = "2006-01-02 15:04:05.000 -0700"
		c.logger.FEEDBACK.Println(time.Now().Format(layout))

		c.changeDetector.PrepareNew()
		if err := c.rebuildSites(dynamicEvents); err != nil {
			c.handleBuildErr(err, "Rebuild failed")
		}

		if doLiveReload {
			if len(partitionedEvents.ContentEvents) == 0 && len(partitionedEvents.AssetEvents) > 0 {
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
				var p *hugolib.Page

				if navigate {
					if onePageName != "" {
						p = c.hugo.GetContentPage(onePageName)
					}
				}

				if p != nil {
					livereload.NavigateToPathForPort(p.RelPermalink(), p.Site.ServerPort())
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

	// Some editors (for example notepad.exe on Windows) triggers a change
	// both for directory and file. So we pick the longest path, which should
	// be the file itself.
	for _, ev := range events {
		if (ev.Op&fsnotify.Write == fsnotify.Write || ev.Op&fsnotify.Create == fsnotify.Create) && len(ev.Name) > len(name) {
			name = ev.Name
		}
	}

	return name
}

// isThemeVsHugoVersionMismatch returns whether the current Hugo version is
// less than any of the themes' min_version.
func (c *commandeer) isThemeVsHugoVersionMismatch(fs afero.Fs) (dir string, mismatch bool, requiredMinVersion string) {
	if !c.hugo.PathSpec.ThemeSet() {
		return
	}

	for _, absThemeDir := range c.hugo.BaseFs.AbsThemeDirs {

		path := filepath.Join(absThemeDir, "theme.toml")

		exists, err := helpers.Exists(path, fs)

		if err != nil || !exists {
			continue
		}

		b, err := afero.ReadFile(fs, path)

		tomlMeta, err := metadecoders.Default.UnmarshalToMap(b, metadecoders.TOML)

		if err != nil {
			continue
		}

		if minVersion, ok := tomlMeta["min_version"]; ok {
			if hugo.CompareVersion(minVersion) > 0 {
				return absThemeDir, true, fmt.Sprint(minVersion)
			}
		}

	}

	return
}
