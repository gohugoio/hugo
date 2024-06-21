---
title: Lastmod
description: Returns the last modification date of site content.
categories: []
keywords: []
action:
  related: []
  returnType: time.Time
  signatures: [SITE.Lastmod]
---

{{< new-in 0.123.0 >}}

The `Lastmod` method on a `Site` object returns a [`time.Time`] value. Use this with time [functions] and [methods]. For example:

```go-html-template
{{ .Site.Lastmod | time.Format ":date_long" }} → January 31, 2024

```

[`time.Time`]: https://pkg.go.dev/time#Time
[functions]: /functions/time/
[methods]: /methods/time/
