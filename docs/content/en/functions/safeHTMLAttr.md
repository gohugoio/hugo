---
title: safeHTMLAttr
# linktitle: safeHTMLAttr
description: Declares the provided string as a safe HTML attribute.
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings]
signature: ["safeHTMLAttr INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

Example: Given a site-wide `config.toml` that contains this menu entry:

{{< code-toggle file="config" >}}
[[menu.main]]
    name = "IRC: #golang at freenode"
    url = "irc://irc.freenode.net/#golang"
{{< /code-toggle >}}

* <span class="bad">`<a href="{{ .URL }}">` &rarr; `<a href="#ZgotmplZ">`</span>
* <span class="good">`<a {{ printf "href=%q" .URL | safeHTMLAttr }}>` &rarr; `<a href="irc://irc.freenode.net/#golang">`</span>
