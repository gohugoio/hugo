---
aliases:
- /doc/supported-formats/
date: 2015-08-01
menu:
  main:
    parent: content
next: /content/front-matter
prev: /content/organization
title: Supported Formats
weight: 15
toc: true
---

  Since 0.14, Hugo has defined a new concept called _external helpers_. It means that you can write your content using Asciidoc[tor], or reStructuredText. If you have files with associated extensions ([details](https://github.com/spf13/hugo/blob/77c60a3440806067109347d04eb5368b65ea0fe8/helpers/general.go#L65)), then Hugo will call external commands to generate the content.

  This means that you will have to install the associated tool on your machine to be able to use those formats.

  For example, for Asciidoc files, Hugo will try to call __asciidoctor__ or __asciidoc__ command.

  To use those formats, just use the standard extension and the front matter exactly as you would do with natively supported _.md_ files.

  Notes:

  * as these are external commands, generation performance for that content will heavily depend on the performance of those external tools.
  * this feature is still in early stage, hence feedback is even more welcome.
