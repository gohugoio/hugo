---
title: Paginate
description: Paginates a collection of pages.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.Pager
    signatures: ['PAGE.Paginate COLLECTION [N]']
---

Pagination is the process of splitting a list page into two or more pagers, where each pager contains a subset of the page collection and navigation links to other pagers.

By default, the number of elements on each pager is determined by your [site configuration]. The default is `10`. Override that value by providing a second argument, an integer, when calling the `Paginate` method.

> [!note]
> There is also a `Paginator` method on `Page` objects, but it can neither filter nor sort the page collection.
>
> The `Paginate` method is more flexible.

You can invoke pagination in [home], [section], [taxonomy], and [term] templates.

```go-html-template {file="layouts/section.html"}
{{ $pages := where .Site.RegularPages "Section" "articles" }}
{{ $pages = $pages.ByTitle }}
{{ range (.Paginate $pages 7).Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .Title }}</a></h2>
{{ end }}
{{ partial "pagination.html" . }}
```

In the example above, we:

1. Build a page collection
1. Sort the collection by title
1. Paginate the collection, with 7 elements per pager
1. Range over the paginated page collection, rendering a link to each page
1. Call the embedded pagination template to create navigation links between pagers

> [!note]
> Please note that the results of pagination are cached. Once you have invoked either the `Paginator` or `Paginate` method, the paginated collection is immutable. Additional invocations of these methods will have no effect.

[home]: /templates/types/#home
[section]: /templates/types/#section
[site configuration]: /configuration/pagination/
[taxonomy]: /templates/types/#taxonomy
[term]: /templates/types/#term
