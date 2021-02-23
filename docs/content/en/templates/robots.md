---
title: Robots.txt File
linktitle: Robots.txt
description: Hugo can generate a customized robots.txt in the same way as any other template.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
keywords: [robots,search engines]
menu:
  docs:
    parent: "templates"
    weight: 165
weight: 165
sections_weight: 165
draft: false
aliases: [/extras/robots-txt/]
toc: false
---

To create your robots.txt as a template, first set the `enableRobotsTXT` value to `true` in your [configuration file][config]. By default, this option generates a robots.txt with the following content, which tells search engines that they are allowed to crawl everything:

```
User-agent: *
```

## Robots.txt Template Lookup Order

The [lookup order][lookup] for the `robots.txt` template is as follows:

* `/layouts/robots.txt`
* `/themes/<THEME>/layouts/robots.txt`

{{% note %}}
If you do not want Hugo to create a default `robots.txt` or leverage the `robots.txt` template, you can hand code your own and place the file in `static`. Remember that everything in the [static directory](/getting-started/directory-structure/) is copied over as-is when Hugo builds your site.
{{% /note %}}

## Robots.txt Template Example

The following is an example `robots.txt` layout:

{{< code file="layouts/robots.txt" download="robots.txt" >}}
User-agent: *

{{range .Pages}}
Disallow: {{.RelPermalink}}
{{end}}
{{< /code >}}

This template disallows all the pages of the site by creating one `Disallow` entry for each page.

[config]: /getting-started/configuration/
[lookup]: /templates/lookup-order/
[robots]: https://www.robotstxt.org/
