---
title: unix
linktitle: Unix
description: Unix returns the local Time corresponding to the given Unix time, sec seconds and nsec nanoseconds since January 1, 1970 UTC.
godocref: https://golang.org/search?q=Unix#Functions
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [dates,time]
categories: [functions]
toc:
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---

`Unix` returns t as a Unix time; i.e., the number of seconds elapsed since January 1, 1970 UTC.

## `Unix` Example: Show Only Upcoming Events

The following assumes you have a content section called `events` (i.e., `content/events/*.md`). The following [partial template][] allows you to only list events that haven't occurred yet.

{{% code file="layouts/partials/upcoming-events.html" download="upcoming-events.html" %}}
```html
<h4>Upcoming Events</h4>
<ul class="upcoming-events">
{{ range where .Data.Pages.ByDate "Section" "events" }}
  {{ if ge .Date.Unix .Now.Unix }}
    <li><span class="event-type">{{ .Type | title }} —</span>
      {{ .Title }}
      on <span class="event-date">
      {{ .Date.Format "2 January at 3:04pm" }}</span>
      at {{ .Params.place }}
    </li>
  {{ end }}
{{ end }}
</ul>
```
{{% /code %}}

[partial template]: /templates/partials/