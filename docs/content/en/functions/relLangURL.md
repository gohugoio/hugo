---
title: relLangURL
description: Adds the relative URL with correct language prefix according to site configuration for multilingual.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
keywords: [multilingual,i18n,urls]
categories: [functions]
menu:
  docs:
    parent: "functions"
signature: ["relLangURL INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
---

`absLangURL` and `relLangURL` functions are similar to their [`absURL`](/functions/absurl/) and [`relURL`](/functions/relurl/) relatives but will add the correct language prefix when the site is configured with more than one language. (See [Configuring Multilingual][multiliconfig].)

So for a site  `baseURL` set to `https://example.com/hugo/` and the current language is `en`:

```
{{ "blog/" | absLangURL }} → "https://example.com/hugo/en/blog/"
{{ "blog/" | relLangURL }} → "/hugo/en/blog/"
```

[multiliconfig]: /content-management/multilingual/#configuring-multilingual-mode
