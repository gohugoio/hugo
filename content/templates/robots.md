---
lastmod: 2016-01-06
date: 2015-12-08
menu:
  main:
    parent: extras
next: /community/mailing-list
prev: /extras/urls
title: Custom robots.txt
weight: 120
---

Hugo can generated a customized [robots.txt](http://www.robotstxt.org/) in the
[same way as any other templates]({{< ref "templates/go-templates.md" >}}).

To enable it, just set `enableRobotsTXT` option to `true` in the [configuration file]({{< ref "overview/configuration.md" >}}). By default, it generates a robots.txt, which allows everything, with the following content:

```http
User-agent: *
```


Hugo will use the template `robots.txt` according to the following list in descending precedence:

* /layouts/robots.txt
* /themes/`THEME`/layout/robots.txt

An example of a robots.txt layout is:

```http
User-agent: *

{{range .Data.Pages}}
Disallow: {{.RelPermalink}}{{end}}
```

This template disallows and all the pages of the site creating one `Disallow` entry for each one.
