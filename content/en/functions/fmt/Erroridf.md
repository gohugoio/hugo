---
title: fmt.Erroridf
description: Log a suppressible ERROR from a template.
categories: []
keywords: []
action:
  aliases: [erroridf]
  related:
    - functions/fmt/Errorf
    - functions/fmt/Warnf
    - functions/fmt/Warnidf
  returnType: string
  signatures: ['fmt.Erroridf ID FORMAT [INPUT]']
aliases: [/functions/erroridf]
---

{{% include "functions/fmt/_common/fmt-layout.md" %}}

The `erroridf` function evaluates the format string, then prints the result to the ERROR log and fails the build. Unlike the [`errorf`] function, you may suppress errors logged by the `erroridf` function by adding the message ID to the `ignoreLogs` array in your site configuration.

This template code:

```go-html-template
{{ erroridf "error-42" "You should consider fixing this." }}
```

Produces this console log:

```text
ERROR You should consider fixing this.
You can suppress this error by adding the following to your site configuration:
ignoreLogs = ['error-42']
```

To suppress this message:

{{< code-toggle file=hugo >}}
ignoreLogs = ["error-42"]
{{< /code-toggle >}}

[`errorf`]: /functions/fmt/errorf/
