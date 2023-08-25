---
title: printf
description: Formats a string using the standard `fmt.Sprintf` function.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: [strings]
signature: ["printf FORMAT INPUT"]
relatedfuncs: []
---

See [the go doc](https://golang.org/pkg/fmt/) for additional information.

```go-html-template
{{ i18n ( printf "combined_%s" $var ) }}
```

```go-html-template
{{ printf "formatted %.2f" 3.1416 }}
```
