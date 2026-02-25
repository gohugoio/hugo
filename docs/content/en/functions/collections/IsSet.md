---
title: collections.IsSet
description: Reports whether a specific key or index exists in the given map or slice.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [isset]
    returnType: bool
    signatures: [collections.IsSet MAP|SLICE KEY|INDEX]
aliases: [/functions/isset]
---

For example, consider this project configuration:

{{< code-toggle file=hugo >}}
[params]
showHeroImage = false
{{< /code-toggle >}}

If the value of `showHeroImage` is `true`, we can detect that it exists using either `if` or `with`:

```go-html-template
{{ if site.Params.showHeroImage }}
  {{ site.Params.showHeroImage }} → true
{{ end }}

{{ with site.Params.showHeroImage }}
  {{ . }} → true
{{ end }}
```

But if the value of `showHeroImage` is `false`, we can't use either `if` or `with` to detect its existence. In this case, you must use the `isset` function:

```go-html-template
{{ if isset site.Params "showheroimage" }}
  <p>The showHeroImage parameter is set to {{ site.Params.showHeroImage }}.<p>
{{ end }}
```

> [!note]
> When using the `isset` function you must reference the key using lower case. See the previous example.
