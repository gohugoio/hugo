// Copyright © 2013-2015 Steve Francia <spf@spf13.com>.
//
// Licensed under the Simple Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://opensource.org/licenses/Simple-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//Package commands defines and implements command-line commands and flags used by Hugo. Commands and flags are implemented using
//cobra.
package commands

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

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

//HugoCmd is Hugo's root command. Every other command attached to HugoCmd is a child command to it.
var HugoCmd = &cobra.Command{
	Use:   "hugo",
	Short: "hugo builds your site",
	Long: `hugo is the main command, used to build your Hugo site. 
	
Hugo is a Fast and Flexible Static Site Generator built with love by spf13 and friends in Go.

Complete documentation is available at http://gohugo.io`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		build()
	},
}

var hugoCmdV *cobra.Command

//Flags that are to be added to commands.
var BuildWatch, IgnoreCache, Draft, Future, UglyURLs, Verbose, Logging, VerboseLog, DisableRSS, DisableSitemap, PluralizeListTitles, NoTimes bool
var Source, CacheDir, Destination, Theme, BaseURL, CfgFile, LogFile, Editor string

//Execute adds all child commands to the root command HugoCmd and sets flags appropriately.
func Execute() {
	AddCommands()
	utils.StopOnErr(HugoCmd.Execute())
}

//AddCommands adds child commands to the root command HugoCmd.
func AddCommands() {
	HugoCmd.AddCommand(serverCmd)
	HugoCmd.AddCommand(version)
	HugoCmd.AddCommand(config)
	HugoCmd.AddCommand(check)
	HugoCmd.AddCommand(benchmark)
	HugoCmd.AddCommand(convertCmd)
	HugoCmd.AddCommand(newCmd)
	HugoCmd.AddCommand(listCmd)
	HugoCmd.AddCommand(undraftCmd)
	HugoCmd.AddCommand(genautocompleteCmd)
	HugoCmd.AddCommand(gendocCmd)
}

//Initializes flags
func init() {
	HugoCmd.PersistentFlags().BoolVarP(&Draft, "buildDrafts", "D", false, "include content marked as draft")
	HugoCmd.PersistentFlags().BoolVarP(&Future, "buildFuture", "F", false, "include content with publishdate in the future")
	HugoCmd.PersistentFlags().BoolVar(&DisableRSS, "disableRSS", false, "Do not build RSS files")
	HugoCmd.PersistentFlags().BoolVar(&DisableSitemap, "disableSitemap", false, "Do not build Sitemap file")
	HugoCmd.PersistentFlags().StringVarP(&Source, "source", "s", "", "filesystem path to read files relative from")
	HugoCmd.PersistentFlags().StringVarP(&CacheDir, "cacheDir", "", "", "filesystem path to cache directory. Defaults: $TMPDIR/hugo_cache/")
	HugoCmd.PersistentFlags().BoolVarP(&IgnoreCache, "ignoreCache", "", false, "Ignores the cache directory for reading but still writes to it")
	HugoCmd.PersistentFlags().StringVarP(&Destination, "destination", "d", "", "filesystem path to write files to")
	HugoCmd.PersistentFlags().StringVarP(&Theme, "theme", "t", "", "theme to use (located in /themes/THEMENAME/)")
	HugoCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	HugoCmd.PersistentFlags().BoolVar(&UglyURLs, "uglyUrls", false, "if true, use /filename.html instead of /filename/")
	HugoCmd.PersistentFlags().StringVarP(&BaseURL, "baseUrl", "b", "", "hostname (and path) to the root eg. http://spf13.com/")
	HugoCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is path/config.yaml|json|toml)")
	HugoCmd.PersistentFlags().StringVar(&Editor, "editor", "", "edit new content with this editor, if provided")
	HugoCmd.PersistentFlags().BoolVar(&Logging, "log", false, "Enable Logging")
	HugoCmd.PersistentFlags().StringVar(&LogFile, "logFile", "", "Log File path (if set, logging enabled automatically)")
	HugoCmd.PersistentFlags().BoolVar(&VerboseLog, "verboseLog", false, "verbose logging")
	HugoCmd.PersistentFlags().BoolVar(&nitro.AnalysisOn, "stepAnalysis", false, "display memory and timing of different steps of the program")
	HugoCmd.PersistentFlags().BoolVar(&PluralizeListTitles, "pluralizeListTitles", true, "Pluralize titles in lists using inflect")
	HugoCmd.Flags().BoolVarP(&BuildWatch, "watch", "w", false, "watch filesystem for changes and recreate as needed")
	HugoCmd.Flags().BoolVarP(&NoTimes, "noTimes", "", false, "Don't sync modification time of files")
	hugoCmdV = HugoCmd

	// for Bash autocomplete
	validConfigFilenames := []string{"json", "js", "yaml", "yml", "toml", "tml"}
	annotation := make(map[string][]string)
	annotation[cobra.BashCompFilenameExt] = validConfigFilenames
	HugoCmd.PersistentFlags().Lookup("config").Annotations = annotation

	// This message will be shown to Windows users if Hugo is opened from explorer.exe
	cobra.MousetrapHelpText = `

  Hugo is a command line tool

  You need to open cmd.exe and run it from there.`
}

func LoadDefaultSettings() {
	viper.SetDefault("Watch", false)
	viper.SetDefault("MetaDataFormat", "toml")
	viper.SetDefault("DisableRSS", false)
	viper.SetDefault("DisableSitemap", false)
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("LayoutDir", "layouts")
	viper.SetDefault("StaticDir", "static")
	viper.SetDefault("ArchetypeDir", "archetypes")
	viper.SetDefault("PublishDir", "public")
	viper.SetDefault("DataDir", "data")
	viper.SetDefault("DefaultLayout", "post")
	viper.SetDefault("BuildDrafts", false)
	viper.SetDefault("BuildFuture", false)
	viper.SetDefault("UglyURLs", false)
	viper.SetDefault("Verbose", false)
	viper.SetDefault("IgnoreCache", false)
	viper.SetDefault("CanonifyURLs", false)
	viper.SetDefault("RelativeURLs", false)
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	viper.SetDefault("Permalinks", make(hugolib.PermalinkOverrides, 0))
	viper.SetDefault("Sitemap", hugolib.Sitemap{Priority: -1})
	viper.SetDefault("PygmentsStyle", "monokai")
	viper.SetDefault("DefaultExtension", "html")
	viper.SetDefault("PygmentsUseClasses", false)
	viper.SetDefault("DisableLiveReload", false)
	viper.SetDefault("PluralizeListTitles", true)
	viper.SetDefault("FootnoteAnchorPrefix", "")
	viper.SetDefault("FootnoteReturnLinkContents", "")
	viper.SetDefault("NewContentEditor", "")
	viper.SetDefault("Paginate", 10)
	viper.SetDefault("PaginatePath", "page")
	viper.SetDefault("Blackfriday", helpers.NewBlackfriday())
	viper.SetDefault("RSSUri", "index.xml")
	viper.SetDefault("SectionPagesMenu", "")
}

// InitializeConfig initializes a config file with sensible default configuration flags.
func InitializeConfig() {
	viper.SetConfigFile(CfgFile)
	viper.AddConfigPath(Source)
	err := viper.ReadInConfig()
	if err != nil {
		jww.ERROR.Println("Unable to locate Config file. Perhaps you need to create a new site. Run `hugo help new` for details")
	}

	viper.RegisterAlias("indexes", "taxonomies")

	LoadDefaultSettings()

	if hugoCmdV.PersistentFlags().Lookup("buildDrafts").Changed {
		viper.Set("BuildDrafts", Draft)
	}

	if hugoCmdV.PersistentFlags().Lookup("buildFuture").Changed {
		viper.Set("BuildFuture", Future)
	}

	if hugoCmdV.PersistentFlags().Lookup("uglyUrls").Changed {
		viper.Set("UglyURLs", UglyURLs)
	}

	if hugoCmdV.PersistentFlags().Lookup("disableRSS").Changed {
		viper.Set("DisableRSS", DisableRSS)
	}

	if hugoCmdV.PersistentFlags().Lookup("disableSitemap").Changed {
		viper.Set("DisableSitemap", DisableSitemap)
	}

	if hugoCmdV.PersistentFlags().Lookup("verbose").Changed {
		viper.Set("Verbose", Verbose)
	}

	if hugoCmdV.PersistentFlags().Lookup("pluralizeListTitles").Changed {
		viper.Set("PluralizeListTitles", PluralizeListTitles)
	}

	if hugoCmdV.PersistentFlags().Lookup("editor").Changed {
		viper.Set("NewContentEditor", Editor)
	}

	if hugoCmdV.PersistentFlags().Lookup("logFile").Changed {
		viper.Set("LogFile", LogFile)
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
		viper.Set("WorkingDir", Source)
	} else {
		dir, _ := os.Getwd()
		viper.Set("WorkingDir", dir)
	}

	if hugoCmdV.PersistentFlags().Lookup("ignoreCache").Changed {
		viper.Set("IgnoreCache", IgnoreCache)
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

	themeVersionMismatch, minVersion := helpers.IsThemeVsHugoVersionMismatch()
	if themeVersionMismatch {
		jww.ERROR.Printf("Current theme does not support Hugo version %s. Minimum version required is %s\n",
			helpers.HugoReleaseVersion(), minVersion)
	}
}

func build(watches ...bool) {
	utils.CheckErr(copyStatic(), fmt.Sprintf("Error copying static files to %s", helpers.AbsPathify(viper.GetString("PublishDir"))))
	watch := false
	if len(watches) > 0 && watches[0] {
		watch = true
	}
	utils.StopOnErr(buildSite(BuildWatch || watch))

	if BuildWatch {
		jww.FEEDBACK.Println("Watching for changes in", helpers.AbsPathify(viper.GetString("ContentDir")))
		jww.FEEDBACK.Println("Press Ctrl+C to stop")
		utils.CheckErr(NewWatcher(0))
	}
}

func copyStatic() error {
	staticDir := helpers.AbsPathify(viper.GetString("StaticDir")) + "/"
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		jww.ERROR.Println("Unable to find Static Directory:", staticDir)
		return nil
	}

	publishDir := helpers.AbsPathify(viper.GetString("PublishDir")) + "/"

	syncer := fsync.NewSyncer()
	syncer.NoTimes = viper.GetBool("notimes")
	syncer.SrcFs = hugofs.SourceFs
	syncer.DestFs = hugofs.DestinationFS

	themeDir, err := helpers.GetThemeStaticDirPath()
	if err != nil {
		jww.ERROR.Println(err)
		return nil
	}

	if themeDir != "" {
		// Copy Static to Destination
		jww.INFO.Println("syncing from", themeDir, "to", publishDir)
		utils.CheckErr(syncer.Sync(publishDir, themeDir), fmt.Sprintf("Error copying static files of theme to %s", publishDir))
	}

	// Copy Static to Destination
	jww.INFO.Println("syncing from", staticDir, "to", publishDir)
	return syncer.Sync(publishDir, staticDir)
}

// getDirList provides NewWatcher() with a list of directories to watch for changes.
func getDirList() []string {
	var a []string
	dataDir := helpers.AbsPathify(viper.GetString("DataDir"))
	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			if path == dataDir && os.IsNotExist(err) {
				jww.WARN.Println("Skip DataDir:", err)
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
		filepath.Walk(helpers.AbsPathify("themes/"+viper.GetString("theme")), walker)
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
		fmt.Println(err)
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
					istemp := strings.HasSuffix(ext, "~") || (ext == ".swp") || (ext == ".swx") || (ext == ".tmp") || strings.HasPrefix(ext, ".goutputstream")
					if istemp {
						continue
					}
					// renames are always followed with Create/Modify
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						continue
					}

					isstatic := strings.HasPrefix(ev.Name, helpers.GetStaticDirPath()) || strings.HasPrefix(ev.Name, helpers.GetThemesDirPath())
					staticChanged = staticChanged || isstatic
					dynamicChanged = dynamicChanged || !isstatic

					if isstatic {
						if staticPath, err := helpers.MakeStaticPathRelative(ev.Name); err == nil {
							staticFilesChanged[staticPath] = true
						}
					}

					// add new directory to watch list
					if s, err := os.Stat(ev.Name); err == nil && s.Mode().IsDir() {
						if ev.Op&fsnotify.Create == fsnotify.Create {
							watcher.Add(ev.Name)
						}
					}
				}

				if staticChanged {
					jww.FEEDBACK.Printf("Static file changed, syncing\n\n")
					utils.StopOnErr(copyStatic(), fmt.Sprintf("Error copying static files to %s", helpers.AbsPathify(viper.GetString("PublishDir"))))

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
