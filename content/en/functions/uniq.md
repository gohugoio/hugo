---
title: uniq
linktitle: uniq
description: Takes in a slice or array and returns a slice with subsequent duplicate elements removed.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [multilingual,i18n,urls]
signature: ["uniq SET"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
aliases: []
needsexamples: false
---

```
{{ uniq (slice 1 2 3 2) }}
{{ slice 1 2 3 2 | uniq }}
<!-- both return [1 2 3] -->
```




