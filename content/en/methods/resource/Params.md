---
title: Params
description: Returns a map of resource parameters as defined in front matter.
categories: []
keywords: []
action:
  related: []
  returnType: map
  signatures: [RESOURCE.Params]
---

Use the `Params` method with [page resources]. It is not applicable to either [global] or [remote] resources.

[global]: /getting-started/glossary/#global-resource
[page resources]: /getting-started/glossary/#page-resource
[remote]: /getting-started/glossary/#remote-resource

With this content structure:

```text
content/
├── posts/
│   ├── cats/
│   │   ├── images/
│   │   │   └── a.jpg
│   │   └── index.md
│   └── _index.md
└── _index.md
```

And this front matter:

{{< code-toggle file=content/posts/cats.md fm=true >}}
title = 'Cats'
[[resources]]
  src = 'images/a.jpg'
  title = 'Felix the cat'
  [resources.params]
    alt = 'Photograph of black cat'
    temperament = 'vicious'
{{< /code-toggle >}}

And this template:

```go-html-template
{{ with .Resources.Get "images/a.jpg" }}
  <figure>
    <img alt="{{ .Params.alt }}" src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}">
    <figcaption>{{ .Title }} is {{ .Params.temperament }}</figcaption>
  </figure>
{{ end }}
```

Hugo renders:

```html
<figure>
  <img alt="Photograph of black cat" src="/posts/post-1/images/a.jpg" width="600" height="400">
  <figcaption>Felix the cat is vicious</figcaption>
</figure>
```

See the [page resources] section for more information.

[page resources]: /content-management/page-resources/
