---
title: Param
description: Returns the site parameter with the given key.
categories: []
keywords: []
action:
  related: []
  returnType: any
  signatures: [SITE.Param KEY]
---

The `Param` method on a `Site` object is a convenience method to return the value of a user-defined parameter in the site configuration.

{{< code-toggle file=hugo >}}
[params]
display_toc = true
{{< /code-toggle >}}


```go-html-template
{{ .Site.Param "display_toc" }} → true
```

The above is equivalent to either of these:

```go-html-template
{{ .Site.Params.display_toc }}
{{ index .Site.Params "display_toc" }}
```
