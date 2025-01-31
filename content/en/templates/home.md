---
title: Home page templates
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

## Introduction

A home page template is used to render your site's home page, and is the only template required for a single-page website.  For example, the home page template below inherits the site's shell from the base template and renders the home page content, such as a list of other pages.

{{< code file=layouts/_default/home.html >}}
{{ define "main" }}
  {{ .Content }}
  {{ range site.RegularPages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
{{< /code >}}

{{% include "templates/_common/filter-sort-group.md" %}}

## Lookup order

Hugo's [template lookup order] determines the template path, allowing you to create unique templates for any page.

[template lookup order]: /templates/lookup-order/#home-templates

{{% note %}}
You must have thorough understanding of the template lookup order when creating templates. Template selection is based on template type, page kind, content type, section, language, and output format.
{{% /note %}}

## Content and front matter

The home page template uses content and front matter from an&nbsp;`_index.md`&nbsp;file located in the root of your content directory.

{{< code-toggle file=content/_index.md fm=true >}}
---
title: The Home Page
date: 2025-01-30T03:36:57-08:00
draft: false
params:
  subtitle: The Subtitle
---
{{< /code-toggle >}}

The home page template below inherits the site's shell from the base template, renders the subtitle and content as defined in the&nbsp;`_index.md`&nbsp;file, then renders of list of the site's [regular pages](g).

{{< code file=layouts/_default/home.html >}}
{{ define "main" }}
  <h3>{{ .Params.Subtitle }}</h3>
  {{ .Content }}
  {{ range site.RegularPages }}
    <h2><a href="{{ .RelPermalink }}">{{ .LinkTitle }}</a></h2>
  {{ end }}
{{ end }}
{{< /code >}}
