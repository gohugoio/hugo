---
title: Section Page Templates
linktitle: Section Page Templates
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [lists,sections]
weight: 40
draft: false
aliases: []
toc: true
needsreview: true
---

## Introduction to the Template Lookup Order

{{< lookupexplanation >}}



## Lookup Order for Section Page Templates

Hugo uses a set of rules to figure out which template to use when rendering a specific page.

Hugo will use the following prioritized list. If a file isn’t present, then the next one in the list will be used. This enables you to craft specific layouts when you want to without creating more templates than necessary. For most sites only the \_default file at the end of the list will be needed.

### Section Lists

A Section will be rendered at /`SECTION`/ (e.g.&nbsp;http://spf13.com/project/)

* /layouts/section/`SECTION`.html
* /layouts/\_default/section.html
* /layouts/\_default/list.html
* /themes/`THEME`/layouts/section/`SECTION`.html
* /themes/`THEME`/layouts/\_default/section.html
* /themes/`THEME`/layouts/\_default/list.html

Note that a sections list page can also have a content file with frontmatter,  see [Source Organization](/overview/source-directory/}}).

### Taxonomy Lists

A Taxonomy will be rendered at /`PLURAL`/`TERM`/ (e.g.&nbsp;http://spf13.com/topics/golang/) from:

* /layouts/taxonomy/`SINGULAR`.html (e.g.&nbsp;`/layouts/taxonomy/topic.html`)
* /layouts/\_default/taxonomy.html
* /layouts/\_default/list.html
* /themes/`THEME`/layouts/taxonomy/`SINGULAR`.html
* /themes/`THEME`/layouts/\_default/taxonomy.html
* /themes/`THEME`/layouts/\_default/list.html

Note that a taxonomy list page can also have a content file with frontmatter,  see [Source Organization](/overview/source-directory/).

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

A list page is a `Page` and have all the [page variables](/templates/variables/)
and [site variables](/templates/variables/) available to use in the templates.

Taxonomy pages will additionally have:

**.Data.`Singular`** The taxonomy itself.<br>

## Example List Template Pages

### Example Section Template: `post.html`

This content template is used for [spf13.com](http://spf13.com/).
It makes use of [partial templates][partials]. All examples use a
[view](/templates/views/) called either "li" or "summary" which this example site
defined.

{{% input "layouts/section/post.html" %}}
```html
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
```
{{% /input %}}

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

In the case of Hugo, each list will render the content based on metadata provided in the [front
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

### Order by ExpiryDate

    {{ range .Data.Pages.ByExpiryDate }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .ExpiryDate.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

### Order by Lastmod

    {{ range .Data.Pages.ByLastmod }}
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

### Order by Parameter

Order based on the specified frontmatter parameter. Pages without that
parameter will use the site's `.Site.Params` default. If the parameter is not
found at all in some entries, those entries will appear together at the end
of the ordering.

The below example sorts a list of posts by their rating.

    {{ range (.Data.Pages.ByParam "rating") }}
      <!-- ... -->
    {{ end }}

If the frontmatter field of interest is nested beneath another field, you can
also get it:

    {{ range (.Date.Pages.ByParam "author.last_name") }}
      <!-- ... -->
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

The ordering of the groups is performed by keys in alphanumeric order (A–Z,
1–100) and in reverse chronological order (newest first) for dates.

While these are logical defaults, they are not always the desired order. There
are two different syntaxes to change the order; they both work the same way, so
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

In this example, I’ve ordered the groups in chronological order and the content
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
request is to only display “Posts” on the homepage. Using the `where` function,
you can do just that.

### `first`

`first` works in a similar manner to the [`limit` keyword in SQL][limitkeyword]. It reduces the array to only the `first N` elements. It takes the array and number of elements as input. `first` takes two arguments:

1. `array` or `slice of maps or structs`
2. `number of elements`

{{% input "layout/_default/section.html" %}}
```golang
{{ range first 10 .Data.Pages }}
  {{ .Render "summary" }}
{{ end }}
```
{{% /input %}}

### `where`

`where` works in a similar manner to the `where` keyword in SQL. It selects all elements of the array or slice that match the provided field and value. `where` takes three arguments:

1. `array` or a `slice of maps or structs`
2. `key` or `field name'
3. `match value`

{{% input "layouts/_default/.html" %}}
```html
{{ range where .Data.Pages "Section" "post" }}
   {{ .Content }}
{{ end }}
```
{{% /input %}}

### `first` and `where` Together

Using `first` and `where` together can be very powerful:

{{% input "first-and-where-together.html" %}}
```golang
{{ range first 5 (where .Data.Pages "Section" "post") }}
   {{ .Content }}
{{ end }}
```
{{% /input %}}

{{% note %}}
If `where` or `first` receives invalid input or a field name that doesn’t exist, it will return an error and stop site generation. `where` and `first` also work on taxonomy list templates *and* taxonomy terms templates. (See [Taxonomy Templates](/templates/taxonomy-templates/).)
{{% /note %}}

## `.Site.GetPage`

Every `Page` in Hugo has a `.Kind` attribute. `Kind` can easily be combined with [`where`](/functions/where/) in your templates to create kind-specific lists of content, but there are times where you may want to fetch the index page of a single section by the section's path.

[`.GetPage`](/function/getpage/) looks up an index page (i.e `_index.md`) of a given `Kind` and `path`. This method is only supported in section page templates but *may* support [single page templates][singlepages] in the future.

`.Site.GetPage` takes two arguments: `kind` and `kind value`.

The valid values for 'kind' are as follows:

1. `home`
2. `section`
3. `taxonomy`
4. `taxonomyTerm`

### `.Site.GetPage` Example

The `.Site.GetPage` example assumes the following project directory structure:

{{% input "grab-blog-section-index-page-title.html" %}}
{{ with .Site.GetPage "section" "blog" }}{{ .Title }}{{ end }}
{{% /input %}}

`.Site.GetPage` will return `nil` if no `_index.md` page is found. If `content/blog/_index.md` does not exist, the template will output a blank section where `{{.Title}}` should have been in the preceding example.

[sections]: /content-management/sections/
[directorystructure]: /getting-started/directory-structure/
[homepage]: /templates/homepage-template/
[limitkeyword]: https://www.techonthenet.com/sql/select_limit.php
[partials]: /templates/partial-templates/
[RSS 2.0]: http://cyber.law.harvard.edu/rss/rss.html "RSS 2.0 Specification"
[singlepages]: /templates/single-page-templates/