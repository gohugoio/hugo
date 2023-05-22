---
title: "hugo new content"
slug: hugo_new_content
url: /commands/hugo_new_content/
---
## hugo new content

Create new content for your site

### Synopsis

Create a new content file and automatically set the date and title.
		It will guess which kind of file to create based on the path provided.
		
		You can also specify the kind with `-k KIND`.
		
		If archetypes are provided in your theme or site, they will be used.
		
		Ensure you run this within the root directory of your site.

```
hugo new content [path] [flags]
```

### Options

```
      --editor string   edit new content with this editor, if provided
  -f, --force           overwrite file if it already exists
  -h, --help            help for content
  -k, --kind string     content type to create
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

