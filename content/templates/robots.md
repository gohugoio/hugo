---
title: Robots.txt File
linktitle: Robots.txt
description:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [templates]
tags: [robots,search engines]
weight: 165
draft: false
aliases: [/extras/robots-txt/]
toc: false
wip: true
---

Hugo can generate a customized [robots.txt][robots] in the same way as any other template.

To enable creating your robots.txt as a template, set the `enableRobotsTXT` value to `true` in your [project's configuration file][config]. By default, this option generates a robots.txt with the following content, which tells search engines that they are allowed to crawl everything:

```http
User-agent: *
```

## Robots.txt Template Lookup Order

The [lookup order][lookup] for the `robots.txt` template is as follows:

* `/layouts/robots.txt`
* `/themes/<THEME>/layout/robots.txt`

## Robots. txt Template Example

The following is an example`robots.txt` layout:

{{% code file="layouts/robots.txt" download="robots.txt" %}}
```http
User-agent: *

{{range .Data.Pages}}
Disallow: {{.RelPermalink}}
{{end}}
```
{{% /code %}}

This template disallows all the pages of the site by creating one `Disallow` entry for each page.

[config]: /getting-started/configuration/
[lookup]: /layouts/lookup-order
[robots]: http://www.robotstxt.org/