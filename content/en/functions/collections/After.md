---
title: collections.After
description: Slices an array to the items after the Nth item.
categories: []
keywords: []
action:
  aliases: [after]
  related:
    - functions/collections/First
    - functions/collections/Last
  returnType: any
  signatures: [collections.After INDEX COLLECTION]
aliases: [/functions/after]
---

The following shows `after` being used in conjunction with the [`slice`]function:

```go-html-template
{{ $data := slice "one" "two" "three" "four" }}
<ul>
  {{ range after 2 $data }}
    <li>{{ . }}</li>
  {{ end }}
</ul>
```

The template above is rendered to:

```html
<ul>
  <li>three</li>
  <li>four</li>
</ul>
```

## Example of `after` with `first`: 2nd&ndash;4th most recent articles

You can use `after` in combination with the [`first`] function and Hugo's [powerful sorting methods](/quick-reference/page-collections/#sort). Let's assume you have a `section` page at `example.com/articles`. You have 10 articles, but you want your template to show only two rows:

1. The top row is titled "Featured" and shows only the most recently published article (i.e. by `publishdate` in the content files' front matter).
2. The second row is titled "Recent Articles" and shows only the 2nd- to 4th-most recently published articles.

{{< code file=layouts/section/articles.html >}}
{{ define "main" }}
  <section class="row featured-article">
    <h2>Featured Article</h2>
    {{ range first 1 .Pages.ByPublishDate.Reverse }}
    <header>
      <h3><a href="{{ .RelPermalink }}">{{ .Title }}</a></h3>
    </header>
    <p>{{ .Description }}</p>
  {{ end }}
  </section>
  <div class="row recent-articles">
    <h2>Recent Articles</h2>
    {{ range first 3 (after 1 .Pages.ByPublishDate.Reverse) }}
      <section class="recent-article">
        <header>
          <h3><a href="{{ .RelPermalink }}">{{ .Title }}</a></h3>
        </header>
        <p>{{ .Description }}</p>
      </section>
    {{ end }}
  </div>
{{ end }}
{{< /code >}}

[`first`]: /functions/collections/first/
[`slice`]: /functions/collections/slice/
