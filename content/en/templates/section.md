---
title: Section templates
description: Use section templates to list members of a section.
categories: [templates]
keywords: []
menu:
  docs:
    parent: templates
    weight: 80
weight: 80
toc: true
aliases: [/templates/sections/,/templates/section-templates/]
---

## Add content and front matter to section templates

To effectively leverage section templates, you should first understand Hugo's [content organization](/content-management/organization/) and, specifically, the purpose of `_index.md` for adding content and front matter to section and other list pages.

## Section template lookup order

See [Template Lookup](/templates/lookup-order/).

## Example: creating a default section template

{{< code file=layouts/_default/section.html >}}
{{ define "main" }}
  <main>
    {{ .Content }}

    {{ $pages := where site.RegularPages "Type" "posts" }}
    {{ $paginator := .Paginate $pages }}

    {{ range $paginator.Pages }}
      <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
    {{ end }}

    {{ template "_internal/pagination.html" . }}
  </main>
{{ end }}
{{< /code >}}

### Example: using `.Site.GetPage`

The `.Site.GetPage` example that follows assumes the following project directory structure:

```txt
.
└── content
    ├── blog
    │   ├── _index.md   <-- title: My Hugo Blog
    │   ├── post-1.md
    │   ├── post-2.md
    │   └── post-3.md
    └── events
        ├── event-1.md
        └── event-2.md
```

`.Site.GetPage` will return `nil` if no `_index.md` page is found. Therefore, if `content/blog/_index.md` does not exist, the template will output the section name:

```go-html-template
<h1>{{ with .Site.GetPage "/blog" }}{{ .Title }}{{ end }}</h1>
```

Since `blog` has a section index page with front matter at `content/blog/_index.md`, the above code will return the following result:

```html
<h1>My Hugo Blog</h1>
```

If we try the same code with the `events` section, however, Hugo will default to the section title because there is no `content/events/_index.md` from which to pull content and front matter:

```go-html-template
<h1>{{ with .Site.GetPage "/events" }}{{ .Title }}{{ end }}</h1>
```

Which then returns the following:

```html
<h1>Events</h1>
```

[contentorg]: /content-management/organization/
[lookup]: /templates/lookup-order/
[`where`]: /functions/collections/where/
[sections]: /content-management/sections/
