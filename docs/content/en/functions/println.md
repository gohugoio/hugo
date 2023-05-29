---
title: println
description: Prints the default representation of the given argument using the standard `fmt.Print` function and enforces a linebreak.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings]
signature: ["println INPUT"]
relatedfuncs: []
---

See [the go doc](https://golang.org/pkg/fmt/) for additional information. `\n` denotes the linebreak but isn't printed in the templates as seen below:

```go-html-template
{{ println "foo" }} → "foo\n"
```
