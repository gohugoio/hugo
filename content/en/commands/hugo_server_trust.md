---
title: "hugo server trust"
slug: hugo_server_trust
url: /commands/hugo_server_trust/
---
## hugo server trust

Install the local CA in the system trust store.

```
hugo server trust [flags] [args]
```

### Options

```
  -h, --help        help for trust
      --uninstall   Uninstall the local CA (but do not delete it).
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

* [hugo server](/commands/hugo_server/)	 - A high performance webserver

