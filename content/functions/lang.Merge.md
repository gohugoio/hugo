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


A more practical example is to fill in the missing translations for the "minority languages" with content from the main language:


```bash
 {{ $pages := .Site.RegularPages }}
 {{ .Scratch.Set "pages" $pages }}
 {{ $mainSite := .Sites.First }}
 {{ if ne $mainSite .Site }}
    {{ .Scratch.Set "pages" ($pages | lang.Merge $mainSite.RegularPages) }}
 {{ end }}
 {{ $pages := .Scratch.Get "pages" }} 
 ```

{{% note %}}
Note that the slightly ugly `.Scratch` construct will not be needed once this is fixed: https://github.com/golang/go/issues/10608
{{% /note %}}
