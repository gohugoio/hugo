---
title: Render
description: Renders the given template with the given page as context.
categories: []
keywords: []
action:
  related:
    - functions/partials/Include
    - functions/partials/IncludeCached
  returnType: template.HTML
  signatures: [PAGE.Render NAME]
aliases: [/functions/render]
---

Typically used when ranging over a page collection, the `Render` method on a `Page` object renders the given template, passing the given page as context.

```go-html-template
{{ range site.RegularPages }}
  <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ .Render "summary" }}
{{ end }}
```

In the example above, note that the template ("summary") is identified by its file name without directory or extension.

Although similar to the [`partial`] function, there are key differences.

`Render` method|`partial` function|
:--|:--
The `Page` object is automatically passed to the given template. You cannot pass additional context.| You must specify the context, allowing you to pass a combination of objects, slices, maps, and scalars.
The path to the template is determined by the [content type].|You must specify the path to the template, relative to the layouts/partials directory.

Consider this layout structure:

```text
layouts/
├── _default/
│   ├── baseof.html
│   ├── home.html
│   ├── li.html      <-- used for other content types
│   ├── list.html
│   ├── single.html
│   └── summary.html
└── books/
    ├── li.html      <-- used when content type is "books"
    └── summary.html
```

And this template:

```go-html-template
<ul>
  {{ range site.RegularPages.ByDate }}
    {{ .Render "li" }}
  {{ end }}
</ul>
```

When rendering content of type "books" the `Render` method calls:

```text
layouts/books/li.html
```

For all other content types the `Render` methods calls:

```text
layouts/_default/li.html
```

See [content views] for more examples.

[content views]: /templates/content-view/
[`partial`]: /functions/partials/include/
[content type]: /getting-started/glossary/#content-type
