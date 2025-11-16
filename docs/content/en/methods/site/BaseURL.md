---
title: BaseURL
description: Returns the base URL as defined in the site configuration.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: string
    signatures: [SITE.BaseURL]
---

Site configuration:

{{< code-toggle file=hugo >}}
baseURL = 'https://example.org/docs/'
{{< /code-toggle >}}

Template:

```go-html-template
{{ .Site.BaseURL }} → https://example.org/docs/
```

> [!note]
> There is almost never a good reason to use this method in your templates. Its usage tends to be fragile due to misconfiguration.
>
> Use the [`absURL`], [`absLangURL`], [`relURL`], or [`relLangURL`] functions instead.

[`absLangURL`]: /functions/urls/absLangURL/
[`absURL`]: /functions/urls/absURL/
[`relLangURL`]: /functions/urls/relLangURL/
[`relURL`]: /functions/urls/relURL/
