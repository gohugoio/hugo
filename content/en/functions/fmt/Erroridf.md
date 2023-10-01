---
title: fmt.Erroridf
linkTitle: erroridf
description: Log a suppressable ERROR from a template.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [erroridf]
  returnType: string
  signatures: ['fmt.Erroridf ID FORMAT [INPUT]']
relatedFunctions:
  - fmt.Errorf
  - fmt.Erroridf
  - fmt.Warnf
aliases: [/functions/erroridf]
---

The documentation for [Go's fmt package] describes the structure and content of the format string.

Like the  [`errorf`] function, the `erroridf` function evaluates the format string, prints the result to the ERROR log, then fails the build. Hugo prints each unique message once to avoid flooding the log with duplicate errors.

Unlike the `errorf` function, you may surpress errors logged by the `erroridf` function by adding the messsage ID to the `ignoreErrors` array in your site configuration.

This template code:

```go-html-template
{{ erroridf "error-42" "You should consider fixing this." }}
```

Produces this console log:

```text
ERROR You should consider fixing this.
If you feel that this should not be logged as an ERROR, you can ignore it by adding this to your site config:
ignoreErrors = ["error-42"]
```

To suppress this message:

{{< code-toggle file=hugo copy=false >}}
ignoreErrors = ["error-42"]
{{< /code-toggle >}}

[`errorf`]: /functions/fmt/errorf
[Go's fmt package]: https://pkg.go.dev/fmt
