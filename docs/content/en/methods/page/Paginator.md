---
title: Paginator
description: Paginates the collection of regular pages received in context.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.Pager
    signatures: [PAGE.Paginator]
---

Pagination is the process of splitting a list page into two or more pagers, where each pager contains a subset of the page collection and navigation links to other pagers.

The number of elements on each pager is determined by your [site configuration]. The default is `10`.

You can invoke pagination in [home], [section], [taxonomy], and [term] templates. Each of these receives a collection of regular pages in [context](g). When you invoke the `Paginator` method, it paginates the page collection received in context.

```go-html-template {file="layouts/section.html"}
{{ range .Paginator.Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
{{ end }}
{{ partial "pagination.html" . }}
```

In the example above, the embedded pagination template creates navigation links between pagers.

> [!note]
> Although simple to invoke, with the `Paginator` method you can neither filter nor sort the page collection. It acts upon the page collection received in context.
>
> The [`Paginate`] method is more flexible, and strongly recommended.

> [!note]
> Please note that the results of pagination are cached. Once you have invoked either the `Paginator` or `Paginate` method, the paginated collection is immutable. Additional invocations of these methods will have no effect.

[home]: /templates/types/#home
[section]: /templates/types/#section
[site configuration]: /configuration/pagination/
[taxonomy]: /templates/types/#taxonomy
[term]: /templates/types/#term
[`Paginate`]: /methods/page/paginate/
