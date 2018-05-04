---
title: absLangURL
description: Adds the absolute URL with correct language prefix according to site configuration for multilingual.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [multilingual,i18n,urls]
signature: ["absLangURL INPUT"]
workson: []
hugoversion:
relatedfuncs: [relLangURL]
deprecated: false
aliases: []
---

Both `absLangURL` and [`relLangURL`](/functions/rellangurl/) are similar to their [`absURL`](/functions/absurl/) and [`relURL`](/functions/relurl) relatives but will add the correct language prefix when the site is configured with more than one language.

So for a site  `baseURL` set to `https://example.com/hugo/` and the current language is `en`:

```
{{ "blog/" | absLangURL }} → "https://example.com/hugo/en/blog/"
{{ "blog/" | relLangURL }} → "/hugo/en/blog/"
```
