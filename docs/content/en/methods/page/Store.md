---
title: Store
description: Returns a persistent data structure for storing and manipulating keyed values, scoped to the current page.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: maps.Scratch
    signatures: [PAGE.Store]
aliases: [/functions/store/,/extras/scratch/,/doc/scratch/,/functions/scratch]
---

Use the `Store` method on a `Page` object to create a persistent data structure for storing and manipulating keyed values, scoped to the current page. To create a data structure with a different [scope](g), refer to the [scope](#scope) section below.

{{% include "_common/store-methods.md" %}}

{{% include "_common/store-scope.md" %}}

## Determinate values

The `Store` method is often used to set values within a _shortcode_ template, a _partial_ template called by a _shortcode_ template, or by a _render hook_ template. In all three cases, the stored values are indeterminate until Hugo renders the page content.

If you need to access a stored value from a parent template, and the parent template has not yet rendered the page content, you can trigger content rendering by assigning the returned value to a [noop](g) variable:

```go-html-template
{{ $noop := .Content }}
{{ .Store.Get "mykey" }}
```

You can also trigger content rendering with the `ContentWithoutSummary`, `FuzzyWordCount`, `Len`, `Plain`, `PlainWords`, `ReadingTime`, `Summary`, `Truncated`, and `WordCount` methods. For example:

```go-html-template
{{ $noop := .WordCount }}
{{ .Store.Get "mykey" }}
```
