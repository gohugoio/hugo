---
title: Content
description: Returns the content of the given resource.
categories: []
keywords: []
action:
  related: []
  returnType: any
  signatures: [RESOURCE.Content]
toc:
---

The `Content` method on a `Resource` object returns `template.HTML` when the resource type is `page`, otherwise it returns a `string`.

[resource type]: /methods/resource/resourcetype/

{{< code file=assets/quotations/kipling.txt >}}
He travels the fastest who travels alone.
{{< /code >}}

To get the content:

```go-html-template
{{ with resources.Get "quotations/kipling.txt" }}
  {{ .Content }} → He travels the fastest who travels alone.
{{ end }}
```

To get the size in bytes:

```go-html-template
{{ with resources.Get "quotations/kipling.txt" }}
  {{ .Content | len }} → 42
{{ end }}
```

To create an inline image:

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  <img src="data:{{ .MediaType.Type }};base64,{{ .Content | base64Encode }}">
{{ end }}
```

To create inline CSS:

```go-html-template
{{ with resources.Get "css/style.css" }}
  <style>{{ .Content | safeCSS }}</style>
{{ end }}
```

To create inline JavaScript:

```go-html-template
{{ with resources.Get "js/script.js" }}
  <script>{{ .Content | safeJS }}</script>
{{ end }}
```

{{% include "methods/resource/_common/global-page-remote-resources.md" %}}
