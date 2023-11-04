---
title: urlquery
description: Returns the escaped value of the textual representation of its arguments in a form suitable for embedding in a URL query.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: []
  returnType: string
  signatures: ['urlquery INPUT [INPUT]...']
relatedFunctions:
  - collections.Querify
  - urlquery
aliases: [/functions/urlquery]
---

{{% readfile file="/functions/_common/go-template-functions.md" %}}

This template code:

```go-html-template
{{ $u := urlquery "https://" "example.com" | safeURL }}
<a href="https://example.org?url={{ $u }}">Link</a>
```

Is rendered to:

```html
<a href="https://example.org?url=https%3A%2F%2Fexample.com">Link</a>
```
