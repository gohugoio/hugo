---
title: RSS templates
description: Use the built-in RSS template, or create your own.
keywords: [rss, xml, templates]
categories: [templates]
menu:
  docs:
    parent: templates
    weight: 160
weight: 160
toc: true
---

## Configuration

By default, when you build your site, Hugo generates RSS feeds for home, section, taxonomy, and term pages. Control feed generation in your site configuration. For example, to generate feeds for home and section pages, but not for taxonomy and term pages:

{{< code-toggle file=hugo copy=false >}}
[outputs]
home = ['html', 'rss']
section = ['html', 'rss']
taxonomy = ['html']
term = ['html']
{{< /code-toggle >}}

To disable feed generation for all [page kinds]:

{{< code-toggle file=hugo copy=false >}}
disableKinds = ['rss']
{{< /code-toggle >}}

By default, the number of items in each feed is unlimited. Change this as needed in your site configuration:

{{< code-toggle file=hugo copy=false >}}
[services.rss]
limit = 42
{{< /code-toggle >}}

Set `limit` to `-1` to generate an unlimited number of items per feed.

The built-in RSS template will render the following values, if present, from your site configuration:

{{< code-toggle file=hugo copy=false >}}
copyright = '© 2023 ABC Widgets, Inc.'
[params.author]
name = 'John Doe'
email = 'jdoe@example.org'
{{< /code-toggle >}}

## Include feed reference

To include a feed reference in the `head` element of your rendered pages, place this within the `head` element of your templates:

```go-html-template
{{ with .OutputFormats.Get "rss" -}}
  {{ printf `<link rel=%q type=%q href=%q title=%q>` .Rel .MediaType.Type .Permalink site.Title | safeHTML }}
{{ end }}
```

Hugo will render this to:

```html
<link rel="alternate" type="application/rss+xml" href="https://example.org/index.xml" title="ABC Widgets">
```

## Custom templates

Override Hugo's [built-in RSS template] by creating one or more of your own, following the naming conventions as shown in the [template lookup order table].

For example, to use different templates for home, section, taxonomy, and term pages:

```text
layouts/
└── _default/
    ├── home.rss.xml
    ├── section.rss.xml
    ├── taxonomy.rss.xml
    └── term.rss.xml
```

RSS templates receive the `.Page` and `.Site` objects in context.

[built-in RSS template]: https://github.com/gohugoio/hugo/blob/master/tpl/tplimpl/embedded/templates/_default/rss.xml
[page kinds]: /getting-started/glossary/#page-kind
[template lookup order table]: #template-lookup-order

## Template lookup order

The table below shows the RSS template lookup order for the different page kinds. The first listing shows the lookup order when running with a theme (`demoTheme`).

{{< datatable-filtered "output" "layouts" "OutputFormat == rss" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}
