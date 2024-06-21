---
title: "hugo convert"
slug: hugo_convert
url: /commands/hugo_convert/
---
## hugo convert

Convert your content to different formats

### Synopsis

Convert your content (e.g. front matter) to different formats.

See convert's subcommands toJSON, toTOML and toYAML for more information.

### Options

```
  -h, --help            help for convert
  -o, --output string   filesystem path to write files to
      --unsafe          enable less safe operations, please backup first
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
      --quiet                      build in quiet mode
  -M, --renderToMemory             render to memory (mostly useful when running the server)
  -s, --source string              filesystem path to read files relative from
      --themesDir string           filesystem path to themes directory
  -v, --verbose                    verbose output
```

### SEE ALSO

* [hugo](/commands/hugo/)	 - hugo builds your site
* [hugo convert toJSON](/commands/hugo_convert_tojson/)	 - Convert front matter to JSON
* [hugo convert toTOML](/commands/hugo_convert_totoml/)	 - Convert front matter to TOML
* [hugo convert toYAML](/commands/hugo_convert_toyaml/)	 - Convert front matter to YAML

