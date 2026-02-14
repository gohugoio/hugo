---
title: Role
description: Returns the role object for the given site.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: roles.RoleSite
    signatures: [SITE.Role]
---

{{< new-in 0.153.0 />}}

The `Role` method on a `Site` object returns the role object for the given site. The role object is derived from the role definition in the site configuration.

## Methods

### IsDefault

(`bool`) Reports whether this is the default role object as defined by the [`defaultContentRole`][] setting in the site configuration.

```go-html-template
{{ .Site.Role.IsDefault }} → true
```

### Name

(`string`) Returns the role name. This is the lower cased key from the site configuration.

```go-html-template
{{ .Site.Role.Name }} → guest
```

[`defaultContentRole`]: /configuration/all/#defaultcontentrole
