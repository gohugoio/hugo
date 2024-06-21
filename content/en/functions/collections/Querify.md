---
title: collections.Querify
description: Returns a URL query string composed of the given key-value pairs.
categories: []
keywords: []
action:
  aliases: [querify]
  related:
    - functions/go-template/urlquery.md
  returnType: string
  signatures:
    - collections.Querify [VALUE...]
aliases: [/functions/querify]
---

Specify the key-value pairs as individual arguments, or as a slice. The following are equivalent:


```go-html-template
{{ collections.Querify "a" 1 "b" 2 }}
{{ collections.Querify (slice "a" 1 "b" 2) }}
```

To append a query string to a URL:

```go-html-template
{{ $qs := collections.Querify "a" 1 "b" 2 }}
{{ $href := printf "https://example.org?%s" $qs }}

<a href="{{ $href }}">Link</a>
```

Hugo renders this to:

```html
<a href="https://example.org?a=1&amp;b=2">Link</a>
```
