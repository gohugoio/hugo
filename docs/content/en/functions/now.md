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

For example, building your site on June 24, 2017, with the following templating:

```
<div>
    <small>&copy; {{ now.Format "2006"}}</small>
</div>
```

would produce the following:

```
<div>
    <small>&copy; 2017</small>
</div>
```

The above example uses the [`.Format` function](/functions/format), which page includes a full listing of date formatting using Go's layout string.

{{% note %}}
Older Hugo themes may still be using the obsolete Page’s `.Now` (uppercase with leading dot), which causes build error that looks like the following:

    ERROR ... Error while rendering "..." in "...": ...
    executing "..." at <.Now.Format>:
    can't evaluate field Now in type *hugolib.PageOutput

Be sure to use `now` (lowercase with _**no**_ leading dot) in your templating.
{{% /note %}}
