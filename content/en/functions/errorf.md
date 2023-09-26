---
title: errorf
description: Log an ERROR from a template.
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
  - fmt.Errorf FORMAT [INPUT]
  - errorf FORMAT [INPUT]
---

The documentation for [Go's fmt package] describes the structure and content of the format string.

Like the  [`printf`] function, the `errorf` function evaluates the format string. It then prints the result to the ERROR log and fails the build. Hugo prints each unique message once to avoid flooding the log with duplicate errors.

```go-html-template
{{ errorf "The %q shortcode requires a src parameter. See %s" .Name .Position }}
```

Use the [`erroridf`] function to allow optional supression of specific errors.

[`erroridf`]: /functions/erroridf/
[`printf`]: /functions/printf/
[Go's fmt package]: https://pkg.go.dev/fmt
