---
title: reflect.IsResource
description: Reports whether the given value is a Resource object.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: []
    returnType: bool
    signatures: [reflect.IsResource INPUT]
---

{{< new-in 0.154.0 />}}

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

These are the values returned by the `reflect.IsResource` function:

```go-html-template {file="layouts/page.html"}
{{ with resources.Get "a.json" }}
  {{ reflect.IsResource . }} → true
{{ end }}

{{ with resources.Get "b.avif" }}
  {{ reflect.IsResource . }} → true
{{ end }}

{{ with resources.Get "c.jpg" }}
  {{ reflect.IsResource . }} → true
{{ end }}
```

```go-html-template {file="layouts/page.html"}
{{ with .Resources.Get "d.json" }}
  {{ reflect.IsResource . }} → true
{{ end }}

{{ with .Resources.Get "e.avif" }}
  {{ reflect.IsResource . }} → true
{{ end }}

{{ with .Resources.Get "f.jpg" }}
  {{ reflect.IsResource . }} → true
{{ end }}
```

```go-html-template {file="layouts/page.html"}
{{ with site.GetPage "/example" }}
  {{ reflect.IsResource . }} → true
{{ end }}
```
