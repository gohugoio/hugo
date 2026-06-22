---
title: urls.Anchorize
description: Returns the given string, sanitized for usage in an HTML id attribute.
categories: []
keywords: []
params:
  functions_and_methods:
    aliases: [anchorize]
    returnType: string
    signatures: [urls.Anchorize INPUT]
aliases: [/functions/anchorize]
---

{{% include "/_common/functions/urls/anchorize-vs-urlize.md" %}}

The `ursl.Anchorize` function sanitizes the resulting string per the [`autoIDType`][] setting in your project configuration.

[`autoIDType`]: /configuration/markup/#parserautoidtype
