---
title: PlainWords
description: Calls the Plain method, splits the result into a slice of words, and returns the slice.
categories: []
keywords: []
action:
  related:
    - methods/page/Content
    - methods/page/Summary
    - methods/page/ContentWithoutSummary
    - methods/page/RawContent
    - methods/page/Plain
    - methods/page/RenderShortcodes
  returnType: '[]string'
  signatures: [PAGE.PlainWords]
---

The `PlainWords` method on a `Page` object calls the [`Plain`] method, then uses Go's [`strings.Fields`] function to split the result into words.

{{% note %}}
_Fields splits the string s around each instance of one or more consecutive whitespace characters, as defined by [`unicode.IsSpace`], returning a slice of substrings of s or an empty slice if s contains only whitespace._

[`unicode.IsSpace`]: https://pkg.go.dev/unicode#IsSpace
{{% /note %}}

As a result, elements within the slice may contain leading or trailing punctuation.

```go-html-template
{{ .PlainWords }}
```

To determine the approximate number of unique words on a page:

```go-html-template
{{ .PlainWords | uniq }} → 42
```

[`Plain`]: /methods/page/plain/
[`strings.Fields`]: https://pkg.go.dev/strings#Fields
