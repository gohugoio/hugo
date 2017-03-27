---
title: trim
linktitle:
description: Returns a slice of a passed string with all leading and trailing characters from cutset removed.
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: [strings]
ns:
signature: ["trim INPUT CUTLIST"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

`trim` returns a slice of the string with all leading and trailing characters contained in cutset removed.

```
{{ trim "++Batman--" "+-" }} → "Batman"
```
