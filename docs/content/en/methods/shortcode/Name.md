---
title: Name
description: Returns the shortcode file name, excluding the file extension.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [SHORTCODE.Name]
---

The `Name` method is useful for error reporting. For example, if your shortcode requires a "greeting" argument:

```go-html-template {file="layouts/_shortcodes/myshortcode.html"}
{{ $greeting := "" }}
{{ with .Get "greeting" }}
  {{ $greeting = . }}
{{ else }}
  {{ errorf "The %q shortcode requires a 'greeting' argument. See %s" .Name .Position }}
{{ end }}
```

In the absence of a "greeting" argument, Hugo will throw an error message and fail the build:

```text
ERROR The "myshortcode" shortcode requires a 'greeting' argument. See "/home/user/project/content/about.md:11:1"
```
