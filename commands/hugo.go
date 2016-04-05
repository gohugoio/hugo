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
	"github.com/spf13/hugo/hugofs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/spf13/hugo/parser"
	flag "github.com/spf13/pflag"

	"regexp"

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
	"gopkg.in/fsnotify.v1"
)

// MainSite represents the Hugo site to build. This variable is exported as it
// is used by at least one external library (the Hugo caddy plugin). We should
// provide a cleaner external API, but until then, this is it.
var MainSite *hugolib.Site

// userError is an error used to signal different error situations in command handling.
type commandError struct {
	s         string
	userError bool
}

func (u commandError) Error() string {
	return u.s
}

func (u commandError) isUserError() bool {
	return u.userError
}

func newUserError(a ...interface{}) commandError {
	return commandError{s: fmt.Sprintln(a...), userError: true}
}

func newUserErrorF(format string, a ...interface{}) commandError {
	return commandError{s: fmt.Sprintf(format, a...), userError: true}
}

func newSystemError(a ...interface{}) commandError {
	return commandError{s: fmt.Sprintln(a...), userError: false}
}

func newSystemErrorF(format string, a ...interface{}) commandError {
	return commandError{s: fmt.Sprintf(format, a...), userError: false}
}

// catch some of the obvious user errors from Cobra.
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
			viper.Set("DisableLiveReload", true)
			watchConfig()
		}

		return build()
	},
}

var hugoCmdV *cobra.Command

// Flags that are to be added to commands.
var (
	buildWatch            bool
	canonifyURLs          bool
	cleanDestination      bool
	enableRobotsTXT       bool
	disableRSS            bool
	disableSitemap        bool
	draft                 bool
	forceSync             bool
	future                bool
	ignoreCache           bool
	logging               bool
	noTimes               bool
	pluralizeListTitles   bool
	preserveTaxonomyNames bool
	renderToMemory        bool // for benchmark testing
	uglyURLs              bool
	verbose               bool
	verboseLog            bool
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
	initCoreCommonFlags(cmd)
	initHugoBuildCommonFlags(cmd)
}

// initCoreCommonFlags initializes common flags used by Hugo core commands.
func initCoreCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is path/config.yaml|json|toml)")

	// Set bash-completion
	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	cmd.Flags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)
}

// initHugoBuildCommonFlags initialize common flags related to the Hugo build.
// Called by initHugoBuilderFlags.
func initHugoBuildCommonFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&cleanDestination, "cleanDestinationDir", false, "Remove files from destination not found in static directories")
	cmd.Flags().BoolVarP(&draft, "buildDrafts", "D", false, "include content marked as draft")
	cmd.Flags().BoolVarP(&future, "buildFuture", "F", false, "include content with publishdate in the future")
	cmd.Flags().BoolVar(&disableRSS, "disableRSS", false, "Do not build RSS files")
	cmd.Flags().BoolVar(&disableSitemap, "disableSitemap", false, "Do not build Sitemap file")
	cmd.Flags().StringVarP(&source, "source", "s", "", "filesystem path to read files relative from")
	cmd.Flags().StringVarP(&contentDir, "contentDir", "c", "", "filesystem path to content directory")
	cmd.Flags().StringVarP(&layoutDir, "layoutDir", "l", "", "filesystem path to layout directory")
	cmd.Flags().StringVarP(&cacheDir, "cacheDir", "", "", "filesystem path to cache directory. Defaults: $TMPDIR/hugo_cache/")
	cmd.Flags().BoolVarP(&ignoreCache, "ignoreCache", "", false, "Ignores the cache directory for reading but still writes to it")
	cmd.Flags().StringVarP(&destination, "destination", "d", "", "filesystem path to write files to")
	cmd.Flags().StringVarP(&theme, "theme", "t", "", "theme to use (located in /themes/THEMENAME/)")
	cmd.Flags().BoolVar(&uglyURLs, "uglyURLs", false, "if true, use /filename.html instead of /filename/")
	cmd.Flags().BoolVar(&canonifyURLs, "canonifyURLs", false, "if true, all relative URLs will be canonicalized using baseURL")
	cmd.Flags().StringVarP(&baseURL, "baseURL", "b", "", "hostname (and path) to the root, e.g. http://spf13.com/")

	cmd.Flags().BoolVar(&nitro.AnalysisOn, "stepAnalysis", false, "display memory and timing of different steps of the program")
	cmd.Flags().BoolVar(&pluralizeListTitles, "pluralizeListTitles", true, "Pluralize titles in lists using inflect")
	cmd.Flags().BoolVar(&preserveTaxonomyNames, "preserveTaxonomyNames", false, `Preserve taxonomy names as written ("Gérard Depardieu" vs "gerard-depardieu")`)
	cmd.Flags().BoolVarP(&forceSync, "forceSyncStatic", "", false, "Copy all files when static is changed.")
	cmd.Flags().BoolVarP(&noTimes, "noTimes", "", false, "Don't sync modification time of files")

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

	initHugoBuilderFlags(HugoCmd)
	initBenchmarkBuildingFlags(HugoCmd)

	HugoCmd.Flags().BoolVarP(&buildWatch, "watch", "w", false, "watch filesystem for changes and recreate as needed")
	hugoCmdV = HugoCmd

	// Set bash-completion
	HugoCmd.PersistentFlags().SetAnnotation("logFile", cobra.BashCompFilenameExt, []string{})
}

func loadDefaultSettings() {
	viper.SetDefault("cleanDestinationDir", false)
	viper.SetDefault("Watch", false)
	viper.SetDefault("MetaDataFormat", "toml")
	viper.SetDefault("DisableRSS", false)
	viper.SetDefault("DisableSitemap", false)
	viper.SetDefault("DisableRobotsTXT", false)
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("LayoutDir", "layouts")
	viper.SetDefault("StaticDir", "static")
	viper.SetDefault("ArchetypeDir", "archetypes")
	viper.SetDefault("PublishDir", "public")
	viper.SetDefault("DataDir", "data")
	viper.SetDefault("ThemesDir", "themes")
	viper.SetDefault("DefaultLayout", "post")
	viper.SetDefault("BuildDrafts", false)
	viper.SetDefault("BuildFuture", false)
	viper.SetDefault("UglyURLs", false)
	viper.SetDefault("Verbose", false)
	viper.SetDefault("IgnoreCache", false)
	viper.SetDefault("CanonifyURLs", false)
	viper.SetDefault("RelativeURLs", false)
	viper.SetDefault("RemovePathAccents", false)
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	viper.SetDefault("Permalinks", make(hugolib.PermalinkOverrides, 0))
	viper.SetDefault("Sitemap", hugolib.Sitemap{Priority: -1, Filename: "sitemap.xml"})
	viper.SetDefault("DefaultExtension", "html")
	viper.SetDefault("PygmentsStyle", "monokai")
	viper.SetDefault("PygmentsUseClasses", false)
	viper.SetDefault("PygmentsCodeFences", false)
	viper.SetDefault("PygmentsOptions", "")
	viper.SetDefault("DisableLiveReload", false)
	viper.SetDefault("PluralizeListTitles", true)
	viper.SetDefault("PreserveTaxonomyNames", false)
	viper.SetDefault("ForceSyncStatic", false)
	viper.SetDefault("FootnoteAnchorPrefix", "")
	viper.SetDefault("FootnoteReturnLinkContents", "")
	viper.SetDefault("NewContentEditor", "")
	viper.SetDefault("Paginate", 10)
	viper.SetDefault("PaginatePath", "page")
	viper.SetDefault("Blackfriday", helpers.NewBlackfriday())
	viper.SetDefault("RSSUri", "index.xml")
	viper.SetDefault("SectionPagesMenu", "")
	viper.SetDefault("DisablePathToLower", false)
	viper.SetDefault("HasCJKLanguage", false)
	viper.SetDefault("EnableEmoji", false)
	viper.SetDefault("PygmentsCodeFencesGuessSyntax", false)
}

// InitializeConfig initializes a config file with sensible default configuration flags.
// A Hugo command that calls initCoreCommonFlags() can pass itself
// as an argument to have its command-line flags processed here.
func InitializeConfig(subCmdVs ...*cobra.Command) error {
	viper.SetConfigFile(cfgFile)
	// See https://github.com/spf13/viper/issues/73#issuecomment-126970794
	if source == "" {
		viper.AddConfigPath(".")
	} else {
		viper.AddConfigPath(source)
	}
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return newSystemError(err)
		}
		return newSystemErrorF("Unable to locate Config file. Perhaps you need to create a new site.\n       Run `hugo help new` for details. (%s)\n", err)
	}

	viper.RegisterAlias("indexes", "taxonomies")

	loadDefaultSettings()

	for _, cmdV := range append([]*cobra.Command{hugoCmdV}, subCmdVs...) {

		if flagChanged(cmdV.PersistentFlags(), "verbose") {
			viper.Set("Verbose", verbose)
		}
		if flagChanged(cmdV.PersistentFlags(), "logFile") {
			viper.Set("LogFile", logFile)
		}
		if flagChanged(cmdV.Flags(), "cleanDestinationDir") {
			viper.Set("cleanDestinationDir", cleanDestination)
		}
		if flagChanged(cmdV.Flags(), "buildDrafts") {
			viper.Set("BuildDrafts", draft)
		}
		if flagChanged(cmdV.Flags(), "buildFuture") {
			viper.Set("BuildFuture", future)
		}
		if flagChanged(cmdV.Flags(), "uglyURLs") {
			viper.Set("UglyURLs", uglyURLs)
		}
		if flagChanged(cmdV.Flags(), "canonifyURLs") {
			viper.Set("CanonifyURLs", canonifyURLs)
		}
		if flagChanged(cmdV.Flags(), "disableRSS") {
			viper.Set("DisableRSS", disableRSS)
		}
		if flagChanged(cmdV.Flags(), "disableSitemap") {
			viper.Set("DisableSitemap", disableSitemap)
		}
		if flagChanged(cmdV.Flags(), "enableRobotsTXT") {
			viper.Set("EnableRobotsTXT", enableRobotsTXT)
		}
		if flagChanged(cmdV.Flags(), "pluralizeListTitles") {
			viper.Set("PluralizeListTitles", pluralizeListTitles)
		}
		if flagChanged(cmdV.Flags(), "preserveTaxonomyNames") {
			viper.Set("PreserveTaxonomyNames", preserveTaxonomyNames)
		}
		if flagChanged(cmdV.Flags(), "ignoreCache") {
			viper.Set("IgnoreCache", ignoreCache)
		}
		if flagChanged(cmdV.Flags(), "forceSyncStatic") {
			viper.Set("ForceSyncStatic", forceSync)
		}
		if flagChanged(cmdV.Flags(), "noTimes") {
			viper.Set("NoTimes", noTimes)
		}

	}

	if baseURL != "" {
		if !strings.HasSuffix(baseURL, "/") {
			baseURL = baseURL + "/"
		}
		viper.Set("BaseURL", baseURL)
	}

	if !viper.GetBool("RelativeURLs") && viper.GetString("BaseURL") == "" {
		jww.ERROR.Println("No 'baseurl' set in configuration or as a flag. Features like page menus will not work without one.")
	}

	if theme != "" {
		viper.Set("theme", theme)
	}

	if destination != "" {
		viper.Set("PublishDir", destination)
	}

	if source != "" {
		dir, _ := filepath.Abs(source)
		viper.Set("WorkingDir", dir)
	} else {
		dir, _ := os.Getwd()
		viper.Set("WorkingDir", dir)
	}

	if contentDir != "" {
		viper.Set("ContentDir", contentDir)
	}

	if layoutDir != "" {
		viper.Set("LayoutDir", layoutDir)
	}

	if cacheDir != "" {
		if helpers.FilePathSeparator != cacheDir[len(cacheDir)-1:] {
			cacheDir = cacheDir + helpers.FilePathSeparator
		}
		isDir, err := helpers.DirExists(cacheDir, hugofs.Source())
		utils.CheckErr(err)
		if isDir == false {
			mkdir(cacheDir)
		}
		viper.Set("CacheDir", cacheDir)
	} else {
		viper.Set("CacheDir", helpers.GetTempDir("hugo_cache", hugofs.Source()))
	}

	if verboseLog || logging || (viper.IsSet("LogFile") && viper.GetString("LogFile") != "") {
		if viper.IsSet("LogFile") && viper.GetString("LogFile") != "" {
			jww.SetLogFile(viper.GetString("LogFile"))
		} else {
			jww.UseTempLogFile("hugo")
		}
	} else {
		jww.DiscardLogging()
	}

	if viper.GetBool("verbose") {
		jww.SetStdoutThreshold(jww.LevelInfo)
	}

	if verboseLog {
		jww.SetLogThreshold(jww.LevelInfo)
	}

	jww.INFO.Println("Using config file:", viper.ConfigFileUsed())

	// Init file systems. This may be changed at a later point.
	hugofs.InitDefaultFs()

	themeDir := helpers.GetThemeDir()
	if themeDir != "" {
		if _, err := os.Stat(themeDir); os.IsNotExist(err) {
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
		fmt.Println("Config file changed:", e.Name)
		// Force a full rebuild
		MainSite = nil
		utils.CheckErr(buildSite(true))
		if !viper.GetBool("DisableLiveReload") {
			// Will block forever trying to write to a channel that nobody is reading if livereload isn't initialized
			livereload.ForceRefresh()
		}
	})
}

func build(watches ...bool) error {

	// Hugo writes the output to memory instead of the disk
	// This is only used for benchmark testing. Cause the content is only visible
	// in memory
	if renderToMemory {
		hugofs.SetDestination(new(afero.MemMapFs))
		// Rendering to memoryFS, publish to Root regardless of publishDir.
		viper.Set("PublishDir", "/")
	}

	if err := copyStatic(); err != nil {
		return fmt.Errorf("Error copying static files to %s: %s", helpers.AbsPathify(viper.GetString("PublishDir")), err)
	}
	watch := false
	if len(watches) > 0 && watches[0] {
		watch = true
	}
	if err := buildSite(buildWatch || watch); err != nil {
		return fmt.Errorf("Error building site: %s", err)
	}

	if buildWatch {
		jww.FEEDBACK.Println("Watching for changes in", helpers.AbsPathify(viper.GetString("ContentDir")))
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
		jww.WARN.Println(err)
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
	publishDir := helpers.AbsPathify(viper.GetString("PublishDir")) + helpers.FilePathSeparator

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
	syncer.NoTimes = viper.GetBool("notimes")
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
	err := syncer.Sync(publishDir, helpers.FilePathSeparator)
	if err != nil {
		return err
	}
	return nil
}

// getDirList provides NewWatcher() with a list of directories to watch for changes.
func getDirList() []string {
	var a []string
	dataDir := helpers.AbsPathify(viper.GetString("DataDir"))
	layoutDir := helpers.AbsPathify(viper.GetString("LayoutDir"))
	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			if path == dataDir && os.IsNotExist(err) {
				jww.WARN.Println("Skip DataDir:", err)
				return nil

			}
			if path == layoutDir && os.IsNotExist(err) {
				jww.WARN.Println("Skip LayoutDir:", err)
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
			linkfi, err := os.Stat(link)
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
	helpers.SymbolicWalk(hugofs.Source(), helpers.AbsPathify(viper.GetString("ContentDir")), walker)
	helpers.SymbolicWalk(hugofs.Source(), helpers.AbsPathify(viper.GetString("LayoutDir")), walker)
	helpers.SymbolicWalk(hugofs.Source(), helpers.AbsPathify(viper.GetString("StaticDir")), walker)
	if helpers.ThemeSet() {
		helpers.SymbolicWalk(hugofs.Source(), helpers.AbsPathify(viper.GetString("themesDir")+"/"+viper.GetString("theme")), walker)
	}

	return a
}

func buildSite(watching ...bool) (err error) {
	fmt.Println("Started building site")
	startTime := time.Now()
	if MainSite == nil {
		MainSite = new(hugolib.Site)
	}
	if len(watching) > 0 && watching[0] {
		MainSite.RunMode.Watching = true
	}
	err = MainSite.Build()
	if err != nil {
		return err
	}
	MainSite.Stats()
	jww.FEEDBACK.Printf("in %v ms\n", int(1000*time.Since(startTime).Seconds()))

	return nil
}

func rebuildSite(events []fsnotify.Event) error {
	startTime := time.Now()
	err := MainSite.ReBuild(events)
	if err != nil {
		return err
	}
	MainSite.Stats()
	jww.FEEDBACK.Printf("in %v ms\n", int(1000*time.Since(startTime).Seconds()))

	return nil
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
					istemp := strings.HasSuffix(ext, "~") ||
						(ext == ".swp") ||
						(ext == ".swx") ||
						(ext == ".tmp") ||
						(ext == ".DS_Store") ||
						filepath.Base(ev.Name) == "4913" ||
						strings.HasPrefix(ext, ".goutputstream") ||
						strings.HasSuffix(ext, "jb_old___") ||
						strings.HasSuffix(ext, "jb_bak___")
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
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
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
					publishDir := helpers.AbsPathify(viper.GetString("PublishDir")) + helpers.FilePathSeparator

					// If root, remove the second '/'
					if publishDir == "//" {
						publishDir = helpers.FilePathSeparator
					}

					jww.FEEDBACK.Println("\nStatic file changes detected")
					const layout = "2006-01-02 15:04 -0700"
					fmt.Println(time.Now().Format(layout))

					if viper.GetBool("ForceSyncStatic") {
						jww.FEEDBACK.Printf("Syncing all static files\n")
						err := copyStatic()
						if err != nil {
							utils.StopOnErr(err, fmt.Sprintf("Error copying static files to %s", helpers.AbsPathify(viper.GetString("PublishDir"))))
						}
					} else {
						staticSourceFs := getStaticSourceFs()

						if staticSourceFs == nil {
							jww.WARN.Println("No static directories found to sync")
							return
						}

						syncer := fsync.NewSyncer()
						syncer.NoTimes = viper.GetBool("notimes")
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
								fmt.Println(err)
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

					if !buildWatch && !viper.GetBool("DisableLiveReload") {
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
					fmt.Print("\nChange detected, rebuilding site\n")
					const layout = "2006-01-02 15:04 -0700"
					fmt.Println(time.Now().Format(layout))

					rebuildSite(dynamicEvents)

					if !buildWatch && !viper.GetBool("DisableLiveReload") {
						// Will block forever trying to write to a channel that nobody is reading if livereload isn't initialized
						livereload.ForceRefresh()
					}
				}
			case err := <-watcher.Errors:
				if err != nil {
					fmt.Println("error:", err)
				}
			}
		}
	}()

	if port > 0 {
		if !viper.GetBool("DisableLiveReload") {
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

	f, err := fs.Open(path)

	if err != nil {
		return
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)

	if err != nil {
		return
	}

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
