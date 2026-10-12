---
title: Deploy with Rsync
description: Deploy your site to a web server with Rsync.
categories: []
keywords: []
---

Use these instructions to deploy your site to a web server with [Rsync][].

## Deployment constraints

{{% include "/_common/deployment-constraints.md" %}}

Rsync meets both requirements.

## Prerequisites

Please complete the following tasks before continuing:

1. Create a Hugo project and test it with the `hugo server` command.
1. Set the [`baseURL`][] in your project configuration to the URL of your production site.
1. Configure SSH key authentication to log in to your web server, then log in once with the `ssh` command to add your web server to your `known_hosts` file. See your hosting provider's documentation for details.
1. Grant your SSH user write access to the document root, the directory from which your web server serves files.
1. Install [`rsync`][] on your local machine and on your web server.
1. Use a POSIX shell such as Bash or Zsh on Linux or macOS. On Windows, use [Windows Subsystem for Linux][] (WSL).

## Procedure

Step 1
: Create a `deploy.sh` file in the root of your project, adjusting the configuration values as needed.

  ```sh {file="deploy.sh" copy=true}
  #!/usr/bin/env sh

  #------------------------------------------------------------------------------
  # @file
  # Builds a Hugo project and deploys it to a web server with Rsync.
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
  rsync \
    --checksum \
    --chmod=D755,F644 \
    --compress \
    --delete \
    --perms \
    --recursive \
    --rsh="ssh -p ${SSH_PORT}" \
    --verbose \
    "${SRC}" "${SSH_USER}@${SSH_HOST}:${DEST}"
  ```

  Set `SRC` to the path of your `public` directory, including the trailing slash. Set `DEST` to the path of the document root on your web server. A path that does not begin with a slash is relative to your home directory, which is common with shared hosting.

  > [!IMPORTANT]
  > The trailing slash on `SRC` is required. Without it, Rsync copies the directory itself, rather than its contents, to the document root.

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

  To preview changes without modifying anything on your web server, add the `--dry-run` flag to the `rsync` command.

## Command-line flags

The `hugo build` command uses these flags:

`--cleanDestinationDir`
: Removes orphaned files from the `public` directory. See [configure build][].

`--gc`
: Removes unused files from the file cache. See [configure file caches][].

`--minify`
: Minifies the output. See [configure minify][].

The `rsync` command uses these flags:

`--checksum`
: Compares files by checksum instead of by modification time and size.

`--chmod=D755,F644`
: Makes directories and files readable by your web server.

`--compress`
: Compresses file data during transfer.

`--delete`
: Deletes files on your web server that do not exist in the local `public` directory.

`--perms`
: Applies the `--chmod` permissions to existing files, not only to new files.

`--recursive`
: Copies directories recursively.

`--rsh`
: Specifies the SSH command and port.

`--verbose`
: Lists each file that is transferred or deleted.

## Related resources

For more information on deploying your site with Rsync, consult the official documentation:

- [Rsync documentation][]

[Rsync documentation]: https://download.samba.org/pub/rsync/rsync.1
[Rsync]: https://rsync.samba.org/
[Windows Subsystem for Linux]: https://learn.microsoft.com/en-us/windows/wsl/
[`baseURL`]: /configuration/all/#baseurl
[`rsync`]: https://rsync.samba.org/download.html
[configure build]: /configuration/build/
[configure file caches]: /configuration/caches/
[configure minify]: /configuration/minify/
