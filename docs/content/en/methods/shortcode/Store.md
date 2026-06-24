---
title: Store
description: Returns a persistent data structure for storing and manipulating keyed values, scoped to the current shortcode.
categories: []
keywords: []
params:
  functions_and_methods:
    returnType: maps.Scratch
    signatures: [SHORTCODE.Store]
---

{{< new-in 0.139.0 />}}

Use the `Store` method to create a persistent data structure for storing and manipulating keyed values, scoped to the current shortcode. To create a data structure with a different [scope](g), refer to the [scope](#scope) section below.

> [!NOTE]
> With the introduction of the [`newScratch`][] function, and the ability to [assign values to template variables][] after initialization, the `Store` method within a shortcode is mostly obsolete.

{{% include "_common/store-methods.md" %}}

{{% include "_common/store-scope.md" %}}

[`newScratch`]: /functions/collections/newScratch/
[assign values to template variables]: https://go.dev/doc/go1.11#texttemplatepkgtexttemplate
