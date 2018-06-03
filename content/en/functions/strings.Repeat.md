---
title: strings.Repeat
# linktitle:
description: Returns a string consisting of count copies of the string s.
godocref:
date: 2018-05-31
publishdate: 2018-05-31
lastmod: 2018-05-31
categories: [functions]
menu:
  docs:
    parent: "functions"
keywords: [strings]
signature: ["strings.Repeat INPUT COUNT"]
workson: []
hugoversion:
relatedfuncs: []
deprecated: false
---

`strings.Repeat` provides the Go [`strings.Repeat`](https://golang.org/pkg/strings/#Repeat) function for Hugo templates. It takes a string and a count, and returns a string with consisting of count copies of the string argument.

```
{{ strings.Repeat "yo" 3 }} → "yoyoyo"
```

`strings.Repeat` *requires* the second argument, which tells the function how many times to repeat the first argument; there is no default. However, it can be used as a pipeline:

```
{{ "yo" | strings.Repeat 3 }} → "yoyoyo"
```
