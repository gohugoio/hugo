---
title: "hugo deploy"
slug: hugo_deploy
url: /commands/hugo_deploy/
---
## hugo deploy

Deploy your site to a Cloud provider.

### Synopsis

Deploy your site to a Cloud provider.

See https://gohugo.io/hosting-and-deployment/hugo-deploy/ for detailed
documentation.


```
hugo deploy [flags]
```

### Options

```
      --clock string               set the clock used by Hugo, e.g. --clock 2021-11-06T22:30:00.00+09:00
      --confirm                    ask for confirmation before making changes to the target
      --dryRun                     dry run
  -e, --environment string         build environment
      --force                      force upload of all files
  -h, --help                       help for deploy
      --ignoreVendorPaths string   ignores any _vendor for module paths matching the given Glob pattern
      --invalidateCDN              invalidate the CDN cache listed in the deployment target (default true)
      --maxDeletes int             maximum # of files to delete, or -1 to disable (default 256)
  -s, --source string              filesystem path to read files relative from
      --target string              target deployment from deployments section in config file; defaults to the first one
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

* [hugo](/commands/hugo/)	 - hugo builds your site

