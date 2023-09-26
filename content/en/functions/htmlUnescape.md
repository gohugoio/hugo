---
title: htmlUnescape
description: Returns the given string with HTML escape codes un-escaped.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: transform
relatedFuncs:
  - transform.HTMLEscape
  - transform.HTMLUnescape
signature:
  - transform.HTMLUnescape INPUT
  - htmlUnescape INPUT
---

`htmlUnescape` returns the given string with HTML escape codes un-escaped.

Remember to pass the output of this to `safeHTML` if fully un-escaped characters are desired. Otherwise, the output will be escaped again as normal.

```go-html-template
{{ htmlUnescape "Hugo &amp; Caddy &gt; WordPress &amp; Apache" }} → "Hugo & Caddy > WordPress & Apache"
```
