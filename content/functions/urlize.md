---
title: urlize
linktitle: urlize
description: Takes a string, sanitizes it for usage in URLs, and converts spaces to hyphens.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [urls,strings]
godocref:
signature:
hugoversion:
deprecated: false
workson: []
relatedfuncs: []
---

`urlize` takes a string, sanitizes it for usage in URLs, and converts spaces to hyphens ("`-`").

The following examples pull from a content file with the following front matter:

{{% code file="content/blog/greatest-city.md" copy="false"%}}
```toml
+++
title = "The World's Greatest City"
location = "Chicago IL"
tags = ["pizza","beer","hot dogs"]
+++
```
{{% /code %}}

The following might be used as a partial within a [single page template][singletemplate]:

{{% code file="layouts/partials/content-header.html" download="content-header.html" %}}
```html
<header class="content-header">
    <h1>{{.Title}}</h1>
    {{ with .Params.location }}
        <div class="location"><a href="/locations/{{ . | urlize}}">{{.}}</a></div>
    {{ end }}
    <div class="tags">
    {{range .Params.tags}}
        <a href="/tags/{{ . | urlize }}" class="tag">{{ . }}</a><br>
    {{end}}
    </div>
</header>
```
{{% /code %}}

The preceding partial would then output to the rendered page as follows, assuming the page is being built with Hugo's default pretty URLs.

{{% output file="/blog/greatest-city/index.html" %}}
```html
<header class="content-header">
    <h1>The World's Greatest City</h1>
    <div class="location"><a href="/locations/chicago-il/">Chicago IL</a></div>
    <div class="tags">
        <a href="/tags/pizza" class="tag">pizza</a>
        <a href="/tags/beer" class="tag">beer</a>
        <a href="/tags/hot-dogs" class="tag">hot dogs</a>
    </div>
</header>
```
{{% /output %}}


[singletemplate]: /templates/single-page-templates/