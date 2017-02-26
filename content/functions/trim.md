---
title: trim
linktitle:
description:
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [strings]
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

Trim returns a slice of the string with all leading and trailing characters contained in cutset removed.

```
{{ trim "++Batman--" "+-" }} → "Batman"
```