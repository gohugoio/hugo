---
title: "hugo import"
slug: hugo_import
url: /commands/hugo_import/
---
## hugo import

Import a site from another system

### Synopsis

Import a site from another system.

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
  -d, --destination string         filesystem path to write files to
  -e, --environment string         build environment
      --ignoreVendorPaths string   ignores any _vendor for module paths matching the given Glob pattern
      --logLevel string            log level (debug|info|warn|error)
      --noBuildLock                don't create .hugo_build.lock file
      --quiet                      build in quiet mode
  -M, --renderToMemory             render to memory (mostly useful when running the server)
  -s, --source string              filesystem path to read files relative from
      --themesDir string           filesystem path to themes directory
```

### SEE ALSO

* [hugo](/commands/hugo/)	 - Build your site
* [hugo import jekyll](/commands/hugo_import_jekyll/)	 - hugo import from Jekyll

