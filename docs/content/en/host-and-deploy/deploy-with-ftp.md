---
title: Deploy with FTP
description: Deploy your site to a web server with an FTP client.
categories: []
keywords: []
---

Use these instructions to deploy your site to a web server with an FTP client.

> [!TIP]
> If your hosting provider offers SSH access, deploy with [Rsync][] or [Rclone][] instead. Both tools transfer only the files that changed and remove orphaned files automatically.

## Deployment constraints

{{% include "/_common/deployment-constraints.md" %}}

FTP clients do not meet these requirements. The procedure below removes orphaned files by deleting the contents of the document root before each upload, and uploads every file each time you deploy.

## Prerequisites

Please complete the following tasks before continuing:

1. Create a Hugo project and test it with the `hugo server` command.
1. Set the [`baseURL`][] in your project configuration to the URL of your production site.
1. Obtain the host name, user name, and password for your hosting account, and the path to the document root, the directory from which your web server serves files.
1. Install an FTP client that supports SFTP.

If your hosting provider supports SFTP, use it instead of FTP. FTP sends your password and files without encryption.

## Procedure

Step 1
: Build your site from the root of your project.

  ```sh
  hugo build --cleanDestinationDir --gc --minify
  ```

Step 2
: Connect to your web server with your FTP client and open the document root.

Step 3
: Delete the contents of the document root.

Step 4
: Upload the contents of the `public` directory to the document root.

  > [!IMPORTANT]
  > Upload the files and directories inside the `public` directory, not the `public` directory itself. For example, the `index.html` file for your home page must be at the top level of the document root.

Your site is incomplete from the moment you delete the contents of the document root until the upload finishes.

## Command-line flags

The `hugo build` command uses these flags:

`--cleanDestinationDir`
: Removes orphaned files from the `public` directory. See [configure build][].

`--gc`
: Removes unused files from the file cache. See [configure file caches][].

`--minify`
: Minifies the output. See [configure minify][].

[Rclone]: /host-and-deploy/deploy-with-rclone/
[Rsync]: /host-and-deploy/deploy-with-rsync/
[`baseURL`]: /configuration/all/#baseurl
[configure build]: /configuration/build/
[configure file caches]: /configuration/caches/
[configure minify]: /configuration/minify/
