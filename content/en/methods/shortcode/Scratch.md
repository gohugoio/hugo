---
title: Scratch
description: Returns a "scratch pad" scoped to the shortcode to store and manipulate data. 
categories: []
keywords: []
action:
  related:
    - functions/collections/NewScratch
  returnType: maps.Scratch
  signatures: [SHORTCODE.Scratch]
expiryDate: 2025-11-18 #  deprecated 2024-11-18
---

{{% deprecated-in 0.139.0 %}}
Use the [`SHORTCODE.Store`] method instead.

This is a soft deprecation. This method will be removed in a future release, but the removal date has not been established. Although Hugo will not emit a warning if you continue to use this method, you should begin using `SHORTCODE.Store` as soon as possible.

Beginning with v0.139.0 the `SHORTCODE.Scratch` method is aliased to `SHORTCODE.Store`.

[`SHORTCODE.Store`]: /methods/shortcode/store/
{{% /deprecated-in %}}

The `Scratch` method within a shortcode creates a [scratch pad] to store and manipulate data. The scratch pad is scoped to the shortcode.

{{% note %}}
With the introduction of the [`newScratch`] function, and the ability to [assign values to template variables] after initialization, the `Scratch` method within a shortcode is obsolete.

[assign values to template variables]: https://go.dev/doc/go1.11#text/template
[`newScratch`]: /functions/collections/newscratch/
{{% /note %}}

[scratch pad]: /getting-started/glossary/#scratch-pad

{{% include "methods/page/_common/scratch-methods.md" %}}
