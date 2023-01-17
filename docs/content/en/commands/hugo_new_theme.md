---
title: "hugo new theme"
slug: hugo_new_theme
url: /commands/hugo_new_theme/
---
## hugo new theme

Create a new theme

### Synopsis

Create a new theme (skeleton) called [name] in ./themes.
New theme is a skeleton. Please add content to the touched files. Add your
name to the copyright line in the license and adjust the theme.toml file
as you see fit.

```
hugo new theme [name] [flags]
```

### Options

```
      --clock string               set the clock used by Hugo, e.g. --clock 2021-11-06T22:30:00.00+09:00
  -e, --environment string         build environment
  -h, --help                       help for theme
      --ignoreVendorPaths string   ignores any _vendor for module paths matching the given Glob pattern
  -s, --source string              filesystem path to read files relative from
      --themesDir string           filesystem path to themes directory
```

### Options inherited from parent commands

```
      --config string      config file (default is hugo.yaml|json|toml)
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

