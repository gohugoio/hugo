---
title: "hugo gen doc"
slug: hugo_gen_doc
url: /commands/hugo_gen_doc/
---
## hugo gen doc

Generate Markdown documentation for the Hugo CLI.

### Synopsis

Generate Markdown documentation for the Hugo CLI.
			This command is, mostly, used to create up-to-date documentation
	of Hugo's command-line interface for https://gohugo.io/.

	It creates one Markdown file per command with front matter suitable
	for rendering in Hugo.

```
hugo gen doc [flags] [args]
```

### Options

```
      --dir string   the directory to write the doc. (default "/tmp/hugodoc/")
  -h, --help         help for doc
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

