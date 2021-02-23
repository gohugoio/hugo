---
title: print
linktitle: print
description: Prints the default representation of the given arguments using the standard `fmt.Print` function.
godocref: https://golang.org/pkg/fmt/
date: 2017-02-01
publishdate: 2017-02-01
lastmod: 2017-02-01
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings]
signature: ["print INPUT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

See [the go doc](https://golang.org/pkg/fmt/) for additional information.

```
{{ print "foo" }} → "foo"
{{ print "foo" "bar" }} → "foobar"
{{ print (slice 1 2 3) }} → [1 2 3]
```
