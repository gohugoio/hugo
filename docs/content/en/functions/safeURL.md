---
title: safeURL
description: Declares the provided string as a safe URL or URL substring.
godocref: https://golang.org/pkg/html/template/#HTMLEscape
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [strings,urls]
categories: [functions]
menu:
  docs:
    parent: "functions"
signature: ["safeURL INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`safeURL` declares the provided string as a "safe" URL or URL substring (see [RFC 3986][]). A URL like `javascript:checkThatFormNotEditedBeforeLeavingPage()` from a trusted source should go in the page, but by default dynamic `javascript:` URLs are filtered out since they are a frequently exploited injection vector.

Without `safeURL`, only the URI schemes `http:`, `https:` and `mailto:` are considered safe by Go templates. If any other URI schemes (e.g., `irc:` and `javascript:`) are detected, the whole URL will be replaced with `#ZgotmplZ`. This is to "defang" any potential attack in the URL by rendering it useless.

The following examples use a [site `config.toml`][configuration] with the following [menu entry][menus]:

{{< code file="config.toml" copy="false" >}}
[[menu.main]]
    name = "IRC: #golang at freenode"
    url = "irc://irc.freenode.net/#golang"
{{< /code >}}

The following is an example of a sidebar partial that may be used in conjunction with the preceding front matter example:

{{< code file="layouts/partials/bad-url-sidebar-menu.html" copy="false" >}}
<!-- This unordered list may be part of a sidebar menu -->
<ul>
  {{ range .Site.Menus.main }}
  <li><a href="{{ .URL }}">{{ .Name }}</a></li>
  {{ end }}
</ul>
{{< /code >}}

This partial would produce the following HTML output:

{{< output file="bad-url-sidebar-menu-output.html" >}}
<!-- This unordered list may be part of a sidebar menu -->
<ul>
    <li><a href="#ZgotmplZ">IRC: #golang at freenode</a></li>
</ul>
{{< /output >}}

The odd output can be remedied by adding ` | safeURL` to our `.URL` page variable:

{{< code file="layouts/partials/correct-url-sidebar-menu.html" copy="false" >}}
<!-- This unordered list may be part of a sidebar menu -->
<ul>
    <li><a href="{{ .URL | safeURL }}">{{ .Name }}</a></li>
</ul>
{{< /code >}}

With the `.URL` page variable piped through `safeURL`, we get the desired output:

{{< output file="correct-url-sidebar-menu-output.html" >}}
<ul class="sidebar-menu">
    <li><a href="irc://irc.freenode.net/#golang">IRC: #golang at freenode</a></li>
</ul>
{{< /output >}}

[configuration]: /getting-started/configuration/
[menus]: /content-management/menus/
[RFC 3986]: http://tools.ietf.org/html/rfc3986
