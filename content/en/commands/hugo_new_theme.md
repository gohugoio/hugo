---
title: "hugo new theme"
slug: hugo_new_theme
url: /commands/hugo_new_theme/
---
## hugo new theme

Create a new site (skeleton)

### Synopsis

Create a new site in the provided directory.
The new site will have the correct structure, but no content or theme yet.
Use `hugo new [contentPath]` to create new content.

```
hugo new theme [path] [flags]
```

### Options

```
  -h, --help   help for theme
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

* [hugo new](/commands/hugo_new/)	 - Create new content for your site

