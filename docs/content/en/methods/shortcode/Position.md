---
title: Position
description: Returns the filename and position from which the shortcode was called.
categories: []
keywords: []
action:
  related:
    - methods/shortcode/Name
    - functions/fmt/Errorf
  returnType: text.Position
  signatures: [SHORTCODE.Position]
---

The `Position` method is useful for error reporting. For example, if your shortcode requires a "greeting" argument:

{{< code file=layouts/shortcodes/myshortcode.html  >}}
{{ $greeting := "" }}
{{ with .Get "greeting" }}
  {{ $greeting = . }}
{{ else }}
  {{ errorf "The %q shortcode requires a 'greeting' argument. See %s" .Name .Position }}
{{ end }}
{{< /code >}}

In the absence of a "greeting" argument, Hugo will throw an error message and fail the build:

```text
ERROR The "myshortcode" shortcode requires a 'greeting' argument. See "/home/user/project/content/about.md:11:1"
```

{{% note %}}
The position can be expensive to calculate. Limit its use to error reporting.
{{% /note %}}
