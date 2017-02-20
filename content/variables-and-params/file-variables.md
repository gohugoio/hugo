---
title: File Variables
linktitle:
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
weight:
categories: [variables and params]
tags: [files]
draft: false
slug:
aliases: []
toc: false
needsreview: true
notesforauthors:
---

**.File.Path** The original relative path of the page, e.g. `content/posts/foo.en.md`<br>
**.File.LogicalName** The name of the content file that represents a page, e.g. `foo.en.md`<br>
**.File.TranslationBaseName** The filename without extension or optional language identifier, e.g. `foo`<br>
**.File.BaseFileName** The filename without extension, e.g. `foo.en`<br>
**.File.Ext** or **.File.Extension** The file extension of the content file, e.g. `md`<br>
**.File.Lang** The language associated with the given file if [Multilingual](/content-management/multilingual-mode/ is enabled, e.g. `en`<br>
**.File.Dir** Given the path `content/posts/dir1/dir2/`, the relative directory path of the content file will be returned, e.g. `posts/dir1/dir2/`<br>