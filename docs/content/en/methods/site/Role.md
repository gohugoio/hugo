---
title: Role
description: Returns the Role object for the given site.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: roles.Role
    signatures: [SITE.Role]
---

{{< new-in 0.153.0 />}}

The `Role` method on a `Site` object returns the `Role` object for the given site, derived from the role definition in your project configuration.

## Methods

### IsDefault

(`bool`) Reports whether this is the [default role][].

```go-html-template
{{ .Site.Role.IsDefault }} → true
```

### Name

(`string`) Returns the role name. This is the lowercased key from your project configuration.

```go-html-template
{{ .Site.Role.Name }} → guest
```

[default role]: /quick-reference/glossary/#default-role
