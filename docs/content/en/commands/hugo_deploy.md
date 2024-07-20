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
hugo deploy [flags] [args]
```

### Options

```
      --confirm          ask for confirmation before making changes to the target
      --dryRun           dry run
      --force            force upload of all files
  -h, --help             help for deploy
      --invalidateCDN    invalidate the CDN cache listed in the deployment target (default true)
      --maxDeletes int   maximum # of files to delete, or -1 to disable (default 256)
      --target string    target deployment from deployments section in config file; defaults to the first one
      --workers int      number of workers to transfer files. defaults to 10 (default 10)
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

