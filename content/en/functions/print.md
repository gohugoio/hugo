---
title: print
description: Prints the default representation of the given arguments using the standard `fmt.Print` function.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings]
signature: ["print INPUT"]
relatedfuncs: []
---

See [the go doc](https://golang.org/pkg/fmt/) for additional information.

```go-html-template
{{ print "foo" }} → "foo"
{{ print "foo" "bar" }} → "foobar"
{{ print (slice 1 2 3) }} → [1 2 3]
```
