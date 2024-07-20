---
title: WordCount
description: Returns the number of words in the content of the given page.
categories: []
keywords: []
action:
  related:
    - methods/page/FuzzyWordCount
    - methods/page/ReadingTime
  returnType: int
  signatures: [PAGE.WordCount]
---

```go-html-template
{{ .WordCount }} → 103
```

To round up to nearest multiple of 100, use the [`FuzzyWordCount`] method.

[`FuzzyWordCount`]: /methods/page/fuzzywordcount/
