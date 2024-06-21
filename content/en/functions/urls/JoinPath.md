---
title: urls.JoinPath
description: Joins the provided elements into a URL string and cleans the result of any ./ or ../ elements. If the argument list is empty, JoinPath returns an empty string.
categories: []
keywords: []
action:
  aliases: []
  related:
    - functions/path/Join
  returnType: string
  signatures: [urls.JoinPath ELEMENT...]
aliases: [/functions/urls.joinpath]
---

{{< new-in 0.112.0 >}}

```go-html-template
{{ urls.JoinPath }} → "" (empty string)
{{ urls.JoinPath "" }} → /
{{ urls.JoinPath "a" }} → a
{{ urls.JoinPath "a" "b" }} → a/b
{{ urls.JoinPath "/a" "b" }} → /a/b
{{ urls.JoinPath "https://example.org" "b" }} → https://example.org/b

{{ urls.JoinPath (slice "a" "b") }} → a/b
```

Unlike the [`path.Join`] function, `urls.JoinPath` retains consecutive leading slashes.

[`path.Join`]: /functions/path/join/
