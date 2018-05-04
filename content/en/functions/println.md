---
title: println
linktitle: println
description: Prints the default representation of the given argument using the standard `fmt.Print` function and enforces a linebreak.
godocref: https://golang.org/pkg/fmt/
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings]
signature: ["println INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

See [the go doc](https://golang.org/pkg/fmt/) for additional information. `\n` denotes the linebreak but isn't printed in the templates as seen below:

```
{{ println "foo" }} → "foo\n"
```
