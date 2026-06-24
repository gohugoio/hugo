---
title: Build options
description: Build options help define how Hugo must treat a given page when building the site.
categories: []
keywords: []
aliases: [/content/build-options/]
---

Build options are stored in a reserved front matter object named `build` with these defaults:

{{< code-toggle file=content/example/index.md fm=true >}}
[build]
list = 'always'
publishResources = true
render = 'always'
{{< /code-toggle >}}

`list`
: When to include the page within page collections. Specify one of:

  - `always`: Include the page in _all_ page collections. For example, `site.RegularPages`, `.Pages`, etc. This is the default value.
  - `local`: Include the page in _local_ page collections. For example, `.RegularPages`, `.Pages`, etc. Use this option to create fully navigable but headless content sections.
  - `never`: Do not include the page in _any_ page collection.

`publishResources`
: Applicable to [page bundles][], determines whether to publish the associated [page resources][]. Specify one of:

  - `true`: Always publish resources. This is the default value.
  - `false`: Only publish a resource when invoking its [`Permalink`][], [`RelPermalink`][], or [`Publish`][] method within a template.

`render`
: When to render the page. Specify one of:

  - `always`: Always render the page to disk. This is the default value.
  - `link`: Do not render the page to disk, but assign `Permalink` and `RelPermalink` values.
  - `never`: Never render the page to disk, and exclude it from all page collections.

> [!NOTE]
> Any page, regardless of its build options, will always be available by using the [`.Page.GetPage`][] or [`.Site.GetPage`][] method.

## Example -- headless page

Create a unpublished page whose content and resources can be included in other pages.

```tree
content/
├── headless/
│   ├── a.jpg
│   ├── b.jpg
│   └── index.md  <-- leaf bundle
└── _index.md     <-- home page
```

Set the build options in front matter:

{{< code-toggle file=content/headless/index.md fm=true >}}
title = 'Headless page'
[build]
  list = 'never'
  publishResources = false
  render = 'never'
{{< /code-toggle >}}

To include the content and images on the home page:

```go-html-template {file="layouts/home.html"}
{{ with .Site.GetPage "/headless" }}
  {{ .Content }}
  {{ range .Resources.ByType "image" }}
    <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
  {{ end }}
{{ end }}
```

The published site will have this structure:

```tree
public/
├── headless/
│   ├── a.jpg
│   └── b.jpg
└── index.html
```

In the example above, note that:

1. Hugo did not publish an HTML file for the page.
1. Despite setting `publishResources` to `false` in front matter, Hugo published the [page resources][] because we invoked the [`RelPermalink`][] method on each resource. This is the expected behavior.

## Example -- headless section

Create a unpublished section whose content and resources can be included in other pages.

```tree
content/
├── headless/
│   ├── note-1/
│   │   ├── a.jpg
│   │   ├── b.jpg
│   │   └── index.md  <-- leaf bundle
│   ├── note-2/
│   │   ├── c.jpg
│   │   ├── d.jpg
│   │   └── index.md  <-- leaf bundle
│   └── _index.md     <-- branch bundle
└── _index.md         <-- home page
```

Set the build options in front matter, using the `cascade` keyword to "cascade" the values down to descendant pages.

{{< code-toggle file=content/headless/_index.md fm=true >}}
title = 'Headless section'
[[cascade]]
[cascade.build]
  list = 'local'
  publishResources = false
  render = 'never'
{{< /code-toggle >}}

In the front matter above, note that we have set `list` to `local` to include the descendant pages in local page collections.

To include the content and images on the home page:

```go-html-template {file="layouts/home.html"}
{{ with .Site.GetPage "/headless" }}
  {{ range .Pages }}
    {{ .Content }}
    {{ range .Resources.ByType "image" }}
      <img src="{{ .RelPermalink }}" width="{{ .Width }}" height="{{ .Height }}" alt="">
    {{ end }}
  {{ end }}
{{ end }}
```

The published site will have this structure:

```tree
public/
├── headless/
│   ├── note-1/
│   │   ├── a.jpg
│   │   └── b.jpg
│   └── note-2/
│       ├── c.jpg
│       └── d.jpg
└── index.html
```

In the example above, note that:

1. Hugo did not publish an HTML file for the page.
1. Despite setting `publishResources` to `false` in front matter, Hugo correctly published the [page resources][] because we invoked the [`RelPermalink`][] method on each resource. This is the expected behavior.

## Example -- list without publishing

Publish a section page without publishing the descendant pages. For example, to create a glossary:

```tree
content/
├── glossary/
│   ├── _index.md
│   ├── bar.md
│   ├── baz.md
│   └── foo.md
└── _index.md
```

Set the build options in front matter, using the `cascade` keyword to "cascade" the values down to descendant pages.

{{< code-toggle file=content/glossary/_index.md fm=true >}}
title = 'Glossary'
[build]
render = 'always'
[[cascade]]
[cascade.build]
  list = 'local'
  publishResources = false
  render = 'never'
{{< /code-toggle >}}

To render the glossary:

```go-html-template {file="layouts/glossary/section.html"}
<dl>
  {{ range .Pages }}
    <dt>{{ .Title }}</dt>
    <dd>{{ .Content }}</dd>
  {{ end }}
</dl>
```

The published site will have this structure:

```tree
public/
├── glossary/
│   └── index.html
└── index.html
```

## Example -- publish without listing

Publish a section's descendant pages without publishing the section page itself.

```tree
content/
├── books/
│   ├── _index.md
│   ├── book-1.md
│   └── book-2.md
└── _index.md
```

Set the build options in front matter:

{{< code-toggle file=content/books/_index.md fm=true >}}
title = 'Books'
[build]
render = 'never'
list = 'never'
{{< /code-toggle >}}

The published site will have this structure:

```tree
public/
├── books/
│   ├── book-1/
│   │   └── index.html
│   └── book-2/
│       └── index.html
└── index.html
```

## Example -- conditionally hide section

Consider this example. A documentation site has a team of contributors with access to 20 custom shortcodes. Each shortcode takes several arguments, and requires documentation for the contributors to reference when using them.

Instead of external documentation for the shortcodes, include an `internal` section that is hidden when building the production site.

```tree
content/
├── internal/
│   ├── shortcodes/
│   │   ├── _index.md
│   │   ├── shortcode-1.md
│   │   └── shortcode-2.md
│   └── _index.md
├── reference/
│   ├── _index.md
│   ├── reference-1.md
│   └── reference-2.md
├── tutorials/
│   ├── _index.md
│   ├── tutorial-1.md
│   └── tutorial-2.md
└── _index.md
```

Set the build options in front matter, using the `cascade` keyword to "cascade" the values down to descendant pages, and use the `target` keyword to target the production environment.

{{< code-toggle file=content/internal/_index.md >}}
title = 'Internal'
[[cascade]]
[cascade.build]
render = 'never'
list = 'never'
[cascade.target]
environment = 'production'
{{< /code-toggle >}}

The production site will have this structure:

```tree
public/
├── reference/
│   ├── reference-1/
│   │   └── index.html
│   ├── reference-2/
│   │   └── index.html
│   └── index.html
├── tutorials/
│   ├── tutorial-1/
│   │   └── index.html
│   ├── tutorial-2/
│   │   └── index.html
│   └── index.html
└── index.html
```

[`.Page.GetPage`]: /methods/page/getpage/
[`.Site.GetPage`]: /methods/site/getpage/
[`Permalink`]: /methods/resource/permalink/
[`Publish`]: /methods/resource/publish/
[`RelPermalink`]: /methods/resource/relpermalink/
[page bundles]: /content-management/page-bundles/
[page resources]: /content-management/page-resources/
