// Copyright © 2013 Steve Francia <spf@spf13.com>.
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

package commands

import (
	"fmt"
	"github.com/mostafah/fsync"
	"github.com/spf13/cobra"
	"github.com/spf13/hugo/hugolib"
	"github.com/spf13/hugo/utils"
	"github.com/spf13/hugo/watcher"
	"github.com/spf13/nitro"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var Config *hugolib.Config
var HugoCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
love by spf13 and friends in Go.

Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		build()
	},
}
var hugoCmdV *cobra.Command

var BuildWatch, Draft, UglyUrls, Verbose bool
var Source, Destination, BaseUrl, CfgFile string

func Execute() {
	AddCommands()
	utils.StopOnErr(HugoCmd.Execute())
}

func AddCommands() {
	HugoCmd.AddCommand(serverCmd)
	HugoCmd.AddCommand(version)
	HugoCmd.AddCommand(check)
	HugoCmd.AddCommand(benchmark)
}

func init() {
	HugoCmd.PersistentFlags().BoolVarP(&Draft, "build-drafts", "D", false, "include content marked as draft")
	HugoCmd.PersistentFlags().StringVarP(&Source, "source", "s", "", "filesystem path to read files relative from")
	HugoCmd.PersistentFlags().StringVarP(&Destination, "destination", "d", "", "filesystem path to write files to")
	HugoCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	HugoCmd.PersistentFlags().BoolVar(&UglyUrls, "uglyurls", false, "if true, use /filename.html instead of /filename/")
	HugoCmd.PersistentFlags().StringVarP(&BaseUrl, "base-url", "b", "", "hostname (and path) to the root eg. http://spf13.com/")
	HugoCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is path/config.yaml|json|toml)")
	HugoCmd.PersistentFlags().BoolVar(&nitro.AnalysisOn, "stepAnalysis", false, "display memory and timing of different steps of the program")
	HugoCmd.Flags().BoolVarP(&BuildWatch, "watch", "w", false, "watch filesystem for changes and recreate as needed")
	hugoCmdV = HugoCmd
}

func InitializeConfig() {
	Config = hugolib.SetupConfig(&CfgFile, &Source)

	if hugoCmdV.PersistentFlags().Lookup("build-drafts").Changed {
		Config.BuildDrafts = Draft
	}

	if hugoCmdV.PersistentFlags().Lookup("uglyurls").Changed {
		Config.UglyUrls = UglyUrls
	}

	if hugoCmdV.PersistentFlags().Lookup("verbose").Changed {
		Config.Verbose = Verbose
	}
	if BaseUrl != "" {
		Config.BaseUrl = BaseUrl
	}
	if Destination != "" {
		Config.PublishDir = Destination
	}
}

func build(watches ...bool) {
	utils.CheckErr(copyStatic(), fmt.Sprintf("Error copying static files to %s", Config.GetAbsPath(Config.PublishDir)))
	watch := false
	if len(watches) > 0 && watches[0] {
		watch = true
	}
	utils.StopOnErr(buildSite(BuildWatch || watch))

	if BuildWatch {
		fmt.Println("Watching for changes in", Config.GetAbsPath(Config.ContentDir))
		fmt.Println("Press ctrl+c to stop")
		utils.CheckErr(NewWatcher(0))
	}
}

func copyStatic() error {
	staticDir := Config.GetAbsPath(Config.StaticDir + "/")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		return nil
	}

	// Copy Static to Destination
	return fsync.Sync(Config.GetAbsPath(Config.PublishDir+"/"), Config.GetAbsPath(Config.StaticDir+"/"))
}

func getDirList() []string {
	var a []string
	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Walker: ", err)
			return nil
		}

		if fi.IsDir() {
			a = append(a, path)
		}
		return nil
	}

	filepath.Walk(Config.GetAbsPath(Config.ContentDir), walker)
	filepath.Walk(Config.GetAbsPath(Config.LayoutDir), walker)
	filepath.Walk(Config.GetAbsPath(Config.StaticDir), walker)

	return a
}

func buildSite(watching ...bool) (err error) {
	startTime := time.Now()
	site := &hugolib.Site{Config: *Config}
	if len(watching) > 0 && watching[0] {
		site.RunMode.Watching = true
	}
	err = site.Build()
	if err != nil {
		return
	}
	site.Stats()
	fmt.Printf("in %v ms\n", int(1000*time.Since(startTime).Seconds()))
	return nil
}

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
			_ = watcher.Watch(d)
		}
	}

	go func() {
		for {
			select {
			case evs := <-watcher.Event:
				if Verbose {
					fmt.Println(evs)
				}

				static_changed := false
				dynamic_changed := false

				for _, ev := range evs {
					ext := filepath.Ext(ev.Name)
					istemp := strings.HasSuffix(ext, "~") || (ext == ".swp") || (ext == ".tmp")
					if istemp {
						continue
					}
					// renames are always followed with Create/Modify
					if ev.IsRename() {
						continue
					}

					isstatic := strings.HasPrefix(ev.Name, Config.GetAbsPath(Config.StaticDir))
					static_changed = static_changed || isstatic
					dynamic_changed = dynamic_changed || !isstatic

					// add new directory to watch list
					if s, err := os.Stat(ev.Name); err == nil && s.Mode().IsDir() {
						if ev.IsCreate() {
							watcher.Watch(ev.Name)
						}
					}
				}

				if static_changed {
					fmt.Print("Static file changed, syncing\n\n")
					utils.CheckErr(copyStatic(), fmt.Sprintf("Error copying static files to %s", Config.GetAbsPath(Config.PublishDir)))
				}

				if dynamic_changed {
					fmt.Print("Change detected, rebuilding site\n\n")
					utils.StopOnErr(buildSite(true))
				}
			case err := <-watcher.Error:
				if err != nil {
					fmt.Println("error:", err)
				}
			}
		}
	}()

	if port > 0 {
		go serve(port)
	}

	wg.Wait()
	return nil
}
