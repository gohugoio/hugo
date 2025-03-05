---
title: fmt.Warnidf
description: Log a suppressible WARNING from a template.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [warnidf]
    returnType: string
    signatures: ['fmt.Warnidf ID FORMAT [INPUT]']
aliases: [/functions/warnidf]
---

{{< new-in 0.123.0 />}}

{{% include "/_common/functions/fmt/format-string.md" %}}

The `warnidf` function evaluates the format string, then prints the result to the WARNING log. Unlike the [`warnf`] function, you may suppress warnings logged by the `warnidf` function by adding the message ID to the `ignoreLogs` array in your site configuration.

This template code:

```go-html-template
{{ warnidf "warning-42" "You should consider fixing this." }}
```

Produces this console log:

```text
WARN You should consider fixing this.
You can suppress this warning by adding the following to your site configuration:
ignoreLogs = ['warning-42']
```

To suppress this message:

{{< code-toggle file=hugo >}}
ignoreLogs = ["warning-42"]
{{< /code-toggle >}}

[`warnf`]: /functions/fmt/warnf/
