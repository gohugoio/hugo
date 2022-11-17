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
      --clock string               set the clock used by Hugo, e.g. --clock 2021-11-06T22:30:00.00+09:00
  -e, --environment string         build environment
      --force                      init inside non-empty directory
  -f, --format string              config file format (default "toml")
  -h, --help                       help for site
      --ignoreVendorPaths string   ignores any _vendor for module paths matching the given Glob pattern
  -s, --source string              filesystem path to read files relative from
      --themesDir string           filesystem path to themes directory
```

### Options inherited from parent commands

```
      --config string      config file (default is path/config.yaml|json|toml)
      --configDir string   config dir (default "config")
      --debug              debug output
      --log                enable Logging
      --logFile string     log File path (if set, logging enabled automatically)
      --quiet              build in quiet mode
  -v, --verbose            verbose output
      --verboseLog         verbose logging
```

### SEE ALSO

* [hugo new](/commands/hugo_new/)	 - Create new content for your site

