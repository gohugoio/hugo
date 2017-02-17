---
title: Markdown and Supported Formats
linktitle:
description: Hugo uses the BlackFriday markdown parser for content files but also provides support for additional syntaxes (eg, Asciidoc) via external helpers.
date: 2017-01-10
publishdate: 2017-01-10
lastmod: 2017-01-10
categories: [content management]
tags: [markdown,asciidoc,mmark,content format]
weight: 10
draft: false
slug:
aliases: [/content/markdown-extras/,/content/supported-formats/,/content/markdown/]
toc:
notes:
---

## Markdown

Markdown is the natively supported content format for Hugo and is rendered using the excellent [BlackFriday project][], a markdown parser written in Golang.

{{% note "Deeply Nested Lists" %}}
Hugo uses [BlackFriday](https://github.com/russross/blackfriday), a markdown processor written in Golang. BlackFriday has a known issue [(#329)](https://github.com/russross/blackfriday/issues/329) with handling deeply nested lists, but there is a workaround. If you write lists in markdown, be sure to include 4 spaces (i.e., <kbd>tab</kbd>) rather than 2 to delimit nesting of lists.
{{% /note %}}

## Additional Resources

<!-- Mention shortcodes as markdown extension -->

* [Markdown Tutorial][]

[BlackFriday project]: https://github.com/russross/blackfriday
[Markdown Tutorial]: http://www.markdowntutorial.com/