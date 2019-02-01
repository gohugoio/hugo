---
title: Pagination
linktitle: Pagination
description: Hugo supports pagination for your homepage, section pages, and taxonomies.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
keywords: [lists,sections,pagination]
menu:
  docs:
    parent: "templates"
    weight: 140
weight: 140
sections_weight: 140
draft: false
aliases: [/extras/pagination,/doc/pagination/]
toc: true
---

The real power of Hugo pagination shines when combined with the [`where` function][where] and its SQL-like operators: [`first`][], [`last`][], and [`after`][]. You can even [order the content][lists] the way you've become used to with Hugo.

## Configure Pagination

Pagination can be configured in your [site configuration][configuration]:

`Paginate`
: default = `10`. This setting can be overridden within the template.

`PaginatePath`
: default = `page`. Allows you to set a different path for your pagination pages.

Setting `Paginate` to a positive value will split the list pages for the homepage, sections and taxonomies into chunks of that size. But note that the generation of the pagination pages for sections, taxonomies and homepage is *lazy* --- the pages will not be created if not referenced by a `.Paginator` (see below).

`PaginatePath` is used to adapt the `URL` to the pages in the paginator (the default setting will produce URLs on the form `/page/1/`.

## List Paginator Pages

{{% warning %}}
`.Paginator` is provided to help you build a pager menu. This feature is currently only supported on homepage and list pages (i.e., taxonomies and section lists).
{{% /warning %}}

There are two ways to configure and use a `.Paginator`:

1. The simplest way is just to call `.Paginator.Pages` from a template. It will contain the pages for *that page*.
2. Select a subset of the pages with the available template functions and ordering options, and pass the slice to `.Paginate`, e.g. `{{ range (.Paginate ( first 50 .Pages.ByTitle )).Pages }}`.

For a given **Page**, it's one of the options above. The `.Paginator` is static and cannot change once created.

The global page size setting (`Paginate`) can be overridden by providing a positive integer as the last argument. The examples below will give five items per page:

* `{{ range (.Paginator 5).Pages }}`
* `{{ $paginator := .Paginate (where .Pages "Type" "posts") 5 }}`

It is also possible to use the `GroupBy` functions in combination with pagination:

```
{{ range (.Paginate (.Pages.GroupByDate "2006")).PageGroups  }}
```

## Build the navigation

The `.Paginator` contains enough information to build a paginator interface.

The easiest way to add this to your pages is to include the built-in template (with `Bootstrap`-compatible styles):

```
{{ template "_internal/pagination.html" . }}
```

{{% note "When to Create `.Paginator`" %}}
If you use any filters or ordering functions to create your `.Paginator` *and* you want the navigation buttons to be shown before the page listing, you must create the `.Paginator` before it's used.
{{% /note %}}

The following example shows how to create `.Paginator` before its used:

```
{{ $paginator := .Paginate (where .Pages "Type" "posts") }}
{{ template "_internal/pagination.html" . }}
{{ range $paginator.Pages }}
   {{ .Title }}
{{ end }}
```

Without the `where` filter, the above example is even simpler:

```
{{ template "_internal/pagination.html" . }}
{{ range .Paginator.Pages }}
   {{ .Title }}
{{ end }}
```

If you want to build custom navigation, you can do so using the `.Paginator` object, which includes the following properties:

`PageNumber`
: The current page's number in the pager sequence

`URL`
: The relative URL to the current pager

`Pages`
: The pages in the current pager

`NumberOfElements`
: The number of elements on this page

`HasPrev`
: Whether there are page(s) before the current

`Prev`
: The pager for the previous page

`HasNext`
: Whether there are page(s) after the current

`Next`
: The pager for the next page

`First`
: The pager for the first page

`Last`
: The pager for the last page

`Pagers`
: A list of pagers that can be used to build a pagination menu

`PageSize`
: Size of each pager

`TotalPages`
: The number of pages in the paginator

`TotalNumberOfElements`
: The number of elements on all pages in this paginator

## Additional information

The pages are built on the following form (`BLANK` means no value):

```
[SECTION/TAXONOMY/BLANK]/index.html
[SECTION/TAXONOMY/BLANK]/page/1/index.html => redirect to  [SECTION/TAXONOMY/BLANK]/index.html
[SECTION/TAXONOMY/BLANK]/page/2/index.html
....
```


[`first`]: /functions/first/
[`last`]: /functions/last/
[`after`]: /functions/after/
[configuration]: /getting-started/configuration/
[lists]: /templates/lists/
[where]: /functions/where/
