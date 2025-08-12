---
title: "hugo gen jsonschemas"
slug: hugo_gen_jsonschemas
url: /commands/hugo_gen_jsonschemas/
---
## hugo gen jsonschemas

Generate JSON Schema for Hugo config and page structures

### Synopsis

Generate a JSON Schema for Hugo configuration options and page structures using reflection.

```
hugo gen jsonschemas [flags] [args]
```

### Options

```
      --dir string   output directory for schema files (default "/tmp/hugo-schemas")
  -h, --help         help for jsonschemas
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

* [hugo gen](/commands/hugo_gen/)	 - Generate documentation and syntax highlighting styles

