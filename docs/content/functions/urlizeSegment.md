---
title: urlizeSegment
# linktitle: urlizeSegment
description: Takes a string, sanitizes it for usage in URLs, and converts spaces, slashes, and pound signs to hyphens.
date: 2017-10-21
publishdate: 2017-10-21
lastmod: 2017-10-21
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [urls,strings]
godocref:
signature: ["urlizeSegment INPUT"]
hugoversion:
deprecated: false
workson: []
relatedfuncs: [urlize]
---

The following examples pull from a content file with the following front matter:

{{< code file="content/blog/greatest-city.md" copy="false">}}
+++
title = "The World's Greatest City"
location = "Chicago IL/USA"
tags = ["food/pizza","drink/beer","food/hot dogs","we're #1"]
+++
{{< /code >}}

The following might be used as a partial within a [single page template][singletemplate]:

{{< code file="layouts/partials/content-header.html" download="content-header.html" >}}
<header>
    <h1>{{.Title}}</h1>
    {{ with .Params.location }}
        <div><a href="/locations/{{ . | urlizeSegment}}">{{.}}</a></div>
    {{ end }}
    <!-- Creates a list of tags for the content and links to each of their pages -->
    {{ with .Params.tags }}
    <ul>
        {{range .}}
            <li>
                <a href="/tags/{{ . | urlizeSegment }}">{{ . }}</a>
            </li>
        {{end}}
    </ul>
    {{ end }}
</header>
{{< /code >}}

The preceding partial would then output to the rendered page as follows, assuming the page is being built with Hugo's default pretty URLs.

{{< output file="/blog/greatest-city/index.html" >}}
<header>
    <h1>The World's Greatest City</h1>
    <div><a href="/locations/chicago-il-usa/">Chicago IL/USA</a></div>
    <ul>
        <li>
            <a href="/tags/food-pizza">food/pizza</a>
        </li>
        <li>
            <a href="/tags/drink-beer">drink/beer</a>
        </li>
        <li>
            <a href="/tags/food-hot-dogs">food/hot dogs</a>
        </li>
        <li>
            <a href="/tags/were--1">we're #1</a>
        </li>
    </ul>
</header>
{{< /output >}}


[singletemplate]: /templates/single-page-templates/
