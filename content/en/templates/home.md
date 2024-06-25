---
title: Home templates
description: The home page of a website is often formatted differently than the other pages. For this reason, Hugo makes it easy for you to define your new site's home page as a unique template.
categories: [templates]
keywords: []
menu:
  docs:
    parent: templates
    weight: 60
weight: 60
toc: true
aliases: [/layout/homepage/,/templates/homepage-template/,/templates/homepage/]
---

The home template is the *only* required template for building a site and therefore useful when bootstrapping a new site and template. It is also the only required template if you are developing a single-page website.


{{< youtube ut1xtRZ1QOA >}}

## Home template lookup order

See [Template Lookup](/templates/lookup-order/).

## Add content and front matter to the home page

The home page accepts content and front matter from an `_index.md` file. This file should live at the root of your `content` folder (i.e., `content/_index.md`). You can then add body copy and metadata to your home page the way you would any other content file.

See the home template below or [Content Organization][contentorg] for more information on the role of `_index.md` in adding content and front matter to list pages.

## Example home template

{{< code file=layouts/_default/home.html >}}
{{ define "main" }}
  <main aria-role="main">
    <header class="home-page-header">
      <h1>{{ .Title }}</h1>
      {{ with .Params.subtitle }}
        <span class="subtitle">{{ . }}</span>
      {{ end }}
    </header>
    <div class="home-page-content">
      <!-- Note that the content for index.html, as a sort of list page, will pull from content/_index.md -->
      {{ .Content }}
    </div>
    <div>
      {{ range first 10 .Site.RegularPages }}
        {{ .Render "summary" }}
      {{ end }}
    </div>
  </main>
{{ end }}
{{< /code >}}

[contentorg]: /content-management/organization/
[lookup]: /templates/lookup-order/
