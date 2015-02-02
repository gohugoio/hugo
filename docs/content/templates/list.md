---
aliases:
- /layout/indexes/
date: 2013-07-01
linktitle: List of Content
menu:
  main:
    parent: layout
next: /templates/homepage
prev: /templates/content
title: Content List Template
weight: 40
---

A list template is any template that will be used to render multiple pieces of
content in a single HTML page (with the exception of the [homepage](/layout/homepage/) which has a
dedicated template).

We are using the term list in its truest sense, a sequential arrangement
of material, especially in alphabetical or numerical order. Hugo uses
list templates to render anyplace where content is being listed such as
taxonomies and sections.

## Which Template will be rendered?

Hugo uses a set of rules to figure out which template to use when
rendering a specific page.

Hugo will use the following prioritized list. If a file isn’t present,
then the next one in the list will be used. This enables you to craft
specific layouts when you want to without creating more templates
than necessary. For most sites only the \_default file at the end of
the list will be needed.


### Section Lists

A Section will be rendered at /`SECTION`/ (e.g.&nbsp;http://spf13.com/project/)

* /layouts/section/`SECTION`.html
* /layouts/\_default/section.html
* /layouts/\_default/list.html
* /themes/`THEME`/layouts/section/`SECTION`.html
* /themes/`THEME`/layouts/\_default/section.html
* /themes/`THEME`/layouts/\_default/list.html


### Taxonomy Lists

A Taxonomy will be rendered at /`PLURAL`/`TERM`/ (e.g.&nbsp;http://spf13.com/topics/golang/) from:

* /layouts/taxonomy/`SINGULAR`.html (e.g.&nbsp;`/layouts/taxonomy/topic.html`)
* /layouts/\_default/taxonomy.html
* /layouts/\_default/list.html
* /themes/`THEME`/layouts/taxonomy/`SINGULAR`.html
* /themes/`THEME`/layouts/\_default/taxonomy.html
* /themes/`THEME`/layouts/\_default/list.html

### Section RSS

A Section’s RSS will be rendered at /`SECTION`/index.xml (e.g.&nbsp;http://spf13.com/project/index.xml)

*Hugo ships with its own [RSS 2.0][] template. In most cases this will
be sufficient, and an RSS template will not need to be provided by the
user.*

Hugo provides the ability for you to define any RSS type you wish, and
can have different RSS files for each section and taxonomy.

* /layouts/section/`SECTION`.rss.xml
* /layouts/\_default/rss.xml
* /themes/`THEME`/layouts/section/`SECTION`.rss.xml
* /themes/`THEME`/layouts/\_default/rss.xml

### Taxonomy RSS

A Taxonomy’s RSS will be rendered at /`PLURAL`/`TERM`/index.xml (e.g.&nbsp;http://spf13.com/topics/golang/index.xml)

*Hugo ships with its own [RSS 2.0][] template. In most cases this will
be sufficient, and an RSS template will not need to be provided by the
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
This content template is used for [spf13.com](http://spf13.com/).
It makes use of [partial templates](/templates/partials/). All examples use a
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

    {{ partial "footer.html" . }}

### Example taxonomy template (tag.html)
This content template is used for [spf13.com](http://spf13.com/).
It makes use of [partial templates](/templates/partials/). All examples use a
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

    {{ partial "footer.html" . }}

## Ordering Content

In the case of Hugo each list will render the content based on metadata provided in the [front
matter](/content/front-matter/). See [ordering content](/content/ordering/) for more information.

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

### Order by PublishDate

    {{ range .Data.Pages.ByPublishDate }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .PublishDate.Format "Mon, Jan 2, 2006" }}</div>
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
group pages by Section, Type, Date etc.

Here are a variety of different ways you can group the content items in
your list templates:

### Grouping by Page field

    {{ range .Data.Pages.GroupBy "Section" }}
    <h3>{{ .Key }}</h3>
    <ul>
        {{ range .Pages }}
        <li>
        <a href="{{ .Permalink }}">{{ .Title }}</a>
        <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
        </li>
        {{ end }}
    </ul>
    {{ end }}

### Grouping by Page date

    {{ range .Data.Pages.GroupByDate "2006-01" }}
    <h3>{{ .Key }}</h3>
    <ul>
        {{ range .Pages }}
        <li>
        <a href="{{ .Permalink }}">{{ .Title }}</a>
        <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
        </li>
        {{ end }}
    </ul>
    {{ end }}

### Grouping by Page publish date

    {{ range .Data.Pages.GroupByPublishDate "2006-01" }}
    <h3>{{ .Key }}</h3>
    <ul>
        {{ range .Pages }}
        <li>
        <a href="{{ .Permalink }}">{{ .Title }}</a>
        <div class="meta">{{ .PublishDate.Format "Mon, Jan 2, 2006" }}</div>
        </li>
        {{ end }}
    </ul>
    {{ end }}

### Grouping by Page param

    {{ range .Data.Pages.GroupByParam "param_key" }}
    <h3>{{ .Key }}</h3>
    <ul>
        {{ range .Pages }}
        <li>
        <a href="{{ .Permalink }}">{{ .Title }}</a>
        <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
        </li>
        {{ end }}
    </ul>
    {{ end }}

### Grouping by Page param in date format

    {{ range .Data.Pages.GroupByParamDate "param_key" "2006-01" }}
    <h3>{{ .Key }}</h3>
    <ul>
        {{ range .Pages }}
        <li>
        <a href="{{ .Permalink }}">{{ .Title }}</a>
        <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
        </li>
        {{ end }}
    </ul>
    {{ end }}

### Reversing Key Order

The ordering of the groups is performed by keys in alpha-numeric order (A–Z,
1–100) and in reverse chronological order (newest first) for dates.

While these are logical defaults, they are not always the desired order. There
are two different syntaxes to change the order, they both work the same way, so
it’s really just a matter of preference.

#### Reverse method

    {{ range (.Data.Pages.GroupBy "Section").Reverse }}
    ...

    {{ range (.Data.Pages.GroupByDate "2006-01").Reverse }}
    ...


#### Providing the (alternate) direction

    {{ range .Data.Pages.GroupByDate "2006-01" "asc" }}
    ...

    {{ range .Data.Pages.GroupBy "Section" "desc" }}
    ...

### Ordering Pages within Group

Because Grouping returns a key and a slice of pages, all of the ordering methods listed above are available.

In this example I’ve ordered the groups in chronological order and the content
within each group in alphabetical order by title.

    {{ range .Data.Pages.GroupByDate "2006-01" "asc" }}
    <h3>{{ .Key }}</h3>
    <ul>
        {{ range .Pages.ByTitle }}
        <li>
        <a href="{{ .Permalink }}">{{ .Title }}</a>
        <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
        </li>
        {{ end }}
    </ul>
    {{ end }}

## Filtering & Limiting Content

Sometimes you only want to list a subset of the available content. A common
request is to only display “Posts” on the homepage. Using the `where` function
you can do just that.

### First 

`first` works like the `limit` keyword in SQL. It reduces the array to only the
first X elements. It takes the array and number of elements as input.

    {{ range first 10 .Data.Pages }}
        {{ .Render "summary"}}
    {{ end }}

### Where

`where` works in a similar manner to the `where` keyword in SQL. It selects all
elements of the slice that match the provided field and value. It takes three
arguments 'array or slice of maps or structs', 'key or field name' and 'match
value'

    {{ range where .Data.Pages "Section" "post" }}
       {{ .Content}}
    {{ end }}

### First & Where Together

Using both together can be very powerful.

    {{ range first 5 (where .Data.Pages "Section" "post") }}
       {{ .Content}}
    {{ end }}

If `where` or `first` receives invalid input or a field name that doesn’t exist they will provide an error and stop site generation.

These are both template functions and work on not only
[lists](/templates/list/), but [taxonomies](/taxonomies/displaying/),
[terms](/templates/terms/) and [groups](/templates/list/).


[RSS 2.0]: http://cyber.law.harvard.edu/rss/rss.html "RSS 2.0 Specification"
