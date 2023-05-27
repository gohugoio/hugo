---
title: "hugo version"
slug: hugo_version
url: /commands/hugo_version/
---
## hugo version

Print Hugo version and environment info

### Synopsis

Print Hugo version and environment info. This is useful in Hugo bug reports.

```
hugo version [flags] [args]
```

### Options

```
  -h, --help   help for version
```

### Options inherited from parent commands

```
      --clock string               set the clock used by Hugo, e.g. --clock 2021-11-06T22:30:00.00+09:00
      --config string              config file (default is hugo.yaml|json|toml)
      --configDir string           config dir (default "config")
      --debug                      debug output
  -d, --destination string         filesystem path to write files to
  -e, --environment string         build environment
      --format string              preferred file format (toml, yaml or json) (default "toml")
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

* [hugo](/commands/hugo/)	 - hugo builds your site

