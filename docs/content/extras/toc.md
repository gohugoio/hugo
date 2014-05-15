---
title: "Table of Contents"
date: "2013-07-09"
weight: 70
menu:
  main:
    parent: 'extras'
---

Hugo will automatically parse the markdown for your content and create
a Table of Contents you can use to guide readers to the sections within
your content.

## Usage

Simply create content like you normally would with the appropriate
headers.

Hugo will take this markdown and create a table of contents stored in the
[content variable](/layout/variables) .TableOfContents


## Template Example

This is example code of a [single.html template](/layout/content).

    {{ template "chrome/header.html" . }}
        <div id="toc" class="well col-md-4 col-sm-6">
        {{ .TableOfContents }}
        </div>
        <h1>{{ .Title }}</h1>
        {{ .Content }}
    {{ template "chrome/footer.html" . }}


