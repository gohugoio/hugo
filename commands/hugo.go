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
	"github.com/howeyc/fsnotify"
	"github.com/mostafah/fsync"
	"github.com/spf13/cobra"
	"github.com/spf13/hugo/hugolib"
	"log"
	"os"
	"path/filepath"
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
	Run: build,
}

var Hugo *cobra.Commander
var BuildWatch, Draft, UglyUrls, Verbose bool
var Source, Destination, BaseUrl, CfgFile string

func Execute() {
	AddCommands()
	Hugo := HugoCmd.ToCommander()
	err := Hugo.Execute()
	if err != nil {
		os.Exit(-1)
	}
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
	HugoCmd.Flags().BoolVarP(&BuildWatch, "watch", "w", false, "watch filesystem for changes and recreate as needed")
}

func InitializeConfig() {
	Config = hugolib.SetupConfig(&CfgFile, &Source)
	Config.BuildDrafts = Draft
	Config.UglyUrls = UglyUrls
	Config.Verbose = Verbose
	if BaseUrl != "" {
		Config.BaseUrl = BaseUrl
	}
	if Destination != "" {
		Config.PublishDir = Destination
	}
}

func build(cmd *cobra.Command, args []string) {
	InitializeConfig()

	err := copyStatic()
	if err != nil {
		log.Fatalf("Error copying static files to %s: %v", Config.GetAbsPath(Config.PublishDir), err)
	}
	if _, err := buildSite(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if BuildWatch {
		fmt.Println("Watching for changes in", Config.GetAbsPath(Config.ContentDir))
		fmt.Println("Press ctrl+c to stop")
		err := NewWatcher(0)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func copyStatic() error {
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

func buildSite() (site *hugolib.Site, err error) {
	startTime := time.Now()
	site = &hugolib.Site{Config: *Config}
	err = site.Build()
	if err != nil {
		return
	}
	site.Stats()
	fmt.Printf("in %v ms\n", int(1000*time.Since(startTime).Seconds()))
	return site, nil
}

func NewWatcher(port int) error {
	watcher, err := fsnotify.NewWatcher()
	var wg sync.WaitGroup

	if err != nil {
		fmt.Println(err)
		return err
	}

	defer watcher.Close()

	wg.Add(1)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if Verbose {
					fmt.Println(ev)
				}
				watchChange(ev)
				// TODO add newly created directories to the watch list
			case err := <-watcher.Error:
				if err != nil {
					fmt.Println("error:", err)
				}
			}
		}
	}()

	for _, d := range getDirList() {
		if d != "" {
			_ = watcher.Watch(d)
		}
	}

	if port > 0 {
		go serve(port)
	}

	wg.Wait()
	return nil
}

func watchChange(ev *fsnotify.FileEvent) {
	if strings.HasPrefix(ev.Name, Config.GetAbsPath(Config.StaticDir)) {
		fmt.Println("Static file changed, syncing\n")
		copyStatic()
	} else {
		fmt.Println("Change detected, rebuilding site\n")
		buildSite()
	}
}
