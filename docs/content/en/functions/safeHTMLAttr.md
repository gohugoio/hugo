---
title: safeHTMLAttr
# linktitle: safeHTMLAttr
description: Declares the provided string as a safe HTML attribute.
godocref: https://golang.org/src/html/template/content.go?s=1661:1676#L33
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

```
[[menu.main]]
    name = "IRC: #golang at freenode"
    url = "irc://irc.freenode.net/#golang"
```

* <span class="bad">`<a href="{{ .URL }}">` &rarr; `<a href="#ZgotmplZ">`</span>
* <span class="good">`<a {{ printf "href=%q" .URL | safeHTMLAttr }}>` &rarr; `<a href="irc://irc.freenode.net/#golang">`</span>

