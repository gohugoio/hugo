---
title: Ordering and Grouping Lists
linktitle: Rendering Hugo Lists
description: Hugo assumes that the same structure that works to organize your source content is used to organize the rendered site, but
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: []
weight: 27
draft: false
aliases: [/templates/ordering/,/templates/grouping/]
toc: true
wip: true
---

![Image demonstrating a hierarchical website sitemap.](/images/site-hierarchy.svg)

## Understanding `.Data.Pages`

From this image, we can assume that the "homepage" for Section A---presumably, `/section-a/index.html`---is going to list the content pages 1,2,3. In this way, pages 1,2,3 are *data* made available to the template that renders to the .

## Example List Template Pages

### Example Section Template: `post.html`

This content template is used for [spf13.com](http://spf13.com/). It makes use of [partial templates][partials]. All examples use a [view](/templates/views/) called either "li" or "summary" which this example site defined.

{{% code file="layouts/section/post.html" %}}
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
{{% /code %}}

### Example Taxonomy Template

This content template is used for [spf13.com](http://spf13.com/). It makes use of [partial templates](/templates/partials/). All examples use a [view](/templates/views/) called either "li" or "summary" which this example site defined.

{{% code file="layouts/_default/taxonomies.html" download="taxonomies.html" %}}
```html
{{ define "main" }}
<section id="main">
  <div>
   <h1 id="title">{{ .Title }}</h1>
    {{ range .Data.Pages }}
        {{ .Render "summary"}}
    {{ end }}
  </div>
</section>
{{ end }}
```
{{% /code %}}

## Ordering Content

In the case of Hugo, each list will render the content based on metadata provided in the [front matter](/content/front-matter/). See [ordering content](/content/ordering/) for more information.

Here are a variety of different ways you can order the content items in
your list templates:

### Default Ordering: Weight > Date

{{% code file="layouts/_default/list.html" %}}
```html
<ul class="pages">
    {{ range .Data.Pages }}
        <li>
            <h1><a href="{{ .Permalink }}">{{ .Title }}</a></h1>
            <time>{{ .Date.Format "Mon, Jan 2, 2006" }}</time>
        </li>
    {{ end }}
</ul>
```
{{% /code %}}

### Ordering a list by Weight -> Date

```html
{{ range .Data.Pages.ByWeight }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
{{ end }}
```

### Ordering a List by Date

    {{ range .Data.Pages.ByDate }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

### Ordering a List by PublishDate

    {{ range .Data.Pages.ByPublishDate }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .PublishDate.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

### Ordering a list by ExpiryDate

    {{ range .Data.Pages.ByExpiryDate }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .ExpiryDate.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

### Ordering a List by Lastmod

    {{ range .Data.Pages.ByLastmod }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

### Ordering a List by Length

    {{ range .Data.Pages.ByLength }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}


### Ordering a List by Title

    {{ range .Data.Pages.ByTitle }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

### Ordering a List by LinkTitle

    {{ range .Data.Pages.ByLinkTitle }}
    <li>
    <a href="{{ .Permalink }}">{{ .LinkTitle }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

### Order List by Parameter

Order based on the specified front matter parameter. Pages without that
parameter will use the site's `.Site.Params` default. If the parameter is not
found at all in some entries, those entries will appear together at the end
of the ordering.

The below example sorts a list of posts by their rating.

    {{ range (.Data.Pages.ByParam "rating") }}
      <!-- ... -->
    {{ end }}

If the front matter field of interest is nested beneath another field, you can
also get it:

```
{{ range (.Date.Pages.ByParam "author.last_name") }}
  <!-- ... -->
{{ end }}
```

### Reverse Order
Can be applied to any of the above. Using Date for an example.

```
{{ range .Data.Pages.ByDate.Reverse }}
<li>
<a href="{{ .Permalink }}">{{ .Title }}</a>
<div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
</li>
{{ end }}
```

## Grouping Content

Hugo provides some grouping functions for list pages. You can use them to
group pages by Section, Type, Date etc.

Here are a variety of different ways you can group the content items in
your list templates:

### Grouping by Page field

```
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
```

### Grouping by Page date

```
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
```

### Grouping by Page publish date

```
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
```

### Grouping by Page param

```html
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
```

### Grouping by Page param in date format

```html
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
```

### Reversing Key Order

The ordering of the groups is performed by keys in alphanumeric order (A–Z,
1–100) and in reverse chronological order (newest first) for dates.

While these are logical defaults, they are not always the desired order. There
are two different syntaxes to change the order; they both work the same way, so
it’s really just a matter of preference.

#### Reverse method

```golang
{{ range (.Data.Pages.GroupBy "Section").Reverse }}
```

```golang
{{ range (.Data.Pages.GroupByDate "2006-01").Reverse }}
```


#### Providing the (alternate) direction

```golang
{{ range .Data.Pages.GroupByDate "2006-01" "asc" }}
```

```golang
{{ range .Data.Pages.GroupBy "Section" "desc" }}
```

### Ordering Pages within Group

Because Grouping returns a key and a slice of pages, all of the ordering methods listed above are available.

In this example, I’ve ordered the groups in chronological ordering and the content
within each group in alphabetical order by title.

```html
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
```

## Filtering and Limiting List Content

Sometimes you only want to list a subset of the available content. A common request is to only display “Posts” on the homepage. Using the `where` function, you can do just that.

### `where`

`where` works in a similar manner to the `where` keyword in SQL. It selects all elements of the array or slice that match the provided field and value. `where` takes three arguments:

1. `array` or a `slice of maps or structs`
2. `key` or `field name'
3. `match value`

{{% code file="layouts/_default/.html" %}}
```html
{{ range where .Data.Pages "Section" "post" }}
   {{ .Content }}
{{ end }}
```
{{% /code %}}

### `first`

`first` works in a similar manner to the [`limit` keyword in SQL][limitkeyword]. It reduces the array to only the `first N` elements. It takes the array and number of elements as input. `first` takes two arguments:

1. `array` or `slice of maps or structs`
2. `number of elements`

{{% code file="layout/_default/section.html" %}}
```golang
{{ range first 10 .Data.Pages }}
  {{ .Render "summary" }}
{{ end }}
```
{{% /code %}}

### `first` and `where` Together

Using `first` and `where` together can be very powerful:

{{% code file="first-and-where-together.html" %}}
```golang
{{ range first 5 (where .Data.Pages "Section" "post") }}
   {{ .Content }}
{{ end }}
```
{{% /code %}}
