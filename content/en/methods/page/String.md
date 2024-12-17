---
title: String
description: Returns the absolute path to the file backing the given page.
categories: []
keywords: []
action:
  related: []
  returnType: string
  signatures: [PAGE.String]
toc: true
---

{{< new-in 0.137.0 >}}

With content from the file system:

```go-html-template
{{ .String }} → /home/user/project/content/posts/post-1.md
```

With content from a content adapter:

```go-html-template
{{ .String }} → /home/user/project/content/posts/_content.gotmpl:/posts/post-1.md
```

With content from a module:

```go-html-template
{{ .String }}  → /home/user/.cache/hugo_cache/modules/filecache/modules/pkg/mod/github.com/user/hugo-module-content@v0.1.9/content/posts/post-1.md
```

Use this method to provide useful information when displaying error and warning messages in the console:

{{< code file="layouts/partials/featured-image.html" lang="go-html-template" >}}
{{ with .Resources.GetMatch "*featured*" }}
  {{ with .Resize "300x webp" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ else }}
  {{ errorf "The featured-image partial was unable to find a featured image. See %s" .String }}
{{ end }}
{{< /code >}}

With shortcodes and render hooks use the `Position` method instead. Note that the `Position` method is not available to heading, image, and link render hooks.
