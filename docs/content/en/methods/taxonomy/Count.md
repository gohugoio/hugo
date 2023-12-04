---
title: Count
description: Returns the number of number of weighted pages to which the given term has been assigned.
categories: []
keywords: []
action:
  related: []
  returnType: int
  signatures: [TAXONOMY.Count TERM]
toc: true
---

The `Count` method on a `Taxonomy` object returns the number of number of [weighted pages] to which the given [term] has been assigned.

{{% include "methods/taxonomy/_common/get-a-taxonomy-object.md" %}}

## Count the weighted pages

Now that we have captured the "genres" `Taxonomy` object, let's count the number of weighted pages to which the "suspense" term has been assigned:

```go-html-template
{{ $taxonomyObject.Count "suspense" }} → 3
```

[weighted pages]: /getting-started/glossary/#weighted-page
[term]: /getting-started/glossary/#term
