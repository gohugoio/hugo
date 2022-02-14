---
title: "hugo list"
slug: hugo_list
url: /commands/hugo_list/
---
## hugo list

Listing out various types of content

### Synopsis

Listing out various types of content.

List requires a subcommand, e.g. `hugo list drafts`.

### Options

```
  -h, --help   help for list
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

* [hugo](/commands/hugo/)	 - hugo builds your site
* [hugo list all](/commands/hugo_list_all/)	 - List all posts
* [hugo list drafts](/commands/hugo_list_drafts/)	 - List all drafts
* [hugo list expired](/commands/hugo_list_expired/)	 - List all posts already expired
* [hugo list future](/commands/hugo_list_future/)	 - List all posts dated in the future

