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
    └── a.jpg
```

```go-html-template
{{ with resources.Get "images/a.jpg" }}
  {{ .Title }} → images/a.jpg
{{ end }}
```

## Page resource

With a [page resource], the `Title` method returns the path to the resource, relative to the page bundle.

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
  {{ .Title }} → images/a.jpg
{{ end }}
```

If you create an element in the `resources` array in front matter, the `Title` method returns the value of the `title` parameter:

{{< code-toggle file="content/posts/post-1.md" fm=true >}}
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
  {{ .Title }} →  Felix the cat
{{ end }}
```

If the page resource is a content file, the `Title` methods return the `title` field as defined in front matter.

```text
content/
├── lessons/
│   ├── lesson-1/
│   │   ├── _objectives.md  <-- resource type = page
│   │   └── index.md
│   └── _index.md
└── _index.md
```

## Remote resource

With a [remote resource], the `Title` method returns a hashed file name.

```go-html-template
{{ with resources.GetRemote "https://example.org/images/a.jpg" }}
  {{ .Title }} → a_18432433023265451104.jpg
{{ end }}
```

[global resource]: /getting-started/glossary/#global-resource
[page resource]: /getting-started/glossary/#page-resource
[remote resource]: /getting-started/glossary/#remote-resource
