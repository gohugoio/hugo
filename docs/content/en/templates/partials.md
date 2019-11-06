---
title: Partial Templates
linktitle: Partial Templates
description: Partials are smaller, context-aware components in your list and page templates that can be used economically to keep your templating DRY.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
keywords: [lists,sections,partials]
menu:
  docs:
    parent: "templates"
    weight: 90
weight: 90
sections_weight: 90
draft: false
aliases: [/templates/partial/,/layout/chrome/,/extras/analytics/]
toc: true
---

{{< youtube pjS4pOLyB7c >}}

## Partial Template Lookup Order

Partial templates---like [single page templates][singletemps] and [list page templates][listtemps]---have a specific [lookup order][]. However, partials are simpler in that Hugo will only check in two places:

1. `layouts/partials/*<PARTIALNAME>.html`
2. `themes/<THEME>/layouts/partials/*<PARTIALNAME>.html`

This allows a theme's end user to copy a partial's contents into a file of the same name for [further customization][customize].

## Use Partials in your Templates

All partials for your Hugo project are located in a single `layouts/partials` directory. For better organization, you can create multiple subdirectories within `partials` as well:

```
.
└── layouts
    └── partials
        ├── footer
        │   ├── scripts.html
        │   └── site-footer.html
        ├── head
        │   ├── favicons.html
        │   ├── metadata.html
        │   ├── prerender.html
        │   └── twitter.html
        └── header
            ├── site-header.html
            └── site-nav.html
```

All partials are called within your templates using the following pattern:

```
{{ partial "<PATH>/<PARTIAL>.html" . }}
```

{{% note %}}
One of the most common mistakes with new Hugo users is failing to pass a context to the partial call. In the pattern above, note how "the dot" (`.`) is required as the second argument to give the partial context. You can read more about "the dot" in the [Hugo templating introduction](/templates/introduction/).
{{% /note %}}

{{% note %}}
`<PARTIAL>` including `baseof` is reserved. ([#5373](https://github.com/gohugoio/hugo/issues/5373))
{{% /note %}}

As shown in the above example directory structure, you can nest your directories within `partials` for better source organization. You only need to call the nested partial's path relative to the `partials` directory:

```
{{ partial "header/site-header.html" . }}
{{ partial "footer/scripts.html" . }}
```

### Variable Scoping

The second argument in a partial call is the variable being passed down. The above examples are passing the `.`, which tells the template receiving the partial to apply the current [context][context].

This means the partial will *only* be able to access those variables. The partial is isolated and *has no access to the outer scope*. From within the partial, `$.Var` is equivalent to `.Var`.

## Returning a value from a Partial

In addition to outputting markup, partials can be used to return a value of any type. In order to return a value, a partial must include a lone `return` statement.

### Example GetFeatured
```go-html-template
{{/* layouts/partials/GetFeatured.html */}}
{{ return first . (where site.RegularPages ".Params.featured" true) }}
```

```go-html-template
{{/* layouts/index.html */}}
{{ range partial "GetFeatured.html" 5 }}
  [...]
{{ end }}
```
### Example GetImage
```go-html-template
{{/* layouts/partials/GetImage.html */}}
{{ $image := false }}
{{ with .Params.gallery }}
  {{ $image = index . 0 }}
{{ end }}
{{ with .Params.image }}
  {{ $image = . }}
{{ end }}
{{ return $image }}
```

```go-html-template
{{/* layouts/_default/single.html */}}
{{ with partial "GetImage.html" . }}
  [...]
{{ end }}
```

{{% note %}}
Only one `return` statement is allowed per partial file.
{{% /note %}}

## Cached Partials

The [`partialCached` template function][partialcached] can offer significant performance gains for complex templates that don't need to be re-rendered on every invocation. The simplest usage is as follows:

```
{{ partialCached "footer.html" . }}
```

You can also pass additional parameters to `partialCached` to create *variants* of the cached partial.

For example, you can tell Hugo to only render the partial `footer.html` once per section:

```
{{ partialCached "footer.html" . .Section }}
```

If you need to pass additional parameters to create unique variants, you can pass as many variant parameters as you need:

```
{{ partialCached "footer.html" . .Params.country .Params.province }}
```

Note that the variant parameters are not made available to the underlying partial template. They are only use to create a unique cache key.

### Example `header.html`

The following `header.html` partial template is used for [spf13.com](https://spf13.com/):

{{< code file="layouts/partials/header.html" download="header.html" >}}
<!DOCTYPE html>
<html class="no-js" lang="en-US" prefix="og: http://ogp.me/ns# fb: http://ogp.me/ns/fb#">
<head>
    <meta charset="utf-8">

    {{ partial "meta.html" . }}

    <base href="{{ .Site.BaseURL }}">
    <title> {{ .Title }} : spf13.com </title>
    <link rel="canonical" href="{{ .Permalink }}">
    {{ if .RSSLink }}<link href="{{ .RSSLink }}" rel="alternate" type="application/rss+xml" title="{{ .Title }}" />{{ end }}

    {{ partial "head_includes.html" . }}
</head>
<body lang="en">
{{< /code >}}

{{% note %}}
The `header.html` example partial was built before the introduction of block templates to Hugo. Read more on [base templates and blocks](/templates/base/) for defining the outer chrome or shell of your master templates (i.e., your site's head, header, and footer). You can even combine blocks and partials for added flexibility.
{{% /note %}}

### Example `footer.html`

The following `footer.html` partial template is used for [spf13.com](https://spf13.com/):

{{< code file="layouts/partials/footer.html" download="footer.html" >}}
<footer>
  <div>
    <p>
    &copy; 2013-14 Steve Francia.
    <a href="https://creativecommons.org/licenses/by/3.0/" title="Creative Commons Attribution">Some rights reserved</a>;
    please attribute properly and link back. Hosted by <a href="http://servergrove.com">ServerGrove</a>.
    </p>
  </div>
</footer>
<script type="text/javascript">

  var _gaq = _gaq || [];
  _gaq.push(['_setAccount', 'UA-XYSYXYSY-X']);
  _gaq.push(['_trackPageview']);

  (function() {
    var ga = document.createElement('script');
    ga.src = ('https:' == document.location.protocol ? 'https://ssl' :
        'http://www') + '.google-analytics.com/ga.js';
    ga.setAttribute('async', 'true');
    document.documentElement.firstChild.appendChild(ga);
  })();

</script>
</body>
</html>
{{< /code >}}

[context]: /templates/introduction/ "The most easily overlooked concept to understand about Go templating is how the dot always refers to the current context."
[customize]: /themes/customizing/ "Hugo provides easy means to customize themes as long as users are familiar with Hugo's template lookup order."
[listtemps]: /templates/lists/ "To effectively leverage Hugo's system, see how Hugo handles list pages, where content for sections, taxonomies, and the homepage are listed and ordered."
[lookup order]: /templates/lookup-order/ "To keep your templating dry, read the documentation on Hugo's lookup order."
[partialcached]: /functions/partialcached/ "Use the partial cached function to improve build times in cases where Hugo can cache partials that don't need to be rendered with every page."
[singletemps]: /templates/single-page-templates/ "The most common form of template in Hugo is the single content template. Read the docs on how to create templates for individual pages."
[themes]: /themes/
