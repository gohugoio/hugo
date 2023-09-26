---
title: urls.JoinPath
description: Joins the provided elements into a URL string and cleans the result of any ./ or ../ elements. If the argument list is empty, JoinPath returns an empty string.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: urls
relatedFuncs:
  - path.Join
  - urls.JoinPath
signature: 
  - urls.JoinPath ELEMENT...
---

```go-html-template
{{ urls.JoinPath }} → ""
{{ urls.JoinPath "" }} → "/"
{{ urls.JoinPath "a" }} → "a"
{{ urls.JoinPath "a" "b" }} → "a/b"
{{ urls.JoinPath "/a" "b" }} → "/a/b"
{{ urls.JoinPath "https://example.org" "b" }} → "https://example.org/b"

{{ urls.JoinPath (slice "a" "b") }} → "a/b"
```

Unlike the [`path.Join`] function, `urls.JoinPath` retains consecutive leading slashes.

[`path.Join`]: /functions/path.join/
