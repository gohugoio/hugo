---
title: Ordere and Grouping Hugo Lists
linktitle: List Ordering and Grouping
description: You can group or order your content in both your templating and content front matter.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
keywords: []
menu:
  docs:
    parent: "templates"
    weight: 27
weight: 27
sections_weight: 27
draft: true
aliases: [/templates/ordering/,/templates/grouping/]
toc: true
notes: This was originally going to be a separate page on the new docs site but it now makes more sense to keep everything within the templates/lists page. - rdwatters, 2017-03-12.
---

In Hugo, A list template is any template that will be used to render multiple pieces of content in a single HTML page.

## Example List Templates

### Section Template

This list template is used for [spf13.com](https://spf13.com/). It makes use of [partial templates][partials]. All examples use a [view](/templates/views/) called either "li" or "summary."

{{< code file="layouts/section/post.html" >}}
{{ partial "header.html" . }}
{{ partial "subheader.html" . }}

<section id="main">
  <div>
   <h1 id="title">{{ .Title }}</h1>
        <ul id="list">
            {{ range .Pages }}
                {{ .Render "li"}}
            {{ end }}
        </ul>
  </div>
</section>
{{ partial "footer.html" . }}
{{< /code >}}

### Taxonomy Template

{{< code file="layouts/_default/taxonomies.html" download="taxonomies.html" >}}
{{ define "main" }}
<section id="main">
  <div>
   <h1 id="title">{{ .Title }}</h1>
    {{ range .Pages }}
        {{ .Render "summary"}}
    {{ end }}
  </div>
</section>
{{ end }}
{{< /code >}}

## Order Content

Hugo lists render the content based on metadata provided in the [front matter](/content-management/front-matter/)..

Here are a variety of different ways you can order the content items in
your list templates:

### Default: Weight > Date

{{< code file="layouts/partials/order-default.html" >}}
<ul class="pages">
    {{ range .Pages }}
        <li>
            <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
            <time>{{ .Date.Format "Mon, Jan 2, 2006" }}</time>
        </li>
    {{ end }}
</ul>
{{< /code >}}

### By Weight

{{< code file="layouts/partials/by-weight.html" >}}
{{ range .Pages.ByWeight }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
{{ end }}
{{< /code >}}

### By Date

{{< code file="layouts/partials/by-date.html" >}}
{{ range .Pages.ByDate }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
{{ end }}
{{< /code >}}

### By Publish Date

{{< code file="layouts/partials/by-publish-date.html" >}}
{{ range .Pages.ByPublishDate }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .PublishDate.Format "Mon, Jan 2, 2006" }}</div>
    </li>
{{ end }}
{{< /code >}}

### By Expiration Date

{{< code file="layouts/partials/by-expiry-date.html" >}}
{{ range .Pages.ByExpiryDate }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .ExpiryDate.Format "Mon, Jan 2, 2006" }}</div>
    </li>
{{ end }}
{{< /code >}}

### By Last Modified Date

{{< code file="layouts/partials/by-last-mod.html" >}}
{{ range .Pages.ByLastmod }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
{{ end }}
{{< /code >}}

### By Length

{{< code file="layouts/partials/by-length.html" >}}
{{ range .Pages.ByLength }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
{{ end }}
{{< /code >}}


### By Title

{{< code file="layouts/partials/by-title.html" >}}
{{ range .Pages.ByTitle }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
{{ end }}
{{< /code >}}

### By Link Title

{{< code file="layouts/partials/by-link-title.html" >}}
{{ range .Pages.ByLinkTitle }}
    <li>
    <a href="{{ .Permalink }}">{{ .LinkTitle }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
{{ end }}
{{< /code >}}

### By Parameter

Order based on the specified front matter parameter. Content that does not have the specified front matter field  will use the site's `.Site.Params` default. If the parameter is not found at all in some entries, those entries will appear together at the end of the ordering.

The below example sorts a list of posts by their rating.

{{< code file="layouts/partials/by-rating.html" >}}
{{ range (.Pages.ByParam "rating") }}
  <!-- ... -->
{{ end }}
{{< /code >}}

If the front matter field of interest is nested beneath another field, you can
also get it:

{{< code file="layouts/partials/by-nested-param.html" >}}
{{ range (.Pages.ByParam "author.last_name") }}
  <!-- ... -->
{{ end }}
{{< /code >}}

### Reverse Order

Reversing order can be applied to any of the above methods. The following uses `ByDate` as an example:

{{< code file="layouts/partials/by-date-reverse.html" >}}
{{ range .Pages.ByDate.Reverse }}
<li>
<a href="{{ .Permalink }}">{{ .Title }}</a>
<div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
</li>
{{ end }}
{{< /code >}}

## Group Content

Hugo provides some functions for grouping pages by Section, Type, Date, etc.

### By Page Field

{{< code file="layouts/partials/by-page-field.html" >}}
{{ range .Pages.GroupBy "Section" }}
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
{{< /code >}}

### By Page date

{{< code file="layouts/partials/by-page-date.html" >}}
{{ range .Pages.GroupByDate "2006-01" }}
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
{{< /code >}}

### By Page publish date

{{< code file="layouts/partials/by-page-publish-date.html" >}}
{{ range .Pages.GroupByPublishDate "2006-01" }}
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
{{< /code >}}

### By Page Param

{{< code file="layouts/partials/by-page-param.html" >}}
{{ range .Pages.GroupByParam "param_key" }}
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
{{< /code >}}

### By Page Param in Date Format

{{< code file="layouts/partials/by-page-param-as-date.html" >}}
{{ range .Pages.GroupByParamDate "param_key" "2006-01" }}
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
{{< /code >}}

### Reverse Key Order

The ordering of the groups is performed by keys in alphanumeric order (A–Z, 1–100) and in reverse chronological order (newest first) for dates.

While these are logical defaults, they are not always the desired order. There are two different syntaxes to change the order, both of which work the same way. You can use your preferred syntax.

#### Reverse Method

```
{{ range (.Pages.GroupBy "Section").Reverse }}
```

```
{{ range (.Pages.GroupByDate "2006-01").Reverse }}
```


#### Provide the Alternate Direction

```
{{ range .Pages.GroupByDate "2006-01" "asc" }}
```

```
{{ range .Pages.GroupBy "Section" "desc" }}
```

### Order Within Groups

Because Grouping returns a `{{.Key}}` and a slice of pages, all of the ordering methods listed above are available.

In the following example, groups are ordered chronologically and then content
within each group is ordered alphabetically by title.

{{< code file="layouts/partials/by-group-by-page.html" >}}
{{ range .Pages.GroupByDate "2006-01" "asc" }}
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
{{< /code >}}

## Filter and Limiting Lists

See the [_Lists/Filtering and Limiting Lists_
section][filteringandlimitinglists] for details.


[views]: /templates/views/
[filteringandlimitinglists]: /templates/lists/#filtering-and-limiting-lists
