---
aliases:
- /doc/pagination/
date: 2014-01-01
menu:
  main:
    parent: extras
next: /extras/scratch
prev: /extras/shortcodes
title: Pagination
weight: 80
---

Hugo supports pagination for the home page, sections and taxonomies. It's built to be easy use, but with loads of flexibility when needed. The real power shines when you combine it with [`where`](/templates/functions/), with its SQL-like operators, `first` and others --- you can even [order the content](/templates/list/) the way you've become used to with Hugo.

## Configuration

Pagination can be configured in the site configuration (e.g. `config.toml`):

* `Paginate` (default `10`) 
* `PaginatePath` (default `page`)

Setting `Paginate` to a positive value will split the list pages for the home page, sections and taxonomies into chunks of that size. But note that the generation of the pagination pages for sections, taxonomies and home page is *lazy* --- the pages will not be created if not referenced by a `.Paginator` (see below).

 `PaginatePath` is used to adapt the `Url` to the pages in the paginator (the default setting will produce urls on the form `/page/1/`. 

## List the pages

**A `.Paginator` is provided to help building a pager menu. This is only relevant for the templates for the home page and the list pages (sections and taxonomies).**  

There are two ways to configure and use a `.Paginator`:

1. The simplest way is just to call `.Paginator.Pages` from a template. It will contain the pages for *that page* .
2. Select a sub-set of the pages with the available template functions and ordering options, and pass the slice to `.Paginate`, e.g. `{{ range (.Paginate ( first 50 .Data.Pages.ByTitle )).Pages }}`.

For a given **Node**, it's one of the options above. The `.Paginator` is static and cannot change once created.

## Build the navigation

The `.Paginator` contains enough information to build a paginator interface. 

The easiest way to add this to your pages is to include the built-in template (with `Bootstrap`-compatible styles):

```
{{ template "_internal/pagination.html" . }}
```

**Note:** If you use any filters or ordering functions to create your `.Paginator` **and** you want the navigation buttons to be shown before the page listing, you must create the `.Paginator` before it's used:

```
{{ $paginator := .Paginate (where .Data.Pages "Type" "post") }}
{{ template "_internal/pagination.html" . }}
{{ range $paginator.Pages }}
   {{ .Title }}
{{ end }}
```

Without the where-filter, the above is simpler:

```
{{ template "_internal/pagination.html" . }}
{{ range .Paginator.Pages }}
   {{ .Title }}
{{ end }}
```

If you want to build custom navigation, you can do so using the `.Paginator` object:

* `PageNumber`: The current page's number in the pager sequence
* `Url`: The relative Url to the current pager
* `Pages`: The pages in the current pager
* `NumberOfElements`: The number of elements on this page
* `HasPrev`: Whether there are page(s) before the current
* `Prev`: The pager for the previous page
* `HasNext`: Whether there are page(s) after the current
* `Next`: The pager for the next page
* `First`: The pager for the first page
* `Last`: The pager for the last page
* `Pagers`: A list of pagers that can be used to build a pagination menu
* `PageSize`: Size of each pager
* `TotalPages`: The number of pages in the paginator
* `TotalNumberOfElements`: The number of elements on all pages in this paginator

## Additional information

The pages are built on the following form (`BLANK` means no value):

```
[SECTION/TAXONOMY/BLANK]/index.html
[SECTION/TAXONOMY/BLANK]/page/1/index.html => redirect to  [SECTION/TAXONOMY/BLANK]/index.html
[SECTION/TAXONOMY/BLANK]/page/2/index.html
....
```


