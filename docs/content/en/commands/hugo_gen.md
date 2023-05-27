---
title: "hugo gen"
slug: hugo_gen
url: /commands/hugo_gen/
---
## hugo gen

A collection of several useful generators.

```
hugo gen [command] [flags]
```

### Options

```
  -h, --help   help for gen
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

* [hugo](/commands/hugo/)	 - hugo builds your site
* [hugo gen chromastyles](/commands/hugo_gen_chromastyles/)	 - Generate CSS stylesheet for the Chroma code highlighter
* [hugo gen doc](/commands/hugo_gen_doc/)	 - Generate Markdown documentation for the Hugo CLI.
* [hugo gen man](/commands/hugo_gen_man/)	 - Generate man pages for the Hugo CLI

