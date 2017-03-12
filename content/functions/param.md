---
title: param
linktitle: Param
description: Calls page or site variables into your template.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: []
categories: [functions]
toc:
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
wip: true
---

In Hugo, you can declare [site-wide params][sitevars] (i.e. in your [configuration][]), as well as params for [individual pages][pagevars].

A common use case is to have a general value for the site and a more specific value for some of the pages (e.g., an image).

You can use the `.Param` method to call these values into your template:

```
$.Param "image"
```

[configuration]: /getting-started/configuration/
[pagevars]: /variables/page/
[sitevars]: /variables/site/