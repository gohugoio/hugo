---
title: Deploy with Rclone
description: Deploy your site to a web server with Rclone.
categories: []
keywords: []
---

Use these instructions to deploy your site to a web server with [Rclone][].

## Deployment constraints

{{% include "/_common/deployment-constraints.md" %}}

Rclone meets both requirements.

## Prerequisites

Please complete the following tasks before continuing:

1. Create a Hugo project and test it with the `hugo server` command.
1. Set the [`baseURL`][] in your project configuration to the URL of your production site.
1. Configure SSH key authentication to log in to your web server, then log in once with the `ssh` command to add your web server to your `known_hosts` file. See your hosting provider's documentation for details.
1. Load your private key into `ssh-agent` on your local machine, typically with the `ssh-add` command. Rclone authenticates through the agent and does not read your key files directly. See [SSH authentication][] in the Rclone documentation for details.
1. Grant your SSH user write access to the document root, the directory from which your web server serves files.
1. Install [`rclone`][] on your local machine. You do not need to install `rclone` on your web server.
1. Use a POSIX shell such as Bash or Zsh on Linux or macOS. On Windows, use [Windows Subsystem for Linux][] (WSL).

## Procedure

Step 1
: Create a `deploy.sh` file in the root of your project, adjusting the configuration values as needed.

  ```sh {file="deploy.sh" copy=true}
  #!/usr/bin/env sh

  #------------------------------------------------------------------------------
  # @file
  # Builds a Hugo project and deploys it to a web server with Rclone.
  #------------------------------------------------------------------------------

  # Exit on error or undefined variables
  set -eu

  # Configuration
  SSH_USER="jsmith"
  SSH_HOST="example.org"
  SSH_PORT="22"
  SRC="public/"
  DEST="/var/www/html/"

  # Build
  hugo build \
    --cleanDestinationDir \
    --gc \
    --minify

  # Deploy
  rclone sync \
    --checksum \
    --config "" \
    --no-update-dir-modtime \
    --stats 0 \
    --verbose \
    "${SRC}" ":sftp,host=${SSH_HOST},port=${SSH_PORT},user=${SSH_USER},known_hosts_file=~/.ssh/known_hosts:${DEST}"
  ```

  Set `SRC` to the path of your `public` directory. Set `DEST` to the path of the document root on your web server. A path that does not begin with a slash is relative to your home directory, which is common with shared hosting.

Step 2
: Make the script executable.

  ```sh
  chmod +x deploy.sh
  ```

Step 3
: Run the script to build and deploy your site.

  ```sh
  ./deploy.sh
  ```

  To preview changes without modifying anything on your web server, add the `--dry-run` flag to the `rclone sync` command.

## Command-line flags

The `hugo build` command uses these flags:

`--cleanDestinationDir`
: Removes orphaned files from the `public` directory. See [configure build][].

`--gc`
: Removes unused files from the file cache. See [configure file caches][].

`--minify`
: Minifies the output. See [configure minify][].

The `rclone sync` command deletes files on your web server that do not exist in the local `public` directory, and uses these flags:

`--checksum`
: Compares files by checksum instead of by modification time and size. To compare checksums, Rclone runs commands such as `md5sum` on your web server. If your SSH user cannot run shell commands, Rclone compares files by size only.

`--config ""`
: Ignores the Rclone configuration file. The script defines the connection with an inline [connection string][] instead.

`--no-update-dir-modtime`
: Skips setting directory modification times on your web server.

`--stats 0`
: Disables periodic transfer statistics.

`--verbose`
: Lists each file that is transferred or deleted.

## Related resources

For more information on deploying your site with Rclone, consult the official documentation:

- [Rclone documentation][]
- [Rclone SFTP documentation][]

[Rclone SFTP documentation]: https://rclone.org/sftp/
[Rclone documentation]: https://rclone.org/docs/
[Rclone]: https://rclone.org/
[SSH authentication]: https://rclone.org/sftp/#ssh-authentication
[Windows Subsystem for Linux]: https://learn.microsoft.com/en-us/windows/wsl/
[`baseURL`]: /configuration/all/#baseurl
[`rclone`]: https://rclone.org/install/
[configure build]: /configuration/build/
[configure file caches]: /configuration/caches/
[configure minify]: /configuration/minify/
[connection string]: https://rclone.org/docs/#connection-strings
