---
title: Pagination
description: Split a list page into two or more subsets.
categories: [templates]
keywords: []
menu:
  docs:
    parent: templates
    weight: 190
weight: 190
toc: true
aliases: [/extras/pagination,/doc/pagination/]
---

Displaying a large page collection on a list page is not user-friendly:

- A massive list can be intimidating and difficult to navigate. Users may get lost in the sheer volume of information.
- Large pages take longer to load, which can frustrate users and lead to them abandoning the site.
- Without any filtering or organization, finding a specific item becomes a tedious scrolling exercise.

Improve usability by paginating `home`, `section`, `taxonomy`, and `term` pages.

{{% note %}}
The most common templating mistake related to pagination is invoking pagination more than once for a given list page. See the [caching](#caching) section below.
{{% /note %}}

## Terminology

paginate
: To split a [list page] into two or more subsets.

pagination
: The process of paginating a list page.

pager
: Created during pagination, a pager contains a subset of a list page and navigation links to other pagers.

paginator
: A collection of pagers.

[list page]: /getting-started/glossary/#list-page

## Configuration

Control pagination behavior in your site configuration. These are the default settings:

{{< code-toggle file=hugo config=pagination />}}

disableAliases
: (`bool`) Whether to disable alias generation for the first pager. Default is `false`.

pagerSize
: (`int`) The number of pages per pager. Default is `10`.

path
: (`string`) The segment of each pager URL indicating that the target page is a pager. Default is `page`.

With multilingual sites you can define the pagination behavior for each language:

{{< code-toggle file=hugo >}}
[languages.en]
contentDir = 'content/en'
languageCode = 'en-US'
languageDirection = 'ltr'
languageName = 'English'
weight = 1
[languages.en.pagination]
disableAliases = true
pagerSize = 10
path = 'page'
[languages.de]
contentDir = 'content/de'
languageCode = 'de-DE'
languageDirection = 'ltr'
languageName = 'Deutsch'
weight = 2
[languages.de.pagination]
disableAliases = true
pagerSize = 20
path = 'blatt'
{{< /code-toggle >}}

## Methods

To paginate a `home`, `section`, `taxonomy`, or `term` page, invoke either of these methods on the `Page` object in the corresponding template:

- [`Paginate`]
- [`Paginator`]

The `Paginate` method is more flexible, allowing you to:

- Paginate any page collection
- Filter, sort, and group the page collection
- Override the number of pages per pager as defined in your site configuration

By comparison, the `Paginator` method paginates the page collection passed into the template, and you cannot override the number of pages per pager.

[`Paginate`]: /methods/page/paginate/
[`Paginator`]: /methods/page/paginator/

## Examples

To paginate a list page using the `Paginate` method:

```go-html-template
{{ $pages := where site.RegularPages "Type" "posts" }}
{{ $paginator := .Paginate $pages.ByTitle 7 }}

{{ range $paginator.Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}

{{ template "_internal/pagination.html" . }}
```

In the example above, we:

1. Build a page collection
2. Sort the page collection by title
3. Paginate the page collection, with 7 pages per pager
4. Range over the paginated page collection, rendering a link to each page
5. Call the embedded pagination template to create navigation links between pagers


To paginate a list page using the `Paginator` method:

```go-html-template
{{ range .Paginator.Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}

{{ template "_internal/pagination.html" . }}
```

In the example above, we:

1. Paginate the page collection passed into the template, with the default number of pages per pager
2. Range over the paginated page collection, rendering a link to each page
3. Call the embedded pagination template to create navigation links between pagers

## Caching

{{% note %}}
The most common templating mistake related to pagination is invoking pagination more than once for a given list page.
{{% /note %}}

Regardless of pagination method, the initial invocation is cached and cannot be changed. If you invoke pagination more than once for a given list page, subsequent invocations use the cached result. This means that subsequent invocations will not behave as written.

When paginating conditionally, do not use the `compare.Conditional` function due to its eager evaluation of arguments. Use an `if-else` construct instead.

[`compare.Conditional`]: /functions/compare/conditional/

## Grouping

Use pagination with any of the [grouping methods]. For example:

[grouping methods]: /quick-reference/page-collections/#group

```go-html-template
{{ $pages := where site.RegularPages "Type" "posts" }}
{{ $paginator := .Paginate ($pages.GroupByDate "Jan 2006") }}

{{ range $paginator.PageGroups }}
  <h2>{{ .Key }}</h2>
  {{ range .Pages }}
    <h3><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h3>
  {{ end }}
{{ end }}

{{ template "_internal/pagination.html" . }}
```

[grouping methods]: /quick-reference/page-collections/#group

## Navigation

As shown in the examples above, the easiest way to add navigation between pagers is with Hugo's embedded pagination template:

```go-html-template
{{ template "_internal/pagination.html" . }}
```

The embedded pagination template has two formats: `default` and `terse`. The above is equivalent to:

```go-html-template
{{ template "_internal/pagination.html" (dict "page" . "format" "default") }}
```

The `terse` format has fewer controls and page slots, consuming less space when styled as a horizontal list. To use the `terse` format:

```go-html-template
{{ template "_internal/pagination.html" (dict "page" . "format" "terse") }}
```

{{% note %}}
To override Hugo's embedded pagination template, copy the [source code] to a file with the same name in the layouts/partials directory, then call it from your templates using the [`partial`] function:

`{{ partial "pagination.html" . }}`

[`partial`]: /functions/partials/include/
[source code]: {{% eturl pagination %}}
{{% /note %}}

Create custom navigation components using any of the `Pager` methods:

{{< list-pages-in-section path=/methods/pager >}}

## Structure

The example below depicts the published site structure when paginating a list page.

With this content:

```text
content/
├── posts/
│   ├── _index.md
│   ├── post-1.md
│   ├── post-2.md
│   ├── post-3.md
│   └── post-4.md
└── _index.md
```

And this site configuration:

{{< code-toggle file=hugo >}}
[pagination]
  disableAliases = false
  pagerSize = 2
  path = 'page'
{{< /code-toggle >}}

And this section template:

```go-html-template
{{ range (.Paginate .Pages).Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}

{{ template "_internal/pagination.html" . }}
```

The published site has this structure:

```text
public/
├── posts/
│   ├── page/
│   │   ├── 1/
│   │   │   └── index.html  <-- alias to public/posts/index.html
│   │   └── 2/
│   │       └── index.html
│   ├── post-1/
│   │   └── index.html
│   ├── post-2/
│   │   └── index.html
│   ├── post-3/
│   │   └── index.html
│   ├── post-4/
│   │   └── index.html
│   └── index.html
└── index.html
```

To disable alias generation for the first pager, change your site configuration:

{{< code-toggle file=hugo >}}
[pagination]
  disableAliases = true
  pagerSize = 2
  path = 'page'
{{< /code-toggle >}}

Now the published site will have this structure:

```text
public/
├── posts/
│   ├── page/
│   │   └── 2/
│   │       └── index.html
│   ├── post-1/
│   │   └── index.html
│   ├── post-2/
│   │   └── index.html
│   ├── post-3/
│   │   └── index.html
│   ├── post-4/
│   │   └── index.html
│   └── index.html
└── index.html
```
