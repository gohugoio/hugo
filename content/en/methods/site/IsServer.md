---
title: IsServer
description: Reports whether the built-in development server is running.
categories: []
keywords: []
action:
  related: []
  returnType: bool
  signatures: [SITE.IsServer]
expiryDate: 2024-10-30 # deprecated 2023-10-30
---

{{% deprecated-in 0.120.0 %}}
Use [`hugo.IsServer`] instead.

[`hugo.IsServer`]: /functions/hugo/isserver/
{{% /deprecated-in %}}

```go-html-template
{{ .Site.IsServer }} → true/false
```
