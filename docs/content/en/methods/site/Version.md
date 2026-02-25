---
title: Version
description: Returns the Version object for the given site.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: versions.Version
    signatures: [SITE.Version]
---

{{< new-in 0.153.0 />}}

The `Version` method on a `Site` object returns the `Version` object for the given site, derived from the version definition in your project configuration.

## Methods

### IsDefault

(`bool`) Reports whether this is the [default version][].

```go-html-template
{{ .Site.Version.IsDefault }} → true
```

### Name

(`string`) Returns the version name. This is the lowercased key from your project configuration.

```go-html-template
{{ .Site.Version.Name }} → v1.0.0
```

[default version]: /quick-reference/glossary/#default-version
