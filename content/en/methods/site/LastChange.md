---
title: LastChange
description: Returns the last modification date of site content.
categories: []
keywords: []
action:
  related: []
  returnType: time.Time
  signatures: [SITE.LastChange]
---

The `LastChange` method on a `Site` object returns a [`time.Time`] value. Use this with time [functions] and [methods]. For example:

```go-html-template
{{ .Site.LastChange | time.Format ":date_long" }} → October 16, 2023

```

[`time.Time`]: https://pkg.go.dev/time#Time
[functions]: /functions/time
[methods]: /methods/time
