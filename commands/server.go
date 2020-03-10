// Copyright 2019 The Hugo Authors. All rights reserved.
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
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/gohugoio/hugo/livereload"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

type serverCmd struct {
	// Can be used to stop the server. Useful in tests
	stop <-chan bool

	disableLiveReload bool
	navigateToChanged bool
	renderToDisk      bool
	serverAppend      bool
	serverInterface   string
	serverPort        int
	liveReloadPort    int
	serverWatch       bool
	noHTTPCache       bool

	disableFastRender   bool
	disableBrowserError bool

	*baseBuilderCmd
}

func (b *commandsBuilder) newServerCmd() *serverCmd {
	return b.newServerCmdSignaled(nil)
}

func (b *commandsBuilder) newServerCmdSignaled(stop <-chan bool) *serverCmd {
	cc := &serverCmd{stop: stop}

	cc.baseBuilderCmd = b.newBuilderCmd(&cobra.Command{
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
		RunE: cc.server,
	})

	cc.cmd.Flags().IntVarP(&cc.serverPort, "port", "p", 1313, "port on which the server will listen")
	cc.cmd.Flags().IntVar(&cc.liveReloadPort, "liveReloadPort", -1, "port for live reloading (i.e. 443 in HTTPS proxy situations)")
	cc.cmd.Flags().StringVarP(&cc.serverInterface, "bind", "", "127.0.0.1", "interface to which the server will bind")
	cc.cmd.Flags().BoolVarP(&cc.serverWatch, "watch", "w", true, "watch filesystem for changes and recreate as needed")
	cc.cmd.Flags().BoolVar(&cc.noHTTPCache, "noHTTPCache", false, "prevent HTTP caching")
	cc.cmd.Flags().BoolVarP(&cc.serverAppend, "appendPort", "", true, "append port to baseURL")
	cc.cmd.Flags().BoolVar(&cc.disableLiveReload, "disableLiveReload", false, "watch without enabling live browser reload on rebuild")
	cc.cmd.Flags().BoolVar(&cc.navigateToChanged, "navigateToChanged", false, "navigate to changed content file on live browser reload")
	cc.cmd.Flags().BoolVar(&cc.renderToDisk, "renderToDisk", false, "render to Destination path (default is render to memory & serve from there)")
	cc.cmd.Flags().BoolVar(&cc.disableFastRender, "disableFastRender", false, "enables full re-renders on changes")
	cc.cmd.Flags().BoolVar(&cc.disableBrowserError, "disableBrowserError", false, "do not show build errors in the browser")

	cc.cmd.Flags().String("memstats", "", "log memory usage to this file")
	cc.cmd.Flags().String("meminterval", "100ms", "interval to poll memory usage (requires --memstats), valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")

	return cc
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

var serverPorts []int

func (sc *serverCmd) server(cmd *cobra.Command, args []string) error {
	// If a Destination is provided via flag write to disk
	destination, _ := cmd.Flags().GetString("destination")
	if destination != "" {
		sc.renderToDisk = true
	}

	var serverCfgInit sync.Once

	cfgInit := func(c *commandeer) error {
		c.Set("renderToMemory", !sc.renderToDisk)
		if cmd.Flags().Changed("navigateToChanged") {
			c.Set("navigateToChanged", sc.navigateToChanged)
		}
		if cmd.Flags().Changed("disableLiveReload") {
			c.Set("disableLiveReload", sc.disableLiveReload)
		}
		if cmd.Flags().Changed("disableFastRender") {
			c.Set("disableFastRender", sc.disableFastRender)
		}
		if cmd.Flags().Changed("disableBrowserError") {
			c.Set("disableBrowserError", sc.disableBrowserError)
		}
		if sc.serverWatch {
			c.Set("watch", true)
		}

		// TODO(bep) yes, we should fix.
		if !c.languagesConfigured {
			return nil
		}

		var err error

		// We can only do this once.
		serverCfgInit.Do(func() {
			serverPorts = make([]int, 1)

			if c.languages.IsMultihost() {
				if !sc.serverAppend {
					err = newSystemError("--appendPort=false not supported when in multihost mode")
				}
				serverPorts = make([]int, len(c.languages))
			}

			currentServerPort := sc.serverPort

			for i := 0; i < len(serverPorts); i++ {
				l, err := net.Listen("tcp", net.JoinHostPort(sc.serverInterface, strconv.Itoa(currentServerPort)))
				if err == nil {
					l.Close()
					serverPorts[i] = currentServerPort
				} else {
					if i == 0 && sc.cmd.Flags().Changed("port") {
						// port set explicitly by user -- he/she probably meant it!
						err = newSystemErrorF("Server startup failed: %s", err)
					}
					c.logger.FEEDBACK.Println("port", sc.serverPort, "already in use, attempting to use an available port")
					sp, err := helpers.FindAvailablePort()
					if err != nil {
						err = newSystemError("Unable to find alternative port to use:", err)
					}
					serverPorts[i] = sp.Port
				}

				currentServerPort = serverPorts[i] + 1
			}
		})

		c.serverPorts = serverPorts

		c.Set("port", sc.serverPort)
		if sc.liveReloadPort != -1 {
			c.Set("liveReloadPort", sc.liveReloadPort)
		} else {
			c.Set("liveReloadPort", serverPorts[0])
		}

		isMultiHost := c.languages.IsMultihost()
		for i, language := range c.languages {
			var serverPort int
			if isMultiHost {
				serverPort = serverPorts[i]
			} else {
				serverPort = serverPorts[0]
			}

			baseURL, err := sc.fixURL(language, sc.baseURL, serverPort)
			if err != nil {
				return nil
			}
			if isMultiHost {
				language.Set("baseURL", baseURL)
			}
			if i == 0 {
				c.Set("baseURL", baseURL)
			}
		}

		return err

	}

	if err := memStats(); err != nil {
		jww.WARN.Println("memstats error:", err)
	}

	c, err := initializeConfig(true, true, &sc.hugoBuilderCommon, sc, cfgInit)
	if err != nil {
		return err
	}

	if err := c.serverBuild(); err != nil {
		return err
	}

	for _, s := range c.hugo().Sites {
		s.RegisterMediaTypes()
	}

	// Watch runs its own server as part of the routine
	if sc.serverWatch {

		watchDirs, err := c.getDirList()
		if err != nil {
			return err
		}

		watchGroups := helpers.ExtractAndGroupRootPaths(watchDirs)

		for _, group := range watchGroups {
			jww.FEEDBACK.Printf("Watching for changes in %s\n", group)
		}
		watcher, err := c.newWatcher(watchDirs...)

		if err != nil {
			return err
		}

		defer watcher.Close()

	}

	return c.serve(sc)

}

func getRootWatchDirsStr(baseDir string, watchDirs []string) string {
	relWatchDirs := make([]string, len(watchDirs))
	for i, dir := range watchDirs {
		relWatchDirs[i], _ = helpers.GetRelativePath(dir, baseDir)
	}

	return strings.Join(helpers.UniqueStringsSorted(helpers.ExtractRootPaths(relWatchDirs)), ",")
}

type fileServer struct {
	baseURLs      []string
	roots         []string
	errorTemplate func(err interface{}) (io.Reader, error)
	c             *commandeer
	s             *serverCmd
}

func (f *fileServer) createEndpoint(i int) (*http.ServeMux, string, string, error) {
	baseURL := f.baseURLs[i]
	root := f.roots[i]
	port := f.c.serverPorts[i]

	publishDir := f.c.Cfg.GetString("publishDir")

	if root != "" {
		publishDir = filepath.Join(publishDir, root)
	}

	absPublishDir := f.c.hugo().PathSpec.AbsPathify(publishDir)

	jww.FEEDBACK.Printf("Environment: %q", f.c.hugo().Deps.Site.Hugo().Environment)

	if i == 0 {
		if f.s.renderToDisk {
			jww.FEEDBACK.Println("Serving pages from " + absPublishDir)
		} else {
			jww.FEEDBACK.Println("Serving pages from memory")
		}
	}

	httpFs := afero.NewHttpFs(f.c.destinationFs)
	fs := filesOnlyFs{httpFs.Dir(absPublishDir)}

	if i == 0 && f.c.fastRenderMode {
		jww.FEEDBACK.Println("Running in Fast Render Mode. For full rebuilds on change: hugo server --disableFastRender")
	}

	// We're only interested in the path
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, "", "", errors.Wrap(err, "Invalid baseURL")
	}

	decorate := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if f.c.showErrorInBrowser {
				// First check the error state
				err := f.c.getErrorWithContext()
				if err != nil {
					f.c.wasError = true
					w.WriteHeader(500)
					r, err := f.errorTemplate(err)
					if err != nil {
						f.c.logger.ERROR.Println(err)
					}

					port = 1313
					if !f.c.paused {
						port = f.c.Cfg.GetInt("liveReloadPort")
					}
					fmt.Fprint(w, injectLiveReloadScript(r, port))

					return
				}
			}

			if f.s.noHTTPCache {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
				w.Header().Set("Pragma", "no-cache")
			}

			for _, header := range f.c.serverConfig.Match(r.RequestURI) {
				w.Header().Set(header.Key, header.Value)
			}

			if f.c.fastRenderMode && f.c.buildErr == nil {
				p := r.RequestURI
				if strings.HasSuffix(p, "/") || strings.HasSuffix(p, "html") || strings.HasSuffix(p, "htm") {
					if !f.c.visitedURLs.Contains(p) {
						// If not already on stack, re-render that single page.
						if err := f.c.partialReRender(p); err != nil {
							f.c.handleBuildErr(err, fmt.Sprintf("Failed to render %q", p))
							if f.c.showErrorInBrowser {
								http.Redirect(w, r, p, http.StatusMovedPermanently)
								return
							}
						}
					}

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

	endpoint := net.JoinHostPort(f.s.serverInterface, strconv.Itoa(port))

	return mu, u.String(), endpoint, nil
}

var logErrorRe = regexp.MustCompile(`(?s)ERROR \d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} `)

func removeErrorPrefixFromLog(content string) string {
	return logErrorRe.ReplaceAllLiteralString(content, "")
}
func (c *commandeer) serve(s *serverCmd) error {

	isMultiHost := c.hugo().IsMultihost()

	var (
		baseURLs []string
		roots    []string
	)

	if isMultiHost {
		for _, s := range c.hugo().Sites {
			baseURLs = append(baseURLs, s.BaseURL.String())
			roots = append(roots, s.Language().Lang)
		}
	} else {
		s := c.hugo().Sites[0]
		baseURLs = []string{s.BaseURL.String()}
		roots = []string{""}
	}

	templ, err := c.hugo().TextTmpl().Parse("__default_server_error", buildErrorTemplate)
	if err != nil {
		return err
	}

	srv := &fileServer{
		baseURLs: baseURLs,
		roots:    roots,
		c:        c,
		s:        s,
		errorTemplate: func(ctx interface{}) (io.Reader, error) {
			b := &bytes.Buffer{}
			err := c.hugo().Tmpl().Execute(templ, b, ctx)
			return b, err
		},
	}

	doLiveReload := !c.Cfg.GetBool("disableLiveReload")

	if doLiveReload {
		livereload.Initialize()
	}

	var sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for i := range baseURLs {
		mu, serverURL, endpoint, err := srv.createEndpoint(i)

		if doLiveReload {
			mu.HandleFunc("/livereload.js", livereload.ServeJS)
			mu.HandleFunc("/livereload", livereload.Handler)
		}
		jww.FEEDBACK.Printf("Web Server is available at %s (bind address %s)\n", serverURL, s.serverInterface)
		go func() {
			err = http.ListenAndServe(endpoint, mu)
			if err != nil {
				c.logger.ERROR.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}
		}()
	}

	jww.FEEDBACK.Println("Press Ctrl+C to stop")

	if s.stop != nil {
		select {
		case <-sigs:
		case <-s.stop:
		}
	} else {
		<-sigs
	}

	return nil
}

// fixURL massages the baseURL into a form needed for serving
// all pages correctly.
func (sc *serverCmd) fixURL(cfg config.Provider, s string, port int) (string, error) {
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

	if sc.serverAppend {
		if strings.Contains(u.Host, ":") {
			u.Host, _, err = net.SplitHostPort(u.Host)
			if err != nil {
				return "", errors.Wrap(err, "Failed to split baseURL hostpost")
			}
		}
		u.Host += fmt.Sprintf(":%d", port)
	}

	return u.String(), nil
}

func memStats() error {
	b := newCommandsBuilder()
	sc := b.newServerCmd().getCommand()
	memstats := sc.Flags().Lookup("memstats").Value.String()
	if memstats != "" {
		interval, err := time.ParseDuration(sc.Flags().Lookup("meminterval").Value.String())
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
