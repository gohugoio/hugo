---
title: warnf
description: Log a WARNING from a template.
categories: [functions]
menu:
  docs:
    parent: functions
keywords: []
namespace: fmt
relatedFuncs:
  - fmt.Errorf
  - fmt.Erroridf
  - fmt.Warnf
signature: 
  - fmt.Warnf FORMAT [INPUT]
  - warnf FORMAT [INPUT]
---

The documentation for [Go's fmt package] describes the structure and content of the format string.

Like the  [`printf`] function, the `warnf` function evaluates the format string. It then prints the result to the WARNING log. Hugo prints each unique message once to avoid flooding the log with duplicate warnings.

```go-html-template
{{ warnf "Copyright notice missing from site configuration" }}
```

[`printf`]: /functions/printf/
[Go's fmt package]: https://pkg.go.dev/fmt
