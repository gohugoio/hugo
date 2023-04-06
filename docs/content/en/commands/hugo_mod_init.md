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
hugo mod init [flags]
```

### Options

```
  -h, --help   help for init
```

### Options inherited from parent commands

```
      --clock string               set the clock used by Hugo, e.g. --clock 2021-11-06T22:30:00.00+09:00
      --config string              config file (default is hugo.yaml|json|toml)
      --configDir string           config dir (default "config")
      --debug                      debug output
  -e, --environment string         build environment
      --ignoreVendorPaths string   ignores any _vendor for module paths matching the given Glob pattern
      --log                        enable Logging
      --logFile string             log File path (if set, logging enabled automatically)
      --quiet                      build in quiet mode
  -s, --source string              filesystem path to read files relative from
      --themesDir string           filesystem path to themes directory
  -v, --verbose                    verbose output
      --verboseLog                 verbose logging
```

### SEE ALSO

* [hugo mod](/commands/hugo_mod/)	 - Various Hugo Modules helpers.

