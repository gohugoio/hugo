---
title: Data inspection
linkTitle: Inspection
description: Use template functions to inspect values and data structures.
categories: [troubleshooting]
keywords: []
menu:
  docs:
    parent: troubleshooting
    weight: 40
weight: 40
---

Use the [`debug.Dump`] function to inspect a data structure:

```go-html-template
<pre>{{ debug.Dump .Params }}</pre>
```

```text
{
  "date": "2023-11-10T15:10:42-08:00",
  "draft": false,
  "iscjklanguage": false,
  "lastmod": "2023-11-10T15:10:42-08:00",
  "publishdate": "2023-11-10T15:10:42-08:00",
  "tags": [
    "foo",
    "bar"
  ],
  "title": "My first post"
}
```

Use the [`printf`] function (render) or [`warnf`] function (log to console) to inspect simple data structures. The layout string below displays both value and data type.

```go-html-template
{{ $value := 42 }}
{{ printf "%[1]v (%[1]T)" $value }} → 42 (int)
```

[`debug.Dump`]: /functions/debug/dump/
[`printf`]: /functions/fmt/printf/
[`warnf`]: /functions/fmt/warnf/
