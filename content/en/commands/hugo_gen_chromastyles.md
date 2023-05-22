---
title: "hugo gen chromastyles"
slug: hugo_gen_chromastyles
url: /commands/hugo_gen_chromastyles/
---
## hugo gen chromastyles

Generate CSS stylesheet for the Chroma code highlighter

### Synopsis

Generate CSS stylesheet for the Chroma code highlighter for a given style. This stylesheet is needed if markup.highlight.noClasses is disabled in config.

See https://xyproto.github.io/splash/docs/all.html for a preview of the available styles

```
hugo gen chromastyles [flags] [args]
```

### Options

```
  -h, --help                    help for chromastyles
      --highlightStyle string   style used for highlighting lines (see https://github.com/alecthomas/chroma) (default "bg:#ffffcc")
      --linesStyle string       style used for line numbers (see https://github.com/alecthomas/chroma)
      --style string            highlighter style (see https://xyproto.github.io/splash/docs/) (default "friendly")
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

* [hugo gen](/commands/hugo_gen/)	 - A collection of several useful generators.

