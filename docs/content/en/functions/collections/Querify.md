---
title: collections.Querify
description: Returns a URL query string from the given map, slice, or sequence of key-value pairs.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [querify]
    returnType: string
    signatures: [collections.Querify MAP|SLICE|KEY VALUE...]
aliases: [/functions/querify]
---

Specify the key-value pairs as a map, a slice, or a sequence of scalar values. For example, the following are equivalent:

```go-html-template
{{ collections.Querify (dict "a" 1 "b" 2) }}
{{ collections.Querify (slice "a" 1 "b" 2) }}
{{ collections.Querify "a" 1 "b" 2 }}
```

To append a query string to a URL:

```go-html-template
{{ $qs := collections.Querify (dict "a" 1 "b" 2) }}
{{ $href := printf "https://example.org?%s" $qs }}

<a href="{{ $href }}">Link</a>
```

Hugo renders this to:

```html
<a href="https://example.org?a=1&amp;b=2">Link</a>
```

You can also pass in a map from your project configuration or front matter. For example:

{{< code-toggle file=content/example.md fm=true >}}
title = 'Example'
[params.query]
a = 1
b = 2
{{< /code-toggle >}}

```go-html-template
{{ collections.Querify .Params.query }}
```
