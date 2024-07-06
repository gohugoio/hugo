// Copyright 2024 The Hugo Authors. All rights reserved.
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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/bep/mclib"

	"github.com/bep/debounce"
	"github.com/bep/simplecobra"
	"github.com/fsnotify/fsnotify"
	"github.com/gohugoio/hugo/common/herrors"
	"github.com/gohugoio/hugo/common/hugo"

	"github.com/gohugoio/hugo/common/types"
	"github.com/gohugoio/hugo/common/urls"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/gohugoio/hugo/hugolib"
	"github.com/gohugoio/hugo/hugolib/filesystems"
	"github.com/gohugoio/hugo/livereload"
	"github.com/gohugoio/hugo/tpl"
	"github.com/gohugoio/hugo/transform"
	"github.com/gohugoio/hugo/transform/livereloadinject"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/fsync"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

var (
	logDuplicateTemplateExecuteRe = regexp.MustCompile(`: template: .*?:\d+:\d+: executing ".*?"`)
	logDuplicateTemplateParseRe   = regexp.MustCompile(`: template: .*?:\d+:\d*`)
)

var logReplacer = strings.NewReplacer(
	"can't", "can’t", // Chroma lexer doesn't do well with "can't"
	"*hugolib.pageState", "page.Page", // Page is the public interface.
	"Rebuild failed:", "",
)

const (
	configChangeConfig = "config file"
	configChangeGoMod  = "go.mod file"
	configChangeGoWork = "go work file"
)

func newHugoBuilder(r *rootCommand, s *serverCommand, onConfigLoaded ...func(reloaded bool) error) *hugoBuilder {
	var visitedURLs *types.EvictingStringQueue
	if s != nil && !s.disableFastRender {
		visitedURLs = types.NewEvictingStringQueue(20)
	}
	return &hugoBuilder{
		r:              r,
		s:              s,
		visitedURLs:    visitedURLs,
		fullRebuildSem: semaphore.NewWeighted(1),
		debounce:       debounce.New(4 * time.Second),
		onConfigLoaded: func(reloaded bool) error {
			for _, wc := range onConfigLoaded {
				if err := wc(reloaded); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func newServerCommand() *serverCommand {
	// Flags.
	var uninstall bool

	c := &serverCommand{
		quit: make(chan bool),
		commands: []simplecobra.Commander{
			&simpleCommand{
				name:  "trust",
				short: "Install the local CA in the system trust store.",
				run: func(ctx context.Context, cd *simplecobra.Commandeer, r *rootCommand, args []string) error {
					action := "-install"
					if uninstall {
						action = "-uninstall"
					}
					os.Args = []string{action}
					return mclib.RunMain()
				},
				withc: func(cmd *cobra.Command, r *rootCommand) {
					cmd.ValidArgsFunction = cobra.NoFileCompletions
					cmd.Flags().BoolVar(&uninstall, "uninstall", false, "Uninstall the local CA (but do not delete it).")
				},
			},
		},
	}

	return c
}

func (c *serverCommand) Commands() []simplecobra.Commander {
	return c.commands
}

type countingStatFs struct {
	afero.Fs
	statCounter uint64
}

func (fs *countingStatFs) Stat(name string) (os.FileInfo, error) {
	f, err := fs.Fs.Stat(name)
	if err == nil {
		if !f.IsDir() {
			atomic.AddUint64(&fs.statCounter, 1)
		}
	}
	return f, err
}

// dynamicEvents contains events that is considered dynamic, as in "not static".
// Both of these categories will trigger a new build, but the asset events
// does not fit into the "navigate to changed" logic.
type dynamicEvents struct {
	ContentEvents []fsnotify.Event
	AssetEvents   []fsnotify.Event
}

type fileChangeDetector struct {
	sync.Mutex
	current map[string]uint64
	prev    map[string]uint64

	irrelevantRe *regexp.Regexp
}

func (f *fileChangeDetector) OnFileClose(name string, checksum uint64) {
	f.Lock()
	defer f.Unlock()
	f.current[name] = checksum
}

func (f *fileChangeDetector) PrepareNew() {
	if f == nil {
		return
	}

	f.Lock()
	defer f.Unlock()

	if f.current == nil {
		f.current = make(map[string]uint64)
		f.prev = make(map[string]uint64)
		return
	}

	f.prev = make(map[string]uint64)
	for k, v := range f.current {
		f.prev[k] = v
	}
	f.current = make(map[string]uint64)
}

func (f *fileChangeDetector) changed() []string {
	if f == nil {
		return nil
	}
	f.Lock()
	defer f.Unlock()
	var c []string
	for k, v := range f.current {
		vv, found := f.prev[k]
		if !found || v != vv {
			c = append(c, k)
		}
	}

	return f.filterIrrelevant(c)
}

func (f *fileChangeDetector) filterIrrelevant(in []string) []string {
	var filtered []string
	for _, v := range in {
		if !f.irrelevantRe.MatchString(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

type fileServer struct {
	baseURLs      []urls.BaseURL
	roots         []string
	errorTemplate func(err any) (io.Reader, error)
	c             *serverCommand
}

func (f *fileServer) createEndpoint(i int) (*http.ServeMux, net.Listener, string, string, error) {
	r := f.c.r
	baseURL := f.baseURLs[i]
	root := f.roots[i]
	port := f.c.serverPorts[i].p
	listener := f.c.serverPorts[i].ln
	logger := f.c.r.logger

	if i == 0 {
		r.Printf("Environment: %q\n", f.c.hugoTry().Deps.Site.Hugo().Environment)
		mainTarget := "disk"
		if f.c.r.renderToMemory {
			mainTarget = "memory"
		}
		if f.c.renderStaticToDisk {
			r.Printf("Serving pages from %s and static files from disk\n", mainTarget)
		} else {
			r.Printf("Serving pages from %s\n", mainTarget)
		}
	}

	var httpFs *afero.HttpFs
	f.c.withConf(func(conf *commonConfig) {
		httpFs = afero.NewHttpFs(conf.fs.PublishDirServer)
	})

	fs := filesOnlyFs{httpFs.Dir(path.Join("/", root))}
	if i == 0 && f.c.fastRenderMode {
		r.Println("Running in Fast Render Mode. For full rebuilds on change: hugo server --disableFastRender")
	}

	decorate := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if f.c.showErrorInBrowser {
				// First check the error state
				err := f.c.getErrorWithContext()
				if err != nil {
					f.c.errState.setWasErr(true)
					w.WriteHeader(500)
					r, err := f.errorTemplate(err)
					if err != nil {
						logger.Errorln(err)
					}

					port = 1313
					f.c.withConf(func(conf *commonConfig) {
						if lrport := conf.configs.GetFirstLanguageConfig().BaseURLLiveReload().Port(); lrport != 0 {
							port = lrport
						}
					})
					lr := baseURL.URL()
					lr.Host = fmt.Sprintf("%s:%d", lr.Hostname(), port)
					fmt.Fprint(w, injectLiveReloadScript(r, lr))

					return
				}
			}

			if f.c.noHTTPCache {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
				w.Header().Set("Pragma", "no-cache")
			}

			var serverConfig config.Server
			f.c.withConf(func(conf *commonConfig) {
				serverConfig = conf.configs.Base.Server
			})

			// Ignore any query params for the operations below.
			requestURI, _ := url.PathUnescape(strings.TrimSuffix(r.RequestURI, "?"+r.URL.RawQuery))

			for _, header := range serverConfig.MatchHeaders(requestURI) {
				w.Header().Set(header.Key, header.Value)
			}

			if redirect := serverConfig.MatchRedirect(requestURI); !redirect.IsZero() {
				// fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
				doRedirect := true
				// This matches Netlify's behavior and is needed for SPA behavior.
				// See https://docs.netlify.com/routing/redirects/rewrites-proxies/
				if !redirect.Force {
					path := filepath.Clean(strings.TrimPrefix(requestURI, baseURL.Path()))
					if root != "" {
						path = filepath.Join(root, path)
					}
					var fs afero.Fs
					f.c.withConf(func(conf *commonConfig) {
						fs = conf.fs.PublishDirServer
					})

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
						file, err := fs.Open(strings.TrimPrefix(redirect.To, baseURL.Path()))
						if err == nil {
							defer file.Close()
							io.Copy(w, file)
						} else {
							fmt.Fprintln(w, "<h1>Page Not Found</h1>")
						}
						return
					case 200:
						if r2 := f.rewriteRequest(r, strings.TrimPrefix(redirect.To, baseURL.Path())); r2 != nil {
							requestURI = redirect.To
							r = r2
						}
					default:
						w.Header().Set("Content-Type", "")
						http.Redirect(w, r, redirect.To, redirect.Status)
						return

					}
				}

			}

			if f.c.fastRenderMode && f.c.errState.buildErr() == nil {
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
	if baseURL.Path() == "" || baseURL.Path() == "/" {
		mu.Handle("/", fileserver)
	} else {
		mu.Handle(baseURL.Path(), http.StripPrefix(baseURL.Path(), fileserver))
	}
	if r.IsTestRun() {
		var shutDownOnce sync.Once
		mu.HandleFunc("/__stop", func(w http.ResponseWriter, r *http.Request) {
			shutDownOnce.Do(func() {
				close(f.c.quit)
			})
		})
	}

	endpoint := net.JoinHostPort(f.c.serverInterface, strconv.Itoa(port))

	return mu, listener, baseURL.String(), endpoint, nil
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

type filesOnlyFs struct {
	fs http.FileSystem
}

func (fs filesOnlyFs) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return noDirFile{f}, nil
}

type noDirFile struct {
	http.File
}

func (f noDirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

type serverCommand struct {
	r *rootCommand

	commands []simplecobra.Commander

	*hugoBuilder

	quit         chan bool // Closed when the server should shut down. Used in tests only.
	serverPorts  []serverPortListener
	doLiveReload bool

	// Flags.
	renderStaticToDisk  bool
	navigateToChanged   bool
	serverAppend        bool
	serverInterface     string
	tlsCertFile         string
	tlsKeyFile          string
	tlsAuto             bool
	pprof               bool
	serverPort          int
	liveReloadPort      int
	serverWatch         bool
	noHTTPCache         bool
	disableLiveReload   bool
	disableFastRender   bool
	disableBrowserError bool
}

func (c *serverCommand) Name() string {
	return "server"
}

func (c *serverCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	if c.pprof {
		go func() {
			http.ListenAndServe("localhost:8080", nil)
		}()
	}
	// Watch runs its own server as part of the routine
	if c.serverWatch {

		watchDirs, err := c.getDirList()
		if err != nil {
			return err
		}

		watchGroups := helpers.ExtractAndGroupRootPaths(watchDirs)

		for _, group := range watchGroups {
			c.r.Printf("Watching for changes in %s\n", group)
		}
		watcher, err := c.newWatcher(c.r.poll, watchDirs...)
		if err != nil {
			return err
		}

		defer watcher.Close()

	}

	err := func() error {
		defer c.r.timeTrack(time.Now(), "Built")
		return c.build()
	}()
	if err != nil {
		return err
	}

	return c.serve()
}

func (c *serverCommand) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "A high performance webserver"
	cmd.Long = `Hugo provides its own webserver which builds and serves the site.
While hugo server is high performance, it is a webserver with limited options.

The ` + "`" + `hugo server` + "`" + ` command will by default write and serve files from disk, but
you can render to memory by using the ` + "`" + `--renderToMemory` + "`" + ` flag. This can be
faster in some cases, but it will consume more memory.

By default hugo will also watch your files for any changes you make and
automatically rebuild the site. It will then live reload any open browser pages
and push the latest content to them. As most Hugo sites are built in a fraction
of a second, you will be able to save and see your changes nearly instantly.`
	cmd.Aliases = []string{"serve"}

	cmd.Flags().IntVarP(&c.serverPort, "port", "p", 1313, "port on which the server will listen")
	_ = cmd.RegisterFlagCompletionFunc("port", cobra.NoFileCompletions)
	cmd.Flags().IntVar(&c.liveReloadPort, "liveReloadPort", -1, "port for live reloading (i.e. 443 in HTTPS proxy situations)")
	_ = cmd.RegisterFlagCompletionFunc("liveReloadPort", cobra.NoFileCompletions)
	cmd.Flags().StringVarP(&c.serverInterface, "bind", "", "127.0.0.1", "interface to which the server will bind")
	_ = cmd.RegisterFlagCompletionFunc("bind", cobra.NoFileCompletions)
	cmd.Flags().StringVarP(&c.tlsCertFile, "tlsCertFile", "", "", "path to TLS certificate file")
	_ = cmd.MarkFlagFilename("tlsCertFile", "pem")
	cmd.Flags().StringVarP(&c.tlsKeyFile, "tlsKeyFile", "", "", "path to TLS key file")
	_ = cmd.MarkFlagFilename("tlsKeyFile", "pem")
	cmd.Flags().BoolVar(&c.tlsAuto, "tlsAuto", false, "generate and use locally-trusted certificates.")
	cmd.Flags().BoolVar(&c.pprof, "pprof", false, "enable the pprof server (port 8080)")
	cmd.Flags().BoolVarP(&c.serverWatch, "watch", "w", true, "watch filesystem for changes and recreate as needed")
	cmd.Flags().BoolVar(&c.noHTTPCache, "noHTTPCache", false, "prevent HTTP caching")
	cmd.Flags().BoolVarP(&c.serverAppend, "appendPort", "", true, "append port to baseURL")
	cmd.Flags().BoolVar(&c.disableLiveReload, "disableLiveReload", false, "watch without enabling live browser reload on rebuild")
	cmd.Flags().BoolVarP(&c.navigateToChanged, "navigateToChanged", "N", false, "navigate to changed content file on live browser reload")
	cmd.Flags().BoolVar(&c.renderStaticToDisk, "renderStaticToDisk", false, "serve static files from disk and dynamic files from memory")
	cmd.Flags().BoolVar(&c.disableFastRender, "disableFastRender", false, "enables full re-renders on changes")
	cmd.Flags().BoolVar(&c.disableBrowserError, "disableBrowserError", false, "do not show build errors in the browser")

	r := cd.Root.Command.(*rootCommand)
	applyLocalFlagsBuild(cmd, r)

	return nil
}

func (c *serverCommand) PreRun(cd, runner *simplecobra.Commandeer) error {
	c.r = cd.Root.Command.(*rootCommand)

	c.hugoBuilder = newHugoBuilder(
		c.r,
		c,
		func(reloaded bool) error {
			if !reloaded {
				if err := c.createServerPorts(cd); err != nil {
					return err
				}

				if (c.tlsCertFile == "" || c.tlsKeyFile == "") && c.tlsAuto {
					c.withConfE(func(conf *commonConfig) error {
						return c.createCertificates(conf)
					})
				}
			}

			if err := c.setServerInfoInConfig(); err != nil {
				return err
			}

			if !reloaded && c.fastRenderMode {
				c.withConf(func(conf *commonConfig) {
					conf.fs.PublishDir = hugofs.NewHashingFs(conf.fs.PublishDir, c.changeDetector)
					conf.fs.PublishDirStatic = hugofs.NewHashingFs(conf.fs.PublishDirStatic, c.changeDetector)
				})
			}

			return nil
		},
	)

	destinationFlag := cd.CobraCommand.Flags().Lookup("destination")
	if c.r.renderToMemory && (destinationFlag != nil && destinationFlag.Changed) {
		return fmt.Errorf("cannot use --renderToMemory with --destination")
	}
	c.doLiveReload = !c.disableLiveReload
	c.fastRenderMode = !c.disableFastRender
	c.showErrorInBrowser = c.doLiveReload && !c.disableBrowserError

	if c.fastRenderMode {
		// For now, fast render mode only. It should, however, be fast enough
		// for the full variant, too.
		c.changeDetector = &fileChangeDetector{
			// We use this detector to decide to do a Hot reload of a single path or not.
			// We need to filter out source maps and possibly some other to be able
			// to make that decision.
			irrelevantRe: regexp.MustCompile(`\.map$`),
		}

		c.changeDetector.PrepareNew()

	}

	err := c.loadConfig(cd, true)
	if err != nil {
		return err
	}

	return nil
}

func (c *serverCommand) setServerInfoInConfig() error {
	if len(c.serverPorts) == 0 {
		panic("no server ports set")
	}
	return c.withConfE(func(conf *commonConfig) error {
		for i, language := range conf.configs.Languages {
			isMultihost := conf.configs.IsMultihost
			var serverPort int
			if isMultihost {
				serverPort = c.serverPorts[i].p
			} else {
				serverPort = c.serverPorts[0].p
			}
			langConfig := conf.configs.LanguageConfigMap[language.Lang]
			baseURLStr, err := c.fixURL(langConfig.BaseURL, c.r.baseURL, serverPort)
			if err != nil {
				return err
			}
			baseURL, err := urls.NewBaseURLFromString(baseURLStr)
			if err != nil {
				return fmt.Errorf("failed to create baseURL from %q: %s", baseURLStr, err)
			}

			baseURLLiveReload := baseURL
			if c.liveReloadPort != -1 {
				baseURLLiveReload, _ = baseURLLiveReload.WithPort(c.liveReloadPort)
			}
			langConfig.C.SetServerInfo(baseURL, baseURLLiveReload, c.serverInterface)

		}
		return nil
	})
}

func (c *serverCommand) getErrorWithContext() any {
	errCount := c.errCount()

	if errCount == 0 {
		return nil
	}

	m := make(map[string]any)

	m["Error"] = cleanErrorLog(c.r.logger.Errors())

	m["Version"] = hugo.BuildVersionString()
	ferrors := herrors.UnwrapFileErrorsWithErrorContext(c.errState.buildErr())
	m["Files"] = ferrors

	return m
}

func (c *serverCommand) createCertificates(conf *commonConfig) error {
	hostname := "localhost"
	if c.r.baseURL != "" {
		u, err := url.Parse(c.r.baseURL)
		if err != nil {
			return err
		}
		hostname = u.Hostname()
	}

	// For now, store these in the Hugo cache dir.
	// Hugo should probably introduce some concept of a less temporary application directory.
	keyDir := filepath.Join(conf.configs.LoadingInfo.BaseConfig.CacheDir, "_mkcerts")

	// Create the directory if it doesn't exist.
	if _, err := os.Stat(keyDir); os.IsNotExist(err) {
		if err := os.MkdirAll(keyDir, 0o777); err != nil {
			return err
		}
	}

	c.tlsCertFile = filepath.Join(keyDir, fmt.Sprintf("%s.pem", hostname))
	c.tlsKeyFile = filepath.Join(keyDir, fmt.Sprintf("%s-key.pem", hostname))

	// Check if the certificate already exists and is valid.
	certPEM, err := os.ReadFile(c.tlsCertFile)
	if err == nil {
		rootPem, err := os.ReadFile(filepath.Join(mclib.GetCAROOT(), "rootCA.pem"))
		if err == nil {
			if err := c.verifyCert(rootPem, certPEM, hostname); err == nil {
				c.r.Println("Using existing", c.tlsCertFile, "and", c.tlsKeyFile)
				return nil
			}
		}
	}

	c.r.Println("Creating TLS certificates in", keyDir)

	// Yes, this is unfortunate, but it's currently the only way to use Mkcert as a library.
	os.Args = []string{"-cert-file", c.tlsCertFile, "-key-file", c.tlsKeyFile, hostname}
	return mclib.RunMain()
}

func (c *serverCommand) verifyCert(rootPEM, certPEM []byte, name string) error {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootPEM)
	if !ok {
		return fmt.Errorf("failed to parse root certificate")
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return fmt.Errorf("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %v", err.Error())
	}

	opts := x509.VerifyOptions{
		DNSName: name,
		Roots:   roots,
	}

	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("failed to verify certificate: %v", err.Error())
	}

	return nil
}

func (c *serverCommand) createServerPorts(cd *simplecobra.Commandeer) error {
	flags := cd.CobraCommand.Flags()
	var cerr error
	c.withConf(func(conf *commonConfig) {
		isMultihost := conf.configs.IsMultihost
		c.serverPorts = make([]serverPortListener, 1)
		if isMultihost {
			if !c.serverAppend {
				cerr = errors.New("--appendPort=false not supported when in multihost mode")
				return
			}
			c.serverPorts = make([]serverPortListener, len(conf.configs.Languages))
		}
		currentServerPort := c.serverPort
		for i := 0; i < len(c.serverPorts); i++ {
			l, err := net.Listen("tcp", net.JoinHostPort(c.serverInterface, strconv.Itoa(currentServerPort)))
			if err == nil {
				c.serverPorts[i] = serverPortListener{ln: l, p: currentServerPort}
			} else {
				if i == 0 && flags.Changed("port") {
					// port set explicitly by user -- he/she probably meant it!
					cerr = fmt.Errorf("server startup failed: %s", err)
					return
				}
				c.r.Println("port", currentServerPort, "already in use, attempting to use an available port")
				l, sp, err := helpers.TCPListen()
				if err != nil {
					cerr = fmt.Errorf("unable to find alternative port to use: %s", err)
					return
				}
				c.serverPorts[i] = serverPortListener{ln: l, p: sp.Port}
			}

			currentServerPort = c.serverPorts[i].p + 1
		}
	})

	return cerr
}

// fixURL massages the baseURL into a form needed for serving
// all pages correctly.
func (c *serverCommand) fixURL(baseURLFromConfig, baseURLFromFlag string, port int) (string, error) {
	certsSet := (c.tlsCertFile != "" && c.tlsKeyFile != "") || c.tlsAuto
	useLocalhost := false
	baseURL := baseURLFromFlag
	if baseURL == "" {
		baseURL = baseURLFromConfig
		useLocalhost = true
	}

	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}

	// do an initial parse of the input string
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// if no Host is defined, then assume that no schema or double-slash were
	// present in the url.  Add a double-slash and make a best effort attempt.
	if u.Host == "" && baseURL != "/" {
		baseURL = "//" + baseURL

		u, err = url.Parse(baseURL)
		if err != nil {
			return "", err
		}
	}

	if useLocalhost {
		if certsSet {
			u.Scheme = "https"
		} else if u.Scheme == "https" {
			u.Scheme = "http"
		}
		u.Host = "localhost"
	}

	if c.serverAppend {
		if strings.Contains(u.Host, ":") {
			u.Host, _, err = net.SplitHostPort(u.Host)
			if err != nil {
				return "", fmt.Errorf("failed to split baseURL hostport: %w", err)
			}
		}
		u.Host += fmt.Sprintf(":%d", port)
	}

	return u.String(), nil
}

func (c *serverCommand) partialReRender(urls ...string) error {
	defer func() {
		c.errState.setWasErr(false)
	}()
	c.errState.setBuildErr(nil)
	visited := types.NewEvictingStringQueue(len(urls))
	for _, url := range urls {
		visited.Add(url)
	}

	h, err := c.hugo()
	if err != nil {
		return err
	}
	// Note: We do not set NoBuildLock as the file lock is not acquired at this stage.
	return h.Build(hugolib.BuildCfg{NoBuildLock: false, RecentlyVisited: visited, PartialReRender: true, ErrRecovery: c.errState.wasErr()})
}

func (c *serverCommand) serve() error {
	var (
		baseURLs []urls.BaseURL
		roots    []string
		h        *hugolib.HugoSites
	)
	err := c.withConfE(func(conf *commonConfig) error {
		isMultihost := conf.configs.IsMultihost
		var err error
		h, err = c.r.HugFromConfig(conf)
		if err != nil {
			return err
		}

		// We need the server to share the same logger as the Hugo build (for error counts etc.)
		c.r.logger = h.Log

		if isMultihost {
			for _, l := range conf.configs.ConfigLangs() {
				baseURLs = append(baseURLs, l.BaseURL())
				roots = append(roots, l.Language().Lang)
			}
		} else {
			l := conf.configs.GetFirstLanguageConfig()
			baseURLs = []urls.BaseURL{l.BaseURL()}
			roots = []string{""}
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Cache it here. The HugoSites object may be unavailable later on due to intermittent configuration errors.
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
	errTempl, templHandler = getErrorTemplateAndHandler(h)

	srv := &fileServer{
		baseURLs: baseURLs,
		roots:    roots,
		c:        c,
		errorTemplate: func(ctx any) (io.Reader, error) {
			// hugoTry does not block, getErrorTemplateAndHandler will fall back
			// to cached values if nil.
			templ, handler := getErrorTemplateAndHandler(c.hugoTry())
			b := &bytes.Buffer{}
			err := handler.ExecuteWithContext(context.Background(), templ, b, ctx)
			return b, err
		},
	}

	doLiveReload := !c.disableLiveReload

	if doLiveReload {
		livereload.Initialize()
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	var servers []*http.Server

	wg1, ctx := errgroup.WithContext(context.Background())

	for i := range baseURLs {
		mu, listener, serverURL, endpoint, err := srv.createEndpoint(i)
		var srv *http.Server
		if c.tlsCertFile != "" && c.tlsKeyFile != "" {
			srv = &http.Server{
				Addr:    endpoint,
				Handler: mu,
				TLSConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
				},
			}
		} else {
			srv = &http.Server{
				Addr:    endpoint,
				Handler: mu,
			}
		}

		servers = append(servers, srv)

		if doLiveReload {
			baseURL := baseURLs[i]
			mu.HandleFunc(baseURL.Path()+"livereload.js", livereload.ServeJS)
			mu.HandleFunc(baseURL.Path()+"livereload", livereload.Handler)
		}
		c.r.Printf("Web Server is available at %s (bind address %s) %s\n", serverURL, c.serverInterface, roots[i])
		wg1.Go(func() error {
			if c.tlsCertFile != "" && c.tlsKeyFile != "" {
				err = srv.ServeTLS(listener, c.tlsCertFile, c.tlsKeyFile)
			} else {
				err = srv.Serve(listener)
			}
			if err != nil && err != http.ErrServerClosed {
				return err
			}
			return nil
		})
	}

	if c.r.IsTestRun() {
		// Write a .ready file to disk to signal ready status.
		// This is where the test is run from.
		var baseURLs []string
		for _, baseURL := range srv.baseURLs {
			baseURLs = append(baseURLs, baseURL.String())
		}
		testInfo := map[string]any{
			"baseURLs": baseURLs,
		}

		dir := os.Getenv("WORK")
		if dir != "" {
			readyFile := filepath.Join(dir, ".ready")
			// encode the test info as JSON into the .ready file.
			b, err := json.Marshal(testInfo)
			if err != nil {
				return err
			}
			err = os.WriteFile(readyFile, b, 0o777)
			if err != nil {
				return err
			}
		}

	}

	c.r.Println("Press Ctrl+C to stop")

	err = func() error {
		for {
			select {
			case <-c.quit:
				return nil
			case <-sigs:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}()
	if err != nil {
		c.r.Println("Error:", err)
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

type serverPortListener struct {
	p  int
	ln net.Listener
}

type staticSyncer struct {
	c *hugoBuilder
}

func (s *staticSyncer) isStatic(h *hugolib.HugoSites, filename string) bool {
	return h.BaseFs.SourceFilesystems.IsStatic(filename)
}

func (s *staticSyncer) syncsStaticEvents(staticEvents []fsnotify.Event) error {
	c := s.c

	syncFn := func(sourceFs *filesystems.SourceFilesystem) (uint64, error) {
		publishDir := helpers.FilePathSeparator

		if sourceFs.PublishFolder != "" {
			publishDir = filepath.Join(publishDir, sourceFs.PublishFolder)
		}

		syncer := fsync.NewSyncer()
		c.withConf(func(conf *commonConfig) {
			syncer.NoTimes = conf.configs.Base.NoTimes
			syncer.NoChmod = conf.configs.Base.NoChmod
			syncer.ChmodFilter = chmodFilter
			syncer.SrcFs = sourceFs.Fs
			syncer.DestFs = conf.fs.PublishDir
			if c.s != nil && c.s.renderStaticToDisk {
				syncer.DestFs = conf.fs.PublishDirStatic
			}
		})

		logger := s.c.r.logger

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

			relPath, found := sourceFs.MakePathRelative(fromPath, true)

			if !found {
				// Not member of this virtual host.
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
				if _, err := sourceFs.Fs.Stat(relPath); herrors.IsNotExist(err) {
					// If file doesn't exist in any static dir, remove it
					logger.Println("File no longer exists in static dir, removing", relPath)
					c.withConf(func(conf *commonConfig) {
						_ = conf.fs.PublishDirStatic.RemoveAll(relPath)
					})

				} else if err == nil {
					// If file still exists, sync it
					logger.Println("Syncing", relPath, "to", publishDir)

					if err := syncer.Sync(relPath, relPath); err != nil {
						c.r.logger.Errorln(err)
					}
				} else {
					c.r.logger.Errorln(err)
				}

				continue
			}

			// For all other event operations Hugo will sync static.
			logger.Println("Syncing", relPath, "to", publishDir)
			if err := syncer.Sync(filepath.Join(publishDir, relPath), relPath); err != nil {
				c.r.logger.Errorln(err)
			}
		}

		return 0, nil
	}

	_, err := c.doWithPublishDirs(syncFn)
	return err
}

func chmodFilter(dst, src os.FileInfo) bool {
	// Hugo publishes data from multiple sources, potentially
	// with overlapping directory structures. We cannot sync permissions
	// for directories as that would mean that we might end up with write-protected
	// directories inside /public.
	// One example of this would be syncing from the Go Module cache,
	// which have 0555 directories.
	return src.IsDir()
}

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

func injectLiveReloadScript(src io.Reader, baseURL *url.URL) string {
	var b bytes.Buffer
	chain := transform.Chain{livereloadinject.New(baseURL)}
	chain.Apply(&b, src)

	return b.String()
}

func partitionDynamicEvents(sourceFs *filesystems.SourceFilesystems, events []fsnotify.Event) (de dynamicEvents) {
	for _, e := range events {
		if !sourceFs.IsContent(e.Name) {
			de.AssetEvents = append(de.AssetEvents, e)
		} else {
			de.ContentEvents = append(de.ContentEvents, e)
		}
	}
	return
}

func pickOneWriteOrCreatePath(contentTypes config.ContentTypesProvider, events []fsnotify.Event) string {
	name := ""

	for _, ev := range events {
		if ev.Op&fsnotify.Write == fsnotify.Write || ev.Op&fsnotify.Create == fsnotify.Create {
			if contentTypes.IsIndexContentFile(ev.Name) {
				return ev.Name
			}

			if contentTypes.IsContentFile(ev.Name) {
				name = ev.Name
			}

		}
	}

	return name
}

func formatByteCount(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
