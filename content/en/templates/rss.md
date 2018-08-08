---
title: RSS Templates
linktitle: RSS Templates
description: Hugo ships with its own RSS 2.0 template that requires almost no configuration, or you can create your own RSS templates.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [rss, xml, templates]
categories: [templates]
menu:
  docs:
    parent: "templates"
    weight: 150
weight: 150
sections_weight: 150
draft: false
toc: true
---

## RSS Template Lookup Order

See [Template Lookup Order](/templates/lookup-order/) for the complete reference.

{{% note "Hugo Ships with an RSS Template" %}}
Hugo ships with its own [RSS 2.0 template](#the-embedded-rss-xml). The embedded template will be sufficient for most use cases.
{{% /note %}}

RSS pages are of the type `Page` and have all the [page variables](/variables/page/) available to use in the templates.

### Section RSS

A [section’s][section] RSS will be rendered at `/<SECTION>/index.xml` (e.g., http://spf13.com/project/index.xml).

Hugo provides the ability for you to define any RSS type you wish and can have different RSS files for each section and taxonomy.

## Lookup Order for RSS Templates

The table below shows the RSS template lookup order for the different page kinds. The first listing shows the lookup order when running with a theme (`demoTheme`).

{{< datatable-filtered "output" "layouts" "OutputFormat == RSS" "Example" "OutputFormat" "Suffix" "Template Lookup Order" >}}

## Configure RSS

By default, Hugo will create an unlimited number of RSS entries. You can limit the number of articles included in the built-in RSS templates by assigning a numeric value to `rssLimit:` field in your project's [`config` file][config].

The following values will also be included in the RSS output if specified in your site’s configuration:

```toml
languageCode = "en-us"
copyright = "This work is licensed under a Creative Commons Attribution-ShareAlike 4.0 International License."

[author]
    name = "My Name Here"
```

## The Embedded rss.xml

This is the default RSS template that ships with Hugo. It adheres to the [RSS 2.0 Specification][RSS 2.0].

```xml
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>{{ if eq  .Title  .Site.Title }}{{ .Site.Title }}{{ else }}{{ with .Title }}{{.}} on {{ end }}{{ .Site.Title }}{{ end }}</title>
    <link>{{ .Permalink }}</link>
    <description>Recent content {{ if ne  .Title  .Site.Title }}{{ with .Title }}in {{.}} {{ end }}{{ end }}on {{ .Site.Title }}</description>
    <generator>Hugo -- gohugo.io</generator>{{ with .Site.LanguageCode }}
    <language>{{.}}</language>{{end}}{{ with .Site.Author.email }}
    <managingEditor>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</managingEditor>{{end}}{{ with .Site.Author.email }}
    <webMaster>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</webMaster>{{end}}{{ with .Site.Copyright }}
    <copyright>{{.}}</copyright>{{end}}{{ if not .Date.IsZero }}
    <lastBuildDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 -0700" | safeHTML }}</lastBuildDate>{{ end }}
    {{ with .OutputFormats.Get "RSS" }}
        {{ printf "<atom:link href=%q rel=\"self\" type=%q />" .Permalink .MediaType | safeHTML }}
    {{ end }}
    {{ range .Pages }}
    <item>
      <title>{{ .Title }}</title>
      <link>{{ .Permalink }}</link>
      <pubDate>{{ .Date.Format "Mon, 02 Jan 2006 15:04:05 -0700" | safeHTML }}</pubDate>
      {{ with .Site.Author.email }}<author>{{.}}{{ with $.Site.Author.name }} ({{.}}){{end}}</author>{{end}}
      <guid>{{ .Permalink }}</guid>
      <description>{{ .Summary | html }}</description>
    </item>
    {{ end }}
  </channel>
</rss>
```

{{% warning "XML Header" %}}
Hugo will automatically add the following header line to this file on render. Please do *not* include this in the template as it's not valid HTML.
```
<?xml version="1.0" encoding="utf-8" standalone="yes" ?>
```
{{% /warning %}}

## Reference your RSS Feed in `<head>`

In your `header.html` template, you can specify your RSS feed in your `<head></head>` tag using Hugo's [Output Formats][Output Formats] like this:

```go-html-template
{{ range .AlternativeOutputFormats -}}
    {{ printf `<link rel="%s" type="%s" href="%s" title="%s" />` .Rel .MediaType.Type .Permalink $.Site.Title | safeHTML }}
{{ end -}}
```

If you only want the RSS link, you can query the formats:

```go-html-template
{{ with .OutputFormats.Get "rss" -}}
    {{ printf `<link rel="%s" type="%s" href="%s" title="%s" />` .Rel .MediaType.Type .Permalink $.Site.Title | safeHTML }}
{{ end -}}
```

Either of the two snippets above will generate the below `link` tag on the site homepage for RSS output:

```html
<link rel="alternate" type="application/rss+xml" href="https://example.com/index.xml" title="Site Title">
```

_We are assuming `BaseURL` to be `https://example.com/` and `$.Site.Title` to be `"Site Title"` in this example._

[config]: /getting-started/configuration/
[embedded]: #the-embedded-rss-xml
[RSS 2.0]: http://cyber.law.harvard.edu/rss/rss.html "RSS 2.0 Specification"
[section]: /content-management/sections/
[Output Formats]: /templates/output-formats/#link-to-output-formats
