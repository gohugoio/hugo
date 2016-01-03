// Copyright 2015 The Hugo Authors. All rights reserved.
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

	"github.com/spf13/hugo/parser"

	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/fsync"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/livereload"
	"github.com/spf13/hugo/utils"
	"github.com/spf13/hugo/watcher"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/nitro"
	"github.com/spf13/viper"
	"gopkg.in/fsnotify.v1"
)

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

//HugoCmd is Hugo's root command. Every other command attached to HugoCmd is a child command to it.
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

		watchConfig()

		return build()

	},
}

var hugoCmdV *cobra.Command

// Flags that are to be added to commands.
var BuildWatch, IgnoreCache, Draft, Future, UglyURLs, CanonifyURLs, Verbose, Logging, VerboseLog, DisableRSS, DisableSitemap, DisableRobotsTXT, PluralizeListTitles, PreserveTaxonomyNames, NoTimes, ForceSync bool
var Source, CacheDir, Destination, Theme, BaseURL, CfgFile, LogFile, Editor string

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

// initCoreCommonFlags initializes common flags used by Hugo core commands
// such as hugo itself, server, check, config and benchmark.
func initCoreCommonFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&Draft, "buildDrafts", "D", false, "include content marked as draft")
	cmd.Flags().BoolVarP(&Future, "buildFuture", "F", false, "include content with publishdate in the future")
	cmd.Flags().BoolVar(&DisableRSS, "disableRSS", false, "Do not build RSS files")
	cmd.Flags().BoolVar(&DisableSitemap, "disableSitemap", false, "Do not build Sitemap file")
	cmd.Flags().BoolVar(&DisableRobotsTXT, "disableRobotsTXT", false, "Do not build Robots TXT file")
	cmd.Flags().StringVarP(&Source, "source", "s", "", "filesystem path to read files relative from")
	cmd.Flags().StringVarP(&CacheDir, "cacheDir", "", "", "filesystem path to cache directory. Defaults: $TMPDIR/hugo_cache/")
	cmd.Flags().BoolVarP(&IgnoreCache, "ignoreCache", "", false, "Ignores the cache directory for reading but still writes to it")
	cmd.Flags().StringVarP(&Destination, "destination", "d", "", "filesystem path to write files to")
	cmd.Flags().StringVarP(&Theme, "theme", "t", "", "theme to use (located in /themes/THEMENAME/)")
	cmd.Flags().BoolVar(&UglyURLs, "uglyURLs", false, "if true, use /filename.html instead of /filename/")
	cmd.Flags().BoolVar(&CanonifyURLs, "canonifyURLs", false, "if true, all relative URLs will be canonicalized using baseURL")
	cmd.Flags().StringVarP(&BaseURL, "baseURL", "b", "", "hostname (and path) to the root, e.g. http://spf13.com/")
	cmd.Flags().StringVar(&CfgFile, "config", "", "config file (default is path/config.yaml|json|toml)")
	cmd.Flags().StringVar(&Editor, "editor", "", "edit new content with this editor, if provided")
	cmd.Flags().BoolVar(&nitro.AnalysisOn, "stepAnalysis", false, "display memory and timing of different steps of the program")
	cmd.Flags().BoolVar(&PluralizeListTitles, "pluralizeListTitles", true, "Pluralize titles in lists using inflect")
	cmd.Flags().BoolVar(&PreserveTaxonomyNames, "preserveTaxonomyNames", false, `Preserve taxonomy names as written ("Gérard Depardieu" vs "gerard-depardieu")`)
	cmd.Flags().BoolVarP(&ForceSync, "forceSyncStatic", "", false, "Copy all files when static is changed.")
	// For bash-completion
	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	cmd.Flags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)
	cmd.Flags().SetAnnotation("source", cobra.BashCompSubdirsInDir, []string{})
	cmd.Flags().SetAnnotation("cacheDir", cobra.BashCompSubdirsInDir, []string{})
	cmd.Flags().SetAnnotation("destination", cobra.BashCompSubdirsInDir, []string{})
	cmd.Flags().SetAnnotation("theme", cobra.BashCompSubdirsInDir, []string{"themes"})
}

// init initializes flags.
func init() {
	HugoCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	HugoCmd.PersistentFlags().BoolVar(&Logging, "log", false, "Enable Logging")
	HugoCmd.PersistentFlags().StringVar(&LogFile, "logFile", "", "Log File path (if set, logging enabled automatically)")
	HugoCmd.PersistentFlags().BoolVar(&VerboseLog, "verboseLog", false, "verbose logging")

	initCoreCommonFlags(HugoCmd)

	HugoCmd.Flags().BoolVarP(&BuildWatch, "watch", "w", false, "watch filesystem for changes and recreate as needed")
	HugoCmd.Flags().BoolVarP(&NoTimes, "noTimes", "", false, "Don't sync modification time of files")
	hugoCmdV = HugoCmd

	// For bash-completion
	HugoCmd.PersistentFlags().SetAnnotation("logFile", cobra.BashCompFilenameExt, []string{})
}

func LoadDefaultSettings() {
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
	viper.SetDefault("Multilingual", false)
	viper.SetDefault("DefaultContentLang", "en")
}

// InitializeConfig initializes a config file with sensible default configuration flags.
// A Hugo command that calls initCoreCommonFlags() can pass itself
// as an argument to have its command-line flags processed here.
func InitializeConfig(subCmdVs ...*cobra.Command) error {
	viper.SetConfigFile(CfgFile)
	// See https://github.com/spf13/viper/issues/73#issuecomment-126970794
	if Source == "" {
		viper.AddConfigPath(".")
	} else {
		viper.AddConfigPath(Source)
	}
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return newSystemError(err)
		} else {
			return newSystemErrorF("Unable to locate Config file. Perhaps you need to create a new site.\n       Run `hugo help new` for details. (%s)\n", err)
		}
	}

	viper.RegisterAlias("indexes", "taxonomies")

	LoadDefaultSettings()

	if hugoCmdV.PersistentFlags().Lookup("verbose").Changed {
		viper.Set("Verbose", Verbose)
	}
	if hugoCmdV.PersistentFlags().Lookup("logFile").Changed {
		viper.Set("LogFile", LogFile)
	}

	for _, cmdV := range append([]*cobra.Command{hugoCmdV}, subCmdVs...) {
		if cmdV.Flags().Lookup("buildDrafts").Changed {
			viper.Set("BuildDrafts", Draft)
		}
		if cmdV.Flags().Lookup("buildFuture").Changed {
			viper.Set("BuildFuture", Future)
		}
		if cmdV.Flags().Lookup("uglyURLs").Changed {
			viper.Set("UglyURLs", UglyURLs)
		}
		if cmdV.Flags().Lookup("canonifyURLs").Changed {
			viper.Set("CanonifyURLs", CanonifyURLs)
		}
		if cmdV.Flags().Lookup("disableRSS").Changed {
			viper.Set("DisableRSS", DisableRSS)
		}
		if cmdV.Flags().Lookup("disableSitemap").Changed {
			viper.Set("DisableSitemap", DisableSitemap)
		}
		if cmdV.Flags().Lookup("disableRobotsTXT").Changed {
			viper.Set("DisableRobotsTXT", DisableRobotsTXT)
		}
		if cmdV.Flags().Lookup("pluralizeListTitles").Changed {
			viper.Set("PluralizeListTitles", PluralizeListTitles)
		}
		if cmdV.Flags().Lookup("preserveTaxonomyNames").Changed {
			viper.Set("PreserveTaxonomyNames", PreserveTaxonomyNames)
		}
		if cmdV.Flags().Lookup("forceSyncStatic").Changed {
			viper.Set("ForceSyncStatic", ForceSync)
		}
		if cmdV.Flags().Lookup("editor").Changed {
			viper.Set("NewContentEditor", Editor)
		}
		if cmdV.Flags().Lookup("ignoreCache").Changed {
			viper.Set("IgnoreCache", IgnoreCache)
		}
	}

	if hugoCmdV.Flags().Lookup("noTimes").Changed {
		viper.Set("NoTimes", NoTimes)
	}

	if BaseURL != "" {
		if !strings.HasSuffix(BaseURL, "/") {
			BaseURL = BaseURL + "/"
		}
		viper.Set("BaseURL", BaseURL)
	}

	if !viper.GetBool("RelativeURLs") && viper.GetString("BaseURL") == "" {
		jww.ERROR.Println("No 'baseurl' set in configuration or as a flag. Features like page menus will not work without one.")
	}

	if Theme != "" {
		viper.Set("theme", Theme)
	}

	if Destination != "" {
		viper.Set("PublishDir", Destination)
	}

	if Source != "" {
		dir, _ := filepath.Abs(Source)
		viper.Set("WorkingDir", dir)
	} else {
		dir, _ := os.Getwd()
		viper.Set("WorkingDir", dir)
	}

	if CacheDir != "" {
		if helpers.FilePathSeparator != CacheDir[len(CacheDir)-1:] {
			CacheDir = CacheDir + helpers.FilePathSeparator
		}
		isDir, err := helpers.DirExists(CacheDir, hugofs.SourceFs)
		utils.CheckErr(err)
		if isDir == false {
			mkdir(CacheDir)
		}
		viper.Set("CacheDir", CacheDir)
	} else {
		viper.Set("CacheDir", helpers.GetTempDir("hugo_cache", hugofs.SourceFs))
	}

	if VerboseLog || Logging || (viper.IsSet("LogFile") && viper.GetString("LogFile") != "") {
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

	if VerboseLog {
		jww.SetLogThreshold(jww.LevelInfo)
	}

	jww.INFO.Println("Using config file:", viper.ConfigFileUsed())

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

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		utils.CheckErr(buildSite(true))
		if !viper.GetBool("DisableLiveReload") {
			// Will block forever trying to write to a channel that nobody is reading if livereload isn't initalized
			livereload.ForceRefresh()
		}
	})
}

func build(watches ...bool) error {

	if err := copyStatic(); err != nil {
		return fmt.Errorf("Error copying static files to %s: %s", helpers.AbsPathify(viper.GetString("PublishDir")), err)
	}
	watch := false
	if len(watches) > 0 && watches[0] {
		watch = true
	}
	if err := buildSite(BuildWatch || watch); err != nil {
		return fmt.Errorf("Error building site: %s", err)
	}

	if BuildWatch {
		jww.FEEDBACK.Println("Watching for changes in", helpers.AbsPathify(viper.GetString("ContentDir")))
		jww.FEEDBACK.Println("Press Ctrl+C to stop")
		utils.CheckErr(NewWatcher(0))
	}

	return nil
}

func copyStatic() error {
	publishDir := helpers.AbsPathify(viper.GetString("PublishDir")) + helpers.FilePathSeparator

	// If root, remove the second '/'
	if publishDir == "//" {
		publishDir = helpers.FilePathSeparator
	}

	syncer := fsync.NewSyncer()
	syncer.NoTimes = viper.GetBool("notimes")
	syncer.SrcFs = hugofs.SourceFs
	syncer.DestFs = hugofs.DestinationFS

	themeDir, err := helpers.GetThemeStaticDirPath()
	if err != nil {
		jww.WARN.Println(err)
	}

	// Copy the theme's static directory
	if themeDir != "" {
		jww.INFO.Println("syncing from", themeDir, "to", publishDir)
		utils.CheckErr(syncer.Sync(publishDir, themeDir), fmt.Sprintf("Error copying static files of theme to %s", publishDir))
	}

	// Copy the site's own static directory
	staticDir := helpers.GetStaticDirPath() + helpers.FilePathSeparator
	if _, err := os.Stat(staticDir); err == nil {
		jww.INFO.Println("syncing from", staticDir, "to", publishDir)
		return syncer.Sync(publishDir, staticDir)
	} else if os.IsNotExist(err) {
		jww.WARN.Println("Unable to find Static Directory:", staticDir)
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

	filepath.Walk(dataDir, walker)
	filepath.Walk(helpers.AbsPathify(viper.GetString("ContentDir")), walker)
	filepath.Walk(helpers.AbsPathify(viper.GetString("LayoutDir")), walker)
	filepath.Walk(helpers.AbsPathify(viper.GetString("StaticDir")), walker)
	if helpers.ThemeSet() {
		filepath.Walk(helpers.AbsPathify(viper.GetString("themesDir")+"/"+viper.GetString("theme")), walker)
	}

	return a
}

func buildSite(watching ...bool) (err error) {
	startTime := time.Now()
	site := &hugolib.Site{}
	if len(watching) > 0 && watching[0] {
		site.RunMode.Watching = true
	}
	err = site.Build()
	if err != nil {
		return err
	}
	site.Stats()
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
				jww.INFO.Println("File System Event:", evs)

				staticChanged := false
				dynamicChanged := false
				staticFilesChanged := make(map[string]bool)

				for _, ev := range evs {
					ext := filepath.Ext(ev.Name)
					istemp := strings.HasSuffix(ext, "~") || (ext == ".swp") || (ext == ".swx") || (ext == ".tmp") || strings.HasPrefix(ext, ".goutputstream") || strings.HasSuffix(ext, "jb_old___") || strings.HasSuffix(ext, "jb_bak___")
					if istemp {
						continue
					}
					// renames are always followed with Create/Modify
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
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

					isstatic := strings.HasPrefix(ev.Name, helpers.GetStaticDirPath()) || (len(helpers.GetThemesDirPath()) > 0 && strings.HasPrefix(ev.Name, helpers.GetThemesDirPath()))
					staticChanged = staticChanged || isstatic
					dynamicChanged = dynamicChanged || !isstatic

					if isstatic {
						staticFilesChanged[ev.Name] = true
					}

					// add new directory to watch list
					if s, err := os.Stat(ev.Name); err == nil && s.Mode().IsDir() {
						if ev.Op&fsnotify.Create == fsnotify.Create {
							watcher.Add(ev.Name)
						}
					}
				}

				if staticChanged {
					jww.FEEDBACK.Printf("Static file changed, syncing\n")
					if viper.GetBool("ForceSyncStatic") {
						jww.FEEDBACK.Printf("Syncing all static files\n")
						err := copyStatic()
						if err != nil {
							fmt.Println(err)
							utils.StopOnErr(err, fmt.Sprintf("Error copying static files to %s", helpers.AbsPathify(viper.GetString("PublishDir"))))
						}
					} else {

						syncer := fsync.NewSyncer()
						syncer.NoTimes = viper.GetBool("notimes")
						syncer.SrcFs = hugofs.SourceFs
						syncer.DestFs = hugofs.DestinationFS

						publishDir := helpers.AbsPathify(viper.GetString("PublishDir")) + helpers.FilePathSeparator

						if publishDir == "//" || publishDir == helpers.FilePathSeparator {
							publishDir = ""
						}

						staticDir := helpers.GetStaticDirPath()
						themeStaticDir := helpers.GetThemesDirPath()

						jww.FEEDBACK.Printf("StaticDir '%s'\nThemeStaticDir '%s'\n", staticDir, themeStaticDir)

						for path := range staticFilesChanged {
							var publishPath string

							if strings.HasPrefix(path, staticDir) {
								publishPath = filepath.Join(publishDir, strings.TrimPrefix(path, staticDir))
							} else if strings.HasPrefix(path, themeStaticDir) {
								publishPath = filepath.Join(publishDir, strings.TrimPrefix(path, themeStaticDir))
							}
							jww.FEEDBACK.Printf("Syncing file '%s'", path)

							if _, err := os.Stat(path); err == nil {
								jww.INFO.Println("syncing from ", path, " to ", publishPath)
								err := syncer.Sync(publishPath, path)
								if err != nil {
									jww.FEEDBACK.Printf("Error on syncing file '%s'\n", path)
								}
							}
						}
					}

					if !BuildWatch && !viper.GetBool("DisableLiveReload") {
						// Will block forever trying to write to a channel that nobody is reading if livereload isn't initalized

						// force refresh when more than one file
						if len(staticFilesChanged) == 1 {
							for path := range staticFilesChanged {
								livereload.RefreshPath(path)
							}

						} else {
							livereload.ForceRefresh()
						}
					}
				}

				if dynamicChanged {
					fmt.Print("\nChange detected, rebuilding site\n")
					const layout = "2006-01-02 15:04 -0700"
					fmt.Println(time.Now().Format(layout))
					utils.CheckErr(buildSite(true))

					if !BuildWatch && !viper.GetBool("DisableLiveReload") {
						// Will block forever trying to write to a channel that nobody is reading if livereload isn't initalized
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

	fs := hugofs.SourceFs
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
