---
title: RSS templates
description: Use the embedded RSS template, or create your own.
categories: []
keywords: []
weight: 140
---

## Configuration

By default, when you build your project, Hugo generates RSS feeds for home, section, taxonomy, and term pages. Control feed generation in your project configuration. For example, to generate feeds for home and section pages, but not for taxonomy and term pages:

{{< code-toggle file=hugo >}}
[outputs]
home = ['html', 'rss']
section = ['html', 'rss']
taxonomy = ['html']
term = ['html']
{{< /code-toggle >}}

To disable feed generation for all [page kinds](g):

{{< code-toggle file=hugo >}}
disableKinds = ['rss']
{{< /code-toggle >}}

By default, the number of items in each feed is unlimited. Change this as needed in your project configuration:

{{< code-toggle file=hugo >}}
[services.rss]
limit = 42
{{< /code-toggle >}}

Set `limit` to `-1` to generate an unlimited number of items per feed.

The built-in RSS template will render the following values, if present, from your project configuration:

{{< code-toggle file=hugo >}}
copyright = '© 2023 ABC Widgets, Inc.'
[params.author]
name = 'John Doe'
email = 'jdoe@example.org'
{{< /code-toggle >}}

## Include feed reference

To include a feed reference in the `head` element of your rendered pages, place this within the `head` element of your templates:

```go-html-template
{{ with .OutputFormats.Get "rss" }}
  {{ printf `<link rel=%q type=%q href=%q title=%q>` .Rel .MediaType.Type .Permalink site.Title | safeHTML }}
{{ end }}
```

Hugo will render this to:

```html
<link rel="alternate" type="application/rss+xml" href="https://example.org/index.xml" title="ABC Widgets">
```

## Custom templates

Override Hugo's [embedded RSS template][] by creating one or more of your own. For example, to use different templates for home, section, taxonomy, and term pages:

```tree
layouts/
  ├── home.rss.xml
  ├── section.rss.xml
  ├── taxonomy.rss.xml
  └── term.rss.xml
```

RSS templates receive the `.Page` and `.Site` objects in context.

[embedded RSS template]: <{{% eturl rss %}}>
