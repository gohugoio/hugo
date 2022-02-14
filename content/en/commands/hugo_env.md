---
title: "hugo env"
slug: hugo_env
url: /commands/hugo_env/
---
## hugo env

Print Hugo version and environment info

### Synopsis

Print Hugo version and environment info. This is useful in Hugo bug reports.

If you add the -v flag, you will get a full dependency list.


```
hugo env [flags]
```

### Options

```
  -h, --help   help for env
```

### Options inherited from parent commands

```
      --config string              config file (default is path/config.yaml|json|toml)
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

* [hugo](/commands/hugo/)	 - hugo builds your site

