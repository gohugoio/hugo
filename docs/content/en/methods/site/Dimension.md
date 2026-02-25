---
title: Dimension
description: Returns the dimension object for the given dimension for the given site.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: page.SiteDimension
    signatures: [SITE.Dimension DIMENSION]
---

{{< new-in 0.153.0 />}}

The `Dimension` method on a `Site` object returns the dimension object for the given [dimension](g).

The `DIMENSION` argument must be one of `language`, `version`, or `role`.

Example|Returns|Equivalent to
:--|:--|:--
`{{ .Site.Dimension "language" }}`|`langs.Language`|`{{ .Site.Language }}`
`{{ .Site.Dimension "version" }}`|`version.Version`|`{{ .Site.Version }}`
`{{ .Site.Dimension "role" }}`|`roles.Role`|`{{ .Site.Role }}`

```go-html-template
{{ $languageObject := .Site.Dimension "language" }}
{{ $languageObject.IsDefault }} → true
{{ $languageObject.Name }} → en

{{ $versionObject := .Site.Dimension "version" }}
{{ $versionObject.IsDefault }} → true
{{ $versionObject.Name }} → v1.0.0

{{ $roleObject := .Site.Dimension "role" }}
{{ $roleObject.IsDefault }} → true
{{ $roleObject.Name }} → guest
```
