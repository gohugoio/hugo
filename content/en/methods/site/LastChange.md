---
title: LastChange
description: Returns the last modification date of site content.
categories: []
keywords: []
action:
  related: []
  returnType: time.Time
  signatures: [SITE.LastChange]
expiryDate: 2025-02-19 # deprecated 2024-02-19
---

{{% deprecated-in 0.123.0 %}}
Use [`.Site.Lastmod`] instead.

[`.Site.Lastmod`]: /methods/site/lastmod/
{{% /deprecated-in %}}

The `LastChange` method on a `Site` object returns a [`time.Time`] value. Use this with time [functions] and [methods]. For example:

```go-html-template
{{ .Site.LastChange | time.Format ":date_long" }} → January 31, 2024

```

[`time.Time`]: https://pkg.go.dev/time#Time
[functions]: /functions/time/
[methods]: /methods/time/
