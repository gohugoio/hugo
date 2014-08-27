---
aliases:
- /layout/indexes/
date: 2013-07-01
linktitle: List
menu:
  main:
    parent: layout
next: /templates/homepage
prev: /templates/content
title: Content List Template
weight: 40
---

A list template is any template that will be used to render multiple pieces of
content in a single html page (with the exception of the [homepage](/layout/homepage) which has a
dedicated template).

We are using the term list in its truest sense, a sequential arrangement
of material, especially in alphabetical or numerical order. Hugo uses
list templates to render anyplace where content is being listed such as
taxonomies and sections.

## Which Template will be rendered?

Hugo uses a set of rules to figure out which template to use when
rendering a specific page.

Hugo will use the following prioritized list. If a file isn’t present
than the next one in the list will be used. This enables you to craft
specific layouts when you want to without creating more templates
then necessary. For most sites only the \_default file at the end of
the list will be needed.


### Section Lists

A Section will be rendered at /`SECTION`/

* /layouts/section/`SECTION`.html
* /layouts/\_default/section.html
* /layouts/\_default/list.html
* /themes/`THEME`/layouts/section/`SECTION`.html
* /themes/`THEME`/\_default/section.html
* /themes/`THEME`/layouts/\_default/list.html


### Taxonomy Lists

A Taxonomy will be rendered at /`PLURAL`/`TERM`/

* /layouts/taxonomy/`SINGULAR`.html
* /layouts/\_default/taxonomy.html
* /layouts/\_default/list.html
* /themes/`THEME`/layouts/taxonomy/`SINGULAR`.html
* /themes/`THEME`/\_default/taxonomy.html
* /themes/`THEME`/layouts/\_default/list.html

### Section RSS

A Section’s RSS will be rendered at /`SECTION`/index.xml

*Hugo ships with it’s own ATOM 2.0 RSS template. In most cases this will
be sufficient and an RSS template will not need to be provided by the
user.*

Hugo provides the ability for you to define any RSS type you wish, and
can have different RSS files for each section and taxonomy.

* /layouts/section/`SECTION`.rss.xml
* /layouts/\_default/rss.xml
* /themes/`THEME`/layouts/section/`SECTION`.rss.xml
* /themes/`THEME`/layouts/\_default/rss.xml

### Taxonomy RSS

A Taxonomy’s RSS will be rendered at /`PLURAL`/`TERM`/index.xml

*Hugo ships with it’s own ATOM 2.0 RSS template. In most cases this will
be sufficient and an RSS template will not need to be provided by the
user.*

Hugo provides the ability for you to define any RSS type you wish, and
can have different RSS files for each section and taxonomy.

* /layouts/taxonomy/`SINGULAR`.rss.xml
* /layouts/\_default/rss.xml
* /themes/`THEME`/layouts/taxonomy/`SINGULAR`.rss.xml
* /themes/`THEME`/layouts/\_default/rss.xml


## Variables

List pages are of the type "node" and have all the [node
variables](/templates/variables/) and [site
variables](/templates/variables/) available to use in the templates. 

Taxonomy pages will additionally have:

**.Data.`singular`** The taxonomy itself.<br>

## Example List Template Pages

### Example section template (post.html)
This content template is used for [spf13.com](http://spf13.com).
It makes use of [partial templates](/templates/partials). All examples use a
[view](/templates/views/) called either "li" or "summary" which this example site
defined.

    {{ partial "header.html" . }}
    {{ partial "subheader.html" . }}

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

    {{ partial "footer.html" }}

### Example taxonomy template (tag.html)
This content template is used for [spf13.com](http://spf13.com).
It makes use of [partial templates](/templates/partials). All examples use a
[view](/templates/views/) called either "li" or "summary" which this example site
defined.

    {{ partial "header.html" . }}
    {{ partial "subheader.html" . }}

    <section id="main">
      <div>
       <h1 id="title">{{ .Title }}</h1>
        {{ range .Data.Pages }}
            {{ .Render "summary"}}
        {{ end }}
      </div>
    </section>

    {{ partial "footer.html" }}

## Ordering Content

In the case of Hugo each list will render the content based on metadata provided in the [front
matter](/content/front-matter). See [ordering content](/content/ordering) for more information.

Here are a variety of different ways you can order the content items in
your list templates:

### Order by Weight -> Date (default)

    {{ range .Data.Pages }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

### Order by Weight -> Date

    {{ range .Data.Pages.ByWeight }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

### Order by Date

    {{ range .Data.Pages.ByDate }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

### Order by Length

    {{ range .Data.Pages.ByLength }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}


### Order by Title

    {{ range .Data.Pages.ByTitle }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

### Order by LinkTitle

    {{ range .Data.Pages.ByLinkTitle }}
    <li>
    <a href="{{ .Permalink }}">{{ .LinkTitle }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

### Reverse Order
Can be applied to any of the above. Using Date for an example.

    {{ range .Data.Pages.ByDate.Reverse }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

## Grouping Content

Hugo provides some grouping functions for list pages. You can use them to
group pages by Section, Date etc.

Here are a variety of different ways you can group the content items in
your list templates:

### Grouping by Page field

    {{ range .Data.Pages.GroupBy "Section" "asc" }}
    <h3>{{ .Key }}</h3>
    <ul>
        {{ range .Data }}
        <li>
        <a href="{{ .Permalink }}">{{ .Title }}</a>
        <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
        </li>
        {{ end }}
    </ul>
    {{ end }}

### Grouping by Page date

    {{ range .Data.Pages.GroupByDate "2006-01" "desc" }}
    <h3>{{ .Key }}</h3>
    <ul>
        {{ range .Data }}
        <li>
        <a href="{{ .Permalink }}">{{ .Title }}</a>
        <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
        </li>
        {{ end }}
    </ul>
    {{ end }}
