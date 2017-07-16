---
title: unix
# Set this to draft for now; "unix" is not a function. We need to redo this section with:
# - methods vs functions
# - function namespaces as subsections.
draft: true
description: Unix returns the local Time corresponding to the given Unix time, sec seconds and nsec nanoseconds since January 1, 1970 UTC.
godocref: https://golang.org/search?q=Unix#Functions
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
#tags: [dates,time]
categories: [functions]
menu:
  docs:
    parent: "functions"
signature: ["Unix"]
workson: [times]
hugoversion:
relatedfuncs: [Format,dateFormat,now,time]
deprecated: false
aliases: []
---

`Unix` returns t as a Unix time; i.e., the number of seconds elapsed since January 1, 1970 UTC.

## `Unix` Example 1: Show Only Upcoming Events

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

## `Unix` Example 2: Time Passed Since Last Modification

This very simple one-liner uses `Unix` with `Now` to calculate the amount of time that has passed between the `.LastMod` for the current page and the last build of the current page.

{{% code file="time-passed.html" %}}
```golang
{{ div (sub .Now.Unix .Lastmod.Unix) 86400 }}
```
{{% /code %}}

Since both values are integers, they can be subtracted and then divided by the number of seconds in a day (i.e., `60 * 60 * 24 == 86400`).

{{% note %}}
Hugo's output is *static*. In example 2, a month-old page published on a Hugo site that only publishes monthly could easily misrepresented the last update as *yesterday* rather than 30 days ago.
 {{% /note %}}



[partial template]: /templates/partials/
