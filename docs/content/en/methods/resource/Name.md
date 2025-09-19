---
title: Name
description: Returns the name of the given resource as optionally defined in front matter, falling back to its file path.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [RESOURCE.Name]
---

The value returned by the `Name` method on a `Resource` object depends on the resource type.

## Global resource

With a [global resource](g), the `Name` method returns the path to the resource, relative to the `assets` directory.

```text
assets/
└── images/
    └── Sunrise in Bryce Canyon.jpg
```

```go-html-template
{{ with resources.Get "images/Sunrise in Bryce Canyon.jpg" }}
  {{ .Name }} → /images/Sunrise in Bryce Canyon.jpg
{{ end }}
```

## Page resource

With a [page resource](g), if you create an element in the `resources` array in front matter, the `Name` method returns the value of the `name` parameter.

```text
content/
├── example/
│   ├── images/
│   │   └── a.jpg
│   └── index.md
└── _index.md
```

{{< code-toggle file=content/example/index.md fm=true >}}
title = 'Example'
[[resources]]
src = 'images/a.jpg'
name = 'Sunrise in Bryce Canyon'
{{< /code-toggle >}}

```go-html-template
{{ with .Resources.Get "images/a.jpg" }}
  {{ .Name }} → Sunrise in Bryce Canyon
{{ end }}
```

You can also capture the image by specifying its `name` instead of its path:

```go-html-template
{{ with .Resources.Get "Sunrise in Bryce Canyon" }}
  {{ .Name }} → Sunrise in Bryce Canyon
{{ end }}
```

If you do not create an element in the `resources` array in front matter, the `Name` method returns the file path, relative to the page bundle.

```text
content/
├── example/
│   ├── images/
│   │   └── Sunrise in Bryce Canyon.jpg
│   └── index.md
└── _index.md
```

```go-html-template
{{ with .Resources.Get "images/Sunrise in Bryce Canyon.jpg" }}
  {{ .Name }} → images/Sunrise in Bryce Canyon.jpg
{{ end }}
```

## Remote resource

With a [remote resource](g), the `Name` method returns a hashed file name.

```go-html-template
{{ with resources.GetRemote "https://example.org/images/a.jpg" }}
  {{ .Name }} → /a_18432433023265451104.jpg
{{ end }}
```
