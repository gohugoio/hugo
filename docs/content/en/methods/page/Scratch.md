---
title: Scratch
description: Returns a "scratch pad" on the given page to store and manipulate data.
categories: []
keywords: []
action:
  related:
    - methods/page/Store
    - functions/collections/NewScratch
  returnType: maps.Scratch
  signatures: [PAGE.Scratch]
toc: true
aliases: [/extras/scratch/,/doc/scratch/,/functions/scratch]
---

The `Scratch` method on a `Page` object creates a [scratch pad] to store and manipulate data. To create a scratch pad that is not reset on server rebuilds, use the [`Store`] method instead.

To create a locally scoped scratch pad that is not attached to a `Page` object, use the [`newScratch`] function.

[`Store`]: /methods/page/store/
[`newScratch`]: /functions/collections/newscratch/
[scratch pad]: /getting-started/glossary/#scratch-pad

{{% include "methods/page/_common/scratch-methods.md" %}}

## Determinate values

The `Scratch` method is often used to set scratch pad values within a shortcode, a partial template called by a shortcode, or by a Markdown render hook. In all three cases, the scratch pad values are not determinate until Hugo renders the page content.

If you need to access a scratch pad value from a parent template, and the parent template has not yet rendered the page content, you can trigger content rendering by assigning the returned value to a [noop] variable:

[noop]: /getting-started/glossary/#noop

```go-html-template
{{ $noop := .Content }}
{{ .Store.Get "mykey" }}
```

You can also trigger content rendering with the `FuzzyWordCount`, `Len`, `Plain`, `PlainWords`, `ReadingTime`, `Summary`, `Truncated`, and `WordCount` methods. For example:

```go-html-template
{{ $noop := .WordCount }}
{{ .Store.Get "mykey" }}
```
