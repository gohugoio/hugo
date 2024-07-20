---
title: Deploy with Rclone
description: If you have access to your web host with SFTP/FTP/SSH/HTTP(DAV), you can use rclone to incrementally deploy your entire Hugo website.
categories: [hosting and deployment]
keywords: [deployment,rclone,sftp]
menu:
  docs:
    parent: hosting-and-deployment
toc: true
aliases: [/tutorials/deployment-with-rclone/]
---

## Assumptions

* A web host running a web server. This could be a shared hosting environment or a VPS.
* Access to your web host with any of the [protocols supported by rclone](https://rclone.org/#providers), such as SFTP.
* A functional static website built with Hugo
* Deploying from an [Rclone](https://rclone.org) compatible operating system
* You have [installed Rclone](https://rclone.org/install/).

**NB**: You can remove ``--interactive`` in the commands below once you are comfortable with rclone, if you wish. Also, ``--gc`` and ``--minify`` are optional in the ``hugo`` commands below.

## Getting started

The spoiler is that you can even deploy your entire website from any compatible OS with no configuration. Using SFTP for example:

```txt
hugo --gc --minify
rclone sync --interactive --sftp-host sftp.example.com --sftp-user www-data --sftp-ask-password public/ :sftp:www/
```

## Configure Rclone for even easier usage

The easiest way is simply to run `rclone config`.

The [Rclone docs](https://rclone.org/docs/) provide [an example of configuring Rclone to use SFTP](https://rclone.org/sftp/).

For the next commands, we will assume you configured a remote you named ``hugo-www``

The above 'spoiler' commands could become:

```txt
hugo --gc --minify
rclone sync --interactive public/ hugo-www:www/
```

After you issue the above commands (and respond to any prompts), check your website and you will see that it is deployed.
