---
title: "Taxonomy Templates"
date: "2013-07-01"
linktitle: "Templates"
aliases: ["/indexes/templates/"]
weight: 30
menu:
  main:
    parent: 'taxonomy'
---

There are two different templates that the use of indexes will require you to provide.

The first is a list of all the content assigned to a specific index key. The
second is a [list](/indexes/lists/) of all keys for that index. This document
addresses the template used for the first type.

## Creating index templates
For each index type a template needs to be provided to render the index page.
In the case of tags, this will render the content for `/tags/TAGNAME/`.

The template must be called the singular name of the index and placed in 
layouts/indexes

    .
    └── layouts
        └── indexes
            └── category.html

The template will be provided Data about the index. 

## Variables

The following variables are available to the index template:

**.Title**  The title for the content. <br>
**.Date** The date the content is published on.<br>
**.Permalink** The Permanent link for this page.<br>
**.RSSLink** Link to the indexes' rss link. <br>
**.Data.Pages** The content that is assigned this index.<br>
**.Data.`singular`** The index itself.<br>

## Example
    {{ template "chrome/header.html" . }}
    {{ template "chrome/subheader.html" . }}

    <section id="main">
      <div>
       <h1 id="title">{{ .Title }}</h1>
        {{ range .Data.Pages }}
            {{ .Render "summary"}}
        {{ end }}
      </div>
    </section>

    {{ template "chrome/footer.html" }}
