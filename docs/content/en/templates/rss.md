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

A [section’s][section] RSS will be rendered at `/<SECTION>/index.xml` (e.g., https://spf13.com/project/index.xml).

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

This is the default RSS template that ships with Hugo:

https://github.com/gohugoio/hugo/blob/master/tpl/tplimpl/embedded/templates/_default/rss.xml

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
[RSS 2.0]: https://cyber.harvard.edu/rss/rss.html "RSS 2.0 Specification"
[section]: /content-management/sections/
[Output Formats]: /templates/output-formats/#link-to-output-formats
