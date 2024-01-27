---
title: Weight
description: Returns the weight of the given page as defined in front matter.
categories: []
keywords: []
action:
  related: []
  returnType: int
  signatures: [PAGE.Weight]
---

The `Weight` method on a `Page` object returns the [weight] of the given page as defined in front matter.

[weight]: /getting-started/glossary/#weight

{{< code-toggle file=content/recipes/sushi.md fm=true >}}
title = 'How to make spicy tuna hand rolls'
weight = 42
{{< /code-toggle >}}

Page weight controls the position of a page within a collection that is sorted by weight. Assign weights using non-zero integers. Lighter items float to the top, while heavier items sink to the bottom. Unweighted or zero-weighted elements are placed at the end of the collection.

Although rarely used within a template, you can access the value with:

```go-html-template
{{ .Weight }} → 42
```
