---
title: "Index Templates"
date: "2013-07-01"
weight: 50
menu:
  main:
    parent: 'layout'
---

An index template is any template that will be used to render multiple pieces of
content (with the exception of the [homepage](/layout/homepage) which has a
dedicated template).

We are using the term index in its truest sense, a sequential arrangement of
material, especially in alphabetical or numerical order. In the case of Hugo
each index will render the content in newest first order based on the date
provided in the [front matter](/content/front-matter).

index pages are of the type "node" and have all the [node
variables](/layout/variables/) available to use in the templates.
All index templates live in the layouts/indexes directory. There are 3 different
kinds of indexes that Hugo can produce.

1. A listing of all the content for a given [section](/content/sections)
2. A listing of all the content for a given [index](/extras/indexes)
3. A listing of listings... [meta index](/extras/indexes)

It's critical that the name of the index template matches either:

1. The section name
2. The index singular name
3. "indexes"

The following illustrates the location of one of each of these types.

    ▾ layouts/
      ▾ indexes/
          indexes.html
          post.html
          tag.html

## Example section template (post.html)
This content template is used for [spf13.com](http://spf13.com).
It makes use of [chrome templates](/layout/chrome). All examples use a
[view](/layout/views/) called either "li" or "summary" which this example site
defined.

    {{ template "chrome/header.html" . }}
    {{ template "chrome/subheader.html" . }}

    <section id="main">
      <div>
       <h1 id="title">{{ .Title }}</h1>
            <ul id="list">
                {{ range .Data.Pages }}
                    {{ .Render "li"}}
                {{ end }}
            </ul>
      </div>
    </section>

    {{ template "chrome/footer.html" }}

## Example index template (tag.html)
This content template is used for [spf13.com](http://spf13.com).
It makes use of [chrome templates](/layout/chrome). All examples use a
[view](/layout/views/) called either "li" or "summary" which this example site
defined.

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

## Example listing of indexes template (indexes.html)
This content template is used for [spf13.com](http://spf13.com).
It makes use of [chrome templates](/layout/chrome). The list of indexes
templates cannot use a [content view](/layout/views) as they don't display the content, but
rather information about the content.

This particular template lists all of the Tags used on
[spf13.com](http://spf13.com) and provides a count for the number of pieces of
content tagged with each tag.

This example demonstrates two different approaches. The first uses .Data.Index and
the latter uses .Data.OrderedIndex. .Data.Index is alphabetical by key name, while
.Data.Orderedindex is ordered by the quantity of content assigned to that particular
index key.  In practice you would only use one of these approaches.

    {{ template "chrome/header.html" . }}
    {{ template "chrome/subheader.html" . }}

    <section id="main">
      <div>
       <h1 id="title">{{ .Title }}</h1>

       <ul>
       {{ $data := .Data }}
        {{ range $key, $value := .Data.Index }}
        <li><a href="{{ $data.Plural }}/{{ $key | urlize }}"> {{ $key }} </a> {{ len $value }} </li>
        {{ end }}
       </ul>
      </div>

       <ul>
        {{ range $data.OrderedIndex }}
        <li><a href="{{ $data.Plural }}/{{ .Name | urlize }}"> {{ .Name }} </a> {{ .Count }} </li>
        {{ end }}
       </ul>
    </section>

    {{ template "chrome/footer.html" }}




