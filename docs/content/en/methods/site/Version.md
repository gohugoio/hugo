---
title: Version
description: Returns the version object for the given site.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: versions.Version
    signatures: [SITE.Version]
---

{{< new-in 0.153.0 />}}

The `Version` method on a `Site` object returns the version object for the given site. The version object is derived from the version definition in the site configuration.

## Methods

### IsDefault

(`bool`) Reports whether this is the default version object as defined by the [`defaultContentVersion`][] setting in the site configuration.

```go-html-template
{{ .Site.Version.IsDefault }} → true
```

### Name

(`string`) Returns the version name. This is the lower cased key from the site configuration.

```go-html-template
{{ .Site.Version.Name }} → v1.0.0
```

[`defaultContentVersion`]: /configuration/all/#defaultcontentversion
