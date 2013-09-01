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

package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/mostafah/fsync"
	flag "github.com/ogier/pflag"
	"github.com/spf13/hugo/hugolib"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

var (
	baseUrl     = flag.StringP("base-url", "b", "", "hostname (and path) to the root eg. http://spf13.com/")
	cfgfile     = flag.String("config", "", "config file (default is path/config.yaml|json|toml)")
	checkMode   = flag.Bool("check", false, "analyze content and provide feedback")
	draft       = flag.BoolP("build-drafts", "D", false, "include content marked as draft")
	help        = flag.BoolP("help", "h", false, "show this help")
	source      = flag.StringP("source", "s", "", "filesystem path to read files relative from")
	destination = flag.StringP("destination", "d", "", "filesystem path to write files to")
	verbose     = flag.BoolP("verbose", "v", false, "verbose output")
	version     = flag.Bool("version", false, "which version of hugo")
	cpuprofile  = flag.Int("profile", 0, "Number of times to create the site and profile it")
	watchMode   = flag.BoolP("watch", "w", false, "watch filesystem for changes and recreate as needed")
	server      = flag.BoolP("server", "S", false, "run a (very) simple web server")
	port        = flag.String("port", "1313", "port to run web server on, default :1313")
	uglyUrls    = flag.Bool("uglyurls", false, "if true, use /filename.html instead of /filename/")
)

func usage() {
	PrintErr("usage: hugo [flags]", "")
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
	}

	if *version {
		fmt.Println("Hugo Static Site Generator v0.8")
		return
	}

	config := hugolib.SetupConfig(cfgfile, source)
	config.BuildDrafts = *draft
	config.UglyUrls = *uglyUrls
	config.Verbose = *verbose

	if *baseUrl != "" {
		config.BaseUrl = *baseUrl
	} else if *server {
		config.BaseUrl = "http://localhost:" + *port
	}

	if *destination != "" {
		config.PublishDir = *destination
	}

	if *cpuprofile != 0 {
		f, err := os.Create("/tmp/hugo-cpuprofile")

		if err != nil {
			panic(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

		for i := 0; i < *cpuprofile; i++ {
			_, _ = buildSite(config)
		}
	}

	err := copyStatic(config)
	if err != nil {
		log.Fatalf("Error copying static files to %s: %v", config.GetAbsPath(config.PublishDir), err)
	}

	if *checkMode {
		site := hugolib.Site{Config: *config}
		site.Analyze()
		os.Exit(0)
	}

	if *watchMode {
		fmt.Println("Watching for changes. Press ctrl+c to stop")
		_, err = buildSite(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		err := NewWatcher(config, *port, *server)
		if err != nil {
			fmt.Println(err)
		}
	}

	if _, err = buildSite(config); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if *server {
		serve(*port, config)
	}

}

func copyStatic(config *hugolib.Config) error {
	// Copy Static to Destination
	return fsync.Sync(config.GetAbsPath(config.PublishDir+"/"), config.GetAbsPath(config.StaticDir+"/"))
}

func serve(port string, config *hugolib.Config) {

	if config.Verbose {
		fmt.Println("Serving pages from " + config.GetAbsPath(config.PublishDir))
	}

	fmt.Println("Web Server is available at http://localhost:" + port)
	fmt.Println("Press ctrl+c to stop")
	panic(http.ListenAndServe(":"+port, http.FileServer(http.Dir(config.GetAbsPath(config.PublishDir)))))
}

func buildSite(config *hugolib.Config) (site *hugolib.Site, err error) {
	startTime := time.Now()
	site = &hugolib.Site{Config: *config}
	err = site.Build()
	if err != nil {
		return
	}
	site.Stats()
	fmt.Printf("in %v ms\n", int(1000*time.Since(startTime).Seconds()))
	return site, nil
}

func watchChange(c *hugolib.Config, ev *fsnotify.FileEvent) {
	if strings.HasPrefix(ev.Name, c.GetAbsPath(c.StaticDir)) {
		fmt.Println("Static file changed, syncing\n")
		copyStatic(c)
	} else {
		fmt.Println("Change detected, rebuilding site\n")
		buildSite(c)
	}
}

func NewWatcher(c *hugolib.Config, port string, server bool) error {
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
				if c.Verbose {
					fmt.Println(ev)
				}
				watchChange(c, ev)
				// TODO add newly created directories to the watch list
			case err := <-watcher.Error:
				if err != nil {
					fmt.Println("error:", err)
				}
			}
		}
	}()

	for _, d := range getDirList(c) {
		if d != "" {
			_ = watcher.Watch(d)
		}
	}

	if server {
		go serve(port, c)
	}

	wg.Wait()
	return nil
}

func getDirList(c *hugolib.Config) []string {
	var a []string
	walker := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			PrintErr("Walker: ", err)
			return nil
		}

		if fi.IsDir() {
			a = append(a, path)
		}
		return nil
	}

	filepath.Walk(c.GetAbsPath(c.ContentDir), walker)
	filepath.Walk(c.GetAbsPath(c.LayoutDir), walker)
	filepath.Walk(c.GetAbsPath(c.StaticDir), walker)

	return a
}

func PrintErr(str string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, str, a)
}
