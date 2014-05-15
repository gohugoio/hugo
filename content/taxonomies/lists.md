---
title: "Taxonomy Lists"
date: "2013-07-01"
aliases: ["/indexes/lists/","/doc/indexes/", "/extras/indexes"]
linktitle: "Lists"
weight: 40
menu:
  main:
    parent: 'taxonomy'
---

An index list is a list of all the keys that are contained in the index. When a
template is present, this will be rendered at `/IndexPlural/`

Hugo also supports creating pages that list your values for each index along
with the number of content items associated with the index key. These are
global pages, not attached to any specific content, but rather display the meta
data in aggregate.

To have hugo create these list of indexes pages, simply create a template in
/layouts/indexes/ called indexes.html

Hugo can order the meta data in two different ways. It can be ordered by the
number of content assigned to that key or alphabetically.


## Example indexes.html file (alphabetical)

    {{ template "chrome/header.html" . }}
    {{ template "chrome/subheader.html" . }}

    <section id="main">
      <div>
       <h1 id="title">{{ .Title }}</h1>
       <ul>
       {{ $data := .Data }}
        {{ range $key, $value := .Data.Index.Alphabetical }}
        <li><a href="{{ $data.Plural }}/{{ $value.Name | urlize }}"> {{ $value.Name }} </a> {{ $value.Count }} </li>
        {{ end }}
       </ul>
      </div>
    </section>
    {{ template "chrome/footer.html" }}

## Example indexes.html file (ordered)

    {{ template "chrome/header.html" . }}
    {{ template "chrome/subheader.html" . }}

    <section id="main">
      <div>
       <h1 id="title">{{ .Title }}</h1>
       <ul>
       {{ $data := .Data }}
        {{ range $key, $value := .Data.Index.ByCount }}
        <li><a href="{{ $data.Plural }}/{{ $value.Name | urlize }}"> {{ $value.Name }} </a> {{ $value.Count }} </li>
        {{ end }}
       </ul>
      </div>
    </section>

    {{ template "chrome/footer.html" }}

## Variables available to list of indexes pages.

**.Title**  The title for the content. <br>
**.Date** The date the content is published on.<br>
**.Permalink** The Permanent link for this page.<br>
**.RSSLink** Link to the indexes' rss link. <br>
**.Data.Singular** The singular name of the index <br>
**.Data.Plural** The plural name of the index<br>
**.Data.Index** The Index itself<br>
**.Data.Index.Alphabetical** The Index alphabetized<br>
**.Data.Index.ByCount** The Index ordered by popularity<br>
