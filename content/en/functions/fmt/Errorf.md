---
title: fmt.Errorf
description: Log an ERROR from a template.
categories: []
keywords: []
action:
  aliases: [errorf]
  related:
    - functions/fmt/Erroridf
    - functions/fmt/Warnf
  returnType: string
  signatures: ['fmt.Errorf FORMAT [INPUT]']
aliases: [/functions/errorf]
---

The documentation for [Go's fmt package] describes the structure and content of the format string.

Like the  [`printf`] function, the `errorf` function evaluates the format string. It then prints the result to the ERROR log and fails the build. Hugo prints each unique message once to avoid flooding the log with duplicate errors.

```go-html-template
{{ errorf "The %q shortcode requires a src parameter. See %s" .Name .Position }}
```

Use the [`erroridf`] function to allow optional suppression of specific errors.

[`erroridf`]: /functions/fmt/erroridf
[`printf`]: /functions/fmt/printf
[Go's fmt package]: https://pkg.go.dev/fmt
