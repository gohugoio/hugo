---
title: now
linktitle: now
description: Returns the current local time 
godocref: https://godoc.org/time#Time
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-04-30
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [dates,time]
signature: ["now"]
workson: []
hugoversion:
relatedfuncs: [Unix,dateFormat]
deprecated: false
aliases: []
---

See [`time.Time`](https://godoc.org/time#Time).

For example, building your site on June 24, 2017 with the following templating:

```
<div>
    <small>&copy; {{ now.Format "2006"}}</small>
</div>
```

Which will produce the following:

```
<div>
    <small>&copy; 2017</small>
</div>
```

The above example uses the [`.Format` function](/functions/format), which page includes a full listing of date formatting using Golang's layout string.

{{% note %}}
Older Hugo themes may use the deprecated `.Now` (uppercase). Be sure to use the lowercase `.now` in your templating.
{{% /note %}}
