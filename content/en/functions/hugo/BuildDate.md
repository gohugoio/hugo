---
title: hugo.BuildDate
description: Returns the compile date of the Hugo binary.
categories: []
keywords: []
action:
  aliases: []
  related: []
  returnType: string
  signatures: [hugo.BuildDate]
---

The `hugo.BuildDate` function returns the compile date of the Hugo binary, formatted per [RFC 3339].

[RFC 3339]: https://datatracker.ietf.org/doc/html/rfc3339

```go-html-template
{{ hugo.BuildDate }} → 2023-11-01T17:57:00Z
```
