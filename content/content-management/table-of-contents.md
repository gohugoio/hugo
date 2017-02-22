---
title: Table of Contents
linktitle:
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [content management]
tags: [table of contents, toc]
weight: 130
draft: false
aliases: [/extras/toc/]
toc: false
needsreview: true
---

Hugo will automatically parse the Markdown for your content and create
a Table of Contents you can use to guide readers to the sections within
your content.

{{% note "TOC Heading Levels are Fixed" %}}
{{% /note %}}

## Usage

Simply create content like you normally would with the appropriate headers.

Hugo will take this Markdown and create a table of contents stored in the [content variable](/variables-and-params/page-variables/) `.TableOfContents`.

## Template Example

This is example code of a [single.html template](/templates/single-page-templates/).

```golang
{{ partial "header.html" . }}
    <div id="toc" class="well col-md-4 col-sm-6">
    {{ .TableOfContents }}
    </div>
    <h1>{{ .Title }}</h1>
    {{ .Content }}
{{ partial "footer.html" . }}
```
