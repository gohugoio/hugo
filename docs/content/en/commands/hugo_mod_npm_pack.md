---
title: "hugo mod npm pack"
slug: hugo_mod_npm_pack
url: /commands/hugo_mod_npm_pack/
---
## hugo mod npm pack

Merges module Node.js dependencies into an npm workspace

### Synopsis

Merges Node.js dependencies from all Hugo modules into a "packages/hugoautogen" npm workspace.

The merged dependencies are written to packages/hugoautogen/package.json, and the root package.json
is updated with a "workspaces" entry pointing to "packages/hugoautogen".

The source entries are read from either package.hugo.json or package.json in the module root, with package.hugo.json taking precedence if both exist.

See [Node.js dependencies](/hugo-modules/nodejs-dependencies/) for more information.


```
hugo mod npm pack [flags] [args]
```

### Options

```
  -b, --baseURL string           hostname (and path) to the root, e.g. https://spf13.com/
      --cacheDir string          filesystem path to cache directory
  -c, --contentDir string        filesystem path to content directory
  -h, --help                     help for pack
      --renderSegments strings   named segments to render (configured in the segments config)
  -t, --theme strings            themes to use (located in /themes/THEMENAME/)
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

* [hugo mod npm](/commands/hugo_mod_npm/)	 - Various npm helpers
