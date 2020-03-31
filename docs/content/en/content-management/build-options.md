---
title: Build Options
linktitle: Build Options
description: Build options help define how Hugo must treat a given page when building the site.
date: 2020-03-02
publishdate: 2020-03-02
keywords: [build,content,front matter, page resources]
categories: ["content management"]
menu:
  docs:
    parent: "content-management"
    weight: 31
weight: 31	#rem
draft: false
aliases: [/content/build-options/]
toc: true
---

They are stored in a reserved Front Matter object named `_build` with the following defaults:

```yaml
_build:
  render: true
  list: always
  publishResources: true
```

#### render
If true, the page will be treated as a published page, holding its dedicated output files (`index.html`, etc...) and permalink.

#### list

Note that we extended this property from a boolean to an enum in Hugo 0.58.0.

Valid values are:

never
: The page will not be incldued in any page collection.

always (default)
: The page will be included in all page collections, e.g. `site.RegularPages`, `$page.Pages`.

local
: The page will be included in any _local_ page collection, e.g. `$page.RegularPages`, `$page.Pages`. One use case for this would be to create fully navigable, but headless content sections. {{< new-in "0.58.0" >}}

If true, the page will be treated as part of the project's collections and, when appropriate, returned by Hugo's listing methods (`.Pages`, `.RegularPages` etc...).

#### publishResources

If set to true the [Bundle's Resources]({{< relref "content-management/page-bundles" >}}) will be published. 
Setting this to false will still publish Resources on demand (when a resource's `.Permalink` or `.RelPermalink` is invoked from the templates) but will skip the others.

{{% note %}}
Any page, regardless of their build options, will always be available using the [`.GetPage`]({{< relref "functions/GetPage" >}}) methods.
{{% /note %}}

------

### Illustrative use cases

#### Not publishing a page
Project needs a "Who We Are" content file for Front Matter and body to be used by the homepage but nowhere else.

```yaml
# content/who-we-are.md`
title: Who we are
_build:
 list: false
 render: false
```

```go-html-template
{{/* layouts/index.html */}}
<section id="who-we-are">
{{ with site.GetPage "who-we-are" }}
  {{ .Content }}
{{ end }}
</section>
```

#### Listing pages without publishing them

Website needs to showcase a few of the hundred "testimonials" available as content files without publishing any of them.

To avoid setting the build options on every testimonials, one can use [`cascade`]({{< relref "/content-management/front-matter#front-matter-cascade" >}}) on the testimonial section's content file.

```yaml
#content/testimonials/_index.md
title: Testimonials
# section build options:
_build:
  render: true
# children build options with cascade
cascade:
  _build:
    render: false
    list: true # default
```

```go-html-template
{{/* layouts/_defaults/testimonials.html */}}
<section id="testimonials">
{{ range first 5 .Pages }}
  <blockquote cite="{{ .Params.cite }}">
    {{ .Content }}
  </blockquote>
{{ end }}
</section>
