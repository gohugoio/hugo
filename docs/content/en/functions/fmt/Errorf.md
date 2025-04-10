---
title: fmt.Errorf
description: Log an ERROR from a template.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [errorf]
    returnType: string
    signatures: ['fmt.Errorf FORMAT [INPUT]']
aliases: [/functions/errorf]
---

{{% include "/_common/functions/fmt/format-string.md" %}}

The `errorf` function evaluates the format string, then prints the result to the ERROR log and fails the build.

```go-html-template
{{ errorf "The %q shortcode requires a src argument. See %s" .Name .Position }}
```

Use the [`erroridf`] function to allow optional suppression of specific errors.

[`erroridf`]: /functions/fmt/erroridf/
