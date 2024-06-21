---
title: "hugo gen man"
slug: hugo_gen_man
url: /commands/hugo_gen_man/
---
## hugo gen man

Generate man pages for the Hugo CLI

### Synopsis

This command automatically generates up-to-date man pages of Hugo's
	command-line interface.  By default, it creates the man page files
	in the "man" directory under the current directory.

```
hugo gen man [flags] [args]
```

### Options

```
      --dir string   the directory to write the man pages. (default "man/")
  -h, --help         help for man
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

* [hugo gen](/commands/hugo_gen/)	 - A collection of several useful generators.

