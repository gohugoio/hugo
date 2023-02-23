---
title: "hugo mod"
slug: hugo_mod
url: /commands/hugo_mod/
---
## hugo mod

Various Hugo Modules helpers.

### Synopsis

Various helpers to help manage the modules in your project's dependency graph.

Most operations here requires a Go version installed on your system (>= Go 1.12) and the relevant VCS client (typically Git).
This is not needed if you only operate on modules inside /themes or if you have vendored them via "hugo mod vendor".


Note that Hugo will always start out by resolving the components defined in the site
configuration, provided by a _vendor directory (if no --ignoreVendorPaths flag provided),
Go Modules, or a folder inside the themes directory, in that order.

See https://gohugo.io/hugo-modules/ for more information.



### Options

```
  -b, --baseURL string             hostname (and path) to the root, e.g. https://spf13.com/
  -D, --buildDrafts                include content marked as draft
  -E, --buildExpired               include expired content
  -F, --buildFuture                include content with publishdate in the future
      --cacheDir string            filesystem path to cache directory. Defaults: $TMPDIR/hugo_cache/
      --cleanDestinationDir        remove files from destination not found in static directories
      --clock string               set the clock used by Hugo, e.g. --clock 2021-11-06T22:30:00.00+09:00
  -c, --contentDir string          filesystem path to content directory
  -d, --destination string         filesystem path to write files to
      --disableKinds strings       disable different kind of pages (home, RSS etc.)
      --enableGitInfo              add Git revision, date, author, and CODEOWNERS info to the pages
  -e, --environment string         build environment
      --forceSyncStatic            copy all files when static is changed.
      --gc                         enable to run some cleanup tasks (remove unused cache files) after the build
  -h, --help                       help for mod
      --ignoreCache                ignores the cache directory
      --ignoreVendorPaths string   ignores any _vendor for module paths matching the given Glob pattern
  -l, --layoutDir string           filesystem path to layout directory
      --minify                     minify any supported output format (HTML, XML etc.)
      --noBuildLock                don't create .hugo_build.lock file
      --noChmod                    don't sync permission mode of files
      --noTimes                    don't sync modification time of files
      --panicOnWarning             panic on first WARNING log
      --poll string                set this to a poll interval, e.g --poll 700ms, to use a poll based approach to watch for file system changes
      --printI18nWarnings          print missing translations
      --printMemoryUsage           print memory usage to screen at intervals
      --printPathWarnings          print warnings on duplicate target paths etc.
      --printUnusedTemplates       print warnings on unused templates.
  -s, --source string              filesystem path to read files relative from
      --templateMetrics            display metrics about template executions
      --templateMetricsHints       calculate some improvement hints when combined with --templateMetrics
  -t, --theme strings              themes to use (located in /themes/THEMENAME/)
      --themesDir string           filesystem path to themes directory
      --trace file                 write trace to file (not useful in general)
```

### Options inherited from parent commands

```
      --config string      config file (default is hugo.yaml|json|toml)
      --configDir string   config dir (default "config")
      --debug              debug output
      --log                enable Logging
      --logFile string     log File path (if set, logging enabled automatically)
      --quiet              build in quiet mode
  -v, --verbose            verbose output
      --verboseLog         verbose logging
```

### SEE ALSO

* [hugo](/commands/hugo/)	 - hugo builds your site
* [hugo mod clean](/commands/hugo_mod_clean/)	 - Delete the Hugo Module cache for the current project.
* [hugo mod get](/commands/hugo_mod_get/)	 - Resolves dependencies in your current Hugo Project.
* [hugo mod graph](/commands/hugo_mod_graph/)	 - Print a module dependency graph.
* [hugo mod init](/commands/hugo_mod_init/)	 - Initialize this project as a Hugo Module.
* [hugo mod npm](/commands/hugo_mod_npm/)	 - Various npm helpers.
* [hugo mod tidy](/commands/hugo_mod_tidy/)	 - Remove unused entries in go.mod and go.sum.
* [hugo mod vendor](/commands/hugo_mod_vendor/)	 - Vendor all module dependencies into the _vendor directory.
* [hugo mod verify](/commands/hugo_mod_verify/)	 - Verify dependencies.

