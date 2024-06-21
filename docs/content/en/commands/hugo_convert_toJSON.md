---
title: "hugo convert toJSON"
slug: hugo_convert_toJSON
url: /commands/hugo_convert_tojson/
---
## hugo convert toJSON

Convert front matter to JSON

### Synopsis

toJSON converts all front matter in the content directory
to use JSON for the front matter.

```
hugo convert toJSON [flags] [args]
```

### Options

```
  -h, --help   help for toJSON
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
  -o, --output string              filesystem path to write files to
      --quiet                      build in quiet mode
  -M, --renderToMemory             render to memory (mostly useful when running the server)
  -s, --source string              filesystem path to read files relative from
      --themesDir string           filesystem path to themes directory
      --unsafe                     enable less safe operations, please backup first
  -v, --verbose                    verbose output
```

### SEE ALSO

* [hugo convert](/commands/hugo_convert/)	 - Convert your content to different formats

