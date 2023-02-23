---
title: "hugo import"
slug: hugo_import
url: /commands/hugo_import/
---
## hugo import

Import your site from others.

### Synopsis

Import your site from other web site generators like Jekyll.

Import requires a subcommand, e.g. `hugo import jekyll jekyll_root_path target_path`.

### Options

```
  -h, --help   help for import
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

* [hugo](/commands/hugo/)	 - hugo builds your site
* [hugo import jekyll](/commands/hugo_import_jekyll/)	 - hugo import from Jekyll

