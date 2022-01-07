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
      --force           init inside non-empty directory
  -f, --format string   config file format (default "toml")
  -h, --help            help for site
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
      --quiet                      build in quiet mode
  -s, --source string              filesystem path to read files relative from
      --themesDir string           filesystem path to themes directory
  -v, --verbose                    verbose output
      --verboseLog                 verbose logging
```

### SEE ALSO

* [hugo new](/commands/hugo_new/)	 - Create new content for your site

