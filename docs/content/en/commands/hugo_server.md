---
title: "hugo server"
slug: hugo_server
url: /commands/hugo_server/
---
## hugo server

A high performance webserver

### Synopsis

Hugo provides its own webserver which builds and serves the site.
While hugo server is high performance, it is a webserver with limited options.

The `hugo server` command will by default write and serve files from disk, but
you can render to memory by using the `--renderToMemory` flag. This can be
faster in some cases, but it will consume more memory.

By default hugo will also watch your files for any changes you make and
automatically rebuild the site. It will then live reload any open browser pages
and push the latest content to them. As most Hugo sites are built in a fraction
of a second, you will be able to save and see your changes nearly instantly.

```
hugo server [command] [flags]
```

### Options

```
      --appendPort               append port to baseURL (default true)
  -b, --baseURL string           hostname (and path) to the root, e.g. https://spf13.com/
      --bind string              interface to which the server will bind (default "127.0.0.1")
  -D, --buildDrafts              include content marked as draft
  -E, --buildExpired             include expired content
  -F, --buildFuture              include content with publishdate in the future
      --cacheDir string          filesystem path to cache directory
      --cleanDestinationDir      remove files from destination not found in static directories
  -c, --contentDir string        filesystem path to content directory
      --disableBrowserError      do not show build errors in the browser
      --disableFastRender        enables full re-renders on changes
      --disableKinds strings     disable different kind of pages (home, RSS etc.)
      --disableLiveReload        watch without enabling live browser reload on rebuild
      --enableGitInfo            add Git revision, date, author, and CODEOWNERS info to the pages
      --forceSyncStatic          copy all files when static is changed.
      --gc                       enable to run some cleanup tasks (remove unused cache files) after the build
  -h, --help                     help for server
      --ignoreCache              ignores the cache directory
  -l, --layoutDir string         filesystem path to layout directory
      --liveReloadPort int       port for live reloading (i.e. 443 in HTTPS proxy situations) (default -1)
      --minify                   minify any supported output format (HTML, XML etc.)
  -N, --navigateToChanged        navigate to changed content file on live browser reload
      --noBuildLock              don't create .hugo_build.lock file
      --noChmod                  don't sync permission mode of files
      --noHTTPCache              prevent HTTP caching
      --noTimes                  don't sync modification time of files
      --panicOnWarning           panic on first WARNING log
      --poll string              set this to a poll interval, e.g --poll 700ms, to use a poll based approach to watch for file system changes
  -p, --port int                 port on which the server will listen (default 1313)
      --pprof                    enable the pprof server (port 8080)
      --printI18nWarnings        print missing translations
      --printMemoryUsage         print memory usage to screen at intervals
      --printPathWarnings        print warnings on duplicate target paths etc.
      --printUnusedTemplates     print warnings on unused templates.
      --renderSegments strings   named segments to render (configured in the segments config)
      --renderStaticToDisk       serve static files from disk and dynamic files from memory
      --templateMetrics          display metrics about template executions
      --templateMetricsHints     calculate some improvement hints when combined with --templateMetrics
  -t, --theme strings            themes to use (located in /themes/THEMENAME/)
      --tlsAuto                  generate and use locally-trusted certificates.
      --tlsCertFile string       path to TLS certificate file
      --tlsKeyFile string        path to TLS key file
      --trace file               write trace to file (not useful in general)
  -w, --watch                    watch filesystem for changes and recreate as needed (default true)
```

### Options inherited from parent commands

```
      --clock string               set the clock used by Hugo, e.g. --clock 2021-11-06T22:30:00.00+09:00
      --config string              config file (default is hugo.yaml|json|toml)
      --configDir string           config dir (default "config")
      --debug                      debug output
  -d, --destination string         filesystem path to write files to
  -e, --environment string         build environment
      --ignoreVendorPaths string   ignores any _vendor for module paths matching the given Glob pattern
      --logLevel string            log level (debug|info|warn|error)
      --quiet                      build in quiet mode
  -M, --renderToMemory             render to memory (mostly useful when running the server)
  -s, --source string              filesystem path to read files relative from
      --themesDir string           filesystem path to themes directory
  -v, --verbose                    verbose output
```

### SEE ALSO

* [hugo](/commands/hugo/)	 - hugo builds your site
* [hugo server trust](/commands/hugo_server_trust/)	 - Install the local CA in the system trust store.

