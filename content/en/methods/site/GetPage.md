---
title: GetPage
description: Returns a Page object from the given path.
categories: []
keywords: []
action:
  related:
    - methods/page/GetPage
  returnType: page.Page
  signatures: [SITE.GetPage PATH]
toc: true
---

The `GetPage` method is also available on `Page` objects, allowing you to specify a path relative to the current page. See&nbsp;[details].

[details]: /methods/page/getpage/

When using the `GetPage` method on a `Site` object, specify a path relative to the content directory.

If Hugo cannot resolve the path to a page, the method returns nil.

Consider this content structure:

```text
content/
├── works/
│   ├── paintings/
│   │   ├── _index.md
│   │   ├── starry-night.md
│   │   └── the-mona-lisa.md
│   ├── sculptures/
│   │   ├── _index.md
│   │   ├── david.md
│   │   └── the-thinker.md
│   └── _index.md
└── _index.md
```

This home template:

```go-html-template
{{ with .Site.GetPage "/works/paintings" }}
  <ul>
    {{ range .Pages }}
      <li>{{ .Title }} by {{ .Params.artist }}</li>
    {{ end }}
  </ul>
{{ end }}
```

Is rendered to:

```html
<ul>
  <li>Starry Night by Vincent van Gogh</li>
  <li>The Mona Lisa by Leonardo da Vinci</li>
</ul>
```

To get a regular page instead of a section page:

```go-html-template
{{ with .Site.GetPage "/works/paintings/starry-night" }}
  {{ .Title }} → Starry Night
  {{ .Params.artist }} → Vincent van Gogh
{{ end }}
```

## Multilingual projects

With multilingual projects, the `GetPage` method on a `Site` object resolves the given path to a page in the current language.

To get a page from a different language, query the `Sites` object:

```go-html-template
{{ with where .Site.Sites "Language.Lang" "eq" "de" }}
  {{ with index . 0 }}
    {{ with .GetPage "/works/paintings/starry-night" }}
      {{ .Title }} → Sternenklare Nacht
    {{ end }}
  {{ end }}
{{ end }}
```

## Page bundles

Consider this content structure:

```text
content/
├── headless/    
│   ├── a.jpg
│   ├── b.jpg
│   ├── c.jpg
│   └── index.md  <-- front matter: headless = true
└── _index.md
```

In the home template, use the `GetPage` method on a `Site` object to render all the images in the headless [page bundle]:

```go-html-template
{{ with .Site.GetPage "/headless" }}
  {{ range .Resources.ByType "image" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

[page bundle]: /getting-started/glossary/#page-bundle
