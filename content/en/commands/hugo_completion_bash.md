---
title: "hugo completion bash"
slug: hugo_completion_bash
url: /commands/hugo_completion_bash/
---
## hugo completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(hugo completion bash)

To load completions for every new session, execute once:

#### Linux:

	hugo completion bash > /etc/bash_completion.d/hugo

#### macOS:

	hugo completion bash > $(brew --prefix)/etc/bash_completion.d/hugo

You will need to start a new shell for this setup to take effect.


```
hugo completion bash
```

### Options

```
  -h, --help              help for bash
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

