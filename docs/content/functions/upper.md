---
title: upper
linktitle: upper
description: Converts all characters in a string to uppercase
godocref:
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
tags: []
categories: [functions]
toc:
ns:
signature: ["upper INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
draft: false
aliases: []
---

`upper` converts all characters in string to uppercase. Note that `upper` can be applied in your templates in more than one way:

```
{{ upper "BatMan" }} → "BATMAN"
{{ "BatMan" | upper }} → "BATMAN"
```

