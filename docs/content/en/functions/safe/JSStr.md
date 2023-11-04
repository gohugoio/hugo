---
title: safe.JSStr
linkTitle: safeJSStr
description: Declares the provided string as a known safe JavaScript string.
categories: [functions]
keywords: []
menu:
  docs:
    parent: functions
function:
  aliases: [safeJSStr]
  returnType: template.JSStr
  signatures: [safe.JSStr INPUT]
relatedFunctions:
  - safe.CSS
  - safe.HTML
  - safe.HTMLAttr
  - safe.JS
  - safe.JSStr
  - safe.URL
aliases: [/functions/safejsstr]
---

Encapsulates a sequence of characters meant to be embedded between quotes in a JavaScript expression. Use of this type presents a security risk: the encapsulated content should come from a trusted source, as it will be included verbatim in the template output.
  
Without declaring a variable to be a safe JavaScript string:

```go-html-template
{{ $title := "Lilo & Stitch" }}
<script>
  const a = "Title: " + {{ $title }};
</script>
```

Rendered:


```html
<script>
  const a = "Title: " + "Lilo \u0026 Stitch";
</script>
```

To avoid escaping by Go's [html/template] package:

```go-html-template
{{ $title := "Lilo & Stitch" }}
<script>
  const a = "Title: " + {{ $title | safeJSStr }};
</script>
```

Rendered:

```html
<script>
  const a = "Title: " + "Lilo & Stitch";
</script>
```

[html/template]: https://pkg.go.dev/html/template
