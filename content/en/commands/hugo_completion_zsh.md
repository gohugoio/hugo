---
title: "hugo completion zsh"
slug: hugo_completion_zsh
url: /commands/hugo_completion_zsh/
---
## hugo completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(hugo completion zsh)

To load completions for every new session, execute once:

#### Linux:

	hugo completion zsh > "${fpath[1]}/_hugo"

#### macOS:

	hugo completion zsh > $(brew --prefix)/share/zsh/site-functions/_hugo

You will need to start a new shell for this setup to take effect.


```
hugo completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
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

* [hugo completion](/commands/hugo_completion/)	 - Generate the autocompletion script for the specified shell

