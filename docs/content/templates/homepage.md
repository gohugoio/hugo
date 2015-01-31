---
aliases:
- /layout/homepage/
date: 2013-07-01
menu:
  main:
    parent: layout
next: /templates/terms
notoc: true
prev: /templates/list
title: Homepage
weight: 50
---

The home page of a website is often formatted differently than the other
pages. In Hugo you can define your own homepage template. 

Homepage is of the type "node" and have all the [node
variables](/templates/variables/) and [site
variables](/templates/variables/) available to use in the templates.

*This is the only required template for building a site and useful when
bootstrapping a new site and template. It is also the only required
template when using a single page site.*

In addition to the standard node variables, the homepage has access to
all site content accessible from `.Data.Pages`. Details on how to use the
list of pages can be found in the [Lists Template](/templates/list/).

## Which Template will be rendered?
Hugo uses a set of rules to figure out which template to use when
rendering a specific page.

Hugo will use the following prioritized list. If a file isn’t present,
then the next one in the list will be used. This enables you to craft
specific layouts when you want to without creating more templates
than necessary. For most sites, only the \_default file at the end of
the list will be needed.

* /layouts/index.html
* /layouts/\_default/list.html
* /layouts/\_default/single.html
* /themes/`THEME`/layouts/index.html
* /themes/`THEME`/layouts/\_default/list.html
* /themes/`THEME`/layouts/\_default/single.html

## Example index.html
This content template is used for [spf13.com](http://spf13.com/).

It makes use of [partial templates](/templates/partials/) and uses a similar approach as a [List](/templates/list/).

    <!DOCTYPE html>
    <html class="no-js" lang="en-US" prefix="og: http://ogp.me/ns# fb: http://ogp.me/ns/fb#">
    <head>
        <meta charset="utf-8">

        {{ partial "meta.html" . }}

        <base href="{{ .Site.BaseUrl }}">
        <title>{{ .Site.Title }}</title>
        <link rel="canonical" href="{{ .Permalink }}">
        <link href="{{ .RSSlink }}" rel="alternate" type="application/rss+xml" title="{{ .Site.Title }}" />

        {{ partial "head_includes.html" . }}
    </head>
    <body lang="en">

    {{ partial "subheader.html" . }}

    <section id="main">
      <div>
        {{ range first 10 .Data.Pages }}
            {{ .Render "summary"}}
        {{ end }}
      </div>
    </section>

    {{ partial "footer.html" . }}
