---
title: printf
linktitle: printf
description:
godocref: https://golang.org/pkg/fmt/
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
tags: []
signature:
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
needsexamples: true
---

Format a string using the standard `fmt.Sprintf` function. See [the go
doc](https://golang.org/pkg/fmt/) for additional information.

```golang
{{ i18n ( printf "combined_%s" $var ) }}
```

```
{{ printf "formatted %.2f" 3.1416 }}
```