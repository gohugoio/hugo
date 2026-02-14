---
title: reflect.IsImageResource
description: Reports whether the given value is a Resource object representing a processable image.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: bool
    signatures: [reflect.IsImageResource INPUT]
---

{{< new-in 0.154.0 />}}

{{% glossary-term "processable image" %}}

With this project structure:

```text
project/
├── assets/
│   ├── a.json
│   ├── b.avif
│   └── c.jpg
└── content/
    └── example/
        ├── index.md
        ├── d.json
        ├── e.avif
        └── f.jpg
```

These are the values returned by the `reflect.IsImageResource` function:

```go-html-template {file="layouts/page.html"}
{{ with resources.Get "a.json" }}
  {{ reflect.IsImageResource . }} → false
{{ end }}

{{ with resources.Get "b.avif" }}
  {{ reflect.IsImageResource . }} → false
{{ end }}

{{ with resources.Get "c.jpg" }}
  {{ reflect.IsImageResource . }} → true
{{ end }}
```

In the example above, the `b.avif` image is not a processable image because Hugo can neither decode nor encode the AVIF image format.

```go-html-template {file="layouts/page.html"}
{{ with .Resources.Get "d.json" }}
  {{ reflect.IsImageResource . }} → false
{{ end }}

{{ with .Resources.Get "e.avif" }}
  {{ reflect.IsImageResource . }} → false
{{ end }}

{{ with .Resources.Get "f.jpg" }}
  {{ reflect.IsImageResource . }} → true
{{ end }}
```

In the example above, the `e.avif` image is not a processable image because Hugo can neither decode nor encode the AVIF image format.

```go-html-template {file="layouts/page.html"}
{{ with site.GetPage "/example" }}
  {{ reflect.IsImageResource . }} → false
{{ end }}
```
