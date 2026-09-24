---
title: urls.JoinPath
description: Returns a URL string created by joining the provided elements and cleaning the result of any ./ or ../ elements, or an empty string if the argument list is empty.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: string
    signatures: [urls.JoinPath ELEMENT...]
aliases: [/functions/urls.joinpath]
---

```go-html-template
{{ urls.JoinPath }} → "" (empty string)
{{ urls.JoinPath "" }} → /
{{ urls.JoinPath "a" }} → a
{{ urls.JoinPath "a" "b" }} → a/b
{{ urls.JoinPath "/a" "b" }} → /a/b
{{ urls.JoinPath "https://example.org" "b" }} → https://example.org/b

{{ urls.JoinPath (slice "a" "b") }} → a/b
```

Unlike the [`path.Join`][] function, `urls.JoinPath` retains consecutive leading slashes.

[`path.Join`]: /functions/path/join/
