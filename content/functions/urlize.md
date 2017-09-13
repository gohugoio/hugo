---
title: urlize
# linktitle: urlize
description: Takes a string, sanitizes it for usage in URLs, and converts spaces to hyphens.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [urls,strings]
godocref:
signature: ["urlize INPUT"]
hugoversion:
deprecated: false
workson: []
relatedfuncs: []
---

The following examples pull from a content file with the following front matter:

{{< code file="content/blog/greatest-city.md" copy="false">}}
+++
title = "The World's Greatest City"
location = "Chicago IL"
tags = ["pizza","beer","hot dogs"]
+++
{{< /code >}}

The following might be used as a partial within a [single page template][singletemplate]:

{{< code file="layouts/partials/content-header.html" download="content-header.html" >}}
<header>
    <h1>{{.Title}}</h1>
    {{ with .Params.location }}
        <div><a href="/locations/{{ . | urlize}}">{{.}}</a></div>
    {{ end }}
    <!-- Creates a list of tags for the content and links to each of their pages -->
    {{ with .Params.tags }}
    <ul>
        {{range .}}
            <li>
                <a href="/tags/{{ . | urlize }}">{{ . }}</a>
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
    <div><a href="/locations/chicago-il/">Chicago IL</a></div>
    <ul>
        <li>
            <a href="/tags/pizza">pizza</a>
        </li>
        <li>
            <a href="/tags/beer">beer</a>
        </li>
        <li>
            <a href="/tags/hot-dogs">hot dogs</a>
        </li>
    </ul>
</header>
{{< /output >}}


[singletemplate]: /templates/single-page-templates/
