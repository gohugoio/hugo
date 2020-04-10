---
title: lang.Merge
description: "Merge missing translations from other languages."
godocref: ""
workson: []
date: 2018-03-16
categories: [functions]
keywords: [multilingual]
menu:
  docs:
    parent: "functions"
toc: false
signature: ["lang.Merge FROM TO"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
comments:
---

As an example:

```bash
{{ $pages := .Site.RegularPages | lang.Merge $frSite.RegularPages | lang.Merge $enSite.RegularPages }}
```

Will "fill in the gaps" in the current site with, from left to right, content from the French site, and lastly the English.


A more practical example is to fill in the missing translations from the other languages:

```bash
{{ $pages := .Site.RegularPages }}
{{ range .Site.Home.Translations }}
{{ $pages = $pages | lang.Merge .Site.RegularPages }}
{{ end }}
 ```
