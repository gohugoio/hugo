---
title: "hugo mod npm pack"
slug: hugo_mod_npm_pack
url: /commands/hugo_mod_npm_pack/
---
## hugo mod npm pack

Experimental: Prepares and writes a composite package.json file for your project.

### Synopsis

Prepares and writes a composite package.json file for your project.

On first run it creates a "package.hugo.json" in the project root if not already there. This file will be used as a template file
with the base dependency set. 

This set will be merged with all "package.hugo.json" files found in the dependency tree, picking the version closest to the project.

This command is marked as 'Experimental'. We think it's a great idea, so it's not likely to be
removed from Hugo, but we need to test this out in "real life" to get a feel of it,
so this may/will change in future versions of Hugo.


```
hugo mod npm pack [flags] [args]
```

### Options

```
  -h, --help   help for pack
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

* [hugo mod npm](/commands/hugo_mod_npm/)	 - Various npm helpers.

