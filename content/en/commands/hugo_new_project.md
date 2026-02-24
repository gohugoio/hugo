---
title: "hugo new project"
slug: hugo_new_project
url: /commands/hugo_new_project/
---
## hugo new project

Create a new project

### Synopsis

Create a new project at the specified path.

```
hugo new project [path] [flags]
```

### Options

```
  -f, --force           init inside non-empty directory
      --format string   preferred file format (toml, yaml or json) (default "toml")
  -h, --help            help for project
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

* [hugo new](/commands/hugo_new/)	 - Create new content
