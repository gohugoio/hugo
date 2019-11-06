---
title: File Variables
linktitle:
description: "You can access filesystem-related data for a content file in the `.File` variable."
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [variables and params]
keywords: [files]
draft: false
menu:
  docs:
    parent: "variables"
    weight: 40
weight: 40
sections_weight: 40
aliases: [/variables/file-variables/]
toc: false
---

{{% note "Rendering Local Files" %}}
For information on creating shortcodes and templates that tap into Hugo's file-related feature set, see [Local File Templates](/templates/files/).
{{% /note %}}

The `.File` object contains the following fields:

.File.Path
: the original relative path of the page, relative to the content dir (e.g., `posts/foo.en.md`)

.File.LogicalName
: the name of the content file that represents a page (e.g., `foo.en.md`)

.File.TranslationBaseName
: the filename without extension or optional language identifier (e.g., `foo`)

.File.ContentBaseName 
: is a either TranslationBaseName or name of containing folder if file is a leaf bundle.
  
.File.BaseFileName
: the filename without extension (e.g., `foo.en`)


.File.Ext
: the file extension of the content file (e.g., `md`); this can also be called using `.File.Extension` as well. Note that it is *only* the extension without `.`.

.File.Lang
: the language associated with the given file if Hugo's [Multilingual features][multilingual] are enabled (e.g., `en`)

.File.Dir
: given the path `content/posts/dir1/dir2/`, the relative directory path of the content file will be returned (e.g., `posts/dir1/dir2/`). Note that the path separator (`\` or `/`) could be dependent on the operating system.

[Multilingual]: /content-management/multilingual/
