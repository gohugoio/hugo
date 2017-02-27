---
title: abslangurl and rellangurl
linktitle: absLangURL and relLangURL
description: Similar to absURL, but adds the absolute URL with correct language prefix according to site configuration for multilingual and baseURL.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: [multilingual,i18n,urls]
categories: [functions]
signature:
workson: []
hugoversion:
relatedfuncs: [relLangURL]
deprecated: false
aliases: [/functions/absurl/,/functions/rellangurl/]
needsexamples: true
---

These functions are similar to their [`absURL` and `relURL`](/functions/absurl-and-relurl) relatives but will add the correct language prefix when the site is configured with more than one language.

So for a site  `baseURL` set to `http://yoursite.com/hugo/` and the current language is `en`:

```golang
{{ "blog/" | absLangURL }} → "http://yoursite.com/hugo/en/blog/"
{{ "blog/" | relLangURL }} → "/hugo/en/blog/"
```