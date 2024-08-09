---
title: Paginate
description: Paginates a collection of pages.
categories: []
keywords: []
action:
  related:
    - methods/page/Paginator
  returnType: page.Pager
  signatures: ['PAGE.Paginate COLLECTION [N]']
---

[Pagination] is the process of splitting a list page into two or more pagers, where each pager contains a subset of the page collection and navigation links to other pagers.

By default, the number of elements on each pager is determined by your [site configuration]. The default is `10`. Override that value by providing a second argument, an integer, when calling the `Paginate` method.

[site configuration]: /getting-started/configuration/#pagination

{{% note %}}
There is also a `Paginator` method on `Page` objects, but it can neither filter nor sort the page collection.

The `Paginate` method is more flexible.
{{% /note %}}

You can invoke pagination on the [home template], [section templates], [taxonomy templates], and [term templates].

[home template]: /templates/types/#home
[section templates]: /templates/types/#section
[taxonomy templates]: /templates/types/#taxonomy
[term templates]: /templates/types/#term

{{< code file=layouts/_default/list.html >}}
{{ $pages := where .Site.RegularPages "Section" "articles" }}
{{ $pages = $pages.ByTitle }}
{{ range (.Paginate $pages 7).Pages }}
  <h2><a href="{{ .RelPermalink }}">{{ .Title }}</a></h2>
{{ end }}
{{ template "_internal/pagination.html" . }}
{{< /code >}}

In the example above, we:

1. Build a page collection
2. Sort the collection by title
3. Paginate the collection, with 7 elements per pager
4. Range over the paginated page collection, rendering a link to each page
5. Call the embedded pagination template to create navigation links between pagers

{{% note %}}
Please note that the results of pagination are cached. Once you have invoked either the `Paginator` or `Paginate` method, the paginated collection is immutable. Additional invocations of these methods will have no effect.
{{% /note %}}
