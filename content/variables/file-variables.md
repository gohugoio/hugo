---
title: File Variables
linktitle:
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [variables and params]
tags: [files]
draft: false
weight: 40
aliases: []
toc: false
needsreview: true
notesforauthors:
---

Hugo provides the ability to traverse your website's files on your server, including the local `Hugo server`. You can access file-system-related data for a piece of content via the `.File` variable.

{{% note "Local File Templates" %}}
For information on creating shortcodes and templates that tap into Hugo's file-related feature set, see [Local File Templates](/templates/local-file-templates/).
{{% /note %}}

The `.File` object contains the following fields:

`.File.Path`
: The original relative path of the page (e.g., `content/posts/foo.en.md`)

`.File.LogicalName`
: The name of the content file that represents a page (e.g., `foo.en.md`)

`.File.TranslationBaseName`
: The filename without extension or optional language identifier (e.g., `foo`)

`.File.BaseFileName`
: The filename without extension (e.g., `foo.en`)

`.File.Ext`
: The file extension of the content file (e.g., `md`). This can also be called using `.File.Extension`.

`.File.Lang`
: The language associated with the given file if Hugo's [Multilingual][] features are enabled (e.g., `en`)

`.File.Dir`
: Given the path `content/posts/dir1/dir2/`, the relative directory path of the content file will be returned (e.g., `posts/dir1/dir2/`)

[Multilingual]: /content-management/multilingual/