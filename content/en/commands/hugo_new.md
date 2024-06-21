---
title: "hugo new"
slug: hugo_new
url: /commands/hugo_new/
---
## hugo new

Create new content for your site

### Synopsis

Create a new content file and automatically set the date and title.
It will guess which kind of file to create based on the path provided.

You can also specify the kind with `-k KIND`.

If archetypes are provided in your theme or site, they will be used.

Ensure you run this within the root directory of your site.

### Options

```
  -h, --help   help for new
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
* [hugo new content](/commands/hugo_new_content/)	 - Create new content for your site
* [hugo new site](/commands/hugo_new_site/)	 - Create a new site (skeleton)
* [hugo new theme](/commands/hugo_new_theme/)	 - Create a new theme (skeleton)

