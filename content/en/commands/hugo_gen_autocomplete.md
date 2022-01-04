---
title: "hugo gen autocomplete"
slug: hugo_gen_autocomplete
url: /commands/hugo_gen_autocomplete/
---
## hugo gen autocomplete

Generate shell autocompletion script for Hugo

### Synopsis

Generates a shell autocompletion script for Hugo.

The script is written to the console (stdout).

To write to file, add the `--completionfile=/path/to/file` flag.

Add `--type={bash, zsh, fish or powershell}` flag to set alternative
shell type.

Logout and in again to reload the completion scripts,
or just source them in directly:

	$ . /etc/bash_completion or /path/to/file

```
hugo gen autocomplete [flags]
```

### Options

```
  -f, --completionfile string   autocompletion file, defaults to stdout
  -h, --help                    help for autocomplete
  -t, --type string             autocompletion type (bash, zsh, fish, or powershell) (default "bash")
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

* [hugo gen](/commands/hugo_gen/)	 - A collection of several useful generators.

