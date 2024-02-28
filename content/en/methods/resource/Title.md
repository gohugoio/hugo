---
title: Title
description: Returns the title of the given resource as optionally defined in front matter, falling back to a relative path or hashed file name depending on resource type.
categories: []
keywords: []
action:
  related:
    - methods/resource/Name
  returnType: string
  signatures: [RESOURCE.Title]
toc: true
---

The value returned by the `Title` method on a `Resource` object depends on the resource type.

## Global resource

With a [global resource], the `Title` method returns the path to the resource, relative to the assets directory.

```text
assets/
└── images/
    └── Sunrise in Bryce Canyon.jpg
```

```go-html-template
{{ with resources.Get "images/Sunrise in Bryce Canyon.jpg" }}
  {{ .Title }} → /images/Sunrise in Bryce Canyon.jpg
{{ end }}
```

## Page resource

With a [page resource], if you create an element in the `resources` array in front matter, the `Title` method returns the value of the `title` parameter.

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
title = 'A beautiful sunrise in Bryce Canyon'
{{< /code-toggle >}}

```go-html-template
{{ with .Resources.Get "images/a.jpg" }}
  {{ .Title }} → A beautiful sunrise in Bryce Canyon
{{ end }}
```

If you do not create an element in the `resources` array in front matter, the `Title` method returns the file path, relative to the page bundle.

```text
content/
├── example/
│   ├── images/
│   │   └── Sunrise in Bryce Canyon.jpg
│   └── index.md
└── _index.md
```

```go-html-template
{{ with .Resources.Get "Sunrise in Bryce Canyon.jpg" }}
  {{ .Title }} → images/Sunrise in Bryce Canyon.jpg
{{ end }}
```

## Remote resource

With a [remote resource], the `Title` method returns a hashed file name.

```go-html-template
{{ with resources.GetRemote "https://example.org/images/a.jpg" }}
  {{ .Title }} → /a_18432433023265451104.jpg
{{ end }}
```

[global resource]: /getting-started/glossary/#global-resource
[page resource]: /getting-started/glossary/#page-resource
[remote resource]: /getting-started/glossary/#remote-resource
