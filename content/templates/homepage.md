---
title: Homepage Template
linktitle: Homepage Template
description: The homepage of a website is often formatted differently than the other pages. For this reason, Hugo makes it easy for you to define your new site's homepage as a unique template.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [homepage]
menu:
  main:
    parent: "Templates"
    weight: 30
weight: 30
sections_weight: 30
draft: false
aliases: [/layout/homepage/,/templates/homepage-template/]
toc: false
---

The homepage of a website often has a unique design. In Hugo, you can define your own homepage template.

Homepage is a `Page` and therefore has all the [page variables][pagevars] and [site variables][sitevars] available for use.

{{% note "The Only Required Template" %}}
The homepage template is the *only* required template for building a site and therefore useful when bootstrapping a new site and template. It is also the only required template if you are developing a single-page website.
{{% /note %}}

## Homepage Template Lookup Order

The [lookup order][lookup] for the homepage template is as follows:

1. `/layouts/index.html`
2. `/layouts/_default/list.html`
3. `/layouts/_default/single.html`
4. `/themes/<THEME>/layouts/index.html`
5. `/themes/<THEME>/layouts/_default/list.html`
6. `/themes/<THEME>/layouts/_default/single.html`

## `.Data.Pages` on the Homepage

In addition to the standard [page variables][pagevars], the homepage template has access to *all* site content via `.Data.Pages`.


## Homepage Content and Front Matter

A homepage can also have a content file with front matter: `content/_index.md`. See [Content Organization][contentorg] for more information.

## Example Homepage Template

The following is an example of a homepage template that uses [partial][partials] and [block][] templates.

{{% code file="layouts/index.html" download="index.html" %}}
```html
{{ define "main" }}
    {{ partial "content-header.html" . }}
    <main aria-role="main">
      <div>
        <!-- Note that .Data.Pages is the equivalent of .Site.Pages on the homepage template. -->
        {{ range first 10 .Data.Pages }}
            {{ .Render "summary"}}
        {{ end }}
      </div>
    </main>
    {{ partial "content-footer.html" . }}
{{ end }}
```
{{% /code %}}

[block]: /templates/base/
[contentorg]: /content-management/organization/
[lists]: /templates/lists/
[lookup]: /templates/lookup-order/
[pagevars]: /variables/page/
[partials]: /templates/partials/
[sitevars]: /variables/site/