---
title: Homepage Template
linktitle: Homepage Template
description: The homepage of a website is often formatted differently than the other pages. For this reason, Hugo makes it easy for you to define your new site's homepage as a unique template.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
keywords: [homepage]
menu:
  docs:
    parent: "templates"
    weight: 30
weight: 30
sections_weight: 30
draft: false
aliases: [/layout/homepage/,/templates/homepage-template/]
toc: true
---

Homepage is a `Page` and therefore has all the [page variables][pagevars] and [site variables][sitevars] available for use.

{{% note "The Only Required Template" %}}
The homepage template is the *only* required template for building a site and therefore useful when bootstrapping a new site and template. It is also the only required template if you are developing a single-page website.
{{% /note %}}

{{< youtube ut1xtRZ1QOA >}}

## Homepage Template Lookup Order

See [Template Lookup](/templates/lookup-order/).

## Add Content and Front Matter to the Homepage

The homepage, similar to other [list pages in Hugo][lists], accepts content and front matter from an `_index.md` file. This file should live at the root of your `content` folder (i.e., `content/_index.md`). You can then add body copy and metadata to your homepage the way you would any other content file.

See the homepage template below or [Content Organization][contentorg] for more information on the role of `_index.md` in adding content and front matter to list pages.

## Example Homepage Template

The following is an example of a homepage template that uses [partial][partials], [base][] templates, and a content file at `content/_index.md` to populate the `{{.Title}}` and `{{.Content}}` [page variables][pagevars].

{{< code file="layouts/index.html" download="index.html" >}}
{{ define "main" }}
    <main aria-role="main">
      <header class="homepage-header">
        <h1>{{.Title}}</h1>
        {{ with .Params.subtitle }}
        <span class="subtitle">{{.}}</span>
        {{ end }}
      </header>
      <div class="homepage-content">
        <!-- Note that the content for index.html, as a sort of list page, will pull from content/_index.md -->
        {{.Content}}
      </div>
      <div>
        {{ range first 10 .Site.RegularPages }}
            {{ .Render "summary"}}
        {{ end }}
      </div>
    </main>
{{ end }}
{{< /code >}}

[base]: /templates/base/
[contentorg]: /content-management/organization/
[lists]: /templates/lists/
[lookup]: /templates/lookup-order/
[pagevars]: /variables/page/
[partials]: /templates/partials/
[sitevars]: /variables/site/
