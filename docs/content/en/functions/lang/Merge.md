---
title: lang.Merge
description: Merge missing translations from other languages.
categories: []
keywords: []
action:
  aliases: []
  related: []
  returnType: any
  signatures: [lang.Merge FROM TO]
aliases: [/functions/lang.merge]
---

As an example:

```sh
{{ $pages := .Site.RegularPages | lang.Merge $frSite.RegularPages | lang.Merge $enSite.RegularPages }}
```

Will "fill in the gaps" in the current site with, from left to right, content from the French site, and lastly the English.

A more practical example is to fill in the missing translations from the other languages:

```sh
{{ $pages := .Site.RegularPages }}
{{ range .Site.Home.Translations }}
{{ $pages = $pages | lang.Merge .Site.RegularPages }}
{{ end }}
 ```
