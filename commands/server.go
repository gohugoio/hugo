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

package commands

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gohugoio/hugo/livereload"

	"github.com/gohugoio/hugo/config"

	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var (
	disableLiveReload bool
	navigateToChanged bool
	renderToDisk      bool
	serverAppend      bool
	serverInterface   string
	serverPort        int
	liveReloadPort    int
	serverWatch       bool
	noHTTPCache       bool

	disableFastRender bool
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"serve"},
	Short:   "A high performance webserver",
	Long: `Hugo provides its own webserver which builds and serves the site.
While hugo server is high performance, it is a webserver with limited options.
Many run it in production, but the standard behavior is for people to use it
in development and use a more full featured server such as Nginx or Caddy.

'hugo server' will avoid writing the rendered and served content to disk,
preferring to store it in memory.

By default hugo will also watch your files for any changes you make and
automatically rebuild the site. It will then live reload any open browser pages
and push the latest content to them. As most Hugo sites are built in a fraction
of a second, you will be able to save and see your changes nearly instantly.`,
	//RunE: server,
}

type filesOnlyFs struct {
	fs http.FileSystem
}

type noDirFile struct {
	http.File
}

func (fs filesOnlyFs) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return noDirFile{f}, nil
}

func (f noDirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func init() {
	initHugoBuilderFlags(serverCmd)

	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 1313, "port on which the server will listen")
	serverCmd.Flags().IntVar(&liveReloadPort, "liveReloadPort", -1, "port for live reloading (i.e. 443 in HTTPS proxy situations)")
	serverCmd.Flags().StringVarP(&serverInterface, "bind", "", "127.0.0.1", "interface to which the server will bind")
	serverCmd.Flags().BoolVarP(&serverWatch, "watch", "w", true, "watch filesystem for changes and recreate as needed")
	serverCmd.Flags().BoolVar(&noHTTPCache, "noHTTPCache", false, "prevent HTTP caching")
	serverCmd.Flags().BoolVarP(&serverAppend, "appendPort", "", true, "append port to baseURL")
	serverCmd.Flags().BoolVar(&disableLiveReload, "disableLiveReload", false, "watch without enabling live browser reload on rebuild")
	serverCmd.Flags().BoolVar(&navigateToChanged, "navigateToChanged", false, "navigate to changed content file on live browser reload")
	serverCmd.Flags().BoolVar(&renderToDisk, "renderToDisk", false, "render to Destination path (default is render to memory & serve from there)")
	serverCmd.Flags().BoolVar(&disableFastRender, "disableFastRender", false, "enables full re-renders on changes")

	serverCmd.Flags().String("memstats", "", "log memory usage to this file")
	serverCmd.Flags().String("meminterval", "100ms", "interval to poll memory usage (requires --memstats), valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")

	serverCmd.RunE = server

}

func server(cmd *cobra.Command, args []string) error {
	cfg, err := InitializeConfig(serverCmd)
	if err != nil {
		return err
	}

	c, err := newCommandeer(cfg)
	if err != nil {
		return err
	}

	if cmd.Flags().Changed("disableLiveReload") {
		c.Set("disableLiveReload", disableLiveReload)
	}

	if cmd.Flags().Changed("navigateToChanged") {
		c.Set("navigateToChanged", navigateToChanged)
	}

	if cmd.Flags().Changed("disableFastRender") {
		c.Set("disableFastRender", disableFastRender)
	}

	if serverWatch {
		c.Set("watch", true)
	}

	if c.Cfg.GetBool("watch") {
		serverWatch = true
		c.watchConfig()
	}

	languages := c.languages()
	serverPorts := make([]int, 1)

	if languages.IsMultihost() {
		if !serverAppend {
			return newSystemError("--appendPort=false not supported when in multihost mode")
		}
		serverPorts = make([]int, len(languages))
	}

	currentServerPort := serverPort

	for i := 0; i < len(serverPorts); i++ {
		l, err := net.Listen("tcp", net.JoinHostPort(serverInterface, strconv.Itoa(currentServerPort)))
		if err == nil {
			l.Close()
			serverPorts[i] = currentServerPort
		} else {
			if i == 0 && serverCmd.Flags().Changed("port") {
				// port set explicitly by user -- he/she probably meant it!
				return newSystemErrorF("Server startup failed: %s", err)
			}
			jww.ERROR.Println("port", serverPort, "already in use, attempting to use an available port")
			sp, err := helpers.FindAvailablePort()
			if err != nil {
				return newSystemError("Unable to find alternative port to use:", err)
			}
			serverPorts[i] = sp.Port
		}

		currentServerPort = serverPorts[i] + 1
	}

	c.serverPorts = serverPorts

	c.Set("port", serverPort)
	if liveReloadPort != -1 {
		c.Set("liveReloadPort", liveReloadPort)
	} else {
		c.Set("liveReloadPort", serverPorts[0])
	}

	if languages.IsMultihost() {
		for i, language := range languages {
			baseURL, err = fixURL(language, baseURL, serverPorts[i])
			if err != nil {
				return err
			}
			language.Set("baseURL", baseURL)
		}
	} else {
		baseURL, err = fixURL(c.Cfg, baseURL, serverPorts[0])
		if err != nil {
			return err
		}
		c.Set("baseURL", baseURL)
	}

	if err := memStats(); err != nil {
		jww.ERROR.Println("memstats error:", err)
	}

	// If a Destination is provided via flag write to disk
	if destination != "" {
		renderToDisk = true
	}

	// Hugo writes the output to memory instead of the disk
	if !renderToDisk {
		cfg.Fs.Destination = new(afero.MemMapFs)
		// Rendering to memoryFS, publish to Root regardless of publishDir.
		c.Set("publishDir", "/")
	}

	if err := c.build(serverWatch); err != nil {
		return err
	}

	for _, s := range Hugo.Sites {
		s.RegisterMediaTypes()
	}

	// Watch runs its own server as part of the routine
	if serverWatch {

		watchDirs, err := c.getDirList()
		if err != nil {
			return err
		}

		baseWatchDir := c.Cfg.GetString("workingDir")
		relWatchDirs := make([]string, len(watchDirs))
		for i, dir := range watchDirs {
			relWatchDirs[i], _ = helpers.GetRelativePath(dir, baseWatchDir)
		}

		rootWatchDirs := strings.Join(helpers.UniqueStrings(helpers.ExtractRootPaths(relWatchDirs)), ",")

		jww.FEEDBACK.Printf("Watching for changes in %s%s{%s}\n", baseWatchDir, helpers.FilePathSeparator, rootWatchDirs)
		err = c.newWatcher(true, watchDirs...)

		if err != nil {
			return err
		}
	}

	return nil
}

type fileServer struct {
	baseURLs []string
	roots    []string
	c        *commandeer
}

func (f *fileServer) createEndpoint(i int) (*http.ServeMux, string, string, error) {
	baseURL := f.baseURLs[i]
	root := f.roots[i]
	port := f.c.serverPorts[i]

	publishDir := f.c.Cfg.GetString("publishDir")

	if root != "" {
		publishDir = filepath.Join(publishDir, root)
	}

	absPublishDir := f.c.PathSpec().AbsPathify(publishDir)

	if i == 0 {
		if renderToDisk {
			jww.FEEDBACK.Println("Serving pages from " + absPublishDir)
		} else {
			jww.FEEDBACK.Println("Serving pages from memory")
		}
	}

	httpFs := afero.NewHttpFs(f.c.Fs.Destination)
	fs := filesOnlyFs{httpFs.Dir(absPublishDir)}

	doLiveReload := !buildWatch && !f.c.Cfg.GetBool("disableLiveReload")
	fastRenderMode := doLiveReload && !f.c.Cfg.GetBool("disableFastRender")

	if i == 0 && fastRenderMode {
		jww.FEEDBACK.Println("Running in Fast Render Mode. For full rebuilds on change: hugo server --disableFastRender")
	}

	// We're only interested in the path
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, "", "", fmt.Errorf("Invalid baseURL: %s", err)
	}

	decorate := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if noHTTPCache {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
				w.Header().Set("Pragma", "no-cache")
			}

			if fastRenderMode {
				p := r.RequestURI
				if strings.HasSuffix(p, "/") || strings.HasSuffix(p, "html") || strings.HasSuffix(p, "htm") {
					f.c.visitedURLs.Add(p)
				}
			}
			h.ServeHTTP(w, r)
		})
	}

	fileserver := decorate(http.FileServer(fs))
	mu := http.NewServeMux()

	if u.Path == "" || u.Path == "/" {
		mu.Handle("/", fileserver)
	} else {
		mu.Handle(u.Path, http.StripPrefix(u.Path, fileserver))
	}

	endpoint := net.JoinHostPort(serverInterface, strconv.Itoa(port))

	return mu, u.String(), endpoint, nil
}

func (c *commandeer) serve() {

	isMultiHost := Hugo.IsMultihost()

	var (
		baseURLs []string
		roots    []string
	)

	if isMultiHost {
		for _, s := range Hugo.Sites {
			baseURLs = append(baseURLs, s.BaseURL.String())
			roots = append(roots, s.Language.Lang)
		}
	} else {
		s := Hugo.Sites[0]
		baseURLs = []string{s.BaseURL.String()}
		roots = []string{""}
	}

	srv := &fileServer{
		baseURLs: baseURLs,
		roots:    roots,
		c:        c,
	}

	doLiveReload := !c.Cfg.GetBool("disableLiveReload")

	if doLiveReload {
		livereload.Initialize()
	}

	for i, _ := range baseURLs {
		mu, serverURL, endpoint, err := srv.createEndpoint(i)

		if doLiveReload {
			mu.HandleFunc("/livereload.js", livereload.ServeJS)
			mu.HandleFunc("/livereload", livereload.Handler)
		}
		jww.FEEDBACK.Printf("Web Server is available at %s (bind address %s)\n", serverURL, serverInterface)
		go func() {
			err = http.ListenAndServe(endpoint, mu)
			if err != nil {
				jww.ERROR.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}
		}()
	}

	jww.FEEDBACK.Println("Press Ctrl+C to stop")
}

// fixURL massages the baseURL into a form needed for serving
// all pages correctly.
func fixURL(cfg config.Provider, s string, port int) (string, error) {
	useLocalhost := false
	if s == "" {
		s = cfg.GetString("baseURL")
		useLocalhost = true
	}

	if !strings.HasSuffix(s, "/") {
		s = s + "/"
	}

	// do an initial parse of the input string
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	// if no Host is defined, then assume that no schema or double-slash were
	// present in the url.  Add a double-slash and make a best effort attempt.
	if u.Host == "" && s != "/" {
		s = "//" + s

		u, err = url.Parse(s)
		if err != nil {
			return "", err
		}
	}

	if useLocalhost {
		if u.Scheme == "https" {
			u.Scheme = "http"
		}
		u.Host = "localhost"
	}

	if serverAppend {
		if strings.Contains(u.Host, ":") {
			u.Host, _, err = net.SplitHostPort(u.Host)
			if err != nil {
				return "", fmt.Errorf("Failed to split baseURL hostpost: %s", err)
			}
		}
		u.Host += fmt.Sprintf(":%d", port)
	}

	return u.String(), nil
}

func memStats() error {
	memstats := serverCmd.Flags().Lookup("memstats").Value.String()
	if memstats != "" {
		interval, err := time.ParseDuration(serverCmd.Flags().Lookup("meminterval").Value.String())
		if err != nil {
			interval, _ = time.ParseDuration("100ms")
		}

		fileMemStats, err := os.Create(memstats)
		if err != nil {
			return err
		}

		fileMemStats.WriteString("# Time\tHeapSys\tHeapAlloc\tHeapIdle\tHeapReleased\n")

		go func() {
			var stats runtime.MemStats

			start := time.Now().UnixNano()

			for {
				runtime.ReadMemStats(&stats)
				if fileMemStats != nil {
					fileMemStats.WriteString(fmt.Sprintf("%d\t%d\t%d\t%d\t%d\n",
						(time.Now().UnixNano()-start)/1000000, stats.HeapSys, stats.HeapAlloc, stats.HeapIdle, stats.HeapReleased))
					time.Sleep(interval)
				} else {
					break
				}
			}
		}()
	}
	return nil
}
