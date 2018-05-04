---
date: 2017-04-16T13:53:58-04:00
categories: ["Releases"]
description: "Hugo 0.20.2 adds support for plain text partials included into HTML templates"
link: ""
title: "Hugo 0.20.2"
draft: false
author: bep
aliases: [/0-20-2/]
---

Hugo `0.20.2` adds support for plain text partials included into `HTML` templates. This was a side-effect of the big new [Custom Output Format](https://gohugo.io/extras/output-formats/) feature in `0.20`, and while the change was intentional and there was an ongoing discussion about fixing it in [#3273](//github.com/gohugoio/hugo/issues/3273), it did break some themes. There were valid workarounds for these themes, but we might as well get it right.

The most obvious use case for this is inline `CSS` styles, which you now can do without having to name your partials with a `html` suffix.

A simple example:

In `layouts/partials/mystyles.css`:

    body {
    	background-color: {{ .Param "colors.main" }}
    }

Then in `config.toml` (note that by using the `.Param` lookup func, we can override the color in a page’s front matter if we want):

    [params]
    [params.colors]
    main = "green"
    text = "blue"

And then in `layouts/partials/head.html` (or the partial used to include the head section into your layout):

    <head>
        <style type="text/css">
        {{ partial "mystyles.css" . | safeCSS }}
        </style>
    </head>

Of course, `0.20` also made it super-easy to create external `CSS` stylesheets based on your site and page configuration. A simple example:

Add “CSS” to your home page’s `outputs` list, create the template `/layouts/index.css` using Go template syntax for the dynamic parts, and then include it into your `HTML` template with:

    {{ with  .OutputFormats.Get "css" }}
    <link rel="{{ .Rel }}" type="{{ .MediaType.Type }}" href="{{ .Permalink |  safeURL }}">
    {{ end }}`