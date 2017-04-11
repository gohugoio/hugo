---
title: after
linktitle: after
description: Slices an array to only the items after the Nth item.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [iteration]
ns:
signature: ["after INDEX COLLECTION"]
workson: []
hugoversion:
relatedfuncs: [last,first,seq]
deprecated: false
aliases: []
---

`after` slices an array to only the items after the <em>N</em>th item. The following shows `after` being used in conjunction with the [`slice` function][slice]:

```html
{{ $data := slice "one" "two" "three" "four" }}
{{ range after 2 $data }}
    {{ . }}
{{ end }}
→ ["three", "four"]
```

## Example of `after` with `first`: 2nd&ndash;4th Most Recent Articles

You can use `after` in combination with the [`first` function][] and Hugo's [powerful sorting methods][lists]. Let's assume you have a list page at `yoursite.com/articles`. You have 10 articles, but you want your templating for the [list/section page][] to show only two rows:

1. The top row is titled "Featured" and shows only the most recently published article (i.e. by `publishdate` in the content files' front matter).
2. The second row is titled "Recent Articles" and shows only the 2nd- to 4th-most recently published articles.

{{% code file="layouts/section/articles.html" download="articles.html" %}}
```html
{{ define "main" }}
<section class="row featured-article">
    <h2>Featured Article</h2>
    {{ range first 1 .Data.Pages.ByPublishDate.Reverse }}
     <header>
        <h3><a href="{{.Permalink}}">{{.Title}}</a></h3>
    </header>
    <p>{{.Description}}</p>
    {{ end }}
</section>
<div class="row recent-articles">
    <h2>Recent Articles</h2>
    {{ range first 3 (after 1 .Data.Pages.ByPublishDate.Reverse) }}
        <section class="recent-article">
            <header>
                <h3><a href="{{.Permalink}}">{{.Title}}</a></h3>
            </header>
            <p>{{.Description}}</p>
        </section>
    {{ end }}
</div>
{{ end }}
```
{{% /code %}}

[`first` function]: /functions/first/
[list/section page]: /templates/section-templates/
[lists]: /lists/
[slice]: /functions/slice/
