---
title: "hugo new site"
slug: hugo_new_site
url: /commands/hugo_new_site/
---
## hugo new site

Create a new site (skeleton)

### Synopsis

Create a new site in the provided directory.
The new site will have the correct structure, but no content or theme yet.
Use `hugo new [contentPath]` to create new content.

```
hugo new site [path] [flags]
```

### Options

```
  -f, --force           init inside non-empty directory
      --format string   preferred file format (toml, yaml or json) (default "toml")
  -h, --help            help for site
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

* [hugo new](/commands/hugo_new/)	 - Create new content for your site

