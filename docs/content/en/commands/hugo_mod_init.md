---
title: "hugo mod init"
slug: hugo_mod_init
url: /commands/hugo_mod_init/
---
## hugo mod init

Initialize this project as a Hugo Module.

### Synopsis

Initialize this project as a Hugo Module.
	It will try to guess the module path, but you may help by passing it as an argument, e.g:
	
		hugo mod init github.com/gohugoio/testshortcodes
	
	Note that Hugo Modules supports multi-module projects, so you can initialize a Hugo Module
	inside a subfolder on GitHub, as one example.
	

```
hugo mod init [flags] [args]
```

### Options

```
  -b, --baseURL string           hostname (and path) to the root, e.g. https://spf13.com/
      --cacheDir string          filesystem path to cache directory
  -c, --contentDir string        filesystem path to content directory
  -h, --help                     help for init
      --renderSegments strings   named segments to render (configured in the segments config)
  -t, --theme strings            themes to use (located in /themes/THEMENAME/)
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

* [hugo mod](/commands/hugo_mod/)	 - Various Hugo Modules helpers.

