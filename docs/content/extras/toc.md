---
date: 2013-07-09
menu:
  main:
    parent: extras
next: /extras/urls
prev: /extras/highlighting
title: Table of Contents
weight: 100
---

Hugo will automatically parse the Markdown for your content and create
a Table of Contents you can use to guide readers to the sections within
your content.

## Usage

Simply create content like you normally would with the appropriate
headers.

Hugo will take this Markdown and create a table of contents stored in the
[content variable](/layout/variables/) `.TableOfContents`


## Template Example

This is example code of a [single.html template](/layout/content/).

    {{ partial "header.html" . }}
        <div id="toc" class="well col-md-4 col-sm-6">
        {{ .TableOfContents }}
        </div>
        <h1>{{ .Title }}</h1>
        {{ .Content }}
    {{ partial "footer.html" . }}


