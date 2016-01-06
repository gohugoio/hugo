---
lastmod: 2015-12-08
date: 2013-07-09
menu:
  main:
    parent: extras
next: /community/mailing-list
prev: /extras/urls
title: Table of Contents
weight: 120
---

Hugo can generated customized [robots.txt](http://www.robotstxt.org/) in the
[same way than any other template]({{< ref "templates/go-templates.md" >}}).

By default it generates a robots.txt which allows everything, it looks exactly

  User-agent: *

To disable it just set `disableRobotsTXT` option to false in the [command line]({{< ref "commands/hugo.md" >}}) or [configuration file]({{< ref "overview/configuration.md" >}}).

Hugo will use the template `robots.txt` following the list starting with the one with more priority

* /layouts/robots.txt
* /themes/`THEME`/layout/robots.txt

An example of a robots.txt layout is:

    User-agent: *

    {{range .Data.Pages}}
    Disallow: {{.RelPermalink}}{{end}}


This template disallows and all the pages of the site creating one `Disallow` entry for each one.
