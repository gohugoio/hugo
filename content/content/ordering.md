---
title: "Ordering Content"
date: "2013-07-01"
linktitle: "Ordering"
groups: ['content']
groups_weight: 60
---

In Hugo you have a good degree of control of how your content can be ordered.

By default, content is ordered by weight, then by date with the most recent date first.

_Both the date and weight fields are optional._

Unweighted pages appear at the end of the list.
If no weights are provided (or if weights are the same) date will be used to sort. If neither are provided
content will be ordered based on how it's read off the disk and no order is guaranteed.

Alternative sorting is also available to order content by date (ignoring weight), length and reverse the order.

## Assigning Weight to content

    +++
    weight = "4"
    title = "Three"
    date = "2012-04-06"
    +++
    Front Matter with Ordered Pages 3

## Order by Weight -> Date (default)

{{% highlight html %}}
{{ range .Data.Pages }}
<li>
<a href="{{ .Permalink }}">{{ .Title }}</a>
<div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
</li>
{{ end }}
{{% /highlight %}}

## Order by Weight -> Date

{{% highlight html %}}
{{ range .Data.Pages.ByWeight }}
<li>
<a href="{{ .Permalink }}">{{ .Title }}</a>
<div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
</li>
{{ end }}
{{% /highlight %}}


## Order by Date

{{% highlight html %}}
{{ range .Data.Pages.ByDate }}
<li>
<a href="{{ .Permalink }}">{{ .Title }}</a>
<div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
</li>
{{ end }}
{{% /highlight %}}

## Order by Length

{{% highlight html %}}
{{ range .Data.Pages.ByLength }}
<li>
<a href="{{ .Permalink }}">{{ .Title }}</a>
<div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
</li>
{{ end }}
{{% /highlight %}}

## Reverse Order
Can be applied to any of the above. Using Date for an example.

{{% highlight html %}}
{{ range .Data.Pages.ByDate.Reverse }}
<li>
<a href="{{ .Permalink }}">{{ .Title }}</a>
<div class="meta">{{ .Date.Format "Mon, Jan 2, 2006" }}</div>
</li>
{{ end }}
{{% /highlight %}}

## Ordering Content Within Indexes

Please see the [Index Ordering Documentation](/indexes/ordering/)
