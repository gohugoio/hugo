---
title: "Ordering Content"
date: "2014-03-06"
linktitle: "Ordering"
menu:
  main:
    parent: "content"
weight: 60
---

In Hugo you have a good degree of control of how your content can be ordered.

By default, content is ordered by weight, then by date with the most recent
date first, but alternative sorting (by title and linktitle) is also available.

_Both the date and weight fields are optional._

Unweighted pages appear at the end of the list. If no weights are provided (or
if weights are the same) date will be used to sort. If neither are provided
content will be ordered based on how it's read off the disk and no order is
guaranteed.

## Assigning Weight to content

    +++
    weight = "4"
    title = "Three"
    date = "2012-04-06"
    +++
    Front Matter with Ordered Pages 3

## Order by Weight -> Date (default)

    {{ range .Data.Pages }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

## Order by Weight -> Date

    {{ range .Data.Pages.ByWeight }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

## Order by Date

    {{ range .Data.Pages.ByDate }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

## Order by Length

    {{ range .Data.Pages.ByLength }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

## Reverse Order
Can be applied to any of the above. Using Date for an example.

    {{ range .Data.Pages.ByDate.Reverse }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

## Order by Title

    {{ range .Data.Pages.ByTitle }}
    <li>
    <a href="{{ .Permalink }}">{{ .Title }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}

## Order by LinkTitle

    {{ range .Data.Pages.ByLinkTitle }}
    <li>
    <a href="{{ .Permalink }}">{{ .LinkTitle }}</a>
    <div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
    </li>
    {{ end }}


## Ordering Content Within Indexes

Please see the [Index Ordering Documentation](/indexes/ordering/)
