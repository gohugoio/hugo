---
title: Pagination
description: Hugo supports pagination for your homepage, section pages, and taxonomies.
categories: [templates]
keywords: [lists,sections,pagination]
menu:
  docs:
    parent: templates
    weight: 100
weight: 100
aliases: [/extras/pagination,/doc/pagination/]
toc: true
---

The real power of Hugo pagination shines when combined with the [`where`] function and its SQL-like operators: [`first`], [`last`], and [`after`]. You can even [order the content][lists] the way you've become used to with Hugo.

## Configure pagination

Pagination can be configured in your [site configuration][configuration]:

`paginate`
: default = `10`. This setting can be overridden within the template.

`paginatePath`
: default = `page`. Allows you to set a different path for your pagination pages.

Setting `paginate` to a positive value will split the list pages for the homepage, sections and taxonomies into chunks of that size. But note that the generation of the pagination pages for sections, taxonomies and homepage is *lazy* --- the pages will not be created if not referenced by a `.Paginator` (see below).

`paginatePath` is used to adapt the `URL` to the pages in the paginator (the default setting will produce URLs on the form `/page/1/`.

## List paginator pages

{{% note %}}
`.Paginator` is provided to help you build a pager menu. This feature is currently only supported on homepage and list pages (i.e., taxonomies and section lists).
{{% /note %}}

There are two ways to configure and use a `.Paginator`:

1. The simplest way is just to call `.Paginator.Pages` from a template. It will contain the pages for *that page*.
2. Select another set of pages with the available template functions and ordering options, and pass the slice to `.Paginate`, e.g.
  * `{{ range (.Paginate ( first 50 .Pages.ByTitle )).Pages }}` or
  * `{{ range (.Paginate .RegularPagesRecursive).Pages }}`.

For a given **Page**, it's one of the options above. The `.Paginator` is static and cannot change once created.

If you call `.Paginator` or `.Paginate` multiple times on the same page, you should ensure all the calls are identical. Once *either* `.Paginator` or `.Paginate` is called while generating a page, its result is cached, and any subsequent similar call will reuse the cached result. This means that any such calls which do not match the first one will not behave as written.

(Remember that function arguments are eagerly evaluated, so a call like `$paginator := cond x .Paginator (.Paginate .RegularPagesRecursive)` is an example of what you should *not* do. Use `if`/`else` instead to ensure exactly one evaluation.)

The global page size setting (`Paginate`) can be overridden by providing a positive integer as the last argument. The examples below will give five items per page:

* `{{ range (.Paginator 5).Pages }}`
* `{{ $paginator := .Paginate (where .Pages "Type" "posts") 5 }}`

It is also possible to use the `GroupBy` functions in combination with pagination:

```go-html-template
{{ range (.Paginate (.Pages.GroupByDate "2006")).PageGroups }}
```

## Build the navigation

The `.Paginator` contains enough information to build a paginator interface.

The easiest way to add this to your pages is to include the built-in template (with `Bootstrap`-compatible styles):

```go-html-template
{{ template "_internal/pagination.html" . }}
```

{{% note %}}
If you use any filters or ordering functions to create your `.Paginator` *and* you want the navigation buttons to be shown before the page listing, you must create the `.Paginator` before it's used.
{{% /note %}}

The following example shows how to create `.Paginator` before its used:

```go-html-template
{{ $paginator := .Paginate (where .Pages "Type" "posts") }}
{{ template "_internal/pagination.html" . }}
{{ range $paginator.Pages }}
   {{ .Title }}
{{ end }}
```

Without the `where` filter, the above example is even simpler:

```go-html-template
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

```txt
[SECTION/TAXONOMY/BLANK]/index.html
[SECTION/TAXONOMY/BLANK]/page/1/index.html => redirect to  [SECTION/TAXONOMY/BLANK]/index.html
[SECTION/TAXONOMY/BLANK]/page/2/index.html
....
```

[`first`]: /functions/collections/first/
[`last`]: /functions/collections/last/
[`after`]: /functions/collections/after/
[configuration]: /getting-started/configuration/
[lists]: /templates/lists/
[`where`]: /functions/collections/where
