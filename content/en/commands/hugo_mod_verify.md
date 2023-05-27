---
title: "hugo mod verify"
slug: hugo_mod_verify
url: /commands/hugo_mod_verify/
---
## hugo mod verify

Verify dependencies.

### Synopsis

Verify checks that the dependencies of the current module, which are stored in a local downloaded source cache, have not been modified since being downloaded.

```
hugo mod verify [flags] [args]
```

### Options

```
      --clean   delete module cache for dependencies that fail verification
  -h, --help    help for verify
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

