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
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gohugoio/hugo/common/htime"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/tpl"
	"golang.org/x/sync/errgroup"

	"github.com/gohugoio/hugo/livereload"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

type serverCmd struct {
	// Can be used to stop the server. Useful in tests
	stop chan bool

	disableLiveReload  bool
	navigateToChanged  bool
	renderToDisk       bool
	renderStaticToDisk bool
	serverAppend       bool
	serverInterface    string
	serverPort         int
	liveReloadPort     int
	serverWatch        bool
	noHTTPCache        bool

	disableFastRender   bool
	disableBrowserError bool

	*baseBuilderCmd
}

func (b *commandsBuilder) newServerCmd() *serverCmd {
	return b.newServerCmdSignaled(nil)
}

func (b *commandsBuilder) newServerCmdSignaled(stop chan bool) *serverCmd {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.server(cmd, args)
			if err != nil && cc.stop != nil {
				cc.stop <- true
			}
			return err
		},
	})

	cc.cmd.Flags().IntVarP(&cc.serverPort, "port", "p", 1313, "port on which the server will listen")
	cc.cmd.Flags().IntVar(&cc.liveReloadPort, "liveReloadPort", -1, "port for live reloading (i.e. 443 in HTTPS proxy situations)")
	cc.cmd.Flags().StringVarP(&cc.serverInterface, "bind", "", "127.0.0.1", "interface to which the server will bind")
	cc.cmd.Flags().BoolVarP(&cc.serverWatch, "watch", "w", true, "watch filesystem for changes and recreate as needed")
	cc.cmd.Flags().BoolVar(&cc.noHTTPCache, "noHTTPCache", false, "prevent HTTP caching")
	cc.cmd.Flags().BoolVarP(&cc.serverAppend, "appendPort", "", true, "append port to baseURL")
	cc.cmd.Flags().BoolVar(&cc.disableLiveReload, "disableLiveReload", false, "watch without enabling live browser reload on rebuild")
	cc.cmd.Flags().BoolVar(&cc.navigateToChanged, "navigateToChanged", false, "navigate to changed content file on live browser reload")
	cc.cmd.Flags().BoolVar(&cc.renderToDisk, "renderToDisk", false, "serve all files from disk (default is from memory)")
	cc.cmd.Flags().BoolVar(&cc.renderStaticToDisk, "renderStaticToDisk", false, "serve static files from disk and dynamic files from memory")
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

func (sc *serverCmd) server(cmd *cobra.Command, args []string) error {
	// If a Destination is provided via flag write to disk
	destination, _ := cmd.Flags().GetString("destination")
	if destination != "" {
		sc.renderToDisk = true
	}

	var serverCfgInit sync.Once

	cfgInit := func(c *commandeer) (rerr error) {
		c.Set("renderToMemory", !(sc.renderToDisk || sc.renderStaticToDisk))
		c.Set("renderStaticToDisk", sc.renderStaticToDisk)
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

		// TODO(bep) see issue 9901
		// cfgInit is called twice, before and after the languages have been initialized.
		// The servers (below) can not be initialized before we
		// know if we're configured in a multihost setup.
		if len(c.languages) == 0 {
			return nil
		}

		// We can only do this once.
		serverCfgInit.Do(func() {
			c.serverPorts = make([]serverPortListener, 1)

			if c.languages.IsMultihost() {
				if !sc.serverAppend {
					rerr = newSystemError("--appendPort=false not supported when in multihost mode")
				}
				c.serverPorts = make([]serverPortListener, len(c.languages))
			}

			currentServerPort := sc.serverPort

			for i := 0; i < len(c.serverPorts); i++ {
				l, err := net.Listen("tcp", net.JoinHostPort(sc.serverInterface, strconv.Itoa(currentServerPort)))
				if err == nil {
					c.serverPorts[i] = serverPortListener{ln: l, p: currentServerPort}
				} else {
					if i == 0 && sc.cmd.Flags().Changed("port") {
						// port set explicitly by user -- he/she probably meant it!
						rerr = newSystemErrorF("Server startup failed: %s", err)
						return
					}
					c.logger.Println("port", sc.serverPort, "already in use, attempting to use an available port")
					l, sp, err := helpers.TCPListen()
					if err != nil {
						rerr = newSystemError("Unable to find alternative port to use:", err)
						return
					}
					c.serverPorts[i] = serverPortListener{ln: l, p: sp.Port}
				}

				currentServerPort = c.serverPorts[i].p + 1
			}
		})

		if rerr != nil {
			return
		}

		c.Set("port", sc.serverPort)
		if sc.liveReloadPort != -1 {
			c.Set("liveReloadPort", sc.liveReloadPort)
		} else {
			c.Set("liveReloadPort", c.serverPorts[0].p)
		}

		isMultiHost := c.languages.IsMultihost()
		for i, language := range c.languages {
			var serverPort int
			if isMultiHost {
				serverPort = c.serverPorts[i].p
			} else {
				serverPort = c.serverPorts[0].p
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

		return
	}

	if err := memStats(); err != nil {
		jww.WARN.Println("memstats error:", err)
	}

	// silence errors in cobra so we can handle them here
	cmd.SilenceErrors = true

	c, err := initializeConfig(true, true, true, &sc.hugoBuilderCommon, sc, cfgInit)
	if err != nil {
		cmd.PrintErrln("Error:", err.Error())
		return err
	}

	err = func() error {
		defer c.timeTrack(time.Now(), "Built")
		err := c.serverBuild()
		if err != nil {
			cmd.PrintErrln("Error:", err.Error())
		}
		return err
	}()
	if err != nil {
		return err
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
		watcher, err := c.newWatcher(sc.poll, watchDirs...)
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
		relWatchDirs[i], _ = paths.GetRelativePath(dir, baseDir)
	}

	return strings.Join(helpers.UniqueStringsSorted(helpers.ExtractRootPaths(relWatchDirs)), ",")
}

type fileServer struct {
	baseURLs      []string
	roots         []string
	errorTemplate func(err any) (io.Reader, error)
	c             *commandeer
	s             *serverCmd
}

func (f *fileServer) rewriteRequest(r *http.Request, toPath string) *http.Request {
	r2 := new(http.Request)
	*r2 = *r
	r2.URL = new(url.URL)
	*r2.URL = *r.URL
	r2.URL.Path = toPath
	r2.Header.Set("X-Rewrite-Original-URI", r.URL.RequestURI())

	return r2
}

func (f *fileServer) createEndpoint(i int) (*http.ServeMux, net.Listener, string, string, error) {
	baseURL := f.baseURLs[i]
	root := f.roots[i]
	port := f.c.serverPorts[i].p
	listener := f.c.serverPorts[i].ln

	// For logging only.
	// TODO(bep) consolidate.
	publishDir := f.c.Cfg.GetString("publishDir")
	publishDirStatic := f.c.Cfg.GetString("publishDirStatic")
	workingDir := f.c.Cfg.GetString("workingDir")

	if root != "" {
		publishDir = filepath.Join(publishDir, root)
		publishDirStatic = filepath.Join(publishDirStatic, root)
	}
	absPublishDir := paths.AbsPathify(workingDir, publishDir)
	absPublishDirStatic := paths.AbsPathify(workingDir, publishDirStatic)

	jww.FEEDBACK.Printf("Environment: %q", f.c.hugo().Deps.Site.Hugo().Environment)

	if i == 0 {
		if f.s.renderToDisk {
			jww.FEEDBACK.Println("Serving pages from " + absPublishDir)
		} else if f.s.renderStaticToDisk {
			jww.FEEDBACK.Println("Serving pages from memory and static files from " + absPublishDirStatic)
		} else {
			jww.FEEDBACK.Println("Serving pages from memory")
		}
	}

	httpFs := afero.NewHttpFs(f.c.publishDirServerFs)
	fs := filesOnlyFs{httpFs.Dir(path.Join("/", root))}

	if i == 0 && f.c.fastRenderMode {
		jww.FEEDBACK.Println("Running in Fast Render Mode. For full rebuilds on change: hugo server --disableFastRender")
	}

	// We're only interested in the path
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, nil, "", "", fmt.Errorf("Invalid baseURL: %w", err)
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
						f.c.logger.Errorln(err)
					}

					port = 1313
					if !f.c.paused {
						port = f.c.Cfg.GetInt("liveReloadPort")
					}
					lr := *u
					lr.Host = fmt.Sprintf("%s:%d", lr.Hostname(), port)
					fmt.Fprint(w, injectLiveReloadScript(r, lr))

					return
				}
			}

			if f.s.noHTTPCache {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
				w.Header().Set("Pragma", "no-cache")
			}

			// Ignore any query params for the operations below.
			requestURI, _ := url.PathUnescape(strings.TrimSuffix(r.RequestURI, "?"+r.URL.RawQuery))

			for _, header := range f.c.serverConfig.MatchHeaders(requestURI) {
				w.Header().Set(header.Key, header.Value)
			}

			if redirect := f.c.serverConfig.MatchRedirect(requestURI); !redirect.IsZero() {
				// fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
				doRedirect := true
				// This matches Netlify's behaviour and is needed for SPA behaviour.
				// See https://docs.netlify.com/routing/redirects/rewrites-proxies/
				if !redirect.Force {
					path := filepath.Clean(strings.TrimPrefix(requestURI, u.Path))
					if root != "" {
						path = filepath.Join(root, path)
					}
					fs := f.c.publishDirServerFs

					fi, err := fs.Stat(path)

					if err == nil {
						if fi.IsDir() {
							// There will be overlapping directories, so we
							// need to check for a file.
							_, err = fs.Stat(filepath.Join(path, "index.html"))
							doRedirect = err != nil
						} else {
							doRedirect = false
						}
					}
				}

				if doRedirect {
					switch redirect.Status {
					case 404:
						w.WriteHeader(404)
						file, err := fs.Open(strings.TrimPrefix(redirect.To, u.Path))
						if err == nil {
							defer file.Close()
							io.Copy(w, file)
						} else {
							fmt.Fprintln(w, "<h1>Page Not Found</h1>")
						}
						return
					case 200:
						if r2 := f.rewriteRequest(r, strings.TrimPrefix(redirect.To, u.Path)); r2 != nil {
							requestURI = redirect.To
							r = r2
						}
						fallthrough
					default:
						w.Header().Set("Content-Type", "")
						http.Redirect(w, r, redirect.To, redirect.Status)
						return

					}
				}

			}

			if f.c.fastRenderMode && f.c.buildErr == nil {
				if strings.HasSuffix(requestURI, "/") || strings.HasSuffix(requestURI, "html") || strings.HasSuffix(requestURI, "htm") {
					if !f.c.visitedURLs.Contains(requestURI) {
						// If not already on stack, re-render that single page.
						if err := f.c.partialReRender(requestURI); err != nil {
							f.c.handleBuildErr(err, fmt.Sprintf("Failed to render %q", requestURI))
							if f.c.showErrorInBrowser {
								http.Redirect(w, r, requestURI, http.StatusMovedPermanently)
								return
							}
						}
					}

					f.c.visitedURLs.Add(requestURI)

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

	return mu, listener, u.String(), endpoint, nil
}

var (
	logErrorRe                    = regexp.MustCompile(`(?s)ERROR \d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} `)
	logDuplicateTemplateExecuteRe = regexp.MustCompile(`: template: .*?:\d+:\d+: executing ".*?"`)
	logDuplicateTemplateParseRe   = regexp.MustCompile(`: template: .*?:\d+:\d*`)
)

func removeErrorPrefixFromLog(content string) string {
	return logErrorRe.ReplaceAllLiteralString(content, "")
}

var logReplacer = strings.NewReplacer(
	"can't", "can’t", // Chroma lexer does'nt do well with "can't"
	"*hugolib.pageState", "page.Page", // Page is the public interface.
	"Rebuild failed:", "",
)

func cleanErrorLog(content string) string {
	content = strings.ReplaceAll(content, "\n", " ")
	content = logReplacer.Replace(content)
	content = logDuplicateTemplateExecuteRe.ReplaceAllString(content, "")
	content = logDuplicateTemplateParseRe.ReplaceAllString(content, "")
	seen := make(map[string]bool)
	parts := strings.Split(content, ": ")
	keep := make([]string, 0, len(parts))
	for _, part := range parts {
		if seen[part] {
			continue
		}
		seen[part] = true
		keep = append(keep, part)
	}
	return strings.Join(keep, ": ")
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

	// Cache it here. The HugoSites object may be unavaialble later on due to intermitent configuration errors.
	// To allow the en user to change the error template while the server is running, we use
	// the freshest template we can provide.
	var (
		errTempl     tpl.Template
		templHandler tpl.TemplateHandler
	)
	getErrorTemplateAndHandler := func(h *hugolib.HugoSites) (tpl.Template, tpl.TemplateHandler) {
		if h == nil {
			return errTempl, templHandler
		}
		templHandler := h.Tmpl()
		errTempl, found := templHandler.Lookup("_server/error.html")
		if !found {
			panic("template server/error.html not found")
		}
		return errTempl, templHandler
	}
	errTempl, templHandler = getErrorTemplateAndHandler(c.hugo())

	srv := &fileServer{
		baseURLs: baseURLs,
		roots:    roots,
		c:        c,
		s:        s,
		errorTemplate: func(ctx any) (io.Reader, error) {
			// hugoTry does not block, getErrorTemplateAndHandler will fall back
			// to cached values if nil.
			templ, handler := getErrorTemplateAndHandler(c.hugoTry())
			b := &bytes.Buffer{}
			err := handler.Execute(templ, b, ctx)
			return b, err
		},
	}

	doLiveReload := !c.Cfg.GetBool("disableLiveReload")

	if doLiveReload {
		livereload.Initialize()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	var servers []*http.Server

	wg1, ctx := errgroup.WithContext(context.Background())

	for i := range baseURLs {
		mu, listener, serverURL, endpoint, err := srv.createEndpoint(i)
		srv := &http.Server{
			Addr:    endpoint,
			Handler: mu,
		}
		servers = append(servers, srv)

		if doLiveReload {
			u, err := url.Parse(helpers.SanitizeURL(baseURLs[i]))
			if err != nil {
				return err
			}

			mu.HandleFunc(u.Path+"/livereload.js", livereload.ServeJS)
			mu.HandleFunc(u.Path+"/livereload", livereload.Handler)
		}
		jww.FEEDBACK.Printf("Web Server is available at %s (bind address %s)\n", serverURL, s.serverInterface)
		wg1.Go(func() error {
			err = srv.Serve(listener)
			if err != nil && err != http.ErrServerClosed {
				return err
			}
			return nil
		})
	}

	jww.FEEDBACK.Println("Press Ctrl+C to stop")

	err := func() error {
		if s.stop != nil {
			for {
				select {
				case <-sigs:
					return nil
				case <-s.stop:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		} else {
			for {
				select {
				case <-sigs:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}()

	if err != nil {
		jww.ERROR.Println("Error:", err)
	}

	if h := c.hugoTry(); h != nil {
		h.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wg2, ctx := errgroup.WithContext(ctx)
	for _, srv := range servers {
		srv := srv
		wg2.Go(func() error {
			return srv.Shutdown(ctx)
		})
	}

	err1, err2 := wg1.Wait(), wg2.Wait()
	if err1 != nil {
		return err1
	}
	return err2
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
				return "", fmt.Errorf("Failed to split baseURL hostpost: %w", err)
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

			start := htime.Now().UnixNano()

			for {
				runtime.ReadMemStats(&stats)
				if fileMemStats != nil {
					fileMemStats.WriteString(fmt.Sprintf("%d\t%d\t%d\t%d\t%d\n",
						(htime.Now().UnixNano()-start)/1000000, stats.HeapSys, stats.HeapAlloc, stats.HeapIdle, stats.HeapReleased))
					time.Sleep(interval)
				} else {
					break
				}
			}
		}()
	}
	return nil
}
