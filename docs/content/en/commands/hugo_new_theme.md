---
title: "hugo new theme"
slug: hugo_new_theme
url: /commands/hugo_new_theme/
---
## hugo new theme

Create a new theme (skeleton)

### Synopsis

Create a new theme (skeleton) called [name] in ./themes.
New theme is a skeleton. Please add content to the touched files. Add your
name to the copyright line in the license and adjust the theme.toml file
according to your needs.

```
hugo new theme [name] [flags]
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
      --logLevel string            log level (debug|info|warn|error)
      --quiet                      build in quiet mode
  -M, --renderToMemory             render to memory (mostly useful when running the server)
  -s, --source string              filesystem path to read files relative from
      --themesDir string           filesystem path to themes directory
  -v, --verbose                    verbose output
```

### SEE ALSO

* [hugo new](/commands/hugo_new/)	 - Create new content for your site

