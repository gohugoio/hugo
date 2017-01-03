// Copyright 2016 The Hugo Authors. All rights reserved.
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
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/spf13/hugo/tpl"

	"github.com/spf13/hugo/hugofs"

	"github.com/spf13/hugo/parser"
	flag "github.com/spf13/pflag"

	"regexp"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/fsync"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/livereload"
	"github.com/spf13/hugo/utils"
	"github.com/spf13/hugo/watcher"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/nitro"
	"github.com/spf13/viper"
)

// Hugo represents the Hugo sites to build. This variable is exported as it
// is used by at least one external library (the Hugo caddy plugin). We should
// provide a cleaner external API, but until then, this is it.
var Hugo *hugolib.HugoSites

// Reset resets Hugo ready for a new full build. This is mainly only useful
// for benchmark testing etc. via the CLI commands.
func Reset() error {
	Hugo = nil
	viper.Reset()
	return nil
}

// commandError is an error used to signal different error situations in command handling.
type commandError struct {
	s         string
	userError bool
}

func (c commandError) Error() string {
	return c.s
}

func (c commandError) isUserError() bool {
	return c.userError
}

func newUserError(a ...interface{}) commandError {
	return commandError{s: fmt.Sprintln(a...), userError: true}
}

func newSystemError(a ...interface{}) commandError {
	return commandError{s: fmt.Sprintln(a...), userError: false}
}

func newSystemErrorF(format string, a ...interface{}) commandError {
	return commandError{s: fmt.Sprintf(format, a...), userError: false}
}

// Catch some of the obvious user errors from Cobra.
// We don't want to show the usage message for every error.
// The below may be to generic. Time will show.
var userErrorRegexp = regexp.MustCompile("argument|flag|shorthand")

func isUserError(err error) bool {
	if cErr, ok := err.(commandError); ok && cErr.isUserError() {
		return true
	}

	return userErrorRegexp.MatchString(err.Error())
}

// HugoCmd is Hugo's root command.
// Every other command attached to HugoCmd is a child command to it.
var HugoCmd = &cobra.Command{
	Use:   "hugo",
	Short: "hugo builds your site",
	Long: `hugo is the main command, used to build your Hugo site.

Hugo is a Fast and Flexible Static Site Generator
built with love by spf13 and friends in Go.

Complete documentation is available at http://gohugo.io/.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := InitializeConfig(); err != nil {
			return err
		}

		if buildWatch {
			viper.Set("disableLiveReload", true)
			watchConfig()
		}

		return build()
	},
}

var hugoCmdV *cobra.Command

// Flags that are to be added to commands.
var (
	buildWatch     bool
	logging        bool
	renderToMemory bool // for benchmark testing
	verbose        bool
	verboseLog     bool
	quiet          bool
)

var (
	baseURL     string
	cacheDir    string
	contentDir  string
	layoutDir   string
	cfgFile     string
	destination string
	logFile     string
	theme       string
	themesDir   string
	source      string
)

// Execute adds all child commands to the root command HugoCmd and sets flags appropriately.
func Execute() {
	HugoCmd.SetGlobalNormalizationFunc(helpers.NormalizeHugoFlags)

	HugoCmd.SilenceUsage = true

	AddCommands()

	if c, err := HugoCmd.ExecuteC(); err != nil {
		if isUserError(err) {
			c.Println("")
			c.Println(c.UsageString())
		}

		os.Exit(-1)
	}
}

// AddCommands adds child commands to the root command HugoCmd.
func AddCommands() {
	HugoCmd.AddCommand(serverCmd)
	HugoCmd.AddCommand(versionCmd)
	HugoCmd.AddCommand(envCmd)
	HugoCmd.AddCommand(configCmd)
	HugoCmd.AddCommand(checkCmd)
	HugoCmd.AddCommand(benchmarkCmd)
	HugoCmd.AddCommand(convertCmd)
	HugoCmd.AddCommand(newCmd)
	HugoCmd.AddCommand(listCmd)
	HugoCmd.AddCommand(undraftCmd)
	HugoCmd.AddCommand(importCmd)

	HugoCmd.AddCommand(genCmd)
	genCmd.AddCommand(genautocompleteCmd)
	genCmd.AddCommand(gendocCmd)
	genCmd.AddCommand(genmanCmd)
}

// initHugoBuilderFlags initializes all common flags, typically used by the
// core build commands, namely hugo itself, server, check and benchmark.
func initHugoBuilderFlags(cmd *cobra.Command) {
	initHugoBuildCommonFlags(cmd)
}

func initRootPersistentFlags() {
	HugoCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is path/config.yaml|json|toml)")
	HugoCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "build in quiet mode")

	// Set bash-completion
	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	HugoCmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)
}

// initHugoBuildCommonFlags initialize common flags related to the Hugo build.
// Called by initHugoBuilderFlags.
func initHugoBuildCommonFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("cleanDestinationDir", false, "Remove files from destination not found in static directories")
	cmd.Flags().BoolP("buildDrafts", "D", false, "include content marked as draft")
	cmd.Flags().BoolP("buildFuture", "F", false, "include content with publishdate in the future")
	cmd.Flags().BoolP("buildExpired", "E", false, "include expired content")
	cmd.Flags().Bool("disable404", false, "Do not render 404 page")
	cmd.Flags().Bool("disableRSS", false, "Do not build RSS files")
	cmd.Flags().Bool("disableSitemap", false, "Do not build Sitemap file")
	cmd.Flags().StringVarP(&source, "source", "s", "", "filesystem path to read files relative from")
	cmd.Flags().StringVarP(&contentDir, "contentDir", "c", "", "filesystem path to content directory")
	cmd.Flags().StringVarP(&layoutDir, "layoutDir", "l", "", "filesystem path to layout directory")
	cmd.Flags().StringVarP(&cacheDir, "cacheDir", "", "", "filesystem path to cache directory. Defaults: $TMPDIR/hugo_cache/")
	cmd.Flags().BoolP("ignoreCache", "", false, "Ignores the cache directory")
	cmd.Flags().StringVarP(&destination, "destination", "d", "", "filesystem path to write files to")
	cmd.Flags().StringVarP(&theme, "theme", "t", "", "theme to use (located in /themes/THEMENAME/)")
	cmd.Flags().StringVarP(&themesDir, "themesDir", "", "", "filesystem path to themes directory")
	cmd.Flags().Bool("uglyURLs", false, "if true, use /filename.html instead of /filename/")
	cmd.Flags().Bool("canonifyURLs", false, "if true, all relative URLs will be canonicalized using baseURL")
	cmd.Flags().StringVarP(&baseURL, "baseURL", "b", "", "hostname (and path) to the root, e.g. http://spf13.com/")
	cmd.Flags().Bool("enableGitInfo", false, "Add Git revision, date and author info to the pages")

	cmd.Flags().BoolVar(&nitro.AnalysisOn, "stepAnalysis", false, "display memory and timing of different steps of the program")
	cmd.Flags().Bool("pluralizeListTitles", true, "Pluralize titles in lists using inflect")
	cmd.Flags().Bool("preserveTaxonomyNames", false, `Preserve taxonomy names as written ("Gérard Depardieu" vs "gerard-depardieu")`)
	cmd.Flags().BoolP("forceSyncStatic", "", false, "Copy all files when static is changed.")
	cmd.Flags().BoolP("noTimes", "", false, "Don't sync modification time of files")
	cmd.Flags().BoolP("noChmod", "", false, "Don't sync permission mode of files")
	cmd.Flags().BoolVarP(&tpl.Logi18nWarnings, "i18n-warnings", "", false, "Print missing translations")

	// Set bash-completion.
	// Each flag must first be defined before using the SetAnnotation() call.
	cmd.Flags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})
	cmd.Flags().SetAnnotation("cacheDir", cobra.BashCompSubdirsInDir, []string{})
	cmd.Flags().SetAnnotation("destination", cobra.BashCompSubdirsInDir, []string{})
	cmd.Flags().SetAnnotation("theme", cobra.BashCompSubdirsInDir, []string{"themes"})
}

func initBenchmarkBuildingFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&renderToMemory, "renderToMemory", false, "render to memory (only useful for benchmark testing)")
}

// init initializes flags.
func init() {
	HugoCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	HugoCmd.PersistentFlags().BoolVar(&logging, "log", false, "Enable Logging")
	HugoCmd.PersistentFlags().StringVar(&logFile, "logFile", "", "Log File path (if set, logging enabled automatically)")
	HugoCmd.PersistentFlags().BoolVar(&verboseLog, "verboseLog", false, "verbose logging")

	initRootPersistentFlags()
	initHugoBuilderFlags(HugoCmd)
	initBenchmarkBuildingFlags(HugoCmd)

	HugoCmd.Flags().BoolVarP(&buildWatch, "watch", "w", false, "watch filesystem for changes and recreate as needed")
	hugoCmdV = HugoCmd

	// Set bash-completion
	HugoCmd.PersistentFlags().SetAnnotation("logFile", cobra.BashCompFilenameExt, []string{})
}

// InitializeConfig initializes a config file with sensible default configuration flags.
func InitializeConfig(subCmdVs ...*cobra.Command) error {
	if err := hugolib.LoadGlobalConfig(source, cfgFile); err != nil {
		return err
	}

	for _, cmdV := range append([]*cobra.Command{hugoCmdV}, subCmdVs...) {
		initializeFlags(cmdV)
	}

	if baseURL != "" {
		if !strings.HasSuffix(baseURL, "/") {
			baseURL = baseURL + "/"
		}
		viper.Set("baseURL", baseURL)
	}

	if !viper.GetBool("relativeURLs") && viper.GetString("baseURL") == "" {
		jww.ERROR.Println("No 'baseURL' set in configuration or as a flag. Features like page menus will not work without one.")
	}

	if theme != "" {
		viper.Set("theme", theme)
	}

	if themesDir != "" {
		viper.Set("themesDir", themesDir)
	}

	if destination != "" {
		viper.Set("publishDir", destination)
	}

	var dir string
	if source != "" {
		dir, _ = filepath.Abs(source)
	} else {
		dir, _ = os.Getwd()
	}
	viper.Set("workingDir", dir)

	if contentDir != "" {
		viper.Set("contentDir", contentDir)
	}

	if layoutDir != "" {
		viper.Set("layoutDir", layoutDir)
	}

	if cacheDir != "" {
		viper.Set("cacheDir", cacheDir)
	}

	cacheDir = viper.GetString("cacheDir")
	if cacheDir != "" {
		if helpers.FilePathSeparator != cacheDir[len(cacheDir)-1:] {
			cacheDir = cacheDir + helpers.FilePathSeparator
		}
		isDir, err := helpers.DirExists(cacheDir, hugofs.Source())
		utils.CheckErr(err)
		if !isDir {
			mkdir(cacheDir)
		}
		viper.Set("cacheDir", cacheDir)
	} else {
		viper.Set("cacheDir", helpers.GetTempDir("hugo_cache", hugofs.Source()))
	}

	jww.INFO.Println("Using config file:", viper.ConfigFileUsed())

	// Init file systems. This may be changed at a later point.
	hugofs.InitDefaultFs()

	themeDir := helpers.GetThemeDir()
	if themeDir != "" {
		if _, err := hugofs.Source().Stat(themeDir); os.IsNotExist(err) {
			return newSystemError("Unable to find theme Directory:", themeDir)
		}
	}

	themeVersionMismatch, minVersion := isThemeVsHugoVersionMismatch()

	if themeVersionMismatch {
		jww.ERROR.Printf("Current theme does not support Hugo version %s. Minimum version required is %s\n",
			helpers.HugoReleaseVersion(), minVersion)
	}

	return nil

}

func createLogger() *jww.NotePad {
	var (
		logHandle         = ioutil.Discard
		outHandle = os.Stdout
		stdoutThreshold = jww.LevelError
		logThreshold    = jww.LevelWarn
	)

	if verboseLog || logging || (viper.IsSet("logFile") && viper.GetString("logFile") != "") {

		var err error
		if viper.IsSet("logFile") && viper.GetString("logFile") != "" {
			path := viper.GetString("logFile")
			logHandle, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				return newSystemError("Failed to open log file:", path, err)
			}
		} else {
			logHandle, err = ioutil.TempFile(os.TempDir(), "hugo")
			if err != nil {
				return newSystemError(err)
			}
		}
	}

	} else if !quiet && viper.GetBool("verbose") {
		stdoutThreshold = jww.LevelInfo
	}

	if verboseLog {
		logThreshold = jww.LevelInfo
	}
	
	return jww.NewNotepad(stdoutThreshold, logThreshold, outHandle, logHandle, "", log.Ldate|log.Ltime)
}

func initializeFlags(cmd *cobra.Command) {
	persFlagKeys := []string{"verbose", "logFile"}
	flagKeys := []string{
		"cleanDestinationDir",
		"buildDrafts",
		"buildFuture",
		"buildExpired",
		"uglyURLs",
		"canonifyURLs",
		"disable404",
		"disableRSS",
		"disableSitemap",
		"enableRobotsTXT",
		"enableGitInfo",
		"pluralizeListTitles",
		"preserveTaxonomyNames",
		"ignoreCache",
		"forceSyncStatic",
		"noTimes",
		"noChmod",
	}

	for _, key := range persFlagKeys {
		setValueFromFlag(cmd.PersistentFlags(), key)
	}
	for _, key := range flagKeys {
		setValueFromFlag(cmd.Flags(), key)
	}
}

func setValueFromFlag(flags *flag.FlagSet, key string) {
	if flagChanged(flags, key) {
		f := flags.Lookup(key)
		viper.Set(key, f.Value.String())
	}
}

func flagChanged(flags *flag.FlagSet, key string) bool {
	flag := flags.Lookup(key)
	if flag == nil {
		return false
	}
	return flag.Changed
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		jww.FEEDBACK.Println("Config file changed:", e.Name)
		// Force a full rebuild
		utils.CheckErr(recreateAndBuildSites(true))
		if !viper.GetBool("disableLiveReload") {
			// Will block forever trying to write to a channel that nobody is reading if livereload isn't initialized
			livereload.ForceRefresh()
		}
	})
}

func build(watches ...bool) error {
	// Hugo writes the output to memory instead of the disk.
	// This is only used for benchmark testing. Cause the content is only visible
	// in memory.
	if renderToMemory {
		hugofs.SetDestination(new(afero.MemMapFs))
		// Rendering to memoryFS, publish to Root regardless of publishDir.
		viper.Set("publishDir", "/")
	}

	if err := copyStatic(); err != nil {
		return fmt.Errorf("Error copying static files to %s: %s", helpers.AbsPathify(viper.GetString("publishDir")), err)
	}
	watch := false
	if len(watches) > 0 && watches[0] {
		watch = true
	}
	if err := buildSites(buildWatch || watch); err != nil {
		return fmt.Errorf("Error building site: %s", err)
	}

	if buildWatch {
		jww.FEEDBACK.Println("Watching for changes in", helpers.AbsPathify(viper.GetString("contentDir")))
		jww.FEEDBACK.Println("Press Ctrl+C to stop")
		utils.CheckErr(NewWatcher(0))
	}

	return nil
}

func getStaticSourceFs() afero.Fs {
	source := hugofs.Source()
	themeDir, err := helpers.GetThemeStaticDirPath()
	staticDir := helpers.GetStaticDirPath() + helpers.FilePathSeparator

	useTheme := true
	useStatic := true

	if err != nil {
		if err != helpers.ErrThemeUndefined {
			jww.WARN.Println(err)
		}
		useTheme = false
	} else {
		if _, err := source.Stat(themeDir); os.IsNotExist(err) {
			jww.WARN.Println("Unable to find Theme Static Directory:", themeDir)
			useTheme = false
		}
	}

	if _, err := source.Stat(staticDir); os.IsNotExist(err) {
		jww.WARN.Println("Unable to find Static Directory:", staticDir)
		useStatic = false
	}

	if !useStatic && !useTheme {
		return nil
	}

	if !useStatic {
		jww.INFO.Println(themeDir, "is the only static directory available to sync from")
		return afero.NewReadOnlyFs(afero.NewBasePathFs(source, themeDir))
	}

	if !useTheme {
		jww.INFO.Println(staticDir, "is the only static directory available to sync from")
		return afero.NewReadOnlyFs(afero.NewBasePathFs(source, staticDir))
	}

	jww.INFO.Println("using a UnionFS for static directory comprised of:")
	jww.INFO.Println("Base:", themeDir)
	jww.INFO.Println("Overlay:", staticDir)
	base := afero.NewReadOnlyFs(afero.NewBasePathFs(hugofs.Source(), themeDir))
	overlay := afero.NewReadOnlyFs(afero.NewBasePathFs(hugofs.Source(), staticDir))
	return afero.NewCopyOnWriteFs(base, overlay)
}

func copyStatic() error {
	publishDir := helpers.AbsPathify(viper.GetString("publishDir")) + helpers.FilePathSeparator

	// If root, remove the second '/'
	if publishDir == "//" {
		publishDir = helpers.FilePathSeparator
	}

	// Includes both theme/static & /static
	staticSourceFs := getStaticSourceFs()

	if staticSourceFs == nil {
		jww.WARN.Println("No static directories found to sync")
		return nil
	}

	syncer := fsync.NewSyncer()
	syncer.NoTimes = viper.GetBool("noTimes")
	syncer.NoChmod = viper.GetBool("noChmod")
	syncer.SrcFs = staticSourceFs
	syncer.DestFs = hugofs.Destination()
	// Now that we are using a unionFs for the static directories
	// We can effectively clean the publishDir on initial sync
	syncer.Delete = viper.GetBool("cleanDestinationDir")
	if syncer.Delete {
		jww.INFO.Println("removing all files from destination that don't exist in static dirs")
	}
	jww.INFO.Println("syncing static files to", publishDir)

	// because we are using a baseFs (to get the union right).
	// set sync src to root
	return syncer.Sync(publishDir, helpers.FilePathSeparator)
}

// getDirList provides NewWatcher() with a list of directories to watch for changes.
func getDirList() []string {
	var a []string
	dataDir := helpers.AbsPathify(viper.GetString("dataDir"))
	i18nDir := helpers.AbsPathify(viper.GetString("i18nDir"))
	layoutDir := helpers.AbsPathify(viper.GetString("layoutDir"))
	staticDir := helpers.AbsPathify(viper.GetString("staticDir"))
	var themesDir string

	if helpers.ThemeSet() {
		themesDir = helpers.AbsPathify(viper.GetString("themesDir") + "/" + viper.GetString("theme"))
	}

	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			if path == dataDir && os.IsNotExist(err) {
				jww.WARN.Println("Skip dataDir:", err)
				return nil
			}

			if path == i18nDir && os.IsNotExist(err) {
				jww.WARN.Println("Skip i18nDir:", err)
				return nil
			}

			if path == layoutDir && os.IsNotExist(err) {
				jww.WARN.Println("Skip layoutDir:", err)
				return nil
			}

			if path == staticDir && os.IsNotExist(err) {
				jww.WARN.Println("Skip staticDir:", err)
				return nil
			}

			if os.IsNotExist(err) {
				// Ignore.
				return nil
			}

			jww.ERROR.Println("Walker: ", err)
			return nil
		}

		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			link, err := filepath.EvalSymlinks(path)
			if err != nil {
				jww.ERROR.Printf("Cannot read symbolic link '%s', error was: %s", path, err)
				return nil
			}
			linkfi, err := hugofs.Source().Stat(link)
			if err != nil {
				jww.ERROR.Printf("Cannot stat '%s', error was: %s", link, err)
				return nil
			}
			if !linkfi.Mode().IsRegular() {
				jww.ERROR.Printf("Symbolic links for directories not supported, skipping '%s'", path)
			}
			return nil
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

	helpers.SymbolicWalk(hugofs.Source(), dataDir, walker)
	helpers.SymbolicWalk(hugofs.Source(), helpers.AbsPathify(viper.GetString("contentDir")), walker)
	helpers.SymbolicWalk(hugofs.Source(), i18nDir, walker)
	helpers.SymbolicWalk(hugofs.Source(), helpers.AbsPathify(viper.GetString("layoutDir")), walker)

	helpers.SymbolicWalk(hugofs.Source(), staticDir, walker)
	if helpers.ThemeSet() {
		helpers.SymbolicWalk(hugofs.Source(), filepath.Join(themesDir, "layouts"), walker)
		helpers.SymbolicWalk(hugofs.Source(), filepath.Join(themesDir, "static"), walker)
		helpers.SymbolicWalk(hugofs.Source(), filepath.Join(themesDir, "i18n"), walker)
		helpers.SymbolicWalk(hugofs.Source(), filepath.Join(themesDir, "data"), walker)

	}

	return a
}

func recreateAndBuildSites(watching bool) (err error) {
	if err := initSites(); err != nil {
		return err
	}
	if !quiet {
		jww.FEEDBACK.Println("Started building sites ...")
	}
	return Hugo.Build(hugolib.BuildCfg{CreateSitesFromConfig: true, Watching: watching, PrintStats: !quiet})
}

func resetAndBuildSites(watching bool) (err error) {
	if err := initSites(); err != nil {
		return err
	}
	if !quiet {
		jww.FEEDBACK.Println("Started building sites ...")
	}
	return Hugo.Build(hugolib.BuildCfg{ResetState: true, Watching: watching, PrintStats: !quiet})
}

func initSites() error {
	if Hugo != nil {
		return nil
	}

	h, err := hugolib.NewHugoSitesFromConfiguration()

	if err != nil {
		return err
	}
	Hugo = h

	return nil
}

func buildSites(watching bool) (err error) {
	if err := initSites(); err != nil {
		return err
	}
	if !quiet {
		jww.FEEDBACK.Println("Started building sites ...")
	}
	return Hugo.Build(hugolib.BuildCfg{Watching: watching, PrintStats: !quiet})
}

func rebuildSites(events []fsnotify.Event) error {
	if err := initSites(); err != nil {
		return err
	}
	return Hugo.Build(hugolib.BuildCfg{PrintStats: !quiet, Watching: true}, events...)
}

// NewWatcher creates a new watcher to watch filesystem events.
func NewWatcher(port int) error {
	if runtime.GOOS == "darwin" {
		tweakLimit()
	}

	watcher, err := watcher.New(1 * time.Second)
	var wg sync.WaitGroup

	if err != nil {
		return err
	}

	defer watcher.Close()

	wg.Add(1)

	for _, d := range getDirList() {
		if d != "" {
			_ = watcher.Add(d)
		}
	}

	go func() {
		for {
			select {
			case evs := <-watcher.Events:
				jww.INFO.Println("Received System Events:", evs)

				staticEvents := []fsnotify.Event{}
				dynamicEvents := []fsnotify.Event{}

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
							jww.FEEDBACK.Println("adding created directory to watchlist", path)
							watcher.Add(path)
						}
						return nil
					}

					// recursively add new directories to watch list
					// When mkdir -p is used, only the top directory triggers an event (at least on OSX)
					if ev.Op&fsnotify.Create == fsnotify.Create {
						if s, err := hugofs.Source().Stat(ev.Name); err == nil && s.Mode().IsDir() {
							helpers.SymbolicWalk(hugofs.Source(), ev.Name, walkAdder)
						}
					}

					isstatic := strings.HasPrefix(ev.Name, helpers.GetStaticDirPath()) || (len(helpers.GetThemesDirPath()) > 0 && strings.HasPrefix(ev.Name, helpers.GetThemesDirPath()))

					if isstatic {
						staticEvents = append(staticEvents, ev)
					} else {
						dynamicEvents = append(dynamicEvents, ev)
					}
				}

				if len(staticEvents) > 0 {
					publishDir := helpers.AbsPathify(viper.GetString("publishDir")) + helpers.FilePathSeparator

					// If root, remove the second '/'
					if publishDir == "//" {
						publishDir = helpers.FilePathSeparator
					}

					jww.FEEDBACK.Println("\nStatic file changes detected")
					const layout = "2006-01-02 15:04 -0700"
					jww.FEEDBACK.Println(time.Now().Format(layout))

					if viper.GetBool("forceSyncStatic") {
						jww.FEEDBACK.Printf("Syncing all static files\n")
						err := copyStatic()
						if err != nil {
							utils.StopOnErr(err, fmt.Sprintf("Error copying static files to %s", publishDir))
						}
					} else {
						staticSourceFs := getStaticSourceFs()

						if staticSourceFs == nil {
							jww.WARN.Println("No static directories found to sync")
							return
						}

						syncer := fsync.NewSyncer()
						syncer.NoTimes = viper.GetBool("noTimes")
						syncer.NoChmod = viper.GetBool("noChmod")
						syncer.SrcFs = staticSourceFs
						syncer.DestFs = hugofs.Destination()

						// prevent spamming the log on changes
						logger := helpers.NewDistinctFeedbackLogger()

						for _, ev := range staticEvents {
							// Due to our approach of layering both directories and the content's rendered output
							// into one we can't accurately remove a file not in one of the source directories.
							// If a file is in the local static dir and also in the theme static dir and we remove
							// it from one of those locations we expect it to still exist in the destination
							//
							// If Hugo generates a file (from the content dir) over a static file
							// the content generated file should take precedence.
							//
							// Because we are now watching and handling individual events it is possible that a static
							// event that occupies the same path as a content generated file will take precedence
							// until a regeneration of the content takes places.
							//
							// Hugo assumes that these cases are very rare and will permit this bad behavior
							// The alternative is to track every single file and which pipeline rendered it
							// and then to handle conflict resolution on every event.

							fromPath := ev.Name

							// If we are here we already know the event took place in a static dir
							relPath, err := helpers.MakeStaticPathRelative(fromPath)
							if err != nil {
								jww.ERROR.Println(err)
								continue
							}

							// Remove || rename is harder and will require an assumption.
							// Hugo takes the following approach:
							// If the static file exists in any of the static source directories after this event
							// Hugo will re-sync it.
							// If it does not exist in all of the static directories Hugo will remove it.
							//
							// This assumes that Hugo has not generated content on top of a static file and then removed
							// the source of that static file. In this case Hugo will incorrectly remove that file
							// from the published directory.
							if ev.Op&fsnotify.Rename == fsnotify.Rename || ev.Op&fsnotify.Remove == fsnotify.Remove {
								if _, err := staticSourceFs.Stat(relPath); os.IsNotExist(err) {
									// If file doesn't exist in any static dir, remove it
									toRemove := filepath.Join(publishDir, relPath)
									logger.Println("File no longer exists in static dir, removing", toRemove)
									hugofs.Destination().RemoveAll(toRemove)
								} else if err == nil {
									// If file still exists, sync it
									logger.Println("Syncing", relPath, "to", publishDir)
									if err := syncer.Sync(filepath.Join(publishDir, relPath), relPath); err != nil {
										jww.ERROR.Println(err)
									}
								} else {
									jww.ERROR.Println(err)
								}

								continue
							}

							// For all other event operations Hugo will sync static.
							logger.Println("Syncing", relPath, "to", publishDir)
							if err := syncer.Sync(filepath.Join(publishDir, relPath), relPath); err != nil {
								jww.ERROR.Println(err)
							}
						}
					}

					if !buildWatch && !viper.GetBool("disableLiveReload") {
						// Will block forever trying to write to a channel that nobody is reading if livereload isn't initialized

						// force refresh when more than one file
						if len(staticEvents) > 0 {
							for _, ev := range staticEvents {
								path, _ := helpers.MakeStaticPathRelative(ev.Name)
								livereload.RefreshPath(path)
							}

						} else {
							livereload.ForceRefresh()
						}
					}
				}

				if len(dynamicEvents) > 0 {
					jww.FEEDBACK.Println("\nChange detected, rebuilding site")
					const layout = "2006-01-02 15:04 -0700"
					jww.FEEDBACK.Println(time.Now().Format(layout))

					rebuildSites(dynamicEvents)

					if !buildWatch && !viper.GetBool("disableLiveReload") {
						// Will block forever trying to write to a channel that nobody is reading if livereload isn't initialized
						livereload.ForceRefresh()
					}
				}
			case err := <-watcher.Errors:
				if err != nil {
					jww.ERROR.Println(err)
				}
			}
		}
	}()

	if port > 0 {
		if !viper.GetBool("disableLiveReload") {
			livereload.Initialize()
			http.HandleFunc("/livereload.js", livereload.ServeJS)
			http.HandleFunc("/livereload", livereload.Handler)
		}

		go serve(port)
	}

	wg.Wait()
	return nil
}

// isThemeVsHugoVersionMismatch returns whether the current Hugo version is
// less than the theme's min_version.
func isThemeVsHugoVersionMismatch() (mismatch bool, requiredMinVersion string) {
	if !helpers.ThemeSet() {
		return
	}

	themeDir := helpers.GetThemeDir()

	fs := hugofs.Source()
	path := filepath.Join(themeDir, "theme.toml")

	exists, err := helpers.Exists(path, fs)

	if err != nil || !exists {
		return
	}

	b, err := afero.ReadFile(fs, path)

	c, err := parser.HandleTOMLMetaData(b)

	if err != nil {
		return
	}

	config := c.(map[string]interface{})

	if minVersion, ok := config["min_version"]; ok {
		switch minVersion.(type) {
		case float32:
			return helpers.HugoVersionNumber < minVersion.(float32), fmt.Sprint(minVersion)
		case float64:
			return helpers.HugoVersionNumber < minVersion.(float64), fmt.Sprint(minVersion)
		default:
			return
		}

	}

	return
}
