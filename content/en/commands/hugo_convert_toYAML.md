---
title: "hugo convert toYAML"
slug: hugo_convert_toYAML
url: /commands/hugo_convert_toyaml/
---
## hugo convert toYAML

Convert front matter to YAML

### Synopsis

toYAML converts all front matter in the content directory
to use YAML for the front matter.

```
hugo convert toYAML [flags]
```

### Options

```
  -h, --help   help for toYAML
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
  -o, --output string              filesystem path to write files to
      --quiet                      build in quiet mode
  -s, --source string              filesystem path to read files relative from
      --themesDir string           filesystem path to themes directory
      --unsafe                     enable less safe operations, please backup first
  -v, --verbose                    verbose output
      --verboseLog                 verbose logging
```

### SEE ALSO

* [hugo convert](/commands/hugo_convert/)	 - Convert your content to different formats

