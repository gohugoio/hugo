---
title: transform.HTMLUnescape
linkTitle: htmlUnescape
description: Returns the given string with HTML escape codes un-escaped.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [htmlUnescape]
  returnType: string
  signatures: [transform.HTMLUnescape INPUT]
relatedFunctions:
  - transform.HTMLEscape
  - transform.HTMLUnescape
aliases: [/functions/htmlunescape]
---

Remember to pass the output of this to `safeHTML` if fully un-escaped characters are desired. Otherwise, the output will be escaped again as normal.

```go-html-template
{{ htmlUnescape "Hugo &amp; Caddy &gt; WordPress &amp; Apache" }} → "Hugo & Caddy > WordPress & Apache"
```
