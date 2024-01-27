---
title: Name
description: Returns the name of the given resource as optionally defined in front matter, falling back to a relative path or hashed file name depending on resource type.
categories: []
keywords: []
action:
  related:
    - methods/resource/Title
  returnType: string
  signatures: [RESOURCE.Name]
toc: true
---

The value returned by the `Name` method on a `Resource` object depends on the resource type.

## Global resource

With a [global resource], the `Name` method returns the path to the resource, relative to the assets directory.

```text
assets/
└── images/
    └── a.jpg
```

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ .Name }} → images/a.jpg
{{ end }}
```

## Page resource

With a [page resource], the `Name` method returns the path to the resource, relative to the page bundle.

```text
content/
├── posts/
│   ├── post-1/
│   │   ├── images/
│   │   │   └── a.jpg
│   │   └── index.md
│   └── _index.md
└── _index.md
```

```go-html-template
{{ with .Resources.Get "images/a.jpg" }}
  {{ .Name }} → images/a.jpg
{{ end }}
```

If you create an element in the `resources` array in front matter, the `Name` method returns the value of the `name` parameter:

{{< code-toggle file=content/posts/post-1.md fm=true >}}
title = 'Post 1'
[[resources]]
src = 'images/a.jpg'
name = 'cat'
title = 'Felix the cat'
[resources.params]
temperament = 'malicious'
{{< /code-toggle >}}

```go-html-template
{{ with .Resources.Get "cat" }}
  {{ .Name }} →  cat
{{ end }}
```
## Remote resource

With a [remote resource], the `Name` method returns a hashed file name.

```go-html-template
{{ with resources.GetRemote "https://example.org/images/a.jpg" }}
  {{ .Name }} → a_18432433023265451104.jpg
{{ end }}
```

[global resource]: /getting-started/glossary/#global-resource
[page resource]: /getting-started/glossary/#page-resource
[remote resource]: /getting-started/glossary/#remote-resource
