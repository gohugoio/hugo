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


Attempting to use the `url` value directly in an attribute like this:

- `<a href="{{ .URL }}"></a>` will produce the following: `<a href="#ZgotmplZ"></a>`.

The `ZgotmplZ` value indicates that you're trying to output content at a spot
where `template/html` deems to be unsafe. To correct the output, use the
`safeHTMLAttr` function like so:

- `<a {{ printf "href=%q" .URL | safeHTMLAttr }}></a>` which produces: `<a href="irc://irc.freenode.net/#golang"></a>`
