---
title: time.Now
linkTitle: now
description: Returns the current local time 
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [now]
  returnType: time.Time
  signatures: [time.Now]
relatedFunctions:
  - time.AsTime
  - time.Duration
  - time.Format
  - time.Now
  - time.ParseDuration
aliases: [/functions/now]
---

See [`time.Time`](https://godoc.org/time#Time).

For example, building your site on June 24, 2017, with the following templating:

```go-html-template
<div>
    <small>&copy; {{ now.Format "2006" }}</small>
</div>
```

would produce the following:

```html
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
