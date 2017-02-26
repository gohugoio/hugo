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
needsreview: true
---

Hugo can generated a customized [robots.txt](http://www.robotstxt.org/) in the
[same way as any other templates]({{< ref "templates/go-templates.md" >}}).

To enable it, just set `enableRobotsTXT` option to `true` in the [configuration file]({{< ref "overview/configuration.md" >}}). By default, it generates a robots.txt, which allows everything, with the following content:

```http
User-agent: *
```

## Robots.txt Template Lookup Order

The [lookup order][lookup] for the `robots.txt` template is as follows:

* `/layouts/robots.txt`
* `/themes/<THEME>/layout/robots.txt`

An example of a `robots.txt` layout is:

```http
User-agent: *

{{range .Data.Pages}}
Disallow: {{.RelPermalink}}{{end}}
```

This template disallows all the pages of the site creating one `Disallow` entry for each one.

[lookup]: /layouts/lookup-order
