---
title: Configure page
linkTitle: Page
description: Configure page behavior.
categories: []
keywords: []
---

{{% glossary-term "default sort order" %}}

Hugo uses the default sort order to determine the _next_ and _previous_ page relative to the current page when calling these methods on a `Page` object:

- [`Next`][] and [`Prev`][]
- [`NextInSection`][] and [`PrevInSection`][]

This is based on this default project configuration:

{{< code-toggle config=page />}}

`nextPrevInSectionSortOrder`
: (`string`) The sort order used to determine the _next_ and _previous_ page within the same section when calling [`NextInSection`][] or [`PrevInSection`][] on a `Page` object. Valid values are `asc` (ascending) or `desc` (descending). Default is `desc`.

`nextPrevSortOrder`
: (`string`) The sort order used to determine the _next_ and _previous_ page when calling [`Next`][] or [`Prev`][] on a `Page` object. Valid values are `asc` (ascending) or `desc` (descending). Default is `desc`.

To reverse the meaning of _next_ and _previous_:

{{< code-toggle file=hugo >}}
[page]
  nextPrevInSectionSortOrder = 'asc'
  nextPrevSortOrder = 'asc'
{{< /code-toggle >}}

> [!NOTE]
> These settings do not apply to the [`Next`][next-pages] or [`Prev`][prev-pages] methods on a `Pages` object.

[`NextInSection`]: /methods/page/nextinsection/
[`Next`]: /methods/page/next/
[`PrevInSection`]: /methods/page/previnsection/
[`Prev`]: /methods/page/prev/
[next-pages]: /methods/pages/next/
[prev-pages]: /methods/pages/prev/
