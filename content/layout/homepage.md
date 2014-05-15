---
title: "Homepage Templates"
date: "2013-07-01"
weight: 30
notoc: true
menu:
  main:
    parent: 'layout'
---

Home pages are of the type "node" and have all the [node
variables](/layout/variables/) available to use in the templates.

*This is the only required template for building a site and useful when
bootstrapping a new site and template.*

In addition to the standard node variables, the homepage has access to
all site content accessible from .Data.Pages . Details on how to use this 
list of pages can be found in [Lists](/indexes/lists/)


    ▾ layouts/
        index.html


## example index.html
This content template is used for [spf13.com](http://spf13.com).

It makes use of [chrome templates](/layout/chrome) and uses a [List](/indexes/lists/)

    <!DOCTYPE html>
    <html class="no-js" lang="en-US" prefix="og: http://ogp.me/ns# fb: http://ogp.me/ns/fb#">
    <head>
        <meta charset="utf-8">

        {{ template "chrome/meta.html" . }}

        <base href="{{ .Site.BaseUrl }}">
        <title>{{ .Site.Title }}</title>
        <link rel="canonical" href="{{ .Permalink }}">
        <link href="{{ .RSSlink }}" rel="alternate" type="application/rss+xml" title="{{ .Site.Title }}" />

        {{ template "chrome/head_includes.html" . }}
    </head>
    <body lang="en">

    {{ template "chrome/subheader.html" . }}

    <section id="main">
      <div>
        {{ range first 10 .Data.Pages }}
            {{ .Render "summary"}}
        {{ end }}
      </div>
    </section>

    {{ template "chrome/footer.html" }}
