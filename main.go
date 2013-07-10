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
	"./hugolib"
	"flag"
	"fmt"
	"github.com/howeyc/fsnotify"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sync"
	"time"
)

const (
	cfgFiledefault = "config.yaml"
)

var (
	baseUrl    = flag.String("b", "", "hostname (and path) to the root eg. http://spf13.com/")
	cfgfile    = flag.String("c", cfgFiledefault, "config file (default is path/config.yaml)")
	checkMode  = flag.Bool("k", false, "analyze content and provide feedback")
	draft      = flag.Bool("d", false, "include content marked as draft")
	help       = flag.Bool("h", false, "show this help")
	path       = flag.String("p", "", "filesystem path to read files relative from")
	verbose    = flag.Bool("v", false, "verbose output")
	version    = flag.Bool("version", false, "which version of hugo")
	cpuprofile = flag.Int("cpuprofile", 0, "Number of times to create the site and profile it")
	watchMode  = flag.Bool("w", false, "watch filesystem for changes and recreate as needed")
	server     = flag.Bool("s", false, "run a (very) simple web server")
	port       = flag.String("port", "1313", "port to run web server on, default :1313")
)

func usage() {
	PrintErr("usage: hugo [flags]", "")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
	}

	config := hugolib.SetupConfig(cfgfile, path)
	config.BuildDrafts = *draft

	if *baseUrl != "" {
		config.BaseUrl = *baseUrl
	} else if *server {
        config.BaseUrl = "http://localhost:" + *port
    }

	if *version {
		fmt.Println("Hugo Static Site Generator v0.8")
	}

	if *cpuprofile != 0 {
		f, err := os.Create("/tmp/hugo-cpuprofile")

		if err != nil {
			panic(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

		for i := 0; i < *cpuprofile; i++ {
			_ = buildSite(config)
		}
	}

	if *checkMode {
		site := hugolib.NewSite(config)
		site.Analyze()
		os.Exit(2)
	}

	if *watchMode {
		fmt.Println("Watching for changes. Press ctrl+c to stop")
		_ = buildSite(config)
		err := NewWatcher(config, *port, *server)

		if err != nil {
			fmt.Println(err)
		}
	}

	_ = buildSite(config)

	if *server {
		serve(*port, config)
	}

}

func serve(port string, config *hugolib.Config) {
	fmt.Println("Web Server is available at http://localhost:" + port)
	fmt.Println("Press ctrl+c to stop")
	panic(http.ListenAndServe(":"+port, http.FileServer(http.Dir(config.PublishDir))))
}

func buildSite(config *hugolib.Config) *hugolib.Site {
    startTime := time.Now()
	site := hugolib.NewSite(config)
	site.Build()
	site.Stats()
	fmt.Printf("in %v ms\n", int(1000 * time.Since(startTime).Seconds()))
	return site
}

func watchChange(c *hugolib.Config) {
	fmt.Println("Change detected, rebuilding site\n")
	buildSite(c)
}

func NewWatcher(c *hugolib.Config, port string, server bool) error {
	watcher, err := fsnotify.NewWatcher()
	var wg sync.WaitGroup

	if err != nil {
		return err
		fmt.Println(err)
	}

	defer watcher.Close()

	wg.Add(1)
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				var _ = ev
				watchChange(c)
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

	filepath.Walk(c.GetAbsPath(c.SourceDir), walker)
	filepath.Walk(c.GetAbsPath(c.LayoutDir), walker)

	return a
}

func PrintErr(str string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, str, a)
}
