---
title: IsDevelopment
description: Reports whether the current running environment is “development”.
categories: []
keywords: []
action:
  related: []
  returnType: bool
  signatures: [SITE.IsDevelopment]
expiryDate: 2024-10-30 # deprecated 2023-10-30
---

{{% deprecated-in 0.120.0 %}}
Use [`hugo.IsDevelopment`] instead.

[`hugo.IsDevelopment`]: /functions/hugo/isdevelopment/
{{% /deprecated-in %}}

```go-html-template
{{ .Site.IsDevelopment }} → true/false
```
