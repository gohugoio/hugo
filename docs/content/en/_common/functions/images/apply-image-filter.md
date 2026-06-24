---
_comment: Do not remove front matter.
---

Apply the filter using the [`images.Filter`][] function:

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with . | images.Filter $filter }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

You can also apply the filter using the [`Filter`][] method on a `Resource` object:

```go-html-template
{{ with resources.Get "images/original.jpg" }}
  {{ with .Filter $filter }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

[`Filter`]: /methods/resource/filter/
[`images.Filter`]: /functions/images/filter/
